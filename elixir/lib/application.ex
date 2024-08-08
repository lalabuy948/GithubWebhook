defmodule WebhookServer.Application do
  use Application

  @impl true
  def start(_type, _args) do
    children = [
      {Plug.Cowboy, scheme: :http, plug: WebhookServer.Webhook, options: [port: 4567]}
    ]

    IO.puts("Starting server on port 4567...")

    opts = [strategy: :one_for_one, name: WebhookServer.Supervisor]
    {:ok, _pid} = Supervisor.start_link(children, opts)

    {:ok, self()}
  end
end
