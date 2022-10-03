// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: query.sql

package users

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users(
    username, password, email,first_name,last_name
) VALUES (
             $1, $2, $3,$4,$5
         )
    RETURNING id, username, email, password, first_name, last_name
`

type CreateUserParams struct {
	Username  string
	Password  string
	Email     string
	FirstName string
	LastName  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserById = `-- name: GetUserById :one
SELECT id, username, email, password, first_name, last_name FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, username, email, password, first_name, last_name FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}

const getUserPassword = `-- name: GetUserPassword :one
SELECT password FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUserPassword(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserPassword, username)
	var password string
	err := row.Scan(&password)
	return password, err
}

const getUserView = `-- name: GetUserView :one
SELECT username, email,first_name,last_name FROM users
WHERE id = $1 LIMIT 1
`

type GetUserViewRow struct {
	Username  string
	Email     string
	FirstName string
	LastName  string
}

func (q *Queries) GetUserView(ctx context.Context, id uuid.UUID) (GetUserViewRow, error) {
	row := q.db.QueryRowContext(ctx, getUserView, id)
	var i GetUserViewRow
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, password, first_name, last_name FROM users
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.FirstName,
			&i.LastName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
set username = $2,
    password = $3,
    email =$4,
    first_name = $5,
    last_name = $6
WHERE id = $1
    RETURNING id, username, email, password, first_name, last_name
`

type UpdateUserParams struct {
	ID        uuid.UUID
	Username  string
	Password  string
	Email     string
	FirstName string
	LastName  string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Password,
		&i.FirstName,
		&i.LastName,
	)
	return i, err
}
