package sharding

import (
	"errors"
	"sync"
)

type Replica struct {
	Address string
	Number  int
}

type Shard struct {
	Address string
	Number  int
}

type Manager struct {
	size int
	ss   *sync.Map
	sr   *sync.Map
}

var (
	ErrorShardNotFound = errors.New("shard not found")
)

func NewManager(size int) *Manager {
	return &Manager{
		size: size,
		ss:   &sync.Map{},
		sr:   &sync.Map{},
	}
}

func (m *Manager) AddShard(s *Shard) {
	m.ss.Store(s.Number, s)
}

func (m *Manager) AddReplica(s *Replica) {
	m.sr.Store(s.Number, s)
}

func (m *Manager) ShardById(entityId int) (*Shard, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId / m.size
	if s, ok := m.ss.Load(n); ok {
		return s.(*Shard), nil
	}
	return nil, ErrorShardNotFound
}

func (m *Manager) ReplicaById(entityId int) (*Replica, error) {
	if entityId < 0 {
		return nil, ErrorShardNotFound
	}
	n := entityId / m.size
	if s, ok := m.sr.Load(n); ok {
		return s.(*Replica), nil
	}
	return nil, ErrorShardNotFound
}
