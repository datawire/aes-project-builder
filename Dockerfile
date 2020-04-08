FROM golang:1.14
WORKDIR /go/src/github.com/github.com/datawire/aes-project-builder
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-w -s' .

# gcr.io/kaniko-project/executor:v0.16.0, plus a bugfix
# https://github.com/GoogleContainerTools/kaniko/pull/1184 Avoiding
# newer Kaniko versions because of
# https://github.com/GoogleContainerTools/kaniko/issues/1162
FROM quay.io/datawire/kaniko:cd6e822
COPY --from=0 /go/bin/aes-project-builder /kaniko/aes-project-builder
ENTRYPOINT ["/kaniko/aes-project-builder"]
