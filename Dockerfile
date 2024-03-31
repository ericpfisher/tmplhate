# syntax: docker/dockerfile:1


FROM --platform=$BUILDPLATFORM golang:1.22-alpine as build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY cmd cmd
COPY core core
COPY main.go ./
ARG TARGETOS TARGETARCH VERSION
RUN --mount=type=cache,target=/root/.cache/go-build \
		--mount=type=cache,target=/go/pkg \
		CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags="-X 'tmplhate/cmd.Version=$VERSION'" \
    -o /tmplhate .

FROM alpine
COPY --from=build /tmplhate /tmplhate
ENTRYPOINT ["/tmplhate"]
CMD ["--help"]
