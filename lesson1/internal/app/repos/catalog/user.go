package catalog

import (
	"context"
)

type User struct {
	ID           int
	Name         string
	Environments []Environment
}

type UserStore interface {
	Create(ctx context.Context, u User) (*User, error)
	Read(ctx context.Context, id int) (*User, error)
	Delete(ctx context.Context, id int) error
	SearchByName(ctx context.Context, name string) ([]User, error)
	SearchByEnvironment(ctx context.Context, environmentName string) ([]User, error)
}

type Users struct {
	ustore UserStore
}

func NewUsers(ustore UserStore) *Users {
	return &Users{
		ustore: ustore,
	}
}

func (us *Users) Create(ctx context.Context, u User) (*User, error)
func (us *Users) Read(ctx context.Context, id int) (*User, error)
func (us *Users) Delete(ctx context.Context, id int) error
func (us *Users) SearchByName(ctx context.Context, name string) ([]User, error)
func (us *Users) SearchByEnvironment(ctx context.Context, environmentName string) ([]User, error)
