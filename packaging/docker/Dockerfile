# DiliVet – ML-DSA diagnostics and vetting toolkit
# Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

FROM gcr.io/distroless/static:nonroot

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY dist/${TARGETOS}-${TARGETARCH}/dilivet /usr/local/bin/dilivet

USER nonroot:nonroot

ENTRYPOINT ["/usr/local/bin/dilivet"]

