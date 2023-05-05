FROM golang:alpine AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o logging ./cmd/api

FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/logging .
CMD ["./logging"]
