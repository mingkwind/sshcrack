package cmd

import (
	"sshcrack/utils"

	"github.com/urfave/cli"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to crack wwak password",
	Description: "start to crack weak password",
	Action:      utils.Scan,
	Flags: []cli.Flag{
		boolFlag("debug,d", "debug mode"),
		intFlag("time, t", 3, "timeout"),
		intFlag("thread_num,n", 50, "thread num"),
		stringFlag("ip_list,i", "ip_list.txt", "iplist"),
		stringFlag("user_dict,u", "user.txt", "user dict"),
		stringFlag("pass_dict,p", "pass.txt", "password dict"),
		stringFlag("outfile,o", "pass_crack.txt", "scan result file"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
