FROM golang:1.17.8-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .
FROM alpine:latest
COPY --from=builder /dist/main .
# VOLUME [ “/sys/fs/bpf” ]
EXPOSE 8090
ENTRYPOINT ["/main"]
