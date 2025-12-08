package storage

import (
	"context"
	"yus/internal/models"

	// "github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type InMemoryStore interface {
	CreateClient(ctx context.Context) (*redis.Client, error)
	GetConnection(ctx context.Context) (*redis.Client, error)

	GenerateSessionID(ctx context.Context) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error

	AdminExists(ctx context.Context) (bool, error)
	CreateAdminSession(ctx context.Context, adminEmail string) (string, error)
	CheckAdminSession(ctx context.Context, sessionID string) (bool, error)
	AdminLogin(ctx context.Context, email string, password string) (bool, error)
	StoreAdmin(ctx context.Context, name string, email string, password string) (string, error)

	CreateDriverSession(ctx context.Context, driverID int) (string, error)
	CheckDriverSession(ctx context.Context, sessionID string) (bool, int, error)

	GetOtp(email string) string
	SetOtp(ctx context.Context, email string, otp string) error

	StoreArrivalStatus(ctx context.Context, driverID int, arrivalStatus map[int]string) error
	GetArrivalStatus(ctx context.Context, driverID int) (map[int]string, error)
	CacheBusRoute(ctx context.Context) error
	GetCachedRoute(ctx context.Context) ([]models.CurrentRoute, error)
}

type DBStore interface {
}

type Store struct {
	InMemoryDB InMemoryStore
	DB         DBStore
}

// func (s *Storage) NewStorage(redis *redis.Client, pg *pgxpool.Pool) *Storage {

// 	return &Storage{
// 		InMemoryDB: redis,
// 	}
// }
