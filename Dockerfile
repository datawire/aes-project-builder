FROM alpine:3.10
RUN apk --no-cache add bash ca-certificates git
COPY --from=gcr.io/kaniko-project/executor:v0.19.0 /kaniko/executor /usr/local/bin/kaniko-executor
COPY run.sh /usr/local/bin/run.sh
ENTRYPOINT ["/usr/local/bin/run.sh"]
