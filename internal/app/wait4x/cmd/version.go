// Copyright 2019 Mohammad Abdolirad
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"text/template"

	"github.com/spf13/cobra"
)

var versionTemplate = `Version:           {{.AppVersion}}
Go version:        {{.GoVersion}}
Git commit:        {{.GitCommit}}
Built:             {{.BuildTime}}
OS/Arch:           {{.GoOs}}/{{.GoArch}}`

var (
	// AppVersion represents Wait4X version
	AppVersion = "$Format:%(describe:tags=true,exclude=*-[0-9]*-g[0-9a-f]*)$"
	// GitCommit represents Wait4X commit sha1 hash from git, output of $(git rev-parse HEAD)
	GitCommit = "$Format:%H$"
	// BuildTime represents Wait4X build time in ISO8601 format, output of $(date -u '+%FT%TZ')
	BuildTime = "1970-01-01T00:00:00Z"
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
