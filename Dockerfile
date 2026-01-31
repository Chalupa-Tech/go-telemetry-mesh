FROM golang:1.24.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY api/ api/
COPY internal/ internal/
COPY cmd/ cmd/

# Build using vendor directory
RUN go build -mod=vendor -o go-telemetry-mesh ./cmd/go-telemetry-mesh

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/go-telemetry-mesh .
CMD ["./go-telemetry-mesh"]
