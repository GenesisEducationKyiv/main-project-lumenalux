# **gses2-app BTC to UAH exchange API**

![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## Translation

- [Українська](README_ua.md).

## Contents

- [About](#about)
- [Installation](#installation)
- [Usage](#usage)
- [Description](#description)
- [How It Works](#how-it-works)
- [Architecture](#architecture-diagram)
- [Project tree](#project-tree)

## About

This is an API that provides the current exchange rate between Bitcoin and the Ukrainian Hryvnia (UAH). It allows users to subscribe to rate updates and receive those updates via email.

## Installation

1. **Clone the repository to your desired location:**

   ```bash
   git clone https://github.com/lumenalux/gses2-app.git gses2-app
   ```

   ```bash
   cd gses2-app
   ```

2. **Configure the environment variables:**

   The application uses a .env file for configuration. Copy the contents of .env.example into a new file named .env. Set up the following environment variables for the SMTP server and email settings:

   ```bash
   GSES2_APP_SMTP_HOST="<smtp server host>"

   GSES2_APP_SMTP_USER="<smtp username>"

   GSES2_APP_SMTP_PASSWORD="<smtp password>"`
   ```

   The rest of the environment variables have default values as listed below, but can be overridden if necessary:

   ```bash
   GSES2_APP_SMTP_PORT=465

   GSES2_APP_EMAIL_FROM=no.reply@test.info.api
   GSES2_APP_EMAIL_SUBJECT=BTC to UAH exchange rate
   GSES2_APP_EMAIL_BODY=The BTC to UAH exchange rate is {{.Rate}} UAH per BTC

   GSES2_APP_STORAGE_PATH=./storage/storage.csv

   GSES2_APP_HTTP_PORT=8080
   GSES2_APP_HTTP_TIMEOUT=10s

   GSES2_APP_KUNAAPI_URL=https://api.kuna.io/v3/tickers?symbols=btcuah

   GSES2_APP_BINANCEAPI_URL=https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT

   GSES2_APP_COINGECKOAPI_URL=https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah

   GSES2_APP_RABBITMQ_URL=amqp://guest:guest@amqp/
   ```

The environment variables include settings for the SMTP server and the content of the email messages sent to subscribers. The body of the email is designed as a template using Go's text/template syntax. The application replaces `{{.Rate}}` with the current BTC to UAH exchange rate before sending the email.

**For the** `email` **settings:**

- `GSES2_APP_EMAIL_FROM`: This variable specifies the email address that will be displayed as the sender of the email.
- `GSES2_APP_EMAIL_SUBJECT`: This variable contains the subject line of the email.
- `GSES2_APP_EMAIL_BODY`: This variable contains the body of the email. Any occurrence of `{{.Rate}}` in this field will be replaced with the current BTC to UAH exchange rate when the email is sent.

If you want to change the content of the email, simply set new values for `GSES2_APP_EMAIL_SUBJECT` and/or `GSES2_APP_EMAIL_BODY` as desired.

> **Note**
> If you wish to modify the content of the email, simply set new values for `GSES2_APP_EMAIL_SUBJECT` and/or `GSES2_APP_EMAIL_BODY` as desired. Remember to up again your `docker-compose` to apply the new settings after making changes to these variables.

> **Warning**
> It's important to keep the `{{.Rate}}` placeholder in the `GSES2_APP_EMAIL_BODY` field if you want to include the current exchange rate in the email.

## Usage

1. **Up the docker compose:**

```bash
docker-compose up --build --detach
```

2. **Use the API:**

   Get the current BTC to UAH rate:

   ```bash
   curl localhost:8080/api/rate
   ```

   **Subscribe to rate updates:**

   ```bash
   curl -X POST -d "email=subscriber@email.com" localhost:8080/api/subscribe
   ```

   **Send rate updates to all subscribers:**

   ```bash
   curl -X POST localhost:8080/api/sendEmails
   ```

## Detailed API Usage

For detailed examples of how the API works including screenshots, please see [API_USAGE.md](./docs/API_USAGE.md).

## Description

This API exposes three endpoints that perform different operations:

1.  **GET** `/api/rate`: This endpoint is used to retrieve the current exchange rate from BTC to UAH.

2.  **POST** `/api/subscribe`: This endpoint is used to add a new email address to the subscriber list.

3.  **POST** `/api/sendEmails`: This endpoint sends an email with the current BTC to UAH rate to all the subscribers.

## How It Works

The `main.go` file is the entry point for the Go application. It creates instances of the above services and injects them into the `controller`. It then maps the controller's methods to the HTTP endpoints and starts the server.

## Architecture diagram

![](docs/images/architecture-diagram.png)

## Project Tree

```
📦 gses2-app
├── 📂build
│   └── 📂package
│       ├── 📜Dockerfile
│       └── 📜entrypoint.sh
├── 📂cmd
│   └── 📂gses2-app
│       └── 📜main.go
├── 📜docker-compose.yml
├── 📂docs
│   ├── 📜API_USAGE.md
│   └── 📂images
├── 📜go.mod
├── 📜go.sum
├── 📂internal
│   ├── 📂core
│   │   ├── 📂port
│   │   │   ├── 📜logger.go
│   │   │   ├── 📜rate.go
│   │   │   ├── 📜user.go
│   │   │   └── 📜user_test.go
│   │   └── 📂service
│   │       ├── 📂rate
│   │       │   ├── 📜rate.go
│   │       │   └── 📜rate_test.go
│   │       ├── 📂sender
│   │       │   ├── 📜sender.go
│   │       │   └── 📜sender_test.go
│   │       └── 📂subscription
│   │           ├── 📜subscription.go
│   │           └── 📜subscription_test.go
│   ├── 📂handler
│   │   ├── 📂httpcontroller
│   │   │   ├── 📜httpcontroller.go
│   │   │   └── 📜httpcontroller_test.go
│   │   └── 📂router
│   │       ├── 📜router.go
│   │       └── 📜router_test.go
│   └── 📂repository
│       ├── 📂config
│       │   ├── 📜config.go
│       │   ├── 📜config_test.go
│       │   └── 📜model.go
│       ├── 📂logger
│       │   └── 📂rabbit
│       │       └── 📜logger.go
│       ├── 📂rate
│       │   └── 📂rest
│       │       ├── 📂binance
│       │       │   ├── 📜binance.go
│       │       │   └── 📜binance_test.go
│       │       ├── 📂coingecko
│       │       │   ├── 📜coingecko.go
│       │       │   └── 📜coingecko_test.go
│       │       ├── 📂kuna
│       │       │   ├── 📜kuna.go
│       │       │   └── 📜kuna_test.go
│       │       ├── 📜rest.go
│       │       └── 📜rest_test.go
│       ├── 📂sender
│       │   ├── 📂email
│       │   │   ├── 📜email.go
│       │   │   ├── 📜email_test.go
│       │   │   └── 📂send
│       │   │       ├── 📜message.go
│       │   │       ├── 📜message_test.go
│       │   │       ├── 📜send.go
│       │   │       └── 📜send_test.go
│       │   └── 📂smtp
│       │       ├── 📜smtp.go
│       │       ├── 📜smtp_test.go
│       │       └── 📜stub.go
│       └── 📂storage
│           ├── 📜csv.go
│           └── 📜csv_test.go
├── 📜LICENSE
├── 📜README.md
├── 📜README_ua.md
└── 📂test
    ├── 📂E2E
    │   ├── 📂build
    │   │   ├── 📜docker-compose.e2e.yml
    │   │   ├── 📜Dockerfile
    │   │   └── 📜entrypoint.e2e.sh
    │   ├── 📂fake
    │   │   ├── 📂kunaapi
    │   │   │   ├── 📜Dockerfile
    │   │   │   └── main.go
    │   │   └── 📂smtp
    │   │       ├── 📜Dockerfile
    │   │       ├── 📜main.go
    │   │       └── 📜san.cnf
    │   └── 📂postman
    │       └── 📜tests.e2e.json
    └── 📂integration
        ├── 📜httpcontroller_integration_test.go
        └── 📜subscription_integration_test.go
```
