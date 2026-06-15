# syntax=docker/dockerfile:1.7

FROM golang:1.26-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" \
    -o /out/outbound-api ./cmd/outbound-api

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=build /out/outbound-api /app/outbound-api
COPY configs/default.yaml /app/configs/default.yaml

ENV OUTBOUND_ENV=productive
ENV OUTBOUND_CONFIG_DIR=/app/configs
ENV OUTBOUND_HTTP_HOST=0.0.0.0
ENV OUTBOUND_HTTP_PORT=8080

EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/outbound-api"]
