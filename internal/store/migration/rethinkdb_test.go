package migration_test

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store/migration"
)

func TestRethinkDB(t *testing.T) {
	err := migration.RethinkDB("mctest", "localhost:30815")
	assert.Okf(t, err, "Failed to perform migration for RethinkDB: %s", err)
}
