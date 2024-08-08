defmodule WebhookServer.MixProject do
  use Mix.Project

  def project do
    [
      app: :webhook_server,
      version: "0.1.0",
      elixir: "~> 1.17",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      mod: {WebhookServer.Application, []},
      extra_applications: [:logger]
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      {:plug_cowboy, "~> 2.7.1"},
      {:plug, "~> 1.16.1"},
      {:jason, "~> 1.4"}
    ]
  end
end
