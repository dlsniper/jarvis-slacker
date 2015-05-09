/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package show enables Jarvis to show various thinkgs back to the slack channel
package show

import (
	"errors"
	"fmt"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"github.com/dlsniper/jarvis-slacker/api"
	"github.com/dlsniper/jarvis-slacker/commands/tell"
	"github.com/dlsniper/jarvis-slacker/commands/tell/gif"
	"github.com/dlsniper/jarvis-slacker/slack"
)

type shw struct{}

func (s *shw) Initialize(keys map[string]string) {
	var ok bool
	if gif.GiphyKey, ok = keys["giphy"]; !ok {
		panic("Giphy key is missing")
	}
}

func (s *shw) Execute(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	cmd := strings.TrimLeft(command.LowerText, "show ")
	user := cmd[:strings.Index(cmd, " ")]
	what := command.Text[len("show "+user+" "):]

	client := urlfetch.Client(ctx)

	if what != "a gif" {
		slack.SendUserMessage(ctx, client, command.UserName, fmt.Sprintf(`I don't know how to handle %q`, what))
		return "", errors.New(fmt.Sprintf(`I don't know how to handle %q`, what))
	}

	ctx.Infof("Should show %q %q", user, what)

	var (
		code         int
		body         []byte
		err          error
		computedWhat string
	)

	switch {
	case strings.HasPrefix(user, "#"):
		computedWhat, err = tell.ComputeWhatResult(ctx, command, what)
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to compute %q, got %q", what, err))
		}

		code, body, err = slack.SendChannelMessage(ctx, client, user, computedWhat)
		if code != 200 && err == nil {
			err = errors.New(fmt.Sprintf("Got %d with body %q", code, string(body)))
		}
	case strings.HasPrefix(user, "@"):
		computedWhat, err = tell.ComputeWhatResult(ctx, command, what)
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to compute %q, got %q", what, err))
		}

		code, body, err = slack.SendUserMessage(ctx, client, strings.TrimLeft(user, "@"), computedWhat)
		if code != 200 && err == nil {
			err = errors.New(fmt.Sprintf("Got %d with body %q", code, string(body)))
		}
	case user == "me":
		computedWhat, err = tell.ComputeWhatResult(ctx, command, what)
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to compute %q, got %q", what, err))
		}

		code, body, err = slack.SendUserMessage(ctx, client, command.UserName, computedWhat)
		if code != 200 && err == nil {
			err = errors.New(fmt.Sprintf("Got %d with body %q", code, string(body)))
		}
	default:
		err = errors.New(fmt.Sprintf("I don't know who %q is", user))
	}

	return "", err
}

func (s *shw) Usage(ctx appengine.Context, command *api.SlackCommand, myCommand string) string {
	return `- show [#channel / @user] [something] -> will compute [something] and show it to #channel / @user
`
}

func (s *shw) Help(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	result := `show will compute the operation you want to be told and show it to the desired target.

A target can be a user, if you prefix it with @, or a channel, if you prefix it with #.

`
	result += gif.Help(ctx, command) + "\n"

	return result, nil
}

func NewCommand() ([]string, api.Command) {
	return []string{"show"}, &shw{}
}
