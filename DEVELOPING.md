How do I update this?
---------------------

just push the desired changes to Git; Quay will automatically build a
`quay.io/datawire/aes-project-builder` image from any branches or tags
pushed.  Then go to
https://quay.io/repository/datawire/aes-project-builder?tab=tags and
find the Docker tag that it made, click the download button on the
right and select "Docker Pull (by digest)".  Copy the image name, and
put it in `apro.git/cmd/amb-sidecar/kale/kale.go`.

How do I iterate more quickly?
------------------------------

When developing locally, you can speed up builds by running `go mod
vendor` before running `docker build .`; this allows the build
happening in Docker to take advantage of the host's module cache.
