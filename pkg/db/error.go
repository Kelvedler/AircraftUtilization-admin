package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const uniqueViolation = "23505"

type UniqueViolation struct {
	table  string
	column string
}

type DoesNotExist struct{}

type ContextCanceled struct{}

func getColumn(pgErr *pgconn.PgError) string {
	columnRe := regexp.MustCompile(fmt.Sprintf("%s_([a-z_]+)_(?:[a-z]+)", pgErr.TableName))
	column := columnRe.FindStringSubmatch(pgErr.ConstraintName)[1]
	return cases.Title(language.English).String(column)
}

func ErrorAsStruct(err error) interface{} {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolation:
			return UniqueViolation{
				table:  pgErr.TableName,
				column: getColumn(pgErr),
			}
		default:
			panic(fmt.Sprintf("Unforseen case - %s code", pgErr.Code))
		}
	}
	switch err {
	case pgx.ErrNoRows:
		return DoesNotExist{}
	case context.Canceled:
		return ContextCanceled{}
	default:
		panic(fmt.Sprintf("Unforseen case - %s", err))
	}
}
