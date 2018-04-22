package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/frozzare/gpx/gpx"
	"github.com/spf13/pflag"
)

var (
	dryRunFlag  = pflag.BoolP("dry-run", "n", false, "run in dry run mode but will not actually do anything")
	pkgFlag     = pflag.StringP("package", "p", "", "define the package name instead of guessing")
	rmFlag      = pflag.BoolP("rm", "r", false, "remove binary after execution")
	verboseFlag = pflag.BoolP("verbose", "v", false, "log verbose")
)

const usage = `
Execute go package binaries.

Examples:

  gpx honnef.co/go/tools/cmd/megacheck ./...

Options:

`

func main() {
	pflag.CommandLine.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{
		UnknownFlags: true,
	}
	pflag.Usage = func() {
		log.Print(usage)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	gpx.DryRun = *dryRunFlag
	gpx.Verbose = *verboseFlag || *dryRunFlag

	index := 1
	for _, p := range os.Args[1:] {
		if p[0] == '-' {
			index++
		}
	}

	if len(os.Args) <= index {
		log.Fatal("Missing package argument")
	}

	// Find out package and package name.
	pkg := gpx.URL(os.Args[index])
	if len(*pkgFlag) == 0 {
		parts := strings.Split(pkg, "/")
		*pkgFlag = parts[len(parts)-1]
	}

	// Install package if not existing.
	if !gpx.Exists(*pkgFlag) {
		if err := gpx.Run("go get -u " + pkg); err != nil {
			log.Fatal(err)
		}
	}

	// Add arguments after package string.
	args := []string{}
	if len(os.Args) > index {
		args = os.Args[index+1:]
	}

	args = append([]string{*pkgFlag}, args...)
	cmd := strings.Join(args, " ")

	// Run package.
	if err := gpx.Run(cmd); err != nil {
		log.Fatal(err)
	}

	// Remove package.
	if *rmFlag {
		if err := gpx.Remove(fmt.Sprintf("bin/%s", *pkgFlag)); err != nil {
			log.Fatal(err)
		}

		if err := gpx.Remove(fmt.Sprintf("src/%s", pkg)); err != nil {
			log.Fatal(err)
		}
	}
}
