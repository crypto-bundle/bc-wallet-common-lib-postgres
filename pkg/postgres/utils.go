package postgres

import (
	"database/sql"
	"fmt"
)

func EmptyOrError(err error, errorMessage string) error {
	if err == sql.ErrNoRows {
		return nil
	}
	return fmt.Errorf("%w:%s", err, errorMessage)
}
