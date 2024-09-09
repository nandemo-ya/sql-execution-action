FROM golang:1.23 AS builder
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/bin/main main.go

FROM gcr.io/distroless/base-debian12:latest
COPY --from=builder /app/bin/main /usr/local/bin/

ENTRYPOINT ["main"]
