package model

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

func DuplicateError(err error) bool {
	var e *mysql.MySQLError

	if errors.As(err, &e) {
		// mariadb-10.5.9/libmariadb/include/mysqld_error.h:69:#define ER_DUP_ENTRY 1062
		if e.Number == 1062 {
			return true
		}
	}

	return false
}

func RecordNotFound(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	return false
}
