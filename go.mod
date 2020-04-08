module github.com/datawire/aes-project-builder

go 1.14

// These replacements are copied from the kaniko go.mod
replace (
	github.com/containerd/containerd v1.4.0-0.20191014053712-acdcf13d5eaf => github.com/containerd/containerd v0.0.0-20191014053712-acdcf13d5eaf
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v0.0.0-20190319215453-e7b5f7dbe98c
	github.com/tonistiigi/fsutil v0.0.0-20190819224149-3d2716dd0a4d => github.com/tonistiigi/fsutil v0.0.0-20191018213012-0f039a052ca1
)

require (
	github.com/GoogleContainerTools/kaniko v0.19.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/sirupsen/logrus v1.5.0
	gopkg.in/src-d/go-git.v4 v4.13.1
)
