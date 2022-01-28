package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

func initBot(config BotConfig) {

	botToken := config.token
	if botToken == "" {
		log.Fatalln("no bot token specified! specify one in config.toml")
	}

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("error authenticating to discord, ", err)
	}

	// code stolen
	// https://github.com/bwmarrin/discordgo/blob/master/examples/pingpong/main.go

	// Register the messageCreate func as a callback for MessageCreate events.
	discord.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("pit viper is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate, config BotConfig) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "$ping" {
		s.ChannelMessageSend(m.ChannelID, "ay yo fuck off already damn")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "$pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	if m.Content == "$help" {
		s.ChannelMessageSend(m.ChannelID, "nah if you don't know how to use me then you don't need to use me")
	}

	if m.Content == "$genshin" {
		s.ChannelMessageSend(m.ChannelID, "begone, weeaboos! get in the pit where you belong")
		for _, userid := range config.genshinUserIDs {
			s.GuildMemberMove(config.guildID, userid, &config.pitVCID)
		}
	}

	if m.Content == "$ungenshin" {
		s.ChannelMessageSend(m.ChannelID, "unimplemented :( ")
	}
}

type BotConfig struct {
	token          string
	genshinUserIDs []string
	guildID        string
	pitVCID        string
}

func getBotConfig() BotConfig {

	var config BotConfig
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal(err)
	}

	return config
}

func main() {
	log.Println("starting up pit viper...")

	log.Println("loading config...")
	config := getBotConfig()
	log.Println()
	log.Println("bot init...")
	initBot(config)
}
