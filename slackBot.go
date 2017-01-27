package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type Bot struct {
	Token        string
	UserStore    map[string][]string
	rtm          *slack.RTM
	routeFetcher *RouteFetcher
}

func (b *Bot) init() {
	api := slack.New(b.Token)
	b.rtm = api.NewRTM()
	// b.rtm.SetDebug(true)
	go b.rtm.ManageConnection()
	go b.handleIncomingEvents()

	b.routeFetcher = &RouteFetcher{
		Api: os.Getenv("MAPS_API"),
	}
	b.routeFetcher.init()
}

func (b *Bot) write(message string) {
	b.rtm.SendMessage(b.rtm.NewOutgoingMessage(message, slackChannel))
}

func (b *Bot) handleIncomingEvents() {
	for msg := range b.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			{
				// if ev.Channel == slackChannel {
				fmt.Printf("Message: %v\n", ev)
				user, postalCode := getPostalCode(ev.Msg.Text)
				if postalCode != "" {
					routeToWork := b.routeFetcher.GetRouteToWork(postalCode)
					fmt.Println("USER ", user, "ROUTE : ", routeToWork)

					var message string
					if len(routeToWork) > 0 {
						message = fmt.Sprintf("Hi %s, your subway stops to work are %s", user, strings.Join(routeToWork, ","))
						b.UserStore[user] = routeToWork
					} else {
						message = fmt.Sprintf("Hi %s. It doesn't look like you take the subway to work")
						delete(b.UserStore, user)
					}

					b.rtm.SendMessage(b.rtm.NewOutgoingMessage(message, slackChannel))
				}
			}
		default:
		}
	}
}

func getPostalCode(message string) (name string, postalCode string) {
	postalMatcher := regexp.MustCompile(`^(\w+) register ([a-zA-Z][0-9][a-zA-Z][0-9][a-zA-Z][0-9])`)

	matches := postalMatcher.FindStringSubmatch(message)
	if len(matches) != 3 {
		return "", ""
	}

	return matches[1], matches[2]
}
