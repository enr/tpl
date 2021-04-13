package main

/*
>go run main.go commands.go version.go utils.go statusfull
*/

import (
	"fmt"
	"os"

	"github.com/enr/clui"
	"github.com/urfave/cli/v2"

	"github.com/enr/tpl/lib/core"
)

var (
	ui              *clui.Clui
	versionTemplate = `%s
Revision: %s
Build date: %s
`
	appVersion        = fmt.Sprintf(versionTemplate, core.Version, core.GitCommit, core.BuildTime)
	ignoreMissingDirs bool
)

func main() {
	app := cli.NewApp()
	app.Name = "tpl"
	app.Version = appVersion
	app.Usage = "Template processor"
	app.Flags = []cli.Flag{

		&cli.BoolFlag{Name: "debug", Aliases: []string{"D"}, Usage: "operates in debug mode: lot of output"},
		&cli.BoolFlag{Name: "quiet", Aliases: []string{"q"}, Usage: "operates in quiet mode"},

		&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Usage: "path to template file"},
		&cli.StringFlag{Name: "destination", Aliases: []string{"d"}, Usage: "path to write the final output"},

		&cli.StringSliceFlag{Name: "varfile", Aliases: []string{"P"}, Usage: "path to file containing variables (in properties format)"},
		&cli.StringSliceFlag{Name: "var", Aliases: []string{"V"}, Usage: "var"},

		&cli.BoolFlag{Name: "inline", Aliases: []string{"I"}, Usage: "over write template file"},
		&cli.BoolFlag{Name: "stdout", Aliases: []string{"O"}, Usage: "write output to /dev/stdout"},
		&cli.BoolFlag{Name: "no-indent", Aliases: []string{"i"}, Usage: "do not keep indentation"},
		&cli.StringFlag{Name: "separator", Aliases: []string{"Z"}, Value: `:`, Usage: "placeholder separator"},
		&cli.StringFlag{Name: "start-delimiter", Aliases: []string{"S", "start"}, Usage: "start placeholder delimiter ${"},
		&cli.StringFlag{Name: "end-delimiter", Aliases: []string{"E", "end"}, Usage: "end placeholder delimiter }"},
	}
	app.EnableBashCompletion = true

	app.Action = doProcess

	app.Before = func(c *cli.Context) error {
		if ui != nil {
			return nil
		}
		verbosityLevel := clui.VerbosityLevelMedium
		if c.Bool("debug") {
			verbosityLevel = clui.VerbosityLevelHigh
		}
		if c.Bool("quiet") {
			verbosityLevel = clui.VerbosityLevelLow
		}
		if c.Bool("stdout") {
			verbosityLevel = clui.VerbosityLevelMute
		}
		var err error
		ui, err = clui.NewClui(func(ui *clui.Clui) {
			ui.VerbosityLevel = verbosityLevel
		})
		if err != nil {
			return err
		}
		ui, err = core.ConfigureAppUI(func(ui *clui.Clui) {
			ui.VerbosityLevel = verbosityLevel
		})
		return err
	}

	app.Commands = commands

	// returns an error
	app.Run(os.Args)
}
