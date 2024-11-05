package matcher

import (
	"fmt"
	"sort"
	"strings"
)

type Route struct {
	Pattern string
	Method  string
	Path    string
}

func MatchHandlerToRoute(handlerSignature string, routes []string) (string, error) {
	// Find the matching routes based on the number of substrings and total length of matched substrings
	var matchingRoutes []matchedRoute
	for _, route := range routes {
		substrings, totalLength := findMatchedSubstrings(handlerSignature, route)
		fmt.Println("Substrings matched:", substrings)
		if len(substrings) > 0 {
			matchingRoutes = append(matchingRoutes, matchedRoute{
				Route:       route,
				Substrings:  substrings,
				TotalLength: totalLength,
			})
		}
	}

	// Prioritize the matching routes based on the number of substrings and total length of matched substrings
	sortMatchingRoutes(matchingRoutes)

	if len(matchingRoutes) > 0 {
		return matchingRoutes[0].Route, nil
	}

	return "", fmt.Errorf("no route found for handler: %s", handlerSignature)
}

type matchedRoute struct {
	Route       string
	Substrings  []string
	TotalLength int
}

// sortMatchingRoutes sorts the matching routes based on the number of substrings and total length of matched substrings
func sortMatchingRoutes(routes []matchedRoute) {
	sort.Slice(routes, func(i, j int) bool {
		// Sort by the number of substrings in descending order
		if len(routes[i].Substrings) != len(routes[j].Substrings) {
			return len(routes[i].Substrings) > len(routes[j].Substrings)
		}
		// If the number of substrings is the same, sort by the total length of matched substrings in descending order
		return routes[i].TotalLength > routes[j].TotalLength
	})
}

// findMatchedSubstrings finds the substrings in the path that match the handler signature
func findMatchedSubstrings(handlerSignature, route string) ([]string, int) {
	var substrings []string
	totalLength := 0
	// preprocess the handlerSignature and route to lowercase and trim spaces
	handlerSignature = strings.TrimSpace(handlerSignature)
	handlerSignature = strings.ToLower(handlerSignature)
	route = strings.TrimSpace(route)
	route = strings.ToLower(route)

	for i := 0; i < len(route); i++ {
		for j := i + 1; j <= len(route); j++ {
			substring := route[i:j]
			if len(substring) > 3 && strings.Contains(handlerSignature, substring) {
				substrings = append(substrings, substring)
				totalLength += len(substring)
			}
		}
	}

	return substrings, totalLength
}
