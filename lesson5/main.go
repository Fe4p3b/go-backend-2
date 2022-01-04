package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/Fe4p3b/go-backend-2/lesson5/sharding"
	_ "github.com/lib/pq"
)

type Pool struct {
	sync.RWMutex

	cc map[string]*sql.DB
}

func NewPool() *Pool {
	return &Pool{
		cc: map[string]*sql.DB{},
	}
}

func (p *Pool) Connection(addr string) (*sql.DB, error) {
	p.RLock()
	if c, ok := p.cc[addr]; ok {
		defer p.RUnlock()
		return c, nil
	}
	p.RUnlock()

	p.Lock()
	defer p.Unlock()
	if c, ok := p.cc[addr]; ok {
		return c, nil
	}
	var err error

	p.cc[addr], err = sql.Open("postgres", addr)
	return p.cc[addr], err
}

type User struct {
	UserId int
	Name   string
	Age    int
	Spouse int
}

func (u *User) shardConnection() (*sql.DB, error) {
	s, err := m.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (u *User) replicaConnection() (*sql.DB, error) {
	s, err := m.ReplicaById(u.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (u *User) Create() error {
	c, err := u.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "users" VALUES ($1, $2, $3, $4)`, u.UserId, u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Read() error {
	c, err := u.replicaConnection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "name", "age", "spouse" FROM "users" WHERE "user_id" = $1`, u.UserId)
	return r.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}

func (u *User) Update() error {
	c, err := u.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = $4 WHERE "user_id" = $1`, u.UserId,
		u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Delete() error {
	c, err := u.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "users" WHERE "user_id" = $1`, u.UserId)
	return err
}

var (
	m = sharding.NewManager(10)
	p = NewPool()
)

func main() {
	m.AddShard(&sharding.Shard{"port=8100 user=test password=test dbname=test sslmode=disable", 0})
	m.AddShard(&sharding.Shard{"port=8110 user=test password=test dbname=test sslmode=disable", 1})
	m.AddShard(&sharding.Shard{"port=8120 user=test password=test dbname=test sslmode=disable", 2})

	m.AddReplica(&sharding.Replica{"port=8101 user=test password=test dbname=test sslmode=disable", 0})
	m.AddReplica(&sharding.Replica{"port=8111 user=test password=test dbname=test sslmode=disable", 1})
	m.AddReplica(&sharding.Replica{"port=8121 user=test password=test dbname=test sslmode=disable", 2})

	uu := []*User{
		{1, "Joe Biden", 78, 10},
		{10, "Jill Biden", 69, 1},
		{13, "Donald Trump", 74, 25},
		{25, "Melania Trump", 78, 13},
	}
	for _, u := range uu {
		err := u.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
	}

	uu = []*User{
		{UserId: 1},
		{UserId: 10},
		{UserId: 13},
		{UserId: 25},
	}
	for _, u := range uu {
		err := u.Read()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", u, err))
		}
		fmt.Printf("user: %v\n\n", u)
	}

	aa := []*Activities{
		{1, time.Now(), "Jumping jacks"},
		{10, time.Now(), "Break dance"},
		{13, time.Now(), "Freestyle rap"},
		{25, time.Now(), "Beatbox"},
	}

	for _, a := range aa {
		err := a.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create activity %v: %w", a, err))
		}
	}

	for _, a := range aa {
		err := a.Read()
		if err != nil {
			fmt.Println(fmt.Errorf("error on read activity %v: %w", a, err))
		}
		fmt.Printf("activity: %v\n\n", a)
	}
}

type Activities struct {
	UserId int
	Date   time.Time
	Name   string
}

func (a *Activities) shardConnection() (*sql.DB, error) {
	s, err := m.ShardById(a.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (a *Activities) replicaConnection() (*sql.DB, error) {
	s, err := m.ReplicaById(a.UserId)
	if err != nil {
		return nil, err
	}
	return p.Connection(s.Address)
}

func (a *Activities) Create() error {
	c, err := a.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "activities" VALUES ($1, $2, $3)`, a.UserId, a.Date, a.Name)
	return err
}

func (a *Activities) Read() error {
	c, err := a.replicaConnection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "date", "name" FROM "activities" WHERE "user_id" = $1`, a.UserId)
	return r.Scan(
		&a.Date,
		&a.Name,
	)
}

func (a *Activities) Update() error {
	c, err := a.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "activities" SET "date" = $2, "name" = $3 WHERE "user_id" = $1`, a.Name, a.Date)
	return err
}

func (a *Activities) Delete() error {
	c, err := a.shardConnection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "activities" WHERE "user_id" = $1`, a.UserId)
	return err
}
