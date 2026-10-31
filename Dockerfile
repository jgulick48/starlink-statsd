# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder
ARG TARGETOS TARGETARCH VERSION=dev COMMIT=unknown
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" -o /out/starlink-statsd .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /out/starlink-statsd /starlink-statsd
USER 65532:65532
ENTRYPOINT ["/starlink-statsd"]
