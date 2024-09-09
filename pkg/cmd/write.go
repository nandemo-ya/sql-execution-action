package cmd

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/nandemo-ya/sql-execution-action/pkg/cli"
	"github.com/nandemo-ya/sql-execution-action/pkg/config"
)

var WriteCmd = &cobra.Command{
	Use:   "write",
	Short: "Send update type queries to the database",
	RunE:  cli.WithContext(runUpdate),
}

func runUpdate(ctx context.Context, cmd *cobra.Command) error {
	args, err := buildArgs(cmd)
	if err != nil {
		return err
	}

	var writerConfig config.Writer
	if err := yaml.Unmarshal(args.sqlPayload, &writerConfig); err != nil {
		slog.Error("failed unmarshaling sql data", "err", err)
		return err
	}

	db, err := sql.Open(args.engine, args.datasource)
	if err != nil {
		slog.Error("failed opening database", "err", err)
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		slog.Error("failed beginning transaction", "err", err)
		return err
	}

	for _, query := range writerConfig.Queries {
		statement, err := db.Prepare(query.SQL)
		if err != nil {
			slog.Error("failed preparing statement", "err", err)
			tx.Rollback()
			return err
		}
		if _, err := statement.Exec(query.Params...); err != nil {
			slog.Error("failed executing statement", "err", err)
			tx.Rollback()
			return err
		}
		slog.Info("successfully executed statement", slog.String("sql", query.SQL), slog.Any("params", query.Params))
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed committing transaction", "err", err)
		return err
	}
	slog.Info("successfully committed transaction")

	return nil
}
