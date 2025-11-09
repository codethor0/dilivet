FROM gcr.io/distroless/static:nonroot

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY dist/${TARGETOS}-${TARGETARCH}/dilivet /usr/local/bin/dilivet

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/dilivet"]

