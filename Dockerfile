#Builder
FROM golang:1.18-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/auth-service ./cmd/auth

#Auth-service
FROM alpine:3.15.4
WORKDIR /app
COPY --from=builder /app/bin /app/
COPY  . /app/
EXPOSE 8080
CMD ["./auth-service", "-c", "./config.yaml"]
