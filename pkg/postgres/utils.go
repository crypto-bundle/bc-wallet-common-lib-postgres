package postgres

import (
	"database/sql"
	"errors"
	"fmt"
)

func EmptyOrError(err error, errorMessage string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return fmt.Errorf("%w:%s", err, errorMessage)
}
