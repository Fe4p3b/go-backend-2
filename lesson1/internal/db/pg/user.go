package pg

import (
	"context"

	"github.com/Fe4p3b/go-backend-2/lesson1/internal/app/repos/catalog"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ catalog.UserStore = &Users{}

type Users struct {
	db *pgxpool.Pool
}

func NewUsers(db *pgxpool.Pool) *Users {
	return &Users{
		db: db,
	}
}

func (us *Users) Create(ctx context.Context, u catalog.User) (*catalog.User, error)
func (us *Users) Read(ctx context.Context, id int) (*catalog.User, error)
func (us *Users) Delete(ctx context.Context, id int) error
func (us *Users) SearchByName(ctx context.Context, name string) ([]catalog.User, error)
func (us *Users) SearchByEnvironment(ctx context.Context, environmentName string) ([]catalog.User, error)
