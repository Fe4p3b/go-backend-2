package catalog

import "context"

type Environment struct {
	ID          int
	Name        string
	Description string
	Users       []User
}

type EnvironmentStore interface {
	Create(ctx context.Context, e Environment) (*Environment, error)
	Read(ctx context.Context, id int) (*Environment, error)
	Delete(ctx context.Context, id int) error
	DeleteUser(ctx context.Context, userId int) error
	AddUser(ctx context.Context, u User) (*Environment, error)
	SearchByName(ctx context.Context, name string) ([]Environment, error)
	SearchByUsernames(ctx context.Context, usernames []string) ([]Environment, error)
}

type Environments struct {
	estore EnvironmentStore
}

func NewEnvironments(estore EnvironmentStore) *Environments {
	return &Environments{
		estore: estore,
	}
}

func (es *Environments) Create(ctx context.Context, e Environment) (*Environment, error)
func (es *Environments) Read(ctx context.Context, id int) (*Environment, error)
func (es *Environments) Delete(ctx context.Context, id int) (*Environment, error)
func (es *Environments) DeleteUser(ctx context.Context, userId int) error
func (es *Environments) SearchByName(ctx context.Context, name int) ([]Environment, error)
func (es *Environments) SearchByUsernames(ctx context.Context, usernames []string) ([]Environment, error)
