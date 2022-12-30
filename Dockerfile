# This file is used to build an embedded all-in-one image

ARG NODE_VERSION=18.12.1
FROM node:${NODE_VERSION}-alpine3.17 AS uibuilder
WORKDIR /webui/
RUN npm config set update-notifier false

### Copy actual frontend source code
COPY webui .

### Compile frontend
RUN npm run build-embed


FROM golang:1.19.4-alpine3.17 AS builder

### Install gcc in order to be able to use CGO later
### and ca-certificates in order to be able to perform HTTPS requests
RUN apk add --no-cache build-base ca-certificates && \
	update-ca-certificates

### Create "appuser" standard user, later used in final image to run ./cmd/webapi as non-root
RUN adduser \
	--home /app/ \
	--disabled-password \
	--shell /bin/false \
	 appuser

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


### Create final container from scratch
FROM scratch

### Inform Docker about which port is used
EXPOSE 3000

### Copy the CA Certificates in order to run HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

### Copy the "appuser" standard user from the build stage
COPY --from=builder /etc/passwd /etc/passwd
USER appuser

### Copy the build executable from the builder image
WORKDIR /app/
COPY --from=builder /app/* ./

### Executable command
ENV CFG_WEB_API_HOST='0.0.0.0:3000'
ENV CFG_LOG_DEBUG=false
CMD ["/app/webapi"]
