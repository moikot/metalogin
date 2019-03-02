ARG GO_VERSION=1.10

FROM golang:${GO_VERSION}-alpine as build-env

ARG DEP_VERSION=0.4.1
ARG APP_FOLDER=/go/src/github.com/moikot/metalogin/

ADD . ${APP_FOLDER}
WORKDIR ${APP_FOLDER}

# Update root certificates, curl and git
RUN apk add -U --no-cache ca-certificates curl git

# Installs Go Dep
RUN curl \
    -fsSL \
    -o /usr/local/bin/dep \
    https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 \
    && chmod \
    +x /usr/local/bin/dep

RUN dep ensure -vendor-only

# Compile independant executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/main .

FROM scratch

COPY --from=build-env /bin/main /

ENTRYPOINT ["/main"]
