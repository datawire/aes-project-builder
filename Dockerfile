FROM golang:1.14
WORKDIR /go/src/github.com/github.com/datawire/aes-project-builder
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-w -s' .

FROM gcr.io/kaniko-project/executor:v0.19.0
COPY --from=0 /go/bin/aes-project-builder /kaniko/aes-project-builder
ENTRYPOINT ["/kaniko/aes-project-builder"]
