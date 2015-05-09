/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package tell allows Jarvis to tell people what they need
package tell

import (
	"errors"
	"fmt"
	"strings"

	"appengine"
	"appengine/urlfetch"

	"github.com/dlsniper/jarvis-slacker/api"
	"github.com/dlsniper/jarvis-slacker/commands/tell/joke"
	"github.com/dlsniper/jarvis-slacker/slack"
)

type tll struct{}

func (t *tll) Initialize(keys map[string]string) {
}

func (t *tll) Execute(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	cmd := strings.TrimLeft(command.LowerText, myCommand+" ")
	user := cmd[:strings.Index(cmd, " ")]

	what := command.Text[len(myCommand+" "+user+" "):]

	ctx.Infof("Should tell %q %q", user, what)

	var (
		code         int
		body         []byte
		err          error
		computedWhat string
		client       = urlfetch.Client(ctx)
	)

	switch {
	case strings.HasPrefix(user, "#"):
		computedWhat, err = ComputeWhatResult(ctx, command, what)
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to compute %q, got %q", what, err))
		}

		code, body, err = slack.SendChannelMessage(ctx, client, user, computedWhat)
		if code != 200 && err == nil {
			err = errors.New(fmt.Sprintf("Got %d with body %q", code, string(body)))
		}
	case strings.HasPrefix(user, "@"):
		computedWhat, err = ComputeWhatResult(ctx, command, what)
		if err != nil {
			err = errors.New(fmt.Sprintf("Failed to compute %q, got %q", what, err))
		}

		code, body, err = slack.SendUserMessage(ctx, client, strings.TrimLeft(user, "@"), computedWhat)
		if code != 200 && err == nil {
			err = errors.New(fmt.Sprintf("Got %d with body %q", code, string(body)))
		}
	case user == "me":
		computedWhat, err = ComputeWhatResult(ctx, command, what)
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

	if err != nil {
		ctx.Errorf("Got error %q", err)
	}
	return "", err
}

func (t *tll) Usage(ctx appengine.Context, command *api.SlackCommand, myCommand string) string {
	switch {
	case myCommand == "say":
		return `- say [#channel / @user] [something] -> will compute [something] and tell it to #channel / @user
`
	default:
		return `- tell [#channel / @user] [something] -> will compute [something] and tell it to #channel / @user
`
	}
}

func (t *tll) Help(ctx appengine.Context, command *api.SlackCommand, myCommand string) (string, error) {
	result := fmt.Sprintf(`%s will compute the operation you want to be told and say it to the desired target.

A target can be a user, if you prefix it with @, or a channel, if you prefix it with #.

`, myCommand)

	result += joke.Help(ctx, command) + "\n"

	return result, nil
}

func NewCommand() ([]string, api.Command) {
	return []string{"tell", "say"}, &tll{}
}
