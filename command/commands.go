package command

import "github.com/urfave/cli/v2"

// RegisterCommands 注册所有命令
func RegisterCommands(app *cli.App) {
	app.Commands = append(app.Commands,
		&Test,
	)
}
