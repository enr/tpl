package main

import (
	"fmt"
	"strings"

	"github.com/enr/tpl/lib/core"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	&processCommand,
}

var processCommand = cli.Command{
	Name:        "process",
	Aliases:     []string{"p"},
	Usage:       "tpl process [OPTIONS]",
	Description: `Process the template`,
	Action:      doProcess,
}

func doProcess(c *cli.Context) error {
	request, err := requestFromUserOptions(c)
	if err != nil {
		return exitErrorf(1, "Error creating request: %v", err)
	}
	ui.Confidentialf("Command process, request: %v", request)

	err = core.Process(request)
	if err != nil {
		return exitErrorf(1, "Error process: %v", err)
	}
	return nil
}

func requestFromUserOptions(c *cli.Context) (core.ProcessRequest, error) {

	source := c.String(`source`)
	destination := c.String(`destination`)
	varfiles := c.StringSlice(`varfile`)
	format := c.String(`format`)
	inline := c.Bool(`inline`)
	stdout := c.Bool(`stdout`)
	skipIndent := c.Bool(`no-indent`)
	separator := c.String(`separator`)
	startDelimiter := c.String(`start-delimiter`)
	endDelimiter := c.String(`end-delimiter`)
	vars := c.StringSlice(`var`)

	var kv []string
	vm := make(map[string]string)
	for _, v := range vars {
		kv = strings.SplitN(v, `=`, 2)
		vm[kv[0]] = kv[1]
	}
	return core.ProcessRequest{
		Source:               source,
		Destination:          destination,
		Varfiles:             varfiles,
		Vars:                 vm,
		Format:               format,
		Inline:               inline,
		Stdout:               stdout,
		SkipIndent:           skipIndent,
		PlaceholderSeparator: separator,
		StartDelimiter:       startDelimiter,
		EndDelimiter:         endDelimiter,
	}, nil
}

func exitErrorf(exitCode int, template string, args ...interface{}) error {
	ui.Errorf(`Something gone wrong:`)
	return cli.NewExitError(fmt.Sprintf(template, args...), exitCode)
}
