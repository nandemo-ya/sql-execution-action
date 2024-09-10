# sql-execution-action

This action is designed to make it easy to run SQL on the GitHub Actions runner. You can persist data through caching by configuring your database engine as SQLite3.

<!-- action-docs-inputs source="action.yml" -->
## Inputs

| name | description | required | default |
| --- | --- | --- | --- |
| `command` | <p>The command to run the query. e.g. write or read.</p> | `true` | `""` |
| `engine` | <p>The relation database engine name, e.g. sqlite3(default), mysql, postgres, mssql, and oracle.</p> | `false` | `sqlite3` |
| `datasource` | <p>The datasource string of the relational database. e.g. /path/to/sqlite.db, user:password@tcp(localhost:3306)/dbname</p> | `true` | `""` |
| `sql-file` | <p>The SQL config string in YAML format.</p> | `false` | `""` |
| `sql` | <p>The path of SQL config string file.</p> | `false` | `""` |
<!-- action-docs-inputs source="action.yml" -->

<!-- action-docs-outputs source="action.yml" -->
## Outputs

| name | description |
| --- | --- |
| `query-result` | <p>The result of the query</p> |
<!-- action-docs-outputs source="action.yml" -->

<!-- action-docs-runs action="action.yml" -->
## Runs

This action is a `docker` action.
<!-- action-docs-runs action="action.yml" -->

<!-- action-docs-usage action="action.yml" project="nandemo-ya/sql-execution-action" version="main" -->
## Usage

```yaml
- uses: nandemo-ya/sql-execution-action@main
  with:
    command:
    # The command to run the query. e.g. write or read.
    #
    # Required: true
    # Default: ""

    engine:
    # The relation database engine name, e.g. sqlite3(default), mysql, postgres, mssql, and oracle.
    #
    # Required: false
    # Default: sqlite3

    datasource:
    # The datasource string of the relational database. e.g. /path/to/sqlite.db, user:password@tcp(localhost:3306)/dbname
    #
    # Required: true
    # Default: ""

    sql-file:
    # The SQL config string in YAML format.
    #
    # Required: false
    # Default: ""

    sql:
    # The path of SQL config string file.
    #
    # Required: false
    # Default: ""
```
<!-- action-docs-usage action="action.yml" project="nandemo-ya/sql-execution-action" version="main" -->

## License

This project is distributed under the [MIT license](LICENSE).