FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETARCH
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /telemetry-mesh ./cmd/server/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /telemetry-mesh /telemetry-mesh
ENTRYPOINT ["/telemetry-mesh"]
