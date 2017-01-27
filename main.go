package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var slackChannel = "G3W5N8UQG"

var UserStore = make(map[string][]string)

func main() {
	botToken := os.Getenv("BOT_TOKEN")

	slackBot := Bot{Token: botToken, UserStore: UserStore}
	slackBot.init()
	slackBot.write("Bot Online!")

	tweetChannel := make(chan twitter.Tweet)

	go startTwitterChecker(tweetChannel)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Starting Stream...")
	startStream(tweetChannel, &slackBot, ch)
}

func startStream(tweetChannel chan twitter.Tweet, slackBot *Bot, ch chan os.Signal) {
	messageCount := 0
	for {
		select {
		case tweet := <-tweetChannel:
			{
				fmt.Println("RECEIVED NEW TWEET")
				messageCount++
				tweetTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 2006", tweet.CreatedAt)
				slackMessage := fmt.Sprintf(":ttc: %s AT TIME %s", tweet.Text, tweetTime.Local().String())
				fmt.Println(slackMessage)
				slackBot.write(slackMessage)

				station, delayed, cleared := findAffectedStation(tweet.Text)
				fmt.Println("TWEET IS ", station, delayed, cleared)

				if delayed {
					DelayedStations[station] = true
				}
				if cleared {
					delete(DelayedStations, station)
				}

				for user, items := range UserStore {
					for _, item := range items {
						if item == station {
							slackBot.write(fmt.Sprintf("%s will be delayed\n", user))
						}
					}
				}
			}
		case <-ch:
			fmt.Println("Stopping Bot...")
			return
		}
	}
}
func startTwitterChecker(tweetChannel chan twitter.Tweet) {
	apiKey := os.Getenv("CONSUMER_API_KEY")
	apiSecret := os.Getenv("CONSUMER_API_SECRET")
	token := os.Getenv("TOKEN")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	config := oauth1.NewConfig(apiKey, apiSecret)
	twitterToken := oauth1.NewToken(token, tokenSecret)

	httpClient := config.Client(oauth1.NoContext, twitterToken)
	client := twitter.NewClient(httpClient)

	timerChannel := time.Tick(1 * time.Second)
	var sinceId int64
	for now := range timerChannel {
		fmt.Println(now)
		timelineParams := &twitter.UserTimelineParams{
			ScreenName: os.Getenv("SCREEN_NAME"),
			Count:      5,
		}
		if sinceId != 0 {
			timelineParams.SinceID = sinceId
		}

		tweets, _, err := client.Timelines.UserTimeline(timelineParams)
		if err != nil {
			fmt.Println(err.Error)
		}

		for i := len(tweets) - 1; i >= 0; i-- {
			tweet := tweets[i]
			sinceId = tweet.ID

			tweetChannel <- tweet
		}
	}

}

func findAffectedStation(text string) (string, bool, bool) {
	if !isSubway(text) {
		fmt.Println(text, "NOT A SUBWAY")
		return "", false, false
	}

	station := getStation(text)
	if station == "" {
		fmt.Println("NOT A STATION")
		return "", false, false
	}
	fmt.Println(station)

	fmt.Println("IS SUBWAY ", text)
	delay := isDelay(text)
	clear := isClear(text)

	fmt.Println("IS DELAY ", delay)
	fmt.Println("IS CLEAR", clear)

	return station, delay, clear
}
