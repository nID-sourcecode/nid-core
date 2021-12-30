package headers

import (
	"context"
)

// GetJWTToken retrieve jwt token from context
func (m GRPCMetadataHelper) GetJWTToken(ctx context.Context) (string, error) {
	return m.GetValFromCtx(ctx, "authorization")
}
