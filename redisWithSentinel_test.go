package minisentinel

import (
	"errors"
	"testing"
	"time"

	"github.com/FZambia/sentinel/v2"
	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/matryer/is"
)

// TestSomething - shows how-to start up sentinel + redis in your unittest
func TestSomething(t *testing.T) {
	is := is.New(t)
	m := miniredis.NewMiniRedis()
	err := m.StartAddr(":6379")
	is.NoErr(err)
	defer m.Close()
	s := NewSentinel(m, WithReplica(m))
	err = s.StartAddr(":26379")
	is.NoErr(err)
	defer s.Close()
	// all the setup is done.. now just use sentinel/redis like you
	// would normally in your tests via a redis client
}

func poolDialFunc(sntnl *sentinel.Sentinel, pass []byte) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		redisHostAddr, errHostAddr := sntnl.MasterAddr()
		if errHostAddr != nil {
			return nil, errHostAddr
		}
		c, errDial := redis.Dial("tcp", redisHostAddr)
		if errDial != nil {
			return nil, errDial
		}
		if pass != nil { // auth first, before doing anything else
			if _, errDo := c.Do("AUTH", string(pass)); errDo != nil {
				_ = c.Close()
				return nil, errDo
			}
		}
		return c, nil
	}
}

// TestRedisWithSentinel - an example of combining miniRedis with miniSentinel for unittests
func TestRedisWithSentinel(t *testing.T) {
	is := is.New(t)
	m, err := miniredis.Run()
	m.RequireAuth("super-secret") // not required, but demonstrates how-to
	is.NoErr(err)
	defer m.Close()
	s := NewSentinel(m, WithReplica(m))
	err = s.Start()
	is.NoErr(err)
	defer s.Close()

	// use redigo to create a sentinel pool
	sntnl := &sentinel.Sentinel{
		Addrs:      []string{s.Addr()},
		MasterName: s.MasterInfo().Name,
		Dial: func(addr string) (redis.Conn, error) {
			connTimeout := time.Duration(50 * time.Millisecond)
			readWriteTimeout := time.Duration(50 * time.Millisecond)
			c, errDial := redis.Dial("tcp", addr, redis.DialReadTimeout(readWriteTimeout), redis.DialWriteTimeout(readWriteTimeout), redis.DialConnectTimeout(connTimeout))
			if errDial != nil {
				return nil, errDial
			}
			return c, nil
		},
	}

	redisPassword := []byte("super-secret") // required because of m.RequireAuth()

	pool := redis.Pool{
		MaxIdle:     3,
		MaxActive:   64,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial:        poolDialFunc(sntnl, redisPassword),
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			isMaster, errTestRole := sentinel.TestRole(c, "master")
			if errTestRole != nil {
				return errTestRole
			}

			if !isMaster {
				return errors.New("Role check failed")
			}
			return nil
		},
	}

	// Optionally set some keys your code expects:
	err = m.Set("foo", "not-bar")
	is.NoErr(err)
	m.HSet("some", "other", "key")

	// Run your code and see if it behaves.
	// An example using the redigo library from "github.com/gomodule/redigo/redis":
	c := pool.Get()
	defer func() {
		_ = c.Close() //release connection back to the pool
	}()

	// use the pool connection to do things via TCP to our miniRedis
	_, err = c.Do("SET", "foo", "bar")
	is.NoErr(err)

	// Optionally check values in redis...
	got, err := m.Get("foo")
	is.NoErr(err)
	is.Equal(got, "bar")

	// ... or use a miniRedis helper for that:
	m.CheckGet(t, "foo", "bar")

	// TTL and expiration:
	err = m.Set("foo", "bar")
	is.NoErr(err)
	m.SetTTL("foo", 10*time.Second)
	m.FastForward(11 * time.Second)
	is.True(m.Exists("foo") == false) // it shouldn't be there anymore
}
