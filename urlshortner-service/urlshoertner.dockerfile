FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o shortnerApp ./cmd/api

RUN chmod +x /app/shortnerApp

# just build a image and run the code and after it copy the executable to a minimal image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/shortnerApp /app

CMD ["/app/shortnerApp"]