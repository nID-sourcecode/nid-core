package models

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
)

// DefaultModel is a default representation of the model
func (m *ClientDB) DefaultModel() Client {
	s := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	id, err := uuid.FromString(s)
	if err != nil {
		log.WithError(err).Fatal("unable to create default id for seed client")
	}

	return Client{
		ID:       id,
		Color:    "green",
		Icon:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==",
		Logo:     "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z/C/HgAGgwJ/lK3Q6wAAAABJRU5ErkJggg==",
		Name:     "testClient",
		Metadata: postgres.Jsonb{RawMessage: []byte(`{"oin":"00000009000000999991"}`)},
	}
}

// GetClientByID retrieves a client given his ID
func (m *ClientDB) GetClientByID(clientID string) (*Client, error) {
	var client Client
	err := m.Db.Find(&client, "id = ?", clientID).Error
	if err != nil {
		return nil, err
	}

	return &client, nil
}
