/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package help allows Jarvis to have help functionality
package help

import (
	"strings"

	"appengine"

	"github.com/dlsniper/jarvis-slacker/api"
)

type (
	Help struct {
		commandHelp map[string]api.Command
	}
)

func (h *Help) Initialize(keys map[string]string) {
}

func (h *Help) RegisterCommandHelp(command string, handler api.Command) {
	h.commandHelp[command] = handler
}

func (h *Help) Execute(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	cmd := ""

	if command.LowerText == "help" {
		return h.Help(ctx, command, "help")
	}

	if len(command.LowerText) >= 6 {
		cmd = command.LowerText[5:]
	} else {
		return h.Help(ctx, command, "help")
	}

	module := cmd
	if strings.Index(cmd, " ") != -1 {
		module = cmd[:strings.Index(cmd, " ")]
	}

	if c, ok := h.commandHelp[module]; ok {
		return c.Help(ctx, command, cmd)
	}

	return h.Help(ctx, command, "help")
}

func (h *Help) Usage(ctx appengine.Context, command *api.SlackCommand, myCommand string) string {
	return `- help -> General usage of Jarvis and it's various supported command
- help [command] -> Displays usage of a certain command
`
}

func (h *Help) Help(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	result := `List of current supported commands.
If you want to learn more about a commands, use /jarvis help [command]

Supported commands:
`
	for idx := range h.commandHelp {
		result += h.commandHelp[idx].Usage(ctx, command, idx)
	}
	return result, nil
}

func NewCommand() ([]string, api.Command) {
	return []string{"help"}, &Help{
		commandHelp: map[string]api.Command{},
	}
}
