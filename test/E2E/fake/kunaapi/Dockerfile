FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN GO111MODULE=auto CGO_ENABLED=0 go build -v -o main .

FROM scratch
COPY --from=builder /app/main /app/main
WORKDIR /app
EXPOSE 8080
CMD ["./main"]
