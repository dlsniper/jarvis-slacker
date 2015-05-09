/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package jarvis holds the bindings between all the other packages to form the bot
package jarvis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dlsniper/jarvis-slacker/api"
	"github.com/dlsniper/jarvis-slacker/commands"

	"github.com/dlsniper/jarvis-slacker/slack"

	"appengine"
)

var (
	slashCommand = `/jarvis`
	slashToken   = ``
)

func slashJarvis(ctx appengine.Context, w http.ResponseWriter, r *http.Request, command *api.SlackCommand) {
	var (
		response = ""
		err      error
	)

	cmd := ""
	if strings.Index(command.LowerText, " ") == -1 {
		cmd = command.LowerText
	} else {
		cmd = command.LowerText[:strings.Index(command.LowerText, " ")]
	}

	if c := commands.HasCommand(cmd); c != nil {
		response, err = c.Execute(ctx, command, cmd)
	} else {
		response = fmt.Sprintf("Command %q not found. Try %q for a list of available commands", command.Text, slashCommand+" help")
	}

	ctx.Infof("Got command %#v", command)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, response)
}

func SlashCommandHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "noup")
		return
	}

	command := &api.SlackCommand{
		Token:       r.PostForm.Get("token"),
		TeamID:      r.PostForm.Get("team_id"),
		TeamDomain:  r.PostForm.Get("team_domain"),
		ChannelID:   r.PostForm.Get("channel_id"),
		ChannelName: r.PostForm.Get("channel_name"),
		UserID:      r.PostForm.Get("user_id"),
		UserName:    r.PostForm.Get("user_name"),
		Command:     r.PostForm.Get("command"),
		Text:        strings.Trim(r.PostForm.Get("text"), " \t\r\n"),
	}
	command.LowerText = strings.ToLower(command.Text)

	if command.Token != slashToken {
		http.Error(w, "noup", http.StatusUnauthorized)
		return
	}

	if command.Command == slashCommand {
		ctx := appengine.NewContext(r)
		slashJarvis(ctx, w, r, command)
		return
	}

	fmt.Fprintf(w, "Command %q not found. Did you meant %q?", command.Command, slashCommand)
}

func RegisterCommand(command api.JarvisRegisterFunction) {
	commands.AddCommand(command)
}

func Initialize(command, botName string, keys map[string]string) {
	if _, ok := keys["slackToken"]; !ok {
		panic("Slack token (slackToken) is missing")
	}
	if _, ok := keys["slashToken"]; !ok {
		panic("Slack slash token (slashToken) is missing")
	} else {
		slashToken = keys["slashToken"]
	}

	slashCommand = command
	slack.Initialize(botName, keys["slackToken"], keys["slashToken"])

	commands.Initialize(keys)
}
