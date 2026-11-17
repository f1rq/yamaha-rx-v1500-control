FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o yamaha-rx-v1500-control .

FROM alpine:3.20
RUN apk add --no-cache tzdata
COPY --from=builder /app/yamaha-rx-v1500-control /usr/local/bin/yamaha-rx-v1500-control
EXPOSE 8080
CMD ["yamaha-rx-v1500-control"]