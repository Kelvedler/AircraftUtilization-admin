package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AdminInput struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Name      string `json:"name"`
	Password  string `json:"password"`
}

type Admin struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"       validate:"gte=3,lte=30"`
	Password  string    `json:"password"   validate:"gte=6,lte=20"`
}

func (input AdminInput) Bind() (output Admin, err error) {
	if input.ID != "" {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return Admin{}, err
		}
		output.ID = id
	}
	if input.CreatedAt != "" {
		createdAt, err := strconv.Atoi(input.CreatedAt)
		if err != nil {
			return Admin{}, err
		}
		output.CreatedAt = time.UnixMilli(int64(createdAt)).UTC()
	}
	if input.UpdatedAt != "" {
		updatedAt, err := strconv.Atoi(input.UpdatedAt)
		if err != nil {
			return Admin{}, err
		}
		output.UpdatedAt = time.UnixMilli(int64(updatedAt)).UTC()
	}
	output.Name = input.Name
	output.Password = input.Password
	return output, nil
}

func (a Admin) createQueue(batch *pgx.Batch) {
	insertCols := "name, password"
	insertVals := "$1, $2"
	returningCols := "id, created_at, updated_at"
	query := fmt.Sprintf(
		"INSERT into admin(%s) VALUES(%s) RETURNING %s",
		insertCols,
		insertVals,
		returningCols,
	)
	batch.Queue(query, a.Name, a.Password)
}

func (a *Admin) createResult(result pgx.BatchResults) error {
	return result.QueryRow().Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (a *Admin) Create() (BatchOperation, BatchRead) {
	return a.createQueue, a.createResult
}

func (a Admin) getByNameQueue(batch *pgx.Batch) {
	selectCols := "id, created_at, updated_at, name, password"
	whereVals := "name=$1"
	query := fmt.Sprintf("SELECT %s FROM admin WHERE %s", selectCols, whereVals)
	batch.Queue(query, a.Name)
}

func (a *Admin) getByNameResult(result pgx.BatchResults) error {
	return result.QueryRow().Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.Name, &a.Password)
}

func (a *Admin) GetByName() (BatchOperation, BatchRead) {
	return a.getByNameQueue, a.getByNameResult
}

func (a Admin) getByIdQueue(batch *pgx.Batch) {
	selectCols := "id, created_at, updated_at, name, password"
	whereVals := "id=$1"
	query := fmt.Sprintf("SELECT %s FROM admin WHERE %s", selectCols, whereVals)
	batch.Queue(query, a.ID)
}

func (a *Admin) getByIdResult(result pgx.BatchResults) error {
	return result.QueryRow().Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.Name, &a.Password)
}

func (a *Admin) GetById() (BatchOperation, BatchRead) {
	return a.getByIdQueue, a.getByIdResult
}
