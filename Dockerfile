FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /app

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o main ./cmd/server

FROM gcr.io/distroless/static-debian12:latest AS server

WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/main .
COPY --chown=nonroot:nonroot web ./web

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT [ "./main" ]