package minisentinel

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/matryer/is"

	"github.com/gomodule/redigo/redis"
)

func TestNewSentinel(t *testing.T) {
	is := is.New(t)
	m, err := miniredis.Run()
	is.NoErr(err)
	defer m.Close()

	s := NewSentinel(m, WithReplica(m))
	is.Equal(s.Master(), m)  // make sure the master is correctly set
	is.Equal(s.Replica(), m) // make sure the replicas are correctly set

	err = s.Start()
	is.NoErr(err)
	defer s.Close()

	c, err := redis.Dial("tcp", s.Addr())
	is.NoErr(err)

	// just running all the commands... not really testing expected results (see specifc unit tests for that)
	// PING command
	{
		v, err := redis.String(c.Do("PING"))
		is.NoErr(err)
		t.Log(v)
	}
	// MASTERS command
	{
		// results is an []interfaces which points to [][]strings
		results, err := c.Do("SENTINEL", "MASTERS")
		is.NoErr(err)
		info, err := redis.Strings(results.([]any)[0], nil)
		is.NoErr(err)
		t.Log("MASTERS response:")
		t.Logf("%v", info)
		for _, v := range info {
			t.Logf("%v", v)
		}
	}
	// GET-MASTER-ADDR-BY-NAME command
	{
		results, err := redis.Strings(c.Do("SENTINEL", "GET-MASTER-ADDR-BY-NAME", "MYMASTER"))
		is.NoErr(err)
		t.Log(results)
	}

	// SLAVES command
	{
		// results is an []interfaces which points to [][]strings
		results, err := c.Do("SENTINEL", "SLAVES")
		is.NoErr(err)
		info, err := redis.Strings(results.([]any)[0], nil)
		is.NoErr(err)
		t.Log("SLAVES response:")
		t.Logf("%v", info)
		for _, v := range info {
			t.Logf("%v", v)
		}
	}

	// SENTINELS command
	{
		// results is an []interfaces which points to [][]strings
		results, err := c.Do("SENTINEL", "SENTINELS", "MYMASTER")
		is.NoErr(err)
		info, err := redis.Strings(results.([]any), nil)
		is.NoErr(err)
		t.Log("SENTINELS response:")
		t.Logf("%v", info)
		for _, v := range info {
			t.Logf("%v", v)
		}
	}

	// ROLE command
	{
		results, err := c.Do("ROLE")
		is.NoErr(err)
		resultsArr := results.([]any)
		t.Log("ROLE response:")
		role, err := redis.String(resultsArr[0], nil)
		is.NoErr(err)
		t.Log("role:", role)
		masterName, err := redis.String(resultsArr[1].([]any)[0], nil)
		is.NoErr(err)
		t.Log("masterName:", masterName)
	}
}
