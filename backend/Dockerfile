FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /gitserv ./cmd/gitserv/

FROM alpine:3.20

RUN apk add --no-cache git ca-certificates tzdata

COPY --from=builder /gitserv /usr/local/bin/gitserv

RUN mkdir -p /data/repos

EXPOSE 8080

ENV DATA_DIR=/data
ENV LISTEN_ADDR=:8080

CMD ["gitserv"]
