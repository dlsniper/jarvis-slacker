/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package commands howlds the bindings for all the commands Jarvis knows how to handle
package commands

import (
	"github.com/dlsniper/jarvis-slacker/api"

	"github.com/dlsniper/jarvis-slacker/commands/help"
	"github.com/dlsniper/jarvis-slacker/commands/version"
)

var (
	// Append new modules to this variable
	commands = []api.JarvisRegisterFunction{
		version.NewCommand,
		help.NewCommand,
	}

	cmds = map[string]api.Command{}
)

func Initialize(keys map[string]string) {
	for idx := range commands {
		cmd, handler := commands[idx]()
		handler.Initialize(keys)
		for _, cm := range cmd {
			cmds[cm] = handler
		}
	}

	if _, ok := cmds["help"]; !ok {
		panic("you must have a help commands")
	}

	for idx := range cmds {
		cmds["help"].(*help.Help).RegisterCommandHelp(idx, cmds[idx])
	}
}

func AddCommand(command api.JarvisRegisterFunction) {
	commands = append(commands, command)
}

func HasCommand(name string) api.Command {
	if command, ok := cmds[name]; ok {
		return command
	}

	return nil
}
