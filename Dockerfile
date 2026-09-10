FROM golang:1.27-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /consumer ./cmd/consumer

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
	&& adduser -D -H appuser

COPY --from=builder /consumer /usr/local/bin/consumer

USER appuser
ENTRYPOINT ["/usr/local/bin/consumer"]