FROM golang:1.23-alpine AS builder
ARG SERVICE
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/$SERVICE ./$SERVICE

FROM alpine:3.20
ARG SERVICE
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/$SERVICE /app/$SERVICE
COPY --from=builder /src/$SERVICE/etc /app/etc
COPY deploy/docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT ["/app/entrypoint.sh"]
