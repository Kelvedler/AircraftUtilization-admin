package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ApiUserInput struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Secret    string `json:"secret"`
	Uses      string `json:"uses"`
	Active    string `json:"active"`
}

type ApiUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"       validate:"gte=3,lte=30"`
	Secret    string    `json:"secret"`
	Uses      int32     `json:"uses"`
	Active    bool      `json:"active"`
}

type ApiUserRange struct {
	ApiUsers []ApiUser
	Limit    uint8
	Offset   uint8
}

func (input ApiUserInput) Bind() (output ApiUser, err error) {
	if input.ID != "" {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return ApiUser{}, err
		}
		output.ID = id
	}
	if input.CreatedAt != "" {
		createdAt, err := strconv.Atoi(input.CreatedAt)
		if err != nil {
			return ApiUser{}, err
		}
		output.CreatedAt = time.UnixMilli(int64(createdAt)).UTC()
	}
	if input.UpdatedAt != "" {
		updatedAt, err := strconv.Atoi(input.UpdatedAt)
		if err != nil {
			return ApiUser{}, err
		}
		output.UpdatedAt = time.UnixMilli(int64(updatedAt)).UTC()
	}
	output.Name = input.Name
	output.Secret = input.Secret
	if input.Uses != "" {
		uses, err := strconv.Atoi(input.Uses)
		if err != nil {
			return ApiUser{}, err
		}
		output.Uses = int32(uses)
	}
	if input.Active != "" {
		active, err := strconv.ParseBool(input.Active)
		if err != nil {
			return ApiUser{}, err
		}
		output.Active = active
	}
	return output, nil
}

func (a ApiUser) createQueue(batch *pgx.Batch) {
	insertCols := "name, secret"
	insertVals := "$1, $2"
	returningCols := "id, created_at, updated_at, uses, active"
	query := fmt.Sprintf(
		"INSERT  into api_user(%s) VALUES(%s) RETURNING %s",
		insertCols,
		insertVals,
		returningCols,
	)
	batch.Queue(query, a.Name, a.Secret)
}

func (a *ApiUser) createResult(result pgx.BatchResults) error {
	return result.QueryRow().Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Uses,
		&a.Active,
	)
}

func (a *ApiUser) Create() (BatchOperation, BatchRead) {
	return a.createQueue, a.createResult
}

func (a ApiUserRange) getQueue(batch *pgx.Batch) {
	selectCols := "id, created_at, updated_at, name, uses, active"
	orderCols := "created_at"
	query := fmt.Sprintf(
		"SELECT %s FROM api_user ORDER BY %s LIMIT $1 OFFSET $2",
		selectCols,
		orderCols,
	)
	batch.Queue(query, a.Limit, a.Offset)
}

func (a *ApiUserRange) getResult(result pgx.BatchResults) error {
	rows, err := result.Query()
	if err != nil {
		return err
	}
	next := rows.Next()
	if !next {
		return nil
	}
	for next {
		var apiUser ApiUser
		err = rows.Scan(
			&apiUser.ID,
			&apiUser.CreatedAt,
			&apiUser.UpdatedAt,
			&apiUser.Name,
			&apiUser.Uses,
			&apiUser.Active,
		)
		if err != nil {
			return err
		}
		a.ApiUsers = append(a.ApiUsers, apiUser)
		next = rows.Next()
	}
	return nil
}

func (a *ApiUserRange) Get() (BatchOperation, BatchRead) {
	return a.getQueue, a.getResult
}
