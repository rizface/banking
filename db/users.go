package db

import (
	"banking/configs"
	"banking/db/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	config configs.Config
	dbPool *pgxpool.Pool
}

func NewUser(dbPool *pgxpool.Pool, config configs.Config) *User {
	return &User{
		dbPool: dbPool,
		config: config,
	}
}

func (u *User) Register(ctx context.Context, usr entity.User) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), u.config.BcryptSalt)
	if err != nil {
		return entity.User{}, err
	}

	var existingId string

	err = conn.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, usr.Email).Scan(&existingId)
	if existingId != "" {
		return entity.User{}, errors.New("EXISTING_EMAIL")
	}

	sql := `
		INSERT INTO users (name, email, password) VALUES ($1, $2, $3)
	`

	_, err = conn.Exec(ctx, sql, usr.Name, usr.Email, string(hashedPassword))

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, email FROM users WHERE email = $1`, usr.Email).Scan(&result.Id, &result.Name, &result.Email)

	if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		Id:    result.Id,
		Name:  result.Name,
		Email: result.Email,
	}, nil
}

func (u *User) Login(ctx context.Context, email, password string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, email, password FROM users WHERE email = $1`, email).Scan(
		&result.Id, &result.Name, &result.Email, &result.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		return result, errors.New("INVALID_PASSWORD")
	}

	return result, nil
}

func (u *User) GetUserById(ctx context.Context, userID string) (entity.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return entity.User{}, err
	}
	defer conn.Release()

	var result entity.User

	err = conn.QueryRow(ctx, `SELECT id, name, email FROM users WHERE id = $1`, userID).Scan(&result.Id, &result.Name, &result.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, ErrNoRow
	}
	if err != nil {
		return result, err
	}

	return result, nil
}
