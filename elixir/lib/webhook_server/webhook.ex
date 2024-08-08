defmodule WebhookServer.Webhook do
  use Plug.Router
  import Plug.Conn
  import Plug.Crypto, only: [secure_compare: 2]

  @secret_token System.get_env("SECRET_TOKEN")
  @cwd File.cwd!()

  plug(:match)
  plug(:dispatch)

  post "/webhook" do
    {:ok, body, _conn} = Plug.Conn.read_body(conn)

    [_sha, signature | _] =
      get_req_header(conn, "x-hub-signature") |> List.first() |> String.split("=")

    if verify_signature(body, signature) do
      payload = Jason.decode!(body)
      commit_id = get_in(payload, ["head_commit", "id"])
      ref = get_in(payload, ["ref"])

      IO.puts("Webhook received... commit [#{commit_id}]")

      if ref == "refs/heads/master" do
        IO.puts("Starting pipeline...")
        IO.puts(@cwd <> "./pipeline.sh")

        {output, exit_code} = System.cmd("sh", ["./pipeline.sh"], stderr_to_stdout: true)

        IO.puts("Pipeline output: #{output}")
        IO.puts("Pipeline exit code: #{exit_code}")
      end

      send_resp(conn, 200, "OK")
    else
      send_resp(conn, 500, "Signatures didn't match!")
    end
  end

  defp verify_signature(_body, nil), do: false
  defp verify_signature(nil, _signature), do: false

  defp verify_signature(body, signature) do
    if is_binary(body) and is_binary(@secret_token) do
      # Generate the HMAC signature
      expected_signature =
        :crypto.mac(:hmac, :sha, @secret_token, body)
        |> Base.encode16(case: :lower)

      secure_compare(expected_signature, signature)
    else
      false
    end
  end
end
