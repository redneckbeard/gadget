package cmd

import (
	"fmt"
	"os"
)

var commands = make(map[string]Command)

type Command interface {
	SetFlags()
	Run()
	Parse([]string) error
}

func Add(name string, c Command) {
	commands[name] = c
	if len(os.Args) > 1 && name == os.Args[1] {
		c.SetFlags()
		c.Parse(os.Args[2:])
	}
}

func Run() {
	if len(os.Args) == 1 {
		for name, _ := range commands {
			fmt.Println(name)
		}
	} else {
		command := os.Args[1]
		for name, c := range commands {
			if name == command {
				c.Run()
			}
		}
	}
}
