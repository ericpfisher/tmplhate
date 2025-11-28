# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build
WORKDIR /src

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary (GoReleaser injects GOOS, GOARCH, ldflags)
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 go build -o /tmplhate .

# Final minimal image
FROM alpine:latest
COPY --from=build /tmplhate /tmplhate
ENTRYPOINT ["/tmplhate"]
CMD ["--help"]
