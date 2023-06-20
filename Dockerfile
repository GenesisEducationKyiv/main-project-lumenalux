FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -o gses2-app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gses2-app .
COPY --from=builder /app/config.yaml .
RUN touch storage.csv
EXPOSE 8080 465
ENTRYPOINT ["./gses2-app"]
