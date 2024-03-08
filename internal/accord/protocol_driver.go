package accord

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	// ErrConnectionAlreadyConnected is an error when connect() is called but we're already connected
	ErrConnectionAlreadyConnected = errors.New("connection provided to driver already open")
)

// Driver is discord's driver struct for the things
type Driver struct {
	Con DiscordConnection

	botUser  *core2.BotUser
	outgoing *core2.OutgoingResponses
	core     *core2.Service
	guilds   GuildManager
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "discord"
}

// SetCore sets the core2.Service
func (d *Driver) SetCore(svc *core2.Service) {
	d.core = svc
}

// New starts a new core2.ProtocolDriver compatible Driver
func New(con DiscordConnection) (*Driver, error) {
	if con.IsConnected() {
		return nil, ErrConnectionAlreadyConnected
	}

	d := &Driver{Con: con}

	d.initCache()

	return d, nil
}

func (d *Driver) initCache() {
	log.Println("Driver: initCache()")
	d.guilds = make(map[string]*Guild, 2)
}

func (d *Driver) getOrMakeGuild(guildID string) *Guild {
	guild, ok := d.guilds[guildID]
	if !ok {
		guild = &Guild{guildID: guildID}
		guild.init()

		d.guilds[guildID] = guild
	}

	return guild
}

func (d *Driver) registerHandlers() {
	log.Println("Driver: registerHandlers()")

	d.Con.AddHandler(d.handleReadyStateEvent)
	d.Con.AddHandler(d.handleGuildCreateEvent)
	d.Con.AddHandler(d.handleGuildEmojisUpdate)
	// Channel Events
	d.Con.AddHandler(d.handleChannelCreateEvent)
	d.Con.AddHandler(d.handleChannelUpdateEvent)
	d.Con.AddHandler(d.handleChannelDeleteEvent)
	// Thread (temporary channels in a trench-coat) Events
	d.Con.AddHandler(d.handleThreadCreateEvent)
	d.Con.AddHandler(d.handleThreadUpdateEvent)
	d.Con.AddHandler(d.handleThreadDeleteEvent)
	// Guild Member Events
	d.Con.AddHandler(d.handleGuildMemberAddEvent)
	d.Con.AddHandler(d.handleGuildMemberUpdateEvent)
	d.Con.AddHandler(d.handleGuildMemberRemoveEvent)
	// Message Events
	d.Con.AddHandler(d.handleMessageCreateEvent)
	d.Con.AddHandler(d.handleMessageUpdateEvent)
	d.Con.AddHandler(d.handleMessageDeleteEvent)
	// Reaction Events
	d.Con.AddHandler(d.handleReactionAdd)
	d.Con.AddHandler(d.handleReactionRemove)
}

func (d *Driver) handleReadyStateEvent(s *discordgo.Session, ev *discordgo.Ready) {
	d.SetBotUser(&core2.BotUser{
		ID:      ev.User.ID,
		Name:    strings.ToLower(ev.User.Username),
		IsBot:   true,
		IsAdmin: false,
	})
}

func (d *Driver) handleGuildCreateEvent(s *discordgo.Session, ev *discordgo.GuildCreate) {
	logEvent("guild_create_ev", time.Now(), ev)

	if ev.Name != "" && ev.Guild != nil {
		guild := d.getOrMakeGuild(ev.Guild.ID)

		for _, chn := range ev.Channels {
			guild.channelCache[chn.ID] = chn
		}

		for _, thread := range ev.Threads {
			guild.channelCache[thread.ID] = thread
		}

		for _, member := range ev.Members {
			guild.memberCache[member.User.ID] = member
		}

		for _, emoji := range ev.Emojis {
			name := strings.ToLower(emoji.Name)

			guild.emojiCache[name] = emoji
		}
	}
}
