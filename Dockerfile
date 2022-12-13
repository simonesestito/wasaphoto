# This file is used by Docker "build" or "buildah" to create a container image for this Go project
# The build is done using "multi-stage" approach, where a temporary container ("builder") is used to build the Go
# executable, and the final image is from scratch (empty container) for both security and performance reasons.

ARG DOCKER_PREFIX
FROM ${DOCKER_PREFIX}node:lts AS uibuilder
WORKDIR webui
RUN npm config set update-notifier false

### Copy node_modules only (updating only source code won't rerun this part)
COPY webui/package.json .
COPY webui/package-lock.json .
COPY webui/node_modules ./node_modules
RUN npm install

### Copy actual frontend source code
COPY webui .

### Compile frontend
RUN npm run build-embed


ARG DOCKER_PREFIX
FROM ${DOCKER_PREFIX}enrico204/golang:1.19.4-6 AS builder
WORKDIR /app

### Copy Go code, copying required folders only
COPY cmd cmd
COPY service service
COPY vendor vendor
COPY go.mod .
COPY go.sum .
COPY --from=uibuilder webui webui

### Set some build variables
ARG APP_VERSION
ARG BUILD_DATE
ARG REPO_HASH

RUN go generate -mod=vendor ./...

### Build executables, strip debug symbols and compress with UPX
ENV CGO_ENABLED=1
RUN /bin/bash -euo pipefail -c "for ex in \$(ls cmd/); do pushd cmd/\$ex; go build -tags webui,openapi -mod=vendor -ldflags \"-extldflags \\\"-static\\\" -X main.AppVersion=${APP_VERSION} -X main.BuildDate=${BUILD_DATE}\" -a -o /app/\$ex .; popd; done"
RUN strip ./webapi && upx -9 ./webapi

### Create necessary folders that will be used later in the final scratch image
USER root
RUN mkdir -p ./db/ ./static/user_content/ && touch ./wasaphoto.db
USER appuser

### Create final container from scratch
FROM scratch

### Inform Docker about which port is used
EXPOSE 3000

### Populate scratch with CA certificates and Timezone infos from the builder image
ENV ZONEINFO /zoneinfo.zip
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /zoneinfo.zip /
COPY --from=builder /etc/passwd /etc/passwd

### Copy the build executable from the builder image
WORKDIR /app/
COPY --from=builder /app/webapi .

### Set some build variables
ARG APP_VERSION
ARG BUILD_DATE
ARG PROJECT_NAME

### Downgrade to user level (from root)
USER appuser

### Configure volumes
COPY --from=builder --chown=appuser:1000 /app/static/user_content/ ./static/user_content/
COPY --from=builder --chown=appuser:1000 /app/db/ ./db/
COPY --from=builder --chown=appuser:1000 /app/wasaphoto.db .

### Executable command
ENV CFG_WEB_API_HOST='0.0.0.0:3000'
CMD ["/app/webapi"]

### OpenContainers tags
LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="${PROJECT_NAME}" \
      org.opencontainers.image.authors="Simone Sestito <simone@simonesestito.com>" \
      org.opencontainers.image.source="https://github.com/simonesestito/${PROJECT_NAME}" \
      org.opencontainers.image.revision="${REPO_HASH}" \
      org.opencontainers.image.vendor="Simone Sestito" \
	  org.opencontainers.image.version="${APP_VERSION}"
