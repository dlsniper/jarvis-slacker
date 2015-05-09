/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package tell

import (
	"appengine"

	"github.com/dlsniper/jarvis-slacker/api"
	"github.com/dlsniper/jarvis-slacker/commands/tell/gif"
	"github.com/dlsniper/jarvis-slacker/commands/tell/joke"
)

func ComputeWhatResult(ctx appengine.Context, command *api.SlackCommand, what string) (string, error) {
	switch {
	case what == "a joke":
		return joke.Joke(ctx, command)
	case what == "a gif":
		return gif.Gif(ctx, command)
	}

	// If we can't find a match we treat it as a mirror
	return what, nil
}
