FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG TARGETARCH
ARG GITHUB_PAT
WORKDIR /src

ENV GOPRIVATE=github.com/Chalupa-Tech/*
RUN if [ -n "$GITHUB_PAT" ]; then \
      git config --global url."https://${GITHUB_PAT}@github.com/Chalupa-Tech".insteadOf "https://github.com/Chalupa-Tech"; \
    fi

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go test ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o /telemetry-mesh ./cmd/server/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /telemetry-mesh /telemetry-mesh
ENTRYPOINT ["/telemetry-mesh"]
