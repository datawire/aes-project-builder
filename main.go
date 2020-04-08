package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/GoogleContainerTools/kaniko/pkg/constants"
	"github.com/kballard/go-shellquote"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	libgit "gopkg.in/src-d/go-git.v4"
	libgitPlumbing "gopkg.in/src-d/go-git.v4/plumbing"
	libgitPlumbingTransport "gopkg.in/src-d/go-git.v4/plumbing/transport"
	libgitHTTP "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	sgrReset = 0
	sgrBold  = 1
	sgrGreen = 32
)

func sgr(codes ...int) string {
	strs := make([]string, 0, len(codes))
	for _, code := range codes {
		strs = append(strs, strconv.Itoa(code))
	}
	return "\x1B[" + strings.Join(strs, ";") + "m"
}

func comment(str string) {
	fmt.Println(sgr(sgrGreen, sgrBold) + "# " + str + sgr(sgrReset))
}

func mockCmd(args ...string) {
	fmt.Println(shellquote.Join(args...))
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

type Args struct {
	URL         string
	Auth        libgitPlumbingTransport.AuthMethod
	Rev         libgitPlumbing.Hash
	Ref         libgitPlumbing.ReferenceName
	KanikoFlags []string
}

func Main(args Args) error {
	comment("Checking out code...")

	mockCmd("git", "clone", "--no-checkout", args.URL, constants.BuildContextDir)
	repo, err := libgit.PlainClone(constants.BuildContextDir, false, &libgit.CloneOptions{
		URL:        args.URL,
		Progress:   os.Stdout,
		Auth:       args.Auth,
		NoCheckout: true,
	})
	if err != nil {
		return err
	}

	mockCmd("cd", constants.BuildContextDir)
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	onBranch := false
	if args.Ref.IsBranch() {
		mockCmd("git", "checkout", args.Ref.Short())
		err = worktree.Checkout(&libgit.CheckoutOptions{
			Branch: args.Ref,
		})
		if err != nil {
			logrus.Infoln(err)
		} else {
			onBranch = true
			mockCmd("git", "reset", "--hard", args.Rev.String())
			err = worktree.Reset(&libgit.ResetOptions{
				Commit: args.Rev,
				Mode:   libgit.HardReset,
			})
		}
	}
	if !onBranch {
		mockCmd("git", "checkout", args.Rev.String())
		err = worktree.Checkout(&libgit.CheckoutOptions{
			Hash: args.Rev,
		})
	}
	if err != nil {
		return err
	}

	fmt.Println()
	comment("Performing preflight check...")

	if !fileExists(filepath.Join(constants.BuildContextDir, "Dockerfile")) {
		return fmt.Errorf("Git commit %s does not contain a file \"Dockerfile\"", args.Rev)
	}

	fmt.Println()
	comment("Building...")

	kanikoArgs := append([]string{"/kaniko/executor", "--context=dir://" + constants.BuildContextDir}, args.KanikoFlags...)
	mockCmd(kanikoArgs...)
	cmd := exec.Command(kanikoArgs[0], kanikoArgs[1:]...)
	// Run kaniko, filtering out spurious errors.
	//
	// TODO: Figure out how to get kaniko to not spit these out.
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}
	wg, _ := errgroup.WithContext(context.Background())
	wg.Go(cmd.Wait)
	wg.Go(func() error {
		scanner := bufio.NewScanner(reader)
		re := regexp.MustCompile(`^ERROR:\s+logging\s+before\s+flag.Parse:.*$`)
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				continue
			}
			fmt.Println(line)
		}
		return scanner.Err()
	})
	return wg.Wait()
}

func main() {
	args := Args{
		URL:         "https://github.com/" + os.Getenv("KALE_REPO") + ".git",
		Auth:        &libgitHTTP.BasicAuth{Username: os.Getenv("KALE_CREDS")},
		Rev:         libgitPlumbing.NewHash(os.Getenv("KALE_REV")),
		Ref:         libgitPlumbing.ReferenceName(os.Getenv("KALE_REF")),
		KanikoFlags: os.Args[1:],
	}
	os.Unsetenv("KALE_CREDS")
	os.Unsetenv("KALE_REPO")
	os.Unsetenv("KALE_REV")
	os.Unsetenv("KALE_REF")
	if err := Main(args); err != nil {
		logrus.Fatal(err)
	}
}
