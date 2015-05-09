/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package frontend enables all the functionality
package frontend

import (
	"net/http"

	"github.com/dlsniper/jarvis-slacker/commands/decode"
	"github.com/dlsniper/jarvis-slacker/commands/show"
	"github.com/dlsniper/jarvis-slacker/commands/tell"
	"github.com/dlsniper/jarvis-slacker/jarvis"
)

func init() {
	jarvis.RegisterCommand(decode.NewCommand)
	jarvis.RegisterCommand(show.NewCommand)
	jarvis.RegisterCommand(tell.NewCommand)

	jarvis.Initialize("/jarvis", "jarvis", keys)

	http.HandleFunc("/slashCommand", jarvis.SlashCommandHandler)
}
