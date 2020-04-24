package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"runtime"
	"text/template"
)

var versionTemplate = `Version:           {{.AppVersion}}
Go version:        {{.GoVersion}}
Git commit:        {{.GitCommit}}
Built:             {{.BuildTime}}
OS/Arch:           {{.GoOs}}/{{.GoArch}}`

var (
	// AppVersion represents Wait4X version
	AppVersion = "unknown-app-version"
	// GitCommit represents Wait4X commit hash
	GitCommit = "unknown-git-commit"
	// BuildTime represents Wait4X build time
	BuildTime = "unknown-build-time"
)

// Version represents some information which useful in version sub-command
type Version struct {
	AppVersion string
	GoVersion  string
	GoOs       string
	GoArch     string
	GitCommit  string
	BuildTime  string
}

// NewVersionCommand creates the version sub-command
func NewVersionCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Wait4X",
		Long:  "All software has versions. It's mine.",
		Run:   runVersion,
	}

	return versionCommand
}

func runVersion(_ *cobra.Command, _ []string) {
	versionValues := Version{
		AppVersion: AppVersion,
		GoVersion:  runtime.Version(),
		GoOs:       runtime.GOOS,
		GoArch:     runtime.GOARCH,
		GitCommit:  GitCommit,
		BuildTime:  BuildTime,
	}
	var tmplBytes bytes.Buffer

	t := template.Must(template.New("version").Parse(versionTemplate))
	err := t.Execute(&tmplBytes, versionValues)
	if err != nil {
		log.Println("executing template:", err)
	}

	fmt.Println(tmplBytes.String())
}
