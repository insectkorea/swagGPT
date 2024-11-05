package matcher

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMatchHandlerToRoute(t *testing.T) {
	testCases := []struct {
		name             string
		handlerSignature string
		routes           []string
		expectedRoute    string
		expectedError    error
	}{
		{
			name:             "MatchingRoute",
			handlerSignature: "RBFBundleHandler_ListByOrg",
			routes: []string{
				"GET /api/v1/organizations/:organization_id/bundles",
				"GET /api/v1/organizations/:organization_id/invitations",
			},
			expectedRoute: "GET /api/v1/organizations/:organization_id/bundles",
			expectedError: nil,
		},
		{
			name:             "NonMatchingRoute",
			handlerSignature: "Random_DoSomething",
			routes: []string{
				"GET /api/v1/organizations/:organization_id/bundles",
				"GET /api/v1/organizations/:organization_id/invitations",
			},
			expectedRoute: "",
			expectedError: &RouteNotFoundError{HandlerSignature: "Random_DoSomething"},
		},
		{
			name:             "MultipleMatchingRoutes",
			handlerSignature: "OrgHandler_List",
			routes: []string{
				"GET /api/v1/organizations/:organization_id/bundles",
				"GET /api/v1/organizations/:organization_id/invitations",
			},
			expectedRoute: "GET /api/v1/organizations/:organization_id/bundles",
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			route, err := MatchHandlerToRoute(tc.handlerSignature, tc.routes)

			if !reflect.DeepEqual(route, tc.expectedRoute) {
				t.Errorf("Unexpected route. Got %+v, expected %+v", route, tc.expectedRoute)
			}

			if (err == nil) != (tc.expectedError == nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Unexpected error. Got %v, expected %v", err, tc.expectedError)
			}
		})
	}
}

type RouteNotFoundError struct {
	HandlerSignature string
}

func (e *RouteNotFoundError) Error() string {
	return fmt.Sprintf("no route found for handler: %s", e.HandlerSignature)
}
