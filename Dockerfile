FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.1.0 AS xx

FROM --platform=$BUILDPLATFORM golang:1.18.1-alpine3.15 AS base
ENV GO111MODULE=auto
ENV CGO_ENABLED=0

COPY --from=xx / /
RUN apk add --update --no-cache build-base bash coreutils git
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM

RUN --mount=type=bind,target=/src,rw \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=target=/go/pkg/mod,type=cache \
    GO_BINARY=xx-go WAIT4X_BUILD_OUTPUT=/usr/bin make build \
    && xx-verify --static /usr/bin/wait4x

FROM scratch AS binary
COPY --from=build /usr/bin/wait4x /

FROM base AS releaser
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /work
RUN --mount=from=binary,target=/build \
  --mount=type=bind,target=/src \
  mkdir -p /out \
  && cp /build/wait4x /src/README.md /src/LICENSE . \
  && tar -czvf "/out/wait4x-${TARGETOS}-${TARGETARCH}${TARGETVARIANT}.tar.tgz" * \
  && sha256sum -z "/out/wait4x-${TARGETOS}-${TARGETARCH}${TARGETVARIANT}.tar.tgz" | awk '{ print $1 }' > "/out/wait4x-${TARGETOS}-${TARGETARCH}${TARGETVARIANT}.tar.tgz.sha256"

FROM scratch AS artifact
COPY --from=releaser /out /

FROM alpine:3.15
RUN apk add --no-cache ca-certificates

COPY --from=binary /wait4x /usr/bin/wait4x

ENTRYPOINT ["wait4x"]
CMD ["help"]
