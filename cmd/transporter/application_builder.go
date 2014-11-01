package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/MongoHQ/transporter/pkg/application"
	"github.com/MongoHQ/transporter/pkg/node"
)

type ApplicationBuilder struct {
	Nodes []node.Node

	// command to run
	Command *Command

	// Running Config
	Config Config

	// path to the config file
	config_path string
}

/*
 * build the application, parse the flags and run the command
 */
func Build() (application.Application, error) {
	builder := ApplicationBuilder{}

	err := builder.flagParse()
	if err != nil || builder.Command == nil {
		builder.usage()
		return nil, err
	}

	err = builder.loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config Error: %s\n", err)
	}

	builder.Nodes = builder.Config.Nodes
	return builder.Command.Run(builder, builder.Command.Flag.Args())
}

/*
 *
 * flag parsing related functions
 *
 */
func (a *ApplicationBuilder) flagParse() error {
	flag.StringVar(&a.config_path, "config", "", "path to the config yaml")
	flag.Usage = a.usage
	flag.Parse()

	if flag.Arg(0) == "" {
		return fmt.Errorf("no command specified")
	}

	if flag.Arg(0) == "help" {
		a.help(flag.Arg(1))
		return nil
	}

	// make sure we're valid
	for _, c := range commands {
		if c.Name == flag.Arg(0) {
			c.Flag.Parse(flag.Args()[1:])

			a.Command = c
			return nil
		}
	}
	return fmt.Errorf("Command '%s' not found", flag.Arg(0))
}

func (a *ApplicationBuilder) usage() {
	t := template.Must(template.New("usage").Parse(usageTpl))
	if err := t.Execute(os.Stderr, commands); err != nil {
		panic(nil)
	}
	os.Exit(0)
}

func (a *ApplicationBuilder) help(which string) error {
	t := template.Must(template.New("help").Parse(helpTpl))

	// find the command
	for _, c := range commands {
		if c.Name == which {
			if err := t.Execute(os.Stderr, c); err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}
	return fmt.Errorf("no such command '%s'", which)
}

var usageTpl = `
Usage:

transporter [global arguments] command [arguments]

commands:
{{range .}}
    {{.Name | printf "%-8s"}} {{.Short}}{{end}}

Use "transporter help [command]" for more information.
`

var helpTpl = `
{{.Name}}

{{.Help}}
`