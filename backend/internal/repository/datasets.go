package repository

import (
	"context"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/domain"
)

type DatasetsRepo struct {
}

func NewDatasetsRepo() Datasets {
	return &DatasetsRepo{}
}

func (r *DatasetsRepo) GetDatasetsByClientId(ctx context.Context, tx Transaction, id int64) ([]*domain.Dataset, error) {
	// TODO
	return []*domain.Dataset{
		{
			Id:       1,
			ClientId: id,
		},
	}, nil
}

func (r *DatasetsRepo) DeleteByClientId(ctx context.Context, tx Transaction, id int64) error {
	// TODO
	return nil
}
