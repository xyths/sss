package utils

import "github.com/urfave/cli/v2"

var (
	alex = &cli.Author{
		Name:  "Alexander Xing",
		Email: "AlexanderXing@gmail.com",
	}
	Authors = []*cli.Author{alex,}
)

func Info(app *cli.App) {
	app.Copyright = "©2020 The Pangu Foundation of SERO community."
	app.Authors = Authors
}
