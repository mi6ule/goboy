package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	command_model "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/model/command"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
)

type UserRepository interface {
	Create(user *command_model.User) error
	GetById(id int64) (*command_model.User, error)
	Update(user *command_model.User) error
	Delete(id int64) error
	// ... other methods
}

type SQLUserRepository struct {
	db *persistence.Database
}

func NewSQLUserRepository(db *persistence.Database) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) Create(user *command_model.User) error {
	// Perform the SQL INSERT operation to create a new user record in the database
	_, err := r.db.Exec("INSERT INTO users (id, username, email, password, person_id) VALUES (?, ?, ?, ?, ?)",
		user.ID, user.Username, user.Email, user.Password, user.PersonID, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *SQLUserRepository) GetByID(id uuid.UUID) (*command_model.User, error) {
	// Perform the SQL SELECT operation to retrieve a user record by its ID from the database
	row := r.db.QueryRow("SELECT id, username, email, person_id, created_at, updated_at FROM users WHERE id = $1", id)
	user := &command_model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PersonID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) GetByUsername(username string) (*command_model.User, error) {
	// Perform the SQL SELECT operation to retrieve a user record by its username from the database
	row := r.db.QueryRow("SELECT id, username, email, person_id, created_at, updated_at FROM users WHERE username = ?", username)
	user := &command_model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PersonID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) Update(user *command_model.User) error {
	// Perform the SQL UPDATE operation to update a user record in the database
	_, err := r.db.Exec("UPDATE users SET username = ?, email = ?, person_id = ?, updated_at = ? WHERE id = ?",
		user.Username, user.Email, user.PersonID, user.UpdatedAt, user.ID)
	return err
}

func (r *SQLUserRepository) Delete(id uuid.UUID) error {
	// Perform the SQL DELETE operation to delete a user record from the database
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
