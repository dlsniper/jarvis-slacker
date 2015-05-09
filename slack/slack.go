/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"appengine"
)

type (
	imOpenResponse struct {
		OK      bool `json:"ok"`
		Channel struct {
			ID string `json:"id"`
		} `json:"channel"`
	}

	userListResponse struct {
		OK      bool `json:"ok"`
		Members []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"members"`
	}
)

var (
	botName        = `jarvis`
	joinedChannels = map[string]bool{}
	joinedUsers    = map[string]string{}
	slackTkn       = ``
	slashTkn       = ``
)

func Initialize(name, slackToken, slashToken string) {
	botName = name
	slackTkn = slackToken
	slashTkn = slashToken
}

func SlackUrl(ctx appengine.Context, url string, args ...string) string {
	args = append(args, "username="+botName)
	extraArgs := strings.Join(args, "&")
	ur := fmt.Sprintf("https://slack.com/api/%s?token=%s&%s", url, slackTkn, extraArgs)
	ctx.Infof("Generated URL: %q\n", ur)
	return ur
}

func ExecuteGetRequest(ctx appengine.Context, client *http.Client, url string) (int, []byte, error) {
	ctx.Infof("Running request %q", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ctx.Errorf("Got error %q while doing the request", err)
		return 0, []byte{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "dlsniper/jarvis-slacker slack bot")
	resp, err := client.Do(req)
	if err != nil {
		ctx.Errorf("Got error %q while doing the request", err)
		return 0, []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Errorf("Got error %q while reading the body", err)
		return 0, []byte{}, err
	}

	return resp.StatusCode, body, nil
}

func ExecutePostRequest(ctx appengine.Context, client *http.Client, url string, headers map[string]string, payload []byte) (int, []byte, error) {
	ctx.Infof("Running request %q with payload: %q", url, string(payload))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		ctx.Errorf("Got error %q while doing the request", err)
		return 0, []byte{}, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "dlsniper/jarvis-slacker slack bot")
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		ctx.Errorf("Got error %q while doing the request", err)
		return 0, []byte{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.Errorf("Got error %q while reading the body", err)
		return 0, []byte{}, err
	}

	return resp.StatusCode, body, nil
}

func SendChannelMessage(ctx appengine.Context, client *http.Client, channel, message string) (int, []byte, error) {
	channel = url.QueryEscape(channel)
	message = url.QueryEscape(message)

	return ExecuteGetRequest(ctx,
		client,
		SlackUrl(ctx, "chat.postMessage",
			fmt.Sprintf("channel=%s", channel),
			fmt.Sprintf("text=%s", message),
		),
	)
}

func GetSlackUserID(ctx appengine.Context, client *http.Client, username string) (string, error) {
	code, body, err := ExecuteGetRequest(ctx, client, SlackUrl(ctx, "users.list"))
	if err != nil {
		return "", err
	}
	if code != 200 {
		return "", errors.New(fmt.Sprintf("failed to get the username list, got %q", string(body)))
	}

	response := &userListResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", err
	}
	if !response.OK {
		return "", errors.New(fmt.Sprintf("something went wrong when getting the user list, got %q", string(body)))
	}

	for idx := range response.Members {
		if response.Members[idx].Name == username {
			return response.Members[idx].ID, nil
		}
	}

	return "", errors.New(fmt.Sprintf("username %q not found in slack user list", username))
}

func OpenIM(ctx appengine.Context, client *http.Client, username string) (string, error) {
	if channelID, ok := joinedUsers[username]; ok {
		ctx.Infof("returning existing user direct channel %q\n", channelID)
		return channelID, nil
	}

	userID, err := GetSlackUserID(ctx, client, username)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to get the user id for %q", username))
	}

	code, body, err := ExecuteGetRequest(ctx,
		client,
		SlackUrl(ctx, "im.open",
			fmt.Sprintf("user=%s", userID),
		),
	)
	if err != nil {
		return "", err
	}

	if code != 200 {
		return "", errors.New(fmt.Sprintf("unexpected reponse recevived: %q", string(body)))
	}

	response := &imOpenResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return "", err
	}

	if !response.OK {
		return "", errors.New(fmt.Sprintf("unexpected reponse recevived: %q", string(body)))
	}

	joinedChannels[response.Channel.ID] = true
	joinedUsers[username] = response.Channel.ID

	return response.Channel.ID, nil
}

func SendUserMessage(ctx appengine.Context, client *http.Client, username, message string) (int, []byte, error) {
	channel, err := OpenIM(ctx, client, username)
	if err != nil {
		return 0, []byte{}, err
	}
	channel = url.QueryEscape(channel)
	message = url.QueryEscape(message)

	return ExecuteGetRequest(ctx,
		client,
		SlackUrl(ctx, "chat.postMessage",
			"as_user=true",
			fmt.Sprintf("channel=%s", channel),
			fmt.Sprintf("text=%s", message),
		),
	)
}

func Link(url, title string) string {
	if title == "" {
		title = url
	}
	return fmt.Sprintf("<%s|%s>", url, title)
}
