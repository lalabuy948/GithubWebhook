# GitHub Webhook Servers

Recieve webhook from GitHub on master branch commit and execute pipline.sh <- add your logic there.

## How to

go

```sh
go build
SECRET_TOKEN=6829d633f9...0f9e5cf9cb3bc17 ./WebhookServer
```
ruby

```sh
bundle install

SECRET_TOKEN=6829d633f9...0f9e5cf9cb3bc17 bundle exec ruby WebhookServer.rb
```

elixir

```sh
mix deps.get

SECRET_TOKEN=6829d633f9...0f9e5cf9cb3bc17 mix run --no-halt
```
