package models

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// DefaultModelPrimary is a default representation of the model
func (m *AudienceDB) DefaultModelPrimary(namespace string) Audience {
	return Audience{
		ID:        uuid.FromStringOrNil("e9f31397-dfbe-4701-8b11-bd7bc79f1aa8"),
		Audience:  fmt.Sprintf("http://databron.%s/gql", namespace),
		Namespace: namespace,
	}
}

// DefaultModelSecondary is a default representation of the model
func (m *AudienceDB) DefaultModelSecondary(namespace string) Audience {
	return Audience{
		ID:        uuid.FromStringOrNil("482923d1-8022-489e-984a-b52425ba4d51"),
		Audience:  fmt.Sprintf("http://other-databron.%s/gql", namespace),
		Namespace: namespace,
	}
}

// DefaultModelInformationService is a default representation of the information service audience
func (m *AudienceDB) DefaultModelInformationService(namespace string) Audience {
	return Audience{
		ID:        uuid.FromStringOrNil("ae734e1c-08f0-40f0-ba27-9e15e32cded5"),
		Audience:  fmt.Sprintf("http://information.%s", namespace),
		Namespace: namespace,
	}
}

// GetAudienceByURI retrieves an audience given a audience URI
func (m *AudienceDB) GetAudienceByURI(audienceURI string) (*Audience, error) {
	var audience Audience
	err := m.Db.Find(&audience, "audience = ?", audienceURI).Error
	if err != nil {
		return nil, err
	}

	return &audience, nil
}
