package gen

import (
	"os"
	"text/template"
)

const genTemplate = `
Build version: <{{.BuildVersion}}>
Build date: <{{.BuildDate}}>
Build commit: <{{.BuildCommit}}>
`

type Gen struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

var tmpl = template.Must(template.New("gen").Parse(genTemplate))

func InitVersion(buildVersion, buildDate, buildCommit string) error {
	greetings := Gen{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}

	if greetings.BuildVersion == "" {
		greetings.BuildVersion = "N/A"
	}
	if greetings.BuildDate == "" {
		greetings.BuildDate = "N/A"
	}
	if greetings.BuildCommit == "" {
		greetings.BuildCommit = "N/A"
	}

	return tmpl.Execute(os.Stdout, greetings)
}
