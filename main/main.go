package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/may215/anaconda"

	"github.com/dghubble/oauth1"
)

func main() {
	//anaconda.SetConsumerKey("wWaLpH4CrSP3lOOT4vzfZfY9L")
	//anaconda.SetConsumerSecret("wOrL9NUxxepxN34ZZuvLOwhvwkJvpHJrPoeBP9muHzsTfvM9Jp")
	//var api = anaconda.NewTwitterApi("3944636567-TzEHqosZ9y7wzSBVb6fd8Lu1GxPJMO7gpwECblj", "zKJABK8gm2jtyYh6rq1XAA2lNbyohis49roCyzRyUOgQh")

	config := oauth1.NewConfig("wWaLpH4CrSP3lOOT4vzfZfY9L", "wOrL9NUxxepxN34ZZuvLOwhvwkJvpHJrPoeBP9muHzsTfvM9Jp")
	token := oauth1.NewToken("3944636567-TzEHqosZ9y7wzSBVb6fd8Lu1GxPJMO7gpwECblj", "zKJABK8gm2jtyYh6rq1XAA2lNbyohis49roCyzRyUOgQh")
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	var api = anaconda.NewTwitterApi("3944636567-TzEHqosZ9y7wzSBVb6fd8Lu1GxPJMO7gpwECblj", "zKJABK8gm2jtyYh6rq1XAA2lNbyohis49roCyzRyUOgQh", httpClient)
	// Convenience Demux demultiplexed stream messages

	demux := anaconda.NewSwitchDemux()

	demux.Tweet = func(tweet *anaconda.Tweet) {
		fmt.Println("tweet.Text")
		fmt.Println(tweet.Text)
	}
	demux.DM = func(dm *anaconda.DirectMessage) {
		fmt.Println("dm.SenderId")
		fmt.Println(dm.SenderId)
	}
	demux.Event = func(event *anaconda.Event) {
		fmt.Println("event")
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := anaconda.StreamFilterParams{
		Track:         []string{"BDS"},
		StallWarnings: anaconda.Bool(true),
	}
	stream, err := api.Streams.Filter(&filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
