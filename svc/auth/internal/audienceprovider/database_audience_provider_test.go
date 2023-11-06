package audienceprovider

import (
	"context"
	"testing"

	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/internal/config"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseAudienceProvider_GetAudience(t *testing.T) {
	tests := []struct {
		Scenario               string
		Request                *models.TokenClientFlowRequest
		Scopes                 []*models.Scope
		AllowMultipleAudiences bool
		Expected               []string
		Error                  error
	}{
		{
			Scenario:               "AllowSingleAudience, no scopes provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{},
			AllowMultipleAudiences: true,
			Expected:               []string{},
			Error:                  nil,
		},
		{
			Scenario:               "AllowMultipleAudiences, no scopes provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{},
			AllowMultipleAudiences: false,
			Expected:               []string{},
			Error:                  nil,
		},

		// AllowMultipleAudiences is false.
		{
			Scenario:               "AllowSingleAudience, scope with one audience is provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}}}},
			AllowMultipleAudiences: false,
			Expected:               []string{"test"},
			Error:                  nil,
		},
		{
			Scenario:               "AllowSingleAudience, scopes with multiple audiences are provided, no audience provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}, {Audience: "test2"}}}},
			AllowMultipleAudiences: false,
			Expected:               nil,
			Error:                  contract.ErrInvalidArguments,
		},
		{
			Scenario:               "AllowSingleAudience, scopes with multiple audiences are provided, audience provided",
			Request:                &models.TokenClientFlowRequest{Audience: "test2"},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}, {Audience: "test2"}}}},
			AllowMultipleAudiences: false,
			Expected:               []string{"test2"},
			Error:                  nil,
		},
		// AllowMultipleAudiences is true.
		{
			Scenario:               "AllowMultipleAudiences, scope with one audience is provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}}}},
			AllowMultipleAudiences: true,
			Expected:               []string{"test"},
			Error:                  nil,
		},
		{
			Scenario:               "AllowMultipleAudiences, scopes with multiple audiences are provided, no audience provided",
			Request:                &models.TokenClientFlowRequest{},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}, {Audience: "test2"}}}},
			AllowMultipleAudiences: true,
			Expected:               []string{"test", "test2"},
			Error:                  nil,
		},
		{
			Scenario:               "AllowMultipleAudiences, scopes with multiple audiences are provided, audience provided",
			Request:                &models.TokenClientFlowRequest{Audience: "test2"},
			Scopes:                 []*models.Scope{{Audiences: []*models.Audience{{Audience: "test"}, {Audience: "test2"}}}},
			AllowMultipleAudiences: true,
			Expected:               []string{"test2"},
			Error:                  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			provider := &DatabaseAudienceProvider{
				conf: &config.AuthConfig{
					AllowMultipleAudiences: test.AllowMultipleAudiences,
				},
			}
			result, err := provider.GetAudience(context.Background(), test.Request, test.Scopes)

			assert.Equal(t, test.Expected, result)
			assert.ErrorIs(t, err, test.Error)
		})
	}
}
