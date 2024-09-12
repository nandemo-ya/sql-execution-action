# sql-execution-action

This action is designed to make it easy to run SQL on the GitHub Actions runner. You can persist data through caching by configuring your database engine as SQLite3.

<!-- action-docs-inputs source="action.yml" -->
## Inputs

| name | description | required | default |
| --- | --- | --- | --- |
| `command` | <p>The command to run the query. e.g. write or read.</p> | `true` | `""` |
| `engine` | <p>The relation database engine name, e.g. sqlite3(default), mysql, postgres, and mssql.</p> | `false` | `sqlite3` |
| `datasource` | <p>The datasource string of the relational database. e.g. /path/to/sqlite.db, user:password@tcp(localhost:3306)/dbname</p> | `true` | `""` |
| `sql-file` | <p>The SQL config string in YAML format.</p> | `false` | `""` |
| `sql` | <p>The path of SQL config string file.</p> | `false` | `""` |
<!-- action-docs-inputs source="action.yml" -->

<!-- action-docs-outputs source="action.yml" -->
## Outputs

| name | description |
| --- | --- |
| `query-result` | <p>JSON array of query results.</p> |
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
    # The relation database engine name, e.g. sqlite3(default), mysql, postgres, and mssql.
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

## Configuration

The query to be executed is set using the "sql" or "sql-file" parameter. The SQL configuration file has the following schema:

### Write operation schema

The write operation schema is as follows:

```yaml
queries:
- sql: |
    CREATE TABLE IF NOT EXISTS users (
        "id"   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
        "name" TEXT,
        "age"  INTEGER
    );
  params: []
```

In addition, you can set multiple INSERT statements as follows:

```yaml
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
```

### Read operation schema

The read operation schema is as follows:

```yaml
query:
  sql:
    SELECT id, name, age FROM users WHERE name = ?;
  params:
    - fuga
```

The query result will be output in JSON format to the `query-result` variable.

## Examples

Examples of usage with SQLite3 and MySQL are shown below.

### SQLite3

SQLite3 data files are persisted via `actions/cache`.

```yaml
# Setup the database file as a cache
- name: Cache database file
  uses: actions/cache@v4
  with:
    path: sqlite3.db
    key: ${{ runner.os }}-database-
    restore-keys: |
      ${{ runner.os }}-database-
- uses: nandemo-ya/sql-execution-action@main
  with:
    command: write
    engine: sqlite3
    datasource: sqlite3.db
    sql: |
      queries:
        - sql: |
            CREATE TABLE IF NOT EXISTS users (
                "id"   INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,        
                "name" TEXT,
                "age"  INTEGER
            );
          params: []
- uses: nandemo-ya/sql-execution-action@main
  with:
    command: write
    engine: sqlite3
    datasource: sqlite3.db
    sql: |
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
- uses: nandemo-ya/sql-execution-action@main
  # Set the "id" to reference the result in a later step.
  id: execute_query
  with:
    command: read
    engine: sqlite3
    datasource: sqlite3.db
    sql: |
      query:
        sql:
          SELECT id, name, age FROM users WHERE name = ?;
        params:
          - fuga
- name: Show query result
  run: |
    echo "query result is: ${{ steps.execute_query.outputs.query-result }}"
```

### MySQL

For MySQL, start the mysql service first.

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping --silent"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3
        env:
          MYSQL_ROOT_PASSWORD: root_password
          MYSQL_DATABASE: test_db
          MYSQL_USER: test_user
          MYSQL_PASSWORD: test_password
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      - name: Wait for MySQL to be ready
        run: |
          while ! mysqladmin ping --host="127.0.0.1" --user="root" --password="root_password" --silent; do
            echo "Waiting for MySQL..."
            sleep 3
          done
      - name: Execute create statement
        uses: ./
        with:
          command: write
          engine: mysql
          datasource: test_user:test_password@tcp(mysql:3306)/test_db
          sql: |
            queries:
              - sql: |
                  CREATE TABLE IF NOT EXISTS users (
                      `id`   INT NOT NULL AUTO_INCREMENT,
                      `name` TEXT,
                      `age`  INT,
                      PRIMARY KEY (`id`)
                  );
                params: []
      - name: Execute insert statement
        uses: ./
        with:
          command: write
          engine: mysql
          datasource: test_user:test_password@tcp(mysql:3306)/test_db
          sql: |
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
      - name: Execute select statement
        uses: ./
        id: execute_query
        with:
          command: read
          engine: mysql
          datasource: test_user:test_password@tcp(mysql:3306)/test_db
          sql: |
            query:
              sql:
                SELECT id, name, age FROM users WHERE name = ?;
              params:
                - fuga

      - name: Show query result
        run: |
          echo "query result is: ${{ steps.execute_query.outputs.query-result }}"
```

This action resolves the MySQL database by its service name because the sql-execution-action is executed in the container.


## License

This project is distributed under the [MIT license](LICENSE).