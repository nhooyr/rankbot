package main

import (
	"bufio"
	"flag"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nhooyr/color/log"
)

var (
	email     = flag.String("email", "", "account email")
	pass      = flag.String("pass", "", "account password")
	guild     = flag.String("guild", "", "guild (server) to join")
	channel   = flag.String("chan", "", "channel to join")
	message   = flag.String("msg", "_", "message to be sent")
	interval  = flag.Int64("int", 60, "interval between messages in seconds")
	delete    = flag.Bool("del", false, "delete every message as soon as it's been sent")
	idiomFile = flag.String("idiom", "", "file containing a set of messages")
)

func main() {
	flag.Parse()
	if *email == "" || *pass == "" {
		log.Fatal("please provide an email and password")
	}
	if *idiomFile != "" && *message != "_" {
		log.Fatal("provide either -msg or -idiom")
	}

	var idiom []string = []string{*message}
	if *idiomFile != "" {
		f, err := os.Open(*idiomFile)
		if err != nil {
			log.Fatal(err)
		}
		idiom = readIdiom(f)
	}

	s, err := discordgo.New(*email, *pass)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("logged in")

	g := findGuild(s)
	if g == nil {
		log.Fatal("could not find guild")
	}
	id := findChannel(s, g)
	if id == "" {
		log.Fatal("could not find channel")
	}

	rand.Seed(time.Now().Unix())
	for t := time.Tick(time.Duration(*interval) * time.Second); ; <-t {
		m := idiom[rand.Intn(len(idiom))]

		msg, err := s.ChannelMessageSend(id, m)
		if err != nil {
			log.Print(err)
		}
		log.Print("sent message")

		if *delete {
			err = s.ChannelMessageDelete(id, msg.ID)
			if err != nil {
				log.Print(err)
			}
			log.Print("deleted message")
		}
	}
}

func readIdiom(r io.Reader) (idiom []string) {
	ln := bufio.NewScanner(r)
	for ln.Scan() {
		idiom = append(idiom, ln.Text())
	}
	return
}

func findGuild(s *discordgo.Session) *discordgo.UserGuild {
	gs, err := s.UserGuilds(0, "", "")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("got guilds")
	for _, g := range gs {
		if g.Name == *guild {
			log.Print("found guild")
			return g
		}
	}
	return nil
}

func findChannel(s *discordgo.Session, g *discordgo.UserGuild) string {
	chs, err := s.GuildChannels(g.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("got channels")
	for _, ch := range chs {
		if ch.Name == *channel {
			log.Print("found channel")
			return ch.ID
		}
	}
	return ""
}
