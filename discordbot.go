package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"regexp"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

var headCommandRegex = regexp.MustCompile(` (.+)`)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "$") {
		command := strings.Replace(m.Content, "$", "", 1)

		embed := generateEmbed(m, executeCommand(command))
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	}

}
func generateEmbed(m *discordgo.MessageCreate, output string) *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Color:       0x00ff00, // Green
		Description: "```cmd\n" + output + "\n```",
	}
	return embed
}
func executeCommand(command string) string {
	//trimming spaces
	/*
		command = strings.TrimSpace(command)
		temp := headCommandRegex.Split(command, -1)
	*/
	args := strings.Split(command, " ")
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	stdoutStderr, err := cmd.CombinedOutput()
	output := string(stdoutStderr)
	if err != nil {

		output = output + "\nError:" + err.Error()
	}
	return output
}
