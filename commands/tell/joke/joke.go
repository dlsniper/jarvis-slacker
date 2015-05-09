/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package joke returns a random joke from icndb.com
package joke

import (
	"encoding/json"
	"errors"
	"fmt"

	"appengine"
	"appengine/urlfetch"

	"github.com/dlsniper/jarvis-slacker/api"
	"github.com/dlsniper/jarvis-slacker/slack"
)

type (
	joke struct {
		Type  string `json:"type"`
		Value struct {
			Joke string `json:"joke"`
		} `json:"value"`
	}
)

func Joke(ctx appengine.Context, command *api.SlackCommand) (string, error) {
	client := urlfetch.Client(ctx)
	code, body, err := slack.ExecuteGetRequest(ctx, client, "http://api.icndb.com/jokes/random")
	if err != nil {
		return "", err
	}

	if code != 200 {
		return "", errors.New(fmt.Sprintf("Got %d, %q", code, string(body)))
	}

	joke := &joke{}
	err = json.Unmarshal(body, joke)
	if err != nil {
		return "", err
	}

	return joke.Value.Joke, nil
}

func Help(ctx appengine.Context, command *api.SlackCommand) string {
	return `joken -> will return a random joke from the Internet`
}
