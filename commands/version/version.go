/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package version allows people to know which verison of Jarvis it's currently running
package version

import (
	"appengine"

	"github.com/dlsniper/jarvis-slacker/api"
)

type ver struct{}

func (v *ver) Initialize(keys map[string]string) {
}

func (v *ver) Execute(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) (string, error) {
	return `Jarvis version 0.4`, nil
}

func (v *ver) Usage(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) string {
	return `- version -> returns current version of bot
`
}

func (v *ver) Help(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) (string, error) {
	return `Display the current version of Jarvis`, nil
}

func NewCommand() ([]string, api.Command) {
	return []string{"version"}, &ver{}
}
