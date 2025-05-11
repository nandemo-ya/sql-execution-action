// +build !test

// This file is separated from root.go to allow conditional compilation of DuckDB support.
// The DuckDB driver depends on Apache Arrow C++ libraries, which may not be available
// in all environments, especially during testing. By using the build tag '!test',
// this file is excluded when running tests with the 'test' tag.
//
// When building the actual application, this file is included, enabling DuckDB support.
// When running tests with 'go test -tags=test', this file is excluded, avoiding
// dependency errors related to Arrow C++ libraries.
//
// This approach allows the codebase to support DuckDB in production while maintaining
// testability in environments without the required C++ dependencies.
package cmd

import (
	_ "github.com/marcboeker/go-duckdb"
)
