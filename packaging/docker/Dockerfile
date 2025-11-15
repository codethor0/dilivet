# DiliVet — ML-DSA diagnostics toolkit
# Copyright (c) 2025 Thor Thor (codethor0)
# Project: github.com/codethor0/dilivet
# LinkedIn: https://www.linkedin.com/in/thor-thor0

FROM gcr.io/distroless/static:nonroot

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY dist/${TARGETOS}-${TARGETARCH}/dilivet /usr/local/bin/dilivet

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/dilivet"]

