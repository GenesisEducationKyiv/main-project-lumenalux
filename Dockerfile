FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go test -v ./...
RUN CGO_ENABLED=0 go build -o /bin/gses2-app ./cmd/gses2-app/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/gses2-app .
COPY --from=builder /app/config.yaml .
RUN touch storage.csv
EXPOSE 8080 465
ENTRYPOINT ["./gses2-app"]
