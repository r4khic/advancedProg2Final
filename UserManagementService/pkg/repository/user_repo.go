package repository

import (
	"advancedProg2Final/UserManagementService/pkg/entity"
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int64) (*entity.User, error) {
	user := &entity.User{}
	query := `SELECT u.id, u.username, u.password, u.email, r.id as role_id, r.name as role_name 
              FROM users u
              INNER JOIN roles r ON u.role_id = r.id
              WHERE u.id = $1`

	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role.ID, &user.Role.Name)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Save(user *entity.User) (int64, error) {
	query := `INSERT INTO users (username, password, email, role_id) 
          VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role.ID).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *UserRepository) Authenticate(username, password string) (*entity.User, error) {
	log.Printf("Authenticating user: %s\n", username)

	user := &entity.User{}
	query := `SELECT u.id, u.username, u.password, u.email, r.id as role_id, r.name as role_name 
              FROM users u
              INNER JOIN roles r ON u.role_id = r.id
              WHERE u.username = $1`

	row := r.db.QueryRow(query, username)
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role.ID, &user.Role.Name)

	// If no rows are returned, return a more specific error.
	if err == sql.ErrNoRows {
		log.Printf("No user found with username: %s\n", username)
		return nil, fmt.Errorf("no user found with username: %v", username)
	} else if err != nil {
		log.Printf("Error retrieving user: %s\n", err)
		return nil, err
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// If the password doesn't match, return an error.
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Printf("Invalid password for user: %s\n", username)
		return nil, fmt.Errorf("invalid password")
	} else if err != nil {
		log.Printf("Error comparing password hashes: %s\n", err)
		return nil, err
	}

	log.Printf("User authenticated successfully: %s\n", username)
	return user, nil
}

func (r *UserRepository) SaveUser(user *entity.User) error {
	log.Printf("Saving user: %s\n", user.Username)

	query := `UPDATE users SET username = $1, password = $2, email = $3, role_id = $4 
              WHERE id = $5`

	// Prepare a hashed version of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %s\n", err)
		return err
	}

	// Execute the SQL query
	_, err = r.db.Exec(query, user.Username, string(hashedPassword), user.Email, user.Role.ID, user.ID)
	if err != nil {
		log.Printf("Error saving user: %s\n", err)
		return err
	}

	log.Printf("User saved successfully: %s\n", user.Username)
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	log.Printf("Deleting user with ID: %d\n", id)

	query := `DELETE FROM users WHERE id = $1`

	// Execute the SQL query
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting user: %s\n", err)
		return err
	}

	log.Printf("User deleted successfully with ID: %d\n", id)
	return nil
}
