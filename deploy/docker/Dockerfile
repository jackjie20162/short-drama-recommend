FROM golang:1.23-alpine AS builder
ARG SERVICE
WORKDIR /src
RUN apk add --no-cache git protobuf
COPY go.mod go.sum* ./
RUN go mod download
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
COPY . .
RUN ./scripts/gen.sh
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
