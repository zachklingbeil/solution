# Build
FROM golang:latest AS go_builder
WORKDIR /electric

# Cache dependencies
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build for amd64 only
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o /bin/electric .

# Run
FROM alpine:latest AS final
RUN apk --update add \
    ca-certificates \
    tzdata \
    && \
    update-ca-certificates

ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

COPY --from=go_builder /bin/electric /bin/

ENTRYPOINT [ "/bin/electric" ]