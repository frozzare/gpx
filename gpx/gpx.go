package gpx

import (
	"fmt"
	"go/build"
	"log"
	"os"
	"regexp"
	"strings"

	"mvdan.cc/sh/interp"
	"mvdan.cc/sh/syntax"
)

var (
	regDomain = regexp.MustCompile(`^[a-z\.a-z]+`)

	Verbose = false
	DryRun  = false
)

// URL will try to add `github.com` to url if bad url.
func URL(pkg string) string {
	pkg = strings.TrimSpace(pkg)

	if regDomain.MatchString(pkg) {
		return pkg
	}

	return "github.com/" + strings.TrimLeft(pkg, "/")
}

// Exists check if a package binary exists or not.
func Exists(pkg string) bool {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	path := fmt.Sprintf("%s/bin/%s", gopath, pkg)

	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Run runs a command.
func Run(input string) error {
	if Verbose {
		log.Printf("Running: %s\n", input)
	}

	if DryRun {
		return nil
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	p, err := syntax.NewParser().Parse(strings.NewReader(input), "")
	if err != nil {
		return err
	}

	r := interp.Runner{
		Dir:    path,
		Exec:   interp.DefaultExec,
		Open:   interp.OpenDevImpls(interp.DefaultOpen),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err = r.Reset(); err != nil {
		return err
	}

	return r.Run(p)
}

// Remove removes a go package binary.
func Remove(pkg string) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	path := fmt.Sprintf("%s/%s", gopath, strings.TrimLeft(pkg, "/"))

	if Verbose {
		log.Printf("Removing: %s\n", path)
	}

	if DryRun {
		return nil
	}

	return os.RemoveAll(path)
}
