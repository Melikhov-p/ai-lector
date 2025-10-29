package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Storage struct {
	db *sql.DB
}

func NewPostgresStorage(connectionString string) (*Storage, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	// Применяем миграции
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, usr *user.User) error {
	const op = "repository.postgres.SaveUser"

	query := `INSERT INTO users (first_name, phone, pass_hash) VALUES ($1, $2, $3) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, usr.FirstName(), usr.Phone(), usr.PassHash())

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	usr.SetID(userID)
	return nil
}

func (s *Storage) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	const op = "repository.postgres.GetUserByID"

	query := `SELECT uuid, trial_requests, email, phone, pass_hash,
       first_name, last_name, age, created_at, updated_at FROM users WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, userID)

	var (
		userUUID      uuid.UUID
		trialRequests int
		email         string
		phone         string
		passHash      []byte
		firstName     string
		lastName      string
		age           int
		createdAt     time.Time
		updatedAt     time.Time
	)

	if err := row.Scan(&userUUID, &trialRequests, &email, &phone, &passHash,
		&firstName, &lastName, &age, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.NewUserFromDB(
		userID,
		userUUID,
		trialRequests,
		email,
		phone,
		passHash,
		firstName,
		lastName,
		age,
		createdAt,
		updatedAt,
	), nil
}

func (s *Storage) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	const op = "repository.postgres.GetUserByPhone"

	query := `SELECT id, uuid, trial_requests, email, pass_hash,
       first_name, last_name, age, created_at, updated_at FROM users WHERE phone = $1`

	row := s.db.QueryRowContext(ctx, query, phone)

	var (
		id            int64
		userUUID      uuid.UUID
		trialRequests int
		email         string
		passHash      []byte
		firstName     string
		lastName      string
		age           int
		createdAt     time.Time
		updatedAt     time.Time
	)

	if err := row.Scan(&id, &userUUID, &trialRequests, &email, &passHash,
		&firstName, &lastName, &age, &createdAt, &updatedAt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user.NewUserFromDB(
		id,
		userUUID,
		trialRequests,
		email,
		phone,
		passHash,
		firstName,
		lastName,
		age,
		createdAt,
		updatedAt,
	), nil
}

func (s *Storage) SearchUser(ctx context.Context, filter user.UserFilter) ([]*user.User, error) {

}

func (s *Storage) GetAllUsers(ctx context.Context) ([]*user.User, error) {

}

func (s *Storage) AddInterestToUser(ctx context.Context, userID int64, inter *interest.Interest) error {

}

func (s *Storage) GetUserInterestsByID(ctx context.Context, userID int64) ([]*interest.Interest, error) {

}

func (s *Storage) UpdateUser(ctx context.Context, usr *user.User) error {

}
