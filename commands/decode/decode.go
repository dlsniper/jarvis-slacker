/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package decode does decoding for base64 encoded strings
package decode

import (
	"encoding/base64"

	"appengine"

	"github.com/dlsniper/jarvis-slacker/api"
)

type (
	ci struct{}
)

func (c *ci) Initialize(keys map[string]string) {
}

func (c *ci) Execute(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) (string, error) {
	what := slackCommand.Text[len(myCommand)+1:]

	data, err := base64.URLEncoding.DecodeString(what)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *ci) Usage(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) string {
	return `- decode base64EncodedString -> returns the base64 decoded string
`
}

func (c *ci) Help(ctx appengine.Context, slackCommand *api.SlackCommand, myCommand string) (string, error) {
	result := `Displays the base64 decode.

The following options are available:
`
	result += c.Usage(ctx, slackCommand, myCommand)

	return result, nil
}

func NewCommand() ([]string, api.Command) {
	return []string{"decode"}, &ci{}
}
