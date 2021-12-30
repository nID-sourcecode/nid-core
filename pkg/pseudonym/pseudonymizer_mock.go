package pseudonym

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockPseudonymizer mocks the Pseudonymizer interface
type MockPseudonymizer struct {
	mock.Mock
}

// GetPseudonym implements Pseudonymizer interface
func (m *MockPseudonymizer) GetPseudonym(ctx context.Context, myPseudo, targetNamespace string) (string, error) {
	args := m.Called(ctx, myPseudo, targetNamespace)

	return args.String(0), args.Error(1)
}

// GeneratePseudonym implements Pseudonymizer interface
func (m *MockPseudonymizer) GeneratePseudonym(ctx context.Context, amount uint32) ([]string, error) {
	args := m.Called(ctx, amount)

	return args.Get(0).([]string), args.Error(1)
}
