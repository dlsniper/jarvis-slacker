#Jarvis Slacker

Jarvis Slacker is a bot for Slack running on Google AppEngine written in Go.

It has an easy to write plugin system where you can register your own commands and aliases.

## Development

Developing a custom plugin for Jarvis is very easy.

Unfortunately, there's no nice way, that I can think of right now to handle the initialization code for it.

As such, you'll need to first copy the ``` frontend/frontend.go ``` to your AppEngine module.
 
### Register a command

Then, you can register your commands like this: ```jarvis.RegisterCommand(decode.NewCommand)```

### Initialize jarvis

A part of the initialization is sending the keys to various modules. This is need as there's no nice way
to have custom variables in AppEngine (as far as I know, yet) and you might want to publish various modules for it.

Once you have the list of all the keys in place, simply call: 

```go
jarvis.Initialize("/jarvis", "jarvis", keys)
http.HandleFunc("/slashCommand", jarvis.SlashCommandHandler)
```

And that's it.

For an example on how this could look like, you can see [jarvis-slacker-3d](https://github.com/dlsniper/jarvis-slacker-3d)

## Deployment

First you need to: ```cp app.yaml.dist app.yaml```

Then you need to replace ```<your AppEngine project ID goes here>``` from app.yaml with your AppEngine project id. 

You can deploy this bot like this:

```bash
cd /repository/root/path
goapp deploy frontend
```

##LICENSE

Jarvis Slacker is distributed under MIT license.

You can see the full license in the [LICENSE](LICENSE.md) file.
