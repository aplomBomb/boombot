package discord

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func help(data *disgord.MessageCreate, args []string) {

	defaultHelp(data)

}

func defaultHelp(data *disgord.MessageCreate) {
	_, err := client.SendMsg(
		ctx,
		data.Message.ChannelID,
		&disgord.CreateMessageParams{
			Embed: &disgord.Embed{
				Title: "**__help__**\n **ALIASES:** h, ?, wtf",
				Description: fmt.Sprintf(
					"This is the help function.\n\n"+
						"Use `%shelp functionName` to find out more about each function\n"+
						"Current available functions : ```\nhelp \n```"+
						"You can also read the source code here : https://github.com/aplombomb/boombot",
					conf.Prefix,
				),
				Footer: &disgord.EmbedFooter{
					Text: fmt.Sprintf("Help requested by %s", data.Message.Author.Username),
				},
				Timestamp: data.Message.Timestamp,
				Color:     0xeec400,
			},
		},
	)
	if err != nil {
		fmt.Println("There was an error sending default help message: ", err)
	}
}
