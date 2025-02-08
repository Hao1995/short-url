package migrationkit

import (
	"context"
	"fmt"

	"github.com/pressly/goose/v3"
)

func GooseMigrate(dbString string, dir string) error {
	db, err := goose.OpenDBWithDriver("mysql", dbString)
	if err != nil {
		return fmt.Errorf("sql connection failed: %s", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := goose.RunContext(ctx, "up", db, dir); err != nil {
		return fmt.Errorf("goose up: %v", err)
	}

	return nil
}
