FROM golang:1.23.5 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o fiber-app
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/fiber-app .
EXPOSE 3000
CMD ["./fiber-app"]