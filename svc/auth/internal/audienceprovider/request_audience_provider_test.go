package audienceprovider

import (
	"context"
	"testing"

	"github.com/nID-sourcecode/nid-core/svc/auth/contract"
	"github.com/nID-sourcecode/nid-core/svc/auth/models"
	"github.com/stretchr/testify/assert"
)

func TestRequestAudienceProvider_GetAudience(t *testing.T) {
	tests := []struct {
		Scenario string
		Request  *models.TokenClientFlowRequest
		Scopes   []*models.Scope
		Expected []string
		Error    error
	}{
		{
			Scenario: "Audience provided in request body",
			Request:  &models.TokenClientFlowRequest{Audience: "test"},
			Scopes:   []*models.Scope{},
			Expected: []string{"test"},
			Error:    nil,
		},
		{
			Scenario: "Audience not provided in request body",
			Request:  &models.TokenClientFlowRequest{},
			Scopes:   []*models.Scope{},
			Expected: nil,
			Error:    contract.ErrInvalidArguments,
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			provider := &RequestAudienceProvider{}
			result, err := provider.GetAudience(context.Background(), test.Request, test.Scopes)

			assert.Equal(t, test.Expected, result)
			assert.ErrorIs(t, err, test.Error)
		})
	}
}
