package cmd

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/nandemo-ya/sql-execution-action/pkg/cli"
	"github.com/nandemo-ya/sql-execution-action/pkg/config"
)

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Send select type query to the database",
	RunE:  cli.WithContext(runSelect),
}

func runSelect(ctx context.Context, cmd *cobra.Command) error {
	args, err := buildArgs(cmd)
	if err != nil {
		return err
	}

	var readerConfig config.Reader
	if err := yaml.Unmarshal(args.sqlPayload, &readerConfig); err != nil {
		slog.Error("failed unmarshaling sql data", "err", err)
		return err
	}

	db, err := sql.Open(args.engine, args.datasource)
	if err != nil {
		slog.Error("failed opening database", "err", err)
		return err
	}
	defer db.Close()

	rows, err := db.Query(readerConfig.Query.SQL, readerConfig.Query.Params...)
	if err != nil {
		slog.Error("failed executing query", "err", err)
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		slog.Error("failed getting columns", "err", err)
		return err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			slog.Error("failed scanning row", "err", err)
			return err
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		results = append(results, rowMap)
	}
	data, err := json.Marshal(results)
	if err != nil {
		slog.Error("failed marshaling results", "err", err)
		return err
	}

	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile != "" {
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			slog.Error("failed opening output file", "err", err)
			return err
		}
		defer f.Close()
		if _, err = fmt.Fprintf(f, "query-result=%s", string(data)); err != nil {
			slog.Error("failed writing output", "err", err)
			return err
		}
	}

	return nil
}
