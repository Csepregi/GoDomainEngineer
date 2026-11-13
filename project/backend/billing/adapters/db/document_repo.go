package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"
	"eats/backend/billing/adapters/db/dbmodels"
	"eats/backend/common"
	"eats/backend/billing/domain"
)

type DocumentRecord struct {
	UUID domain.DocumentUUID
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) CreateDocument(
	ctx context.Context,
	series domain.DocumentSeries,
	createFunc func(documentNumber domain.DocumentNumber) (DocumentRecord, error),
) (domain.DocumentUUID, error) {
	//return domain.DocumentUUID{}, fmt.Errorf("TODO: implement")
	var document DocumentRecord
	err := common.UpdateInReadCommittedTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)
		lastNumber, err := queries.NextDocumentNumber(ctx, series.String())
		if err != nil {
			return err
		}
		docNumber, err := domain.NewDocumentNumber(series, int(lastNumber))
		if err != nil {
			return err
		}
		document, err = createFunc(docNumber)
		if err != nil {
			return err
		}
		err = queries.SaveDocument(ctx, dbmodels.SaveDocumentParams{
			DocumentUuid: document.UUID,
			DocumentNumber: docNumber.String(),
			SeriesPrefix: series.String(),
		})
		if err != nil {
			return fmt.Errorf("failed to save document: %w", err)
		}
		return nil
	})
	if err != nil {
		return domain.DocumentUUID{}, err
	}
	return document.UUID, err
	}
