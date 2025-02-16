package repositories

import (
	"database/sql"
	"errors"
	"eshop/internal/domain/entities"
	"eshop/internal/infrastructure/constants"
	"eshop/internal/infrastructure/errs"
	"eshop/pkg/postgre"
	"fmt"
	"github.com/google/uuid"
)

type UserRepository struct {
	db postgre.DBinteraction
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(user *entities.User) (uuid.UUID, error) {
	var id uuid.UUID

	const query = `INSERT INTO users(id, username, password, email, isadmin, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

	err := u.db.QueryRow(query, user.ID, user.Username, user.Password, user.Email, user.IsAdmin, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return constants.EmptyID, fmt.Errorf("failed creating new user: %w", err)
	}

	return id, nil
}

func (u *UserRepository) GetByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	const query = `SELECT id, username, password, email, isadmin, created_at, updated_at FROM users WHERE id = $1;`

	err := u.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNoUserFound
	} else if err != nil {
		return nil, fmt.Errorf("refused to communicate with users database: %w", err)
	}

	return &user, nil
}

func (u *UserRepository) GetAll() ([]entities.User, error) {
	var users []entities.User
	const query = `SELECT id, username, password, email, isadmin, created_at, updated_at FROM users;`

	rows, err := u.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.User
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan structs: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate over users: %w", err)
	}

	return users, nil
}

func (u *UserRepository) Delete(id uuid.UUID) error {
	var returnedID uuid.UUID
	const query = `DELETE FROM users WHERE id = $1 RETURNING id;`

	err := u.db.QueryRow(query, id).Scan(&returnedID)
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNoUserFound
	}
	if (err != nil) || (returnedID != id) {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (u *UserRepository) Update(user *entities.User) error {
	const query = `UPDATE users SET username = $1, password = $2, email = $3, isadmin = $4, updated_at = $5 WHERE id = $6;`

	_, err := u.db.Exec(query, user.Username, user.Password, user.Email,
		user.IsAdmin, user.UpdatedAt, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return errs.ErrNoUserFound
	}
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}
