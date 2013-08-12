package cmd

import (
	"flag"
	"fmt"
	"strings"
)

func init() {
	Add(&Help{})
}

type Help struct {
	*Flagger
}

func (c *Help) Desc() string {
	return "Display help information for a particular command."
}

func (c *Help) SetFlags() {}

func (c *Help) Run() {
	if args := c.FlagSet.Args(); len(args) == 1 {
		if command, ok := commandSet[args[0]]; ok {
			fmt.Println(command.Desc())

			flagDocs := []string{}
			command.SetFlags()
			command.GetFlagSet().VisitAll(func(f *flag.Flag) {
				var defValue string
				if f.DefValue != "" {
					defValue = " Default value: " + f.DefValue + "."
				}

				flagDocs = append(flagDocs, fmt.Sprintf(
					"-%s%s%s%s\n",
					f.Name,
					strings.Repeat(" ", 10-len(f.Name)),
					f.Usage,
					defValue,
				))
			})
			if len(flagDocs) > 0 {
				fmt.Println("\nOptions:")
				for _, fd := range flagDocs {
					fmt.Print(fd)
				}
			}

			return
		}
	}
	fmt.Printf("Usage: '%s help <command>'\n", appName)
}
