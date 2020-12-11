FROM golang:1.15-alpine as builder

WORKDIR /root

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o exporter ./



FROM alpine:3.12
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /root/exporter .

CMD ["./exporter"]