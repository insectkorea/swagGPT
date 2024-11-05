package matcher

import (
	"reflect"
	"testing"

	"github.com/insectkorea/swagGPT/internal/model"
)

func TestMatchHandlerToRoute(t *testing.T) {
	testCases := []struct {
		name             string
		handlerSignature string
		routes           []model.Route
		expectedRoute    string
		expectedError    error
	}{
		{
			name:             "MatchingRoute",
			handlerSignature: "RBFBundleHandler_ListByOrg",
			routes: []model.Route{
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/bundles", Pattern: "/api/v1/organizations/:organization_id/bundles"},
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/invitations", Pattern: "/api/v1/organizations/:organization_id/invitations"},
			},
			expectedRoute: "/api/v1/organizations/:organization_id/bundles [get]",
			expectedError: nil,
		},
		{
			name:             "NonMatchingRoute",
			handlerSignature: "Random_DoSomething",
			routes: []model.Route{
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/bundles", Pattern: "/api/v1/organizations/:organization_id/bundles"},
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/invitations", Pattern: "/api/v1/organizations/:organization_id/invitations"},
			},
			expectedRoute: "",
			expectedError: nil,
		},
		{
			name:             "MultipleMatchingRoutes",
			handlerSignature: "OrgHandler_List",
			routes: []model.Route{
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/bundles", Pattern: "/api/v1/organizations/:organization_id/bundles"},
				{Method: "GET", Path: "/api/v1/organizations/:organization_id/invitations", Pattern: "/api/v1/organizations/:organization_id/invitations"},
			},
			expectedRoute: "/api/v1/organizations/:organization_id/bundles [get]",
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
