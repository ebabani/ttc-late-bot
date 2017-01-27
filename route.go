package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	routeToWork := []string{}

	directionRequest := &maps.DirectionsRequest{
		Origin:      postalCode,
		Destination: "1 Toronto Street, Toronto, Canada",
		Mode:        maps.TravelModeTransit,
	}

	routes, _, err := r.client.Directions(context.Background(), directionRequest)
	if err != nil {
		return routeToWork
	}

	if len(routes) == 0 {
		fmt.Println("no routes")
		return routeToWork
	}
	legs := routes[0].Legs
	if len(legs) == 0 {
		fmt.Println("no legs")
		return routeToWork
	}
	steps := legs[0].Steps
	if len(steps) == 0 {
		fmt.Println("no steps")
		return routeToWork
	}

	for _, step := range steps {
		if step.TransitDetails == nil {
			fmt.Println("NO TRANSIT DETAILS")
			continue
		}
		if !strings.Contains(step.TransitDetails.Headsign, "Line ") {
			fmt.Println(step.TransitDetails.Headsign)
			continue
		}

		start := getStationFromStop(step.TransitDetails.DepartureStop)
		end := getStationFromStop(step.TransitDetails.ArrivalStop)

		routeToWork = append(routeToWork, getStations(start, end)...)
	}

	// pretty.Println(routes)

	return routeToWork
}

func getStationFromStop(stop maps.TransitStop) string {
	return strings.TrimSuffix(stop.Name, " Station")
}

func getStations(start, end string) []string {
	fmt.Println("SEARCHING FOR " + start + " " + end)
	stations := getStationsForLine(start, end, Line1)
	if len(stations) == 0 {
		return getStationsForLine(start, end, Line2)
	}
	return stations
}

func getStationsForLine(start, end string, line []string) []string {

	startIndex := findIndex(start, line)
	endIndex := findIndex(end, line)
	if startIndex != -1 && endIndex != -1 {
		if startIndex > endIndex {
			startIndex, endIndex = endIndex, startIndex
		}
		return line[startIndex : endIndex+1]
	}

	return []string{}
}

func findIndex(item string, list []string) int {
	for i, _ := range list {
		if list[i] == item {
			return i
		}
	}
	return -1
}
