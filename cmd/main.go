package main

import (
	"context"
	"fmt"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/rueidis"
	"github.com/viniciuschagas/minisentinel"
)

func main() {
	master, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("------ master addr: ", master.Addr())
	defer master.Close()

	replica, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("------ replica addr: ", replica.Addr())
	defer replica.Close()

	s := minisentinel.NewSentinel(master, minisentinel.WithReplica(replica))
	err = s.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("---- sentinel addr: ", s.Addr())
	defer s.Close()

	pass := "super-secret"
	master.RequireAuth(pass)
	replica.RequireAuth(pass)
	s.RequireAuth(pass)

	addrs := []string{master.Addr(), replica.Addr(), s.Addr()}

	clientInstance, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: addrs,
		Sentinel: rueidis.SentinelOption{
			MasterSet: s.MasterInfo().Name,
			Password:  pass,
		},
		Password:              pass,
		ClientName:            "ClientName",
		DisableCache:          true,
		DisableAutoPipelining: true,
	})

	if err != nil {
		panic(err)
	}

	// just run a command to make sure it works
	cmd := clientInstance.B().Ping().Build()
	resp := clientInstance.Do(context.TODO(), cmd)
	if resp.Error() != nil {
		panic(resp.Error())
	}
	fmt.Println("PING response: ", resp.String())

	cmd = clientInstance.B().Role().Build()
	resp = clientInstance.Do(context.TODO(), cmd)
	if resp.Error() != nil {
		panic(resp.Error())
	}
	fmt.Println("ROLE response: ", resp.String())
}
