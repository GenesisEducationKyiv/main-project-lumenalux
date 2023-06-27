FROM golang:1.20 AS builder
WORKDIR /app
COPY . .

RUN apt-get update && apt-get install -y openssl
RUN openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config san.cnf

RUN go mod init smtp
RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main .

FROM scratch
COPY --from=builder /app/main /app/main
COPY --from=builder /app/cert.pem /app/
COPY --from=builder /app/key.pem /app/

WORKDIR /app
EXPOSE 1025

CMD ["./main"]