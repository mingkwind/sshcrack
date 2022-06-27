package main

import (
	"os"
	"runtime"
	"sshcrack/cmd"

	"github.com/urfave/cli"
)

func init() {
	// 让协程使用多核CPU
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "password-crack"
	app.Author = "M1ngkvv1nd"
	app.Usage = "Weak ssh password scanner"
	app.Commands = []cli.Command{cmd.Scan}
	app.Flags = append(app.Flags, cmd.Scan.Flags...)
	err := app.Run(os.Args)
	_ = err
}
