FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go get -d -v ./...

COPY . .

RUN go build -o /TradeSimulator .

FROM alpine:latest

RUN apk --no-cache add tzdata \
    && apk add --no-cache gettext

WORKDIR /root/
COPY --from=builder /app/public ./public

COPY --from=builder /TradeSimulator /TradeSimulator

EXPOSE 8000
CMD ["/TradeSimulator"]