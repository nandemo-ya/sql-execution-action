package cmd

import (
	"fmt"
	"log/slog"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/spf13/cobra"
)

func CreateRoot() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "sql-execution-action",
	}
	rootCmd.PersistentFlags().String("engine", "sqlite3", "The relation database engine name, e.g. sqlite3(default), mysql, postgres, mssql, and oracle.")
	rootCmd.PersistentFlags().String("datasource", "", "The datasource string of the relational database. e.g. /path/to/sqlite.db, user:password@tcp(localhost:3306)/dbname")
	rootCmd.PersistentFlags().String("sql", "", "The SQL config string in YAML format.")
	rootCmd.PersistentFlags().String("sql-file", "", "The path of SQL config string file.")
	rootCmd.MarkPersistentFlagRequired("datasource")

	rootCmd.AddCommand(readCommand())
	rootCmd.AddCommand(writeCommand())
	return rootCmd
}

type args struct {
	engine     string
	datasource string
	sqlPayload []byte
}

func buildArgs(cmd *cobra.Command) (*args, error) {
	engine, err := cmd.Flags().GetString("engine")
	if err != nil {
		slog.Error("failed getting engine flag", "err", err)
		return nil, err
	}

	ds, err := cmd.Flags().GetString("datasource")
	if err != nil {
		slog.Error("failed getting datasource flag", "err", err)
		return nil, err
	}

	sql, err := cmd.Flags().GetString("sql")
	if err != nil {
		slog.Error("failed getting sql flag", "err", err)
		return nil, err
	}

	sqlFile, err := cmd.Flags().GetString("sql-file")
	if err != nil {
		slog.Error("failed getting sql-file flag", "err", err)
		return nil, err
	}

	if sql != "" && sqlFile != "" {
		return nil, fmt.Errorf("set either sql or sql-file attribute")
	}

	var sqlPayload []byte
	if sqlFile != "" {
		data, err := os.ReadFile(sqlFile)
		if err != nil {
			slog.Error("failed reading sql file", "err", err)
			return nil, err
		}
		sqlPayload = data
	} else if sql != "" {
		sqlPayload = []byte(sql)
	} else {
		return nil, fmt.Errorf("set either sql or sql-file attribute")
	}

	return &args{
		engine:     engine,
		datasource: ds,
		sqlPayload: sqlPayload,
	}, nil
}
