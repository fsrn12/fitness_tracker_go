package store

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password"`
	Bio          string `json:"bio"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByID(id int64) (*User, error)
	UpdateUser(*User) error
	DeleteUser(id int64) error
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `
	INSERT INTO users (username, email, password_hash, bio)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	err = tx.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.Bio).Scan(&user.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return user, nil

}

func (pg *PostgresUserStore) GetUserByID(id int64) (*User, error) {
	user := &User{}

	userQuery := `
	SELECT id, username, email, password_hash, bio
	FROM users
	WHERE id = $1
	`

	err := pg.db.QueryRow(userQuery, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Bio)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userQuery := `
	UPDATE users
	SET username = $1, email = $2, bio = $3
	WHERE id = $4
	`

	result, err := tx.Exec(userQuery, user.Username, user.Email, user.Bio, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return tx.Commit()
}

func (pg *PostgresUserStore) DeleteUser(id int64) error {
	query := `DELETE from users WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	if err != nil {
		return err
	}

	fmt.Println("Item Deleted Successfully")

	return nil
}
