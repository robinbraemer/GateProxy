package command

import (
	"fmt"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/color"
	. "go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
	"simple-proxy/minimessage"
)

func newBroadcastCmd(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	announceFirst, _ := color.Hex("#5555ff")
	announceSecond, _ := color.Hex("#5498ff")

	return brigodier.Literal("broadcast").Then(
		brigodier.Argument("message", brigodier.StringPhrase).
			Requires(command.Requires(func(c *command.RequiresContext) bool {
				return c.Source.HasPermission("divine.broadcast")
			})).
			Executes(command.Command(func(c *command.Context) error {
				// Colorize/format message
				component, err := minimessage.LegacyCodec.Unmarshal([]byte(c.String("message")))
				if err != nil {
					return c.Source.SendMessage(&Text{
						Content: fmt.Sprintf("Error formatting message: %v", err)})
				}

				message := &Text{
					Extra: []Component{
						minimessage.Gradient("Announcement", Style{Bold: True}, *announceFirst, *announceSecond),
						&Text{Content: " - ", S: Style{Color: color.Gray}},
						component,
					},
				}

				// Send to all players on this proxy
				for _, player := range p.Players() {
					// Send message in new goroutine to not block
					// this loop if any player has a slow connection.
					go player.SendMessage(message)
				}
				return nil
			})),
	)
}
