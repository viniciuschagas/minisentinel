package minisentinel

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/matryer/is"

	"github.com/gomodule/redigo/redis"
)

func TestPing(t *testing.T) {
	is := is.New(t)

	m, err := miniredis.Run()
	is.NoErr(err)
	s, err := Run(m)
	is.NoErr(err)
	defer s.Close()
	c, err := redis.Dial("tcp", s.Addr())
	is.NoErr(err)

	// PING command
	{
		v, errPingRead := redis.String(c.Do("PING"))
		t.Logf("PING returned: %v", v)
		is.NoErr(errPingRead)
		is.True(v == "PONG")
	}
}

func TestAuth(t *testing.T) {
	is := is.New(t)

	m, err := miniredis.Run()
	is.NoErr(err)
	s, err := Run(m)
	is.NoErr(err)
	defer s.Close()

	m.RequireAuth("mypassword")
	s.RequireAuth("mypassword")

	c, err := redis.Dial("tcp", s.Addr(), redis.DialPassword("mypassword"))
	is.NoErr(err)

	defer func() { _ = c.Close() }()

	v, errPingRead := redis.String(c.Do("PING"))
	t.Logf("PING returned: %v", v)
	is.NoErr(errPingRead)
	is.True(v == "PONG")

	_, err = redis.Dial("tcp", s.Addr(), redis.DialPassword("wrongpassword"))
	if err == nil {
		t.Fatal("Expected error when connecting with wrong password, but got none")
	}
}
