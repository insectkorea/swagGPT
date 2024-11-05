package matcher

import (
	"fmt"
	"sort"
	"strings"

	"github.com/insectkorea/swagGPT/internal/model"
)

type matchedRoute struct {
	Route       string
	Method      string
	Substrings  []string
	TotalLength int
}

func MatchHandlerToRoute(handlerSignature string, routes []model.Route) (string, error) {
	// Find the matching routes based on the number of substrings and total length of matched substrings
	var matchingRoutes []matchedRoute
	for _, route := range routes {
		substrings, totalLength := findMatchedSubstrings(handlerSignature, route)
		if len(substrings) > 0 {
			matchingRoutes = append(matchingRoutes, matchedRoute{
				Route:       route.Path,
				Method:      strings.ToLower(route.Method),
				Substrings:  substrings,
				TotalLength: totalLength,
			})
		}
	}

	// Prioritize the matching routes based on the number of substrings and total length of matched substrings
	sortMatchingRoutes(matchingRoutes)
	routeString := ""
	for i, route := range matchingRoutes {
		if i > 3 {
			// Only consider the top 3 routes
			break
		}
		routeString += fmt.Sprintf("%s [%s]", route.Route, route.Method)
		if i < len(matchingRoutes)-1 {
			routeString += ", "
		}
	}

	return routeString, nil
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
func findMatchedSubstrings(handlerSignature string, route model.Route) ([]string, int) {
	var substrings []string
	totalLength := 0
	// preprocess the handlerSignature and route to lowercase and trim spaces
	handlerSignature = strings.TrimSpace(handlerSignature)
	handlerSignature = strings.ToLower(handlerSignature)
	path := strings.TrimSpace(route.Path)
	path = strings.ToLower(path)

	for i := 0; i < len(path); i++ {
		for j := i + 1; j <= len(path); j++ {
			substring := path[i:j]
			if len(substring) > 3 && strings.Contains(handlerSignature, substring) {
				substrings = append(substrings, substring)
				totalLength += len(substring)
			}
		}
	}

	return substrings, totalLength
}
