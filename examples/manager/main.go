package main

import (
	"github.com/yeswearebot/yes-core/core"

	_ "github.com/yeswearebot/yes-core/examples/manager/plugin_ai"
	_ "github.com/yeswearebot/yes-core/examples/manager/plugin_manager"
	_ "github.com/yeswearebot/yes-core/examples/manager/plugin_onebot"
)

func main() {
	app := core.NewApp()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
