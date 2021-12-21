package pg

import (
	"context"

	"github.com/Fe4p3b/go-backend-2/lesson1/internal/app/repos/catalog"
	"github.com/jackc/pgx/v4/pgxpool"
)

var _ catalog.EnvironmentStore = &Environments{}

type Environments struct {
	db *pgxpool.Pool
}

func NewEnvironments(db *pgxpool.Pool) *Environments {
	return &Environments{
		db: db,
	}
}

func (es *Environments) Create(ctx context.Context, e catalog.Environment) (*catalog.Environment, error)
func (es *Environments) Read(ctx context.Context, id int) (*catalog.Environment, error)
func (es *Environments) Delete(ctx context.Context, id int) error
func (es *Environments) DeleteUser(ctx context.Context, userId int) error
func (es *Environments) AddUser(ctx context.Context, u catalog.User) (*catalog.Environment, error)
func (es *Environments) SearchByName(ctx context.Context, name string) ([]catalog.Environment, error)
func (es *Environments) SearchByUsernames(ctx context.Context, usernames []string) ([]catalog.Environment, error)
