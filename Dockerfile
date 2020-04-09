FROM alpine:3.11

LABEL maintainer="Mohammad Abdolirad <m.abdolirad@gmail.com>" \
    org.label-schema.name="wait4x" \
    org.label-schema.vendor="atkrad" \
    org.label-schema.description="Wait4X is a cli tool to wait for everything! It can be wait for a port to open or enter to rquested state." \
    org.label-schema.vcs-url="https://github.com/atkrad/wait4x" \
    org.label-schema.license="MIT"

COPY .docker/root/ /
COPY wait4x-linux-musl-amd64 /usr/local/bin/wait4x

RUN apk --update --no-cache add ca-certificates

ENTRYPOINT ["entrypoint"]
CMD ["wait4x"]
