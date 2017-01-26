package main

import (
	"context"
	"log"

	"github.com/kr/pretty"

	"googlemaps.github.io/maps"
)

type RouteFetcher struct {
	Api    string
	client *maps.Client
}

func (r *RouteFetcher) init() {
	client, err := maps.NewClient(maps.WithAPIKey(r.Api))
	if err != nil {
		log.Fatal("Can't get a maps client")
	}
	r.client = client

}

func (r *RouteFetcher) GetRouteToWork(postalCode string) []string {
	startStation := r.getStartStation(postalCode)
	if startStation == "" {
		return []string{}
	}

	return getRouteToWork(startStation)
}

func getRouteToWork(startStaton string) []string {
	return []string{}
}

func (r *RouteFetcher) getStartStation(postalCode string) string {
	directionRequest := &maps.DirectionsRequest{
		Origin:      postalCode,
		Destination: "1 Toronto Street, Toronto, Canada",
		Mode:        maps.TravelModeTransit,
	}

	resp, _, err := r.client.Directions(context.Background(), directionRequest)
	if err != nil {
		return ""
	}
	pretty.Println(resp)

	return "Dufferin"
}
