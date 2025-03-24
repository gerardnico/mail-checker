# Generated with JReleaser 1.17.0 at 2025-03-24T18:56:51.561713196+01:00
# Based on https://raw.githubusercontent.com/jreleaser/jreleaser/refs/tags/v1.17.0/core/jreleaser-templates/src/main/resources/META-INF/jreleaser/templates/binary/docker/Dockerfile.tpl
# to resolve https://github.com/jreleaser/jreleaser/discussions/1834
# dockerBaseImage is https://hub.docker.com/_/scratch
FROM scratch

LABEL "org.opencontainers.image.title"="mail-checker"
LABEL "org.opencontainers.image.description"="Mail Checker"
LABEL "org.opencontainers.image.url"="https://github.com/gerardnico/mail-checker"
LABEL "org.opencontainers.image.licenses"="MIT"
LABEL "org.opencontainers.image.version"="0.1.0"
LABEL "org.opencontainers.image.revision"="f99f3a0dc76ca13fe1bb97e26575bcad793c0842"


# / at then means, bin should be a directory
COPY --chmod=0777 assembly/bin/mail-checker /bin/

# dockerBaseImage is scratch and has no sh shell
# RUN CHMOD will not work then
# RUN chmod +x mail-checker-0.1.0-linux-amd64/bin/mail-checker


ENV PATH="${PATH}:/bin"

ENTRYPOINT ["mail-checker"]
