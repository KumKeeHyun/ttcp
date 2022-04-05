FROM golang:1.17.8-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod go.sum main.go bpf_bpfel.go bpf_bpfel.o ./
RUN go mod download
RUN go build -o main .

WORKDIR /dist
RUN cp /build/main .
FROM alpine:latest
COPY --from=builder /dist/main .
ENTRYPOINT ["/main"]
