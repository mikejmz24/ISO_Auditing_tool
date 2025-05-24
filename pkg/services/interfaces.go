package services

import (
	"ISO_Auditing_Tool/pkg/types"
	"context"
)

type DraftServiceInterface interface {
	Create(ctx context.Context, draft types.Draft) (types.Draft, error)
	GetByID(ctx context.Context, draft types.Draft) (types.Draft, error)
	Update(ctx context.Context, draft types.Draft) (types.Draft, error)
	GetAll(ctx context.Context) ([]types.Draft, error)
}
