# {{jreleaserCreationStamp}}
# Based on https://raw.githubusercontent.com/jreleaser/jreleaser/refs/tags/v1.17.0/core/jreleaser-templates/src/main/resources/META-INF/jreleaser/templates/binary/docker/Dockerfile.tpl
# to resolve https://github.com/jreleaser/jreleaser/discussions/1834
# dockerBaseImage is https://hub.docker.com/_/scratch
FROM {{dockerBaseImage}}

{{#dockerLabels}}
LABEL {{.}}
{{/dockerLabels}}

{{#dockerPreCommands}}
{{.}}
{{/dockerPreCommands}}

# / at then means, bin should be a directory
COPY --chmod=0777 assembly/bin/{{distributionExecutableUnix}} /bin/

# dockerBaseImage is scratch and has no sh shell
# RUN CHMOD will not work then
# RUN chmod +x {{distributionArtifactRootEntryName}}/bin/{{distributionExecutableUnix}}

{{#dockerPostCommands}}
{{.}}
{{/dockerPostCommands}}

ENV PATH="${PATH}:/bin"

ENTRYPOINT ["{{distributionExecutableUnix}}"]
