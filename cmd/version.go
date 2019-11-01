package cmd

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionTemplate = `Version:           {{.AppVersion}}
Go version:        {{.GoVersion}}
Git commit:        {{.GitCommit}}
Built:             {{.BuildTime}}
OS/Arch:           {{.GoOs}}/{{.GoArch}}`

var (
	AppVersion = "unknown-app-version"
	GoVersion  = "unknown-go-version"
	GoOs       = "unknown-go-os"
	GoArch     = "unknown-go-arch"
	GitCommit  = "unknown-git-commit"
	BuildTime  = "unknown-build-time"
)

type Version struct {
	AppVersion string
	GoVersion  string
	GoOs       string
	GoArch     string
	GitCommit  string
	BuildTime  string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of wait4x",
	Long:  "All software has versions. It's mine.",
	Run: func(cmd *cobra.Command, args []string) {
		versionValues := Version{
			AppVersion: AppVersion,
			GoVersion: GoVersion,
			GoOs: GoOs,
			GoArch: GoArch,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		}
		var tmplBytes bytes.Buffer

		t := template.Must(template.New("version").Parse(versionTemplate))
		err := t.Execute(&tmplBytes, versionValues)
		if err != nil {
			log.Println("executing template:", err)
		}

		fmt.Println(tmplBytes.String())
	},
}
