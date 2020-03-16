package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo"
)

var messageLogs []string
var bots [2]*discordgo.Session

func runBot(token string, shardID int) (discord *discordgo.Session) {
	discord, err := discordgo.New("Bot " + token)
	checkErr("make bot", err)

	discord.ShardCount = len(bots)
	discord.ShardID = shardID
	discord.State.PresenceAdd("343634182615990273", &discordgo.Presence{
		Nick: "test",
	})

	discord.AddHandler(messageCreate)

	err = discord.Open()
	checkErr("opening", err)

	return
}

func main() {
	bytes, err := ioutil.ReadFile("token")
	checkErr("load token file", err)

	for i := 0; i < len(bots); i++ {
		bots[i] = runBot(string(bytes), i)
	}

	e := echo.New()
	e.GET("/send", func(c echo.Context) error {
		return c.File("./send.html")
	})
	e.POST("/send", func(c echo.Context) error {
		for i := 0; i < 5; i++ {
			bots[0].ChannelMessageSend("343634182615990273", c.FormValue("text"))
			bots[1].ChannelMessageSend("688747716796481597", c.FormValue("text"))
		}
		return c.HTML(http.StatusOK, "<script>window.history.back()</script>")
	})
	e.GET("/log", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<textarea>"+strings.Join(messageLogs, "\n")+"</textarea>")
	})

	go func() {
		e.Logger.Fatal(e.Start(":1232"))
	}()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	for _, bot := range bots {
		bot.Close()
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "ping" {
		for i := 1; i <= 9999999999; i++ {
		}
		s.ChannelMessageSend(m.ChannelID, "ping! "+s.HeartbeatLatency().String())
	}

	if m.Content == "pong" {
		for i := 0; i <= 5; i++ {
			s.ChannelMessageSend(m.ChannelID, "pong!")
		}
	}

	messageLogs = append(messageLogs, m.Content)
}

func checkErr(text string, err error) {
	if err != nil {
		log.Fatalln("error occured in", text, err)
	}
}
