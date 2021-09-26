package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"strings"
	"statuschecking/statuses"
)
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "bot token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("New Login as StatusChecker!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == s.State.User.ID {
		return
	}
	
	if strings.HasPrefix(msg.Content, "!status") {
		args := msg.Content[8:]
		res, err := http.Get(args)
		if err != nil || res.StatusCode >= 400 {
			s.ChannelMessageSend(msg.ChannelID, "Request Failed.")
			return
		}
		content := fmt.Sprintf("%s : %s -> %s : %d", "Finished Checking", args, "Current Status", res.StatusCode)	
		s.ChannelMessageSend(msg.ChannelID, content)
	}

	if strings.HasPrefix(msg.Content, "!check") {
		args := msg.Content[7:]
		dictionary := statuses.Dictionary{
			"100" : "Client is sending the body.",
			"101" : "Switching protocols.",
			"102" : "Currently [rocessing.",
			"103" : "Currently requesting to get back the post.",
			"200" : "Status OK.",
			"201" : "Created a new resource.",
			"202" : "Accepted but 99% reloaded.",
			"203" : "Not fully loaded.",
			"204" : "No data to give.",
			"205" : "Need to reset the content.",
			"206" : "Only got some of the resources due to header problems.",
			"207" : "Mutli-status",
			"208" : "Dodging multiple requests.",
			"226" : "Im used.",
			"300" : "Maximum 5 links to move.",
			"301" : "Website is moved to another website permanetly.",
			"302" : "The website is temporarily moved to another link.",
			"303" : "The requested page can be found in another website.",
			"304" : "After the last request, the page isn't modified.",
			"305" : "This happens when you requested without a proxy.",
			"306" : "Unused code.",
			"307" : "Don't change the http method.",
			"308" : "Same as status 301.",
			"400" : "Bad Request.",
			"401" : "Missing Permissions/Unauthorized.",
			"402" : "Payment required.",
			"403" : "Forbidden.",
			"404" : "Not found.",
			"405" : "Method is not allowed.",
			"406" : "Request not acceptable.",
			"407" : "Proxy authentication is required.",
			"408" : "Request timeout.",
			"409" : "Request conflict.",
			"410" : "The page is not available now.",
			"411" : "Length is required.",
			"412" : "Precondition failed.",
			"413" : "Request entity is too large.",
			"414" : "The request url is too long.",
			"415" : "Unsupported media type.",
			"416" : "Requested range is not satisfiable.",
			"417" : "Expectation failed.",
			"418" : "I'm a teapot.",
			"422" : "Requested sended without errors, but has grammar error.",
			"429" : "Too many requests.",
			"500" : "Internal server error.",
			"501" : "The server could not recgonize the request method, or doesn't have the ability to run the method.",
			"502" : "Bad Gateway.",
			"504" : "Gateway timeout.",
			"505" : "Http version not supported.",
			"511" : "Network authentication required.",
		}
		description, err := dictionary.Search(args)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, "The status doesn't exist.")
		}
		s.ChannelMessageSend(msg.ChannelID, description)
	}
}
