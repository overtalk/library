package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultConnctTimeOut   = 10 * time.Second
	defaultPingWaitingTime = 2 * time.Second
)

// Cfg config for mongodb
type Cfg struct {
	Size          uint64
	Addr          []string
	Username      string
	Password      string
	DBName        string
	ReplicaSet    string
	AuthMechanism string
}

// Connect connect mongo db
func (m *Cfg) Connect() (*mongo.Database, error) {
	if m.Size < 1 {
		return nil, errors.New("invalid mysql pool size")
	}
	option := options.Client()
	option.SetHosts(m.Addr)
	option.SetAuth(options.Credential{
		PasswordSet:   true,
		Username:      m.Username,
		Password:      m.Password,
		AuthMechanism: m.AuthMechanism,
	})
	option.SetMaxPoolSize(m.Size)
	option.SetReplicaSet(m.ReplicaSet)

	pool, err := mongo.NewClient(option)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultConnctTimeOut)
	defer cancel()

	err = pool.Connect(ctx)
	if err != nil {
		return nil, err
	}

	ctx1, cancel1 := context.WithTimeout(context.Background(), defaultPingWaitingTime)
	defer cancel1()

	if err := pool.Ping(ctx1, nil); err != nil {
		return nil, err
	}
	return pool.Database(m.DBName), nil
}
