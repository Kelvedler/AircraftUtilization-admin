package db

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kelvedler/AircraftUtilization-admin/pkg/setting"
)

func NewConnectionPool(
	ctx context.Context,
	mainLogger *slog.Logger,
) *pgxpool.Pool {
	dbpool, err := pgxpool.New(ctx, setting.Setting.DatabaseUrl)
	if err != nil {
		mainLogger.Error(err.Error())
		os.Exit(1)
	}
	return dbpool
}

type (
	BatchOperation func(batch *pgx.Batch)
	BatchRead      func(results pgx.BatchResults) error
	BatchSet       func() (opeartion BatchOperation, read BatchRead)
)

func PerformBatch(ctx context.Context, dbpool *pgxpool.Pool, batchSets []BatchSet) (errs []error) {
	batch := pgx.Batch{}
	var batchReads []BatchRead
	for _, item := range batchSets {
		operation, read := item()
		operation(&batch)
		batchReads = append(batchReads, read)
	}
	results := dbpool.SendBatch(ctx, &batch)
	for _, read := range batchReads {
		errs = append(errs, read(results))
	}
	results.Close()
	return errs
}
