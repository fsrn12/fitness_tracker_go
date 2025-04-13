package store

import "database/sql"

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
	// if err != nil{
	// 	return nil, err
	// }
	return user, nil
}
