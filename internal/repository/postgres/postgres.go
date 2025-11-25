package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/subscription"
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
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	// Проверяем соединение
	if err = db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	// Применяем миграции
	if err = goose.Up(db, "./internal/repository/migrations"); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveInterest(ctx context.Context, interest *interest.Interest) (int64, error) {
	const op = "repository.postgres.SaveInterest"

	query := `INSERT INTO interests (name, emoji) VALUES ($1, $2) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, interest.Name(), interest.Emoji())
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) DeleteInterest(ctx context.Context, id int64) error {
	const op = "repository.postgres.DeleteInterest"

	query := `DELETE FROM interests WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetInterestByID(ctx context.Context, id int64) (*interest.Interest, error) {
	const op = "repository.postgres.GetInterestByID"

	query := `SELECT name, emoji FROM interests WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	var (
		name, emoji string
	)

	if err := row.Scan(&name, &emoji); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interest.NewInterestFromDB(id, name, emoji), nil
}

func (s *Storage) GetInterestByName(ctx context.Context, name string) (*interest.Interest, error) {
	const op = "repository.postgres.GetInterestByName"

	query := `SELECT id, emoji FROM interests WHERE name = $1`
	row := s.db.QueryRowContext(ctx, query, name)

	var (
		id    int64
		emoji string
	)

	if err := row.Scan(&id, &emoji); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return interest.NewInterestFromDB(id, name, emoji), nil
}

func (s *Storage) GetAllInterests(ctx context.Context) ([]*interest.Interest, error) {
	const op = "repository.postgres.GetAllInterests"

	var (
		inters []*interest.Interest
		err    error
	)

	rows, err := s.db.QueryContext(ctx, "SELECT id, name, emoji FROM interests")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var (
			id          int64
			name, emoji string
		)

		if err = rows.Scan(&id, &name, &emoji); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		inters = append(inters, interest.NewInterestFromDB(id, name, emoji))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return inters, nil
}

func (s *Storage) SaveUser(ctx context.Context, usr *user.User) (int64, error) {
	const op = "repository.postgres.SaveUser"

	query := `INSERT INTO users (
                  uuid, trial_requests, email, phone, pass_hash,
       first_name, last_name, age, class, created_at, updated_at
                   ) VALUES (
                             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
                             ) RETURNING id`

	row := s.db.QueryRowContext(
		ctx,
		query,
		usr.UUID(), usr.TrialRequests(), usr.Email(), usr.Phone(), usr.PassHash(), usr.FirstName(), usr.LastName(),
		usr.Age(), usr.Class(), usr.CreatedAt(), usr.UpdatedAt(),
	)

	var userID int64
	if err := row.Scan(&userID); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	const op = "repository.postgres.GetUserByID"

	query := `SELECT uuid, trial_requests, email, phone, pass_hash,
       first_name, last_name, age, class, created_at, updated_at FROM users WHERE id = $1`

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
		class         int
		createdAt     time.Time
		updatedAt     time.Time
	)

	if err := row.Scan(&userUUID, &trialRequests, &email, &phone, &passHash,
		&firstName, &lastName, &age, &class, &createdAt, &updatedAt); err != nil {
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
		class,
		createdAt,
		updatedAt,
	), nil
}

func (s *Storage) GetUserByPhone(ctx context.Context, phone string) (*user.User, error) {
	const op = "repository.postgres.GetUserByPhone"

	query := `SELECT id, uuid, trial_requests, email, pass_hash,
       first_name, last_name, age, class, created_at, updated_at FROM users WHERE phone = $1`

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
		class         int
		createdAt     time.Time
		updatedAt     time.Time
	)

	if err := row.Scan(&id, &userUUID, &trialRequests, &email, &passHash,
		&firstName, &lastName, &age, &class, &createdAt, &updatedAt); err != nil {
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
		class,
		createdAt,
		updatedAt,
	), nil
}

func (s *Storage) SearchUser(ctx context.Context, filter user.UserFilter) ([]*user.User, error) {
	const op = "repository.postgres.SearchUser"

	var (
		conditions []string
		params     []interface{}
		usrs       []*user.User
	)

	// Собираем условия и параметры отдельно
	if filter.Email != "" {
		conditions = append(conditions, "email = $"+strconv.Itoa(len(params)+1))
		params = append(params, filter.Email)
	}
	if filter.Phone != "" {
		conditions = append(conditions, "phone = $"+strconv.Itoa(len(params)+1))
		params = append(params, filter.Phone)
	}
	if filter.FirstName != "" {
		conditions = append(conditions, "first_name = $"+strconv.Itoa(len(params)+1))
		params = append(params, filter.FirstName)
	}
	if filter.LastName != "" {
		conditions = append(conditions, "last_name = $"+strconv.Itoa(len(params)+1))
		params = append(params, filter.LastName)
	}

	// Формируем базовый запрос
	query := `SELECT id, uuid, trial_requests, email, phone, pass_hash,
		first_name, last_name, age, class, created_at, updated_at
		FROM users`

	// Добавляем WHERE только если есть условия
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var (
		id            int64
		userUUID      uuid.UUID
		trialRequests int
		email         string
		phone         string
		passHash      []byte
		firstName     string
		lastName      string
		age           int
		class         int
		createdAt     time.Time
		updatedAt     time.Time
	)

	rows, err := s.db.QueryContext(ctx, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if err = rows.Scan(&id, &userUUID, &trialRequests, &email, &phone, &passHash,
			&firstName, &lastName, &age, &class, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		usrs = append(usrs, user.NewUserFromDB(
			id,
			userUUID,
			trialRequests,
			email,
			phone,
			passHash,
			firstName,
			lastName,
			age,
			class,
			createdAt,
			updatedAt,
		))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usrs, nil
}

func (s *Storage) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	const op = "repository.postgres.GetAllUsers"

	var (
		usrs []*user.User
		err  error
	)

	query := `SELECT id, uuid, trial_requests, email, phone, pass_hash,
		first_name, last_name, age, class, created_at, updated_at FROM users`

	var (
		id            int64
		userUUID      uuid.UUID
		trialRequests int
		email         string
		phone         string
		passHash      []byte
		firstName     string
		lastName      string
		age           int
		class         int
		createdAt     time.Time
		updatedAt     time.Time
	)

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if err = rows.Scan(&id, &userUUID, &trialRequests, &email, &phone, &passHash,
			&firstName, &lastName, &age, &class, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		usrs = append(usrs, user.NewUserFromDB(
			id,
			userUUID,
			trialRequests,
			email,
			phone,
			passHash,
			firstName,
			lastName,
			age,
			class,
			createdAt,
			updatedAt,
		))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return usrs, nil
}

func (s *Storage) AddInterestToUser(ctx context.Context, userID int64, inter *interest.Interest) error {
	const op = "repository.postgres.AddInterestToUser"

	query := `INSERT INTO users_interests (user_id, interest_id) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, userID, inter.ID())
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserInterestsByID(ctx context.Context, userID int64) ([]int64, error) {
	const op = "repository.postgres.GetUserInterestsByID"

	var ids []int64

	query := `SELECT interest_id FROM users_interests WHERE user_id = $1`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, interest.ErrInterestNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return ids, nil
}

func (s *Storage) UpdateUser(ctx context.Context, usr *user.User) error {
	const op = "repository.postgres.UpdateUser"

	query := `UPDATE users
SET
    trial_requests = $2,
    email = $3,
    phone = $4,
    pass_hash = $5,
    first_name = $6,
    last_name = $7,
    age = $8,
    class = $9,
    updated_at = NOW()
WHERE uuid = $1;`

	result, err := s.db.ExecContext(ctx, query,
		usr.UUID(),
		usr.TrialRequests(),
		usr.Email(),
		usr.Phone(),
		usr.PassHash(),
		usr.FirstName(),
		usr.LastName(),
		usr.Age(),
		usr.Class(),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, что хотя бы одна строка была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: user with uuid %s not found", op, usr.UUID())
	}

	return nil
}

func (s *Storage) SubscribeUser(ctx context.Context, usrSub *subscription.UserSubscription) error {
	const op = "repository.postgres.SubscribeUser"

	query := `INSERT INTO users_subscriptions (user_id, sub_id, created_at, end_at, payed) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.ExecContext(
		ctx,
		query,
		usrSub.UserID(), usrSub.SubID(), usrSub.CreatedAt(), usrSub.EndAt(), usrSub.Payed(),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
