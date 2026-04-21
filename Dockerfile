FROM golang:1.25.8-alpine AS builder


RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/data ./data

EXPOSE 8080


CMD ["./app"]