/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package api defines the api for the commands we support
package api

import "appengine"

type (
	SlackCommand struct {
		Token       string
		TeamID      string
		TeamDomain  string
		ChannelID   string
		ChannelName string
		UserID      string
		UserName    string
		Command     string
		Text        string
		LowerText   string
	}

	JarvisRegisterFunction func() ([]string, Command)
	JarvisFunction         func(appengine.Context, *SlackCommand, string) (string, error)

	Command interface {
		Initialize(keys map[string]string)
		Execute(appengine.Context, *SlackCommand, string) (string, error)
		Usage(appengine.Context, *SlackCommand, string) string
		Help(appengine.Context, *SlackCommand, string) (string, error)
	}
)
