# This file is used to build an embedded all-in-one image

ARG NODE_VERSION=18.12.1
FROM node:${NODE_VERSION}-alpine3.17 AS uibuilder
WORKDIR /webui/
RUN npm config set update-notifier false

### Copy actual frontend source code
COPY webui .

### Compile frontend
RUN npm run build-embed

#################

FROM golang:1.19.4-alpine3.17 AS builder

### Install gcc in order to be able to use CGO later
RUN apk add --no-cache build-base

### Copy Go code, copying required folders only (use .dockerignore)
WORKDIR /src/
COPY . .

### Copy built frontend
COPY --from=uibuilder /webui/dist ./webui/dist

### Build executables, strip debug symbols and compress with UPX
# Use the "netgo" tag to resolve DNS queries in Go directly
RUN CGO_ENABLED=1 \
	go build  \
    	-tags netgo,webui  \
    	-mod=vendor  \
    	-ldflags '-extldflags "-static"' \
    	-a \
    	-o /app/webapi \
    	./cmd/webapi \
    && strip /app/*

#############################

### Create final container from scratch
FROM alpine:3.17

### Inform Docker about which port is used
EXPOSE 3000

### Install the CA Certificates in order to run HTTPS requests
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

### Create the "appuser" standard user
RUN adduser \
	--home /app/ \
	--disabled-password \
	--shell /bin/false \
	 appuser
USER appuser

### Copy the build executable from the builder image
WORKDIR /app/
COPY --from=builder /app/* ./

### Create appropriate volume
ENV CFG_DB_FILENAME=/app/db/wasaphoto.db
ENV CFG_USER_CONTENT_DIR=/app/static/user_content
RUN mkdir -p /app/db /app/static/user_content

### Executable command
CMD ["/app/webapi"]
