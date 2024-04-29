ARG ARCH=

FROM ${ARCH}golang:1.22.2 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./ ./

RUN go test ./...

RUN GOOS=linux CGO_ENABLED=0 go build

FROM ${ARCH}alpine:3.19.1

COPY --from=builder /app/starlink-statsd /bin/starlink-statsd
WORKDIR /var/lib/starlink-statsd/

CMD ["/bin/starlink-statsd"]
