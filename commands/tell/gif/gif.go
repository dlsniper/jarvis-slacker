/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package gif displays a random gif from giphy.com
package gif

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
	gif struct {
		Data struct {
			Image string `json:"image_original_url"`
		} `json:"data"`
	}
)

var (
	GiphyKey = ""
	gifURL   = "http://api.giphy.com/v1/gifs/random?api_key="
)

func Gif(ctx appengine.Context, command *api.SlackCommand) (string, error) {
	client := urlfetch.Client(ctx)
	code, body, err := slack.ExecuteGetRequest(ctx, client, gifURL + GiphyKey)
	if err != nil {
		return "", err
	}

	if code != 200 {
		return "", errors.New(fmt.Sprintf("Got %d, %q", code, string(body)))
	}

	gif := &gif{}
	err = json.Unmarshal(body, gif)
	if err != nil {
		return "", err
	}

	return gif.Data.Image, nil
}

func Help(ctx appengine.Context, command *api.SlackCommand) string {
	return `gif -> will display a random gif from the Internet`
}
