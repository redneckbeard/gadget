package cmd

import (
	"flag"
	"fmt"
	"github.com/redneckbeard/gadget/strutil"
	"os"
	"reflect"
)

var (
	appName string
	commandSet = make(map[string]Command)
)

type Flagger struct {
	*flag.FlagSet
}

func (f *Flagger) GetFlagSet() *flag.FlagSet {
	return f.FlagSet
}

type Command interface {
	Desc() string
	SetFlags()
	Run()
	Parse([]string) error
	GetFlagSet() *flag.FlagSet
}

func Add(commands ...Command) {
	for _, c := range commands {
		t := reflect.TypeOf(c).Elem()
		v := reflect.ValueOf(c).Elem()
		flagger := v.FieldByName("Flagger")
		flagger.Set(reflect.ValueOf(&Flagger{&flag.FlagSet{}}))
		name := strutil.Hyphenate(t.Name())
		commandSet[name] = c
		if len(os.Args) > 1 && name == os.Args[1] && name != "help" {
			c.SetFlags()
			c.Parse(os.Args[2:])
		}
	}
}

func Run() {
	appName = os.Args[0]
	if len(os.Args) == 1 {
		fmt.Println("Available commands:")
		for name, _ := range commandSet {
			fmt.Println("  " + name)
		}
		fmt.Printf("Type '%s help <command>' for more information on a specific command.\n", os.Args[0])
	} else {
		command := os.Args[1]
		if command == "help" {
			c := commandSet["help"]
			c.Parse(os.Args[2:])
			c.Run()
		} else {
			for name, c := range commandSet {
				if name == command {
					c.Run()
				}
			}
		}
	}
}
