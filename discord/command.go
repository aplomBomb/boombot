package discord

import (
	"github.com/andersfylling/disgord"
)

//RespondToCommand handles all messages that begin with prefix
func RespondToCommand(s disgord.Session, data *disgord.MessageCreate) {
	cmd, args := ParseMessage(data)

	switch cmd {
	case "help", "h", "?", "wtf":
		help(data, args)
	default:
		unknown(data)
	}

}
