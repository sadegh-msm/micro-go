FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

RUN chmod +x /app/loggerApp

# just build a image and run the code and after it copy the executable to a minimal image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/loggerApp /app

CMD ["/app/loggerApp"]