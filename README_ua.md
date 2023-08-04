# gses2-app API обміну BTC на UAH

![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

## Переклад

- [English](README.md).

## Зміст

- [Про проект](#про-проект)
- [Встановлення](#встановлення)
- [Використання](#використання)
- [Опис API](#опис-api)
- [Як це працює](#як-це-працює)
- [Архітектура](#архітектурна-діаграма)
- [Дерево проєкту](#дерево-проєкту)

## Про проект

Це API, яке надає поточний курс обміну між Bitcoin та українською гривнею (UAH). Воно дозволяє користувачам підписуватися на оновлення курсу та отримувати ці оновлення по електронній пошті.

## Встановлення

1. **Склонуйте репозиторій до бажаного місця:**

   ```bash
   git clone https://github.com/lumenalux/gses2-app.git gses2-api
   cd gses2-api
   ```

2. **Налаштуйте змінні середовища:**

   Додаток використовує файл `.env` для конфігурації. Скопіюйте вміст файлу `.env.example` у новий файл з назвою `.env`. Налаштуйте наступні змінні середовища для налаштування сервера SMTP та електронної пошти:

   ```bash
    GSES2_APP_SMTP_HOST="<адреса сервера smtp>"

    GSES2_APP_SMTP_USER="<ім'я користувача smtp>"

    GSES2_APP_SMTP_PASSWORD="<пароль smtp>"
   ```

   Решта змінних середовища має значення за замовчуванням, що наведено нижче, але їх можна перевизначити, якщо це необхідно:

   ```bash
    GSES2_APP_SMTP_PORT=465

    GSES2_APP_EMAIL_FROM=no.reply@test.info.api
    GSES2_APP_EMAIL_SUBJECT=BTC до курсу UAH
    GSES2_APP_EMAIL_BODY=Курс обміну BTC до UAH становить {{.Rate}} UAH за BTC

    GSES2_APP_STORAGE_PATH=./storage/storage.csv

    GSES2_APP_HTTP_PORT=8080
    GSES2_APP_HTTP_TIMEOUT=10s

    GSES2_APP_KUNAAPI_URL=https://api.kuna.io/v3/tickers?symbols=btcuah

    GSES2_APP_BINANCEAPI_URL=https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT

    GSES2_APP_COINGECKOAPI_URL=https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah

    GSES2_APP_RABBITMQ_URL=amqp://guest:guest@amqp/
   ```

   Змінні середовища включають налаштування сервера SMTP та вміст повідомлень електронної пошти, що відправляються підписникам. Тіло електронного листа розроблено як шаблон з використанням синтаксису text/template Go. Додаток замінює `{{.Rate}}` на поточний курс обміну BTC на UAH перед відправкою електронного листа.

   **Щодо налаштувань** `email`**:**

   - `GSES2_APP_EMAIL_FROM`: ця змінна визначає адресу електронної пошти, яка буде відображатися як відправник електронного листа
   - `GSES2_APP_EMAIL_SUBJECT`: ця змінна містить тему електронного листа
   - `GSES2_APP_EMAIL_BODY`: ця змінна містить тіло електронного листа. Будь-яке входження `{{.Rate}}` у це поле буде замінено на поточний курс обміну BTC на UAH при відправці електронного листа

   Якщо ви хочете змінити вміст електронного листа, просто встановіть нові значення для `GSES2_APP_EMAIL_SUBJECT` та/або `GSES2_APP_EMAIL_BODY` за бажанням

   > **Note**
   > Якщо ви бажаєте змінити вміст електронного листа, просто встановіть нові значення для `GSES2_APP_EMAIL_SUBJECT` та/або `GSES2_APP_EMAIL_BODY` за бажанням. Пам'ятайте, що після внесення змін у ці змінні, вам потрібно буде підняти `docker-compose`, щоб застосувати нові налаштування.

   > **Warning**
   > Важливо зберегти заповнювач `{{.Rate}}` у полі `GSES2_APP_EMAIL_BODY`, якщо ви хочете включити поточний курс обміну в електронний лист.

## Використання

1.  **Підніміть docker compose:**

    ```bash
    docker-compose up --build --detach
    ```

2.  **Використання API:**

    **Отримати поточний курс BTC до UAH:**

    ```bash
    curl localhost:8080/api/rate
    ```

    **Підписатися на оновлення курсу:**

    ```bash
    curl -X POST -d "email=subscriber@email.com" localhost:8080/api/subscribe
    ```

    **Надіслати оновлення курсу всім підписникам:**

    ```bash
    curl -X POST localhost:8080/api/sendEmails
    ```

## Додатково про використання API

Для докладних прикладів того, як працює API, включно зі скріншотами, будь ласка, перегляньте файл [API_USAGE.md](./docs/API_USAGE.md).

## Опис API

Це API надає три ендпоінти, які виконують різні операції:

1.  **GET** `/api/rate`: Цей ендпоінт використовується для отримання поточного обмінного курсу від BTC до UAH.

2.  **POST** `/api/subscribe`: Цей ендпоінт використовується для додавання нової адреси електронної пошти до списку підписників.

3.  **POST** `/api/sendEmails`: Цей ендпоінт надсилає електронний лист з поточним обмінним курсом від BTC до UAH всім підписникам.

## Як це працює

Файл `main.go` є точкою входу для програми Go. Він створює екземпляри вищезазначених сервісів та впроваджує їх у `controller`. Потім він співставляє методи контролера на HTTP-ендпоінти та запускає сервер.

## Архітектурна діаграма

![](docs/images/architecture-diagram.png)

## Дерево проєкту

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
