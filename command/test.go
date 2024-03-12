package command

import (
	"fmt"

	ei18n "github.com/mylukin/easy-i18n/i18n"
	"github.com/urfave/cli/v2"
)

var Test = cli.Command{
	Name:  "test",
	Usage: ei18n.Sprintf("test"),
	Action: func(c *cli.Context) error {
		return ExecuteTest()
	},
}

// ExecuteCmd1 执行命令逻辑
func ExecuteTest() error {
	fmt.Println(ei18n.Sprintf("Executing test command..."))
	// 实现命令逻辑
	return nil
}
