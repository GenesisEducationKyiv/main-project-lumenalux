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
- [Architecture](#architecture)

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

   The application uses environment variables for configuration. Set up the following environment variables for the SMTP server and email settings:

   ```bash
   export GSES2_APP_SMTP_HOST="<smtp server host>"
   ```

   ```bash
   export GSES2_APP_SMTP_USER="<smtp username>"
   ```

   ```bash
   export GSES2_APP_SMTP_PASSWORD="<smtp password>"`
   ```

   The rest of the environment variables have default values as listed below, but can be overridden if necessary:

   - `GSES2_APP_SMTP_PORT="465"`
   - `GSES2_APP_EMAIL_FROM="no.reply@test.info.api"`
   - `GSES2_APP_EMAIL_SUBJECT="BTC to UAH exchange rate"`
   - `GSES2_APP_EMAIL_BODY="The BTC to UAH rate is {{.Rate}}"`
   - `GSES2_APP_STORAGE_PATH="./storage/storage.csv"`
   - `GSES2_APP_HTTP_PORT="8080"`
   - `GSES2_APP_HTTP_TIMEOUT="10s"`
   - `GSES2_APP_KUNA_API_URL="https://api.kuna.io/v3/tickers?symbols=btcuah"`
   - `GSES2_APP_KUNA_API_DEFAULT_RATE="0"`

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

## Architecture

1.  **Application Init** **(**`main.go`**)**: This is the entry point of the application. It initializes all necessary services, injects them into the controller, maps the controller's methods to HTTP endpoints, and starts the server. This signifies the birth of the application's lifecycle

2.  **Controller**: It sits right beneath the **Application Init**. The controller handles incoming HTTP requests, utilizes appropriate services for required operations, and responds to these requests. In other words, it's responsible for coordinating the tasks and directing the flow of the application

3.  **Services**: There are three core services that the controller depends on:

    - **Rate Service**: This service is responsible for communicating with external APIs to fetch the current BTC to UAH exchange rates. It has different implementations including the Kuna provider and a stub provider for testing purposes

    - **Subscription Service**: This service manages the operations of adding to and retrieving subscribers from the CSV storage. It allows new subscribers to be added and also retrieval of the existing subscriber list

    - **Sender Service**: This service is in charge of sending emails. It uses different mechanisms, for example, the Email provider, which communicates with an external SMTP server to send emails, and a stub provider for testing

4.  **Additional Components**:

    - **Storage (CSV)**: A storage mechanism used by the Subscription Service to keep track of subscribers. It uses CSV file for storage

    - **Transport**: Handling the API routing with the help of the **Controller**

This architecture ensures a separation of concerns, where each component focuses on a specific task. This separation allows for easier maintenance and improved scalability as changes to one component shouldn't affect the others. Also, the use of stubs allows for easier testing, by isolating the application logic from external dependencies
