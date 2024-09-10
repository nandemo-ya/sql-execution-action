package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testCreateConfig = `
queries:
- sql: |
    CREATE TABLE IF NOT EXISTS users (
        "id"   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
        "name" TEXT,
        "age"  INTEGER
    );
  params: []
`
	testInsertConfig = `
queries:
- sql:
    INSERT INTO users (name, age) VALUES (?, ?);
  params:
    - hoge
    - 20
- sql:
    INSERT INTO users (name, age) VALUES (?, ?);
  params:
    - fuga
    - 30
- sql:
    INSERT INTO users (name, age) VALUES (?, ?);
  params:
    - piyo
    - 40
`

	testSelectConfig = `
query:
  sql:
    SELECT id, name, age FROM users;
  params: []
`

	testQueryResult = `query-result=[{"age":20,"id":1,"name":"hoge"},{"age":30,"id":2,"name":"fuga"},{"age":40,"id":3,"name":"piyo"}]`
)

func TestBuildArgs(t *testing.T) {
	output := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)

	dbFile, err := os.CreateTemp("", "test.db_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbFile.Name())

	rootCmd := CreateRoot()
	writeCmd := writeCommand()

	rootCmd.SetOut(output)
	rootCmd.SetErr(stdErr)
	rootCmd.AddCommand(writeCmd)
	rootCmd.AddCommand(readCommand())
	rootCmd.SetArgs([]string{
		"write",
		"--engine=sqlite3",
		"--datasource=" + dbFile.Name(),
		"--sql=" + testCreateConfig,
	})
	rootCmd.Execute()

	expected := &args{
		engine:     "sqlite3",
		datasource: dbFile.Name(),
		sqlPayload: []byte(testCreateConfig),
	}

	actual, err := buildArgs(rootCmd)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestWriteCommand_NoArgs(t *testing.T) {
	output := new(bytes.Buffer)

	rootCmd := CreateRoot()
	rootCmd.SetOut(output)
	rootCmd.AddCommand(writeCommand())

	rootCmd.SetArgs([]string{"write"})
	err := rootCmd.Execute()

	assert.Error(t, err)
}

func TestSelectCommand_NoArgs(t *testing.T) {
	output := new(bytes.Buffer)

	rootCmd := CreateRoot()
	rootCmd.SetOut(output)
	rootCmd.AddCommand(readCommand())

	rootCmd.SetArgs([]string{"read"})
	err := rootCmd.Execute()

	assert.Error(t, err)
}

func TestAllCommand_WithArgs(t *testing.T) {
	dbFile, err := os.CreateTemp("", "test.db_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbFile.Name())

	outputFile, err := os.CreateTemp("", "github_output_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(outputFile.Name())
	os.Setenv("GITHUB_OUTPUT", outputFile.Name())

	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)

	rootCmd := CreateRoot()
	rootCmd.SetOut(stdOut)
	rootCmd.SetErr(stdErr)
	rootCmd.AddCommand(writeCommand())
	rootCmd.AddCommand(readCommand())

	// The section of creating a table
	rootCmd.SetArgs([]string{
		"write",
		"--engine=sqlite3",
		"--datasource=" + dbFile.Name(),
		"--sql=" + testCreateConfig,
	})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	// The section of inserting data
	rootCmd.SetArgs([]string{
		"write",
		"--engine=sqlite3",
		"--datasource=" + dbFile.Name(),
		"--sql=" + testInsertConfig,
	})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	// The section of selecting data
	rootCmd.SetArgs([]string{
		"read",
		"--engine=sqlite3",
		"--datasource=" + dbFile.Name(),
		"--sql=" + testSelectConfig,
	})
	err = rootCmd.Execute()
	assert.NoError(t, err)

	outputData, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testQueryResult, string(outputData))
}
