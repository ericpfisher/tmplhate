# syntax=docker/dockerfile:1

# Minimal image via GoReleaser
FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/usr/bin/tmplhate"]
CMD ["--help"]
COPY $TARGETPLATFORM/tmplhate /usr/bin/
