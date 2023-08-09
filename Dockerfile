# See https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
FROM golang:1.21-alpine AS builder

# Git is used for dependencies
RUN apk update \
    && apk add \
        --no-cache \
        git \
        ca-certificates

ENV USER=hackathon
ENV UID=1000
# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    -u "${UID}" \
    "${USER}"

ARG PROJECT
ARG RELEASE
ARG COMMIT

WORKDIR $GOPATH/src/${PROJECT}/

COPY ${PROJECT}/go.mod ${PROJECT}/go.sum ./

RUN go mod download \
    && go mod verify

COPY ${PROJECT}/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /go/bin/hackathon \
    -ldflags="\
        -w -s \
        -X ${PROJECT}/version.Name=${PROJECT} \
        -X ${PROJECT}/version.Release=${RELEASE} \
        -X ${PROJECT}/version.Commit=${COMMIT} \
        -X ${PROJECT}/version.BuildTime=$(date -u -Iseconds)" \
    ${PROJECT}.go

FROM scratch

ARG PROJECT

ENV RABBITMQ "amqp://guest:guest@rabbitmq:5672/"

WORKDIR /go

COPY --from=builder \
    /etc/passwd \
    /etc/group \
    /etc/
COPY --from=builder \
    /etc/ssl/certs/ca-certificates.crt \
    /etc/ssl/certs/
COPY --from=builder \
    /go/bin/hackathon \
    /go/bin/hackathon

USER 1000

ENTRYPOINT [ "/go/bin/hackathon" ]
