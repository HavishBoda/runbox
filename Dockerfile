FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o runbox .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/runbox .
EXPOSE 8080
CMD ["./runbox"]