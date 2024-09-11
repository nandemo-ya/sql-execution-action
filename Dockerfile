FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

ARG TARGETARCH TARGETOS
WORKDIR /go/src/github.com/nandemo-ya/sql-execution-action
COPY . .

RUN go mod download
RUN go mod vendor
RUN GO111MODULE=on GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ./bin/sql-execution-action -mod=vendor main.go

FROM gcr.io/distroless/base-debian12:latest
COPY --from=builder /go/src/github.com/nandemo-ya/sql-execution-action/bin/sql-execution-action /usr/local/bin/

ENTRYPOINT ["sql-execution-action"]