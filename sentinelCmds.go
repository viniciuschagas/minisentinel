package minisentinel

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/alicebob/miniredis/v2/server"
)

const msgInvalidSentinelCommand = "ERR unknown command '%s'"

func commandsSentinel(s *Sentinel) {
	_ = s.srv.Register("SENTINEL", s.cmdsSentinel)
	s.setupRoleCommand()
}

func (s *Sentinel) setupRoleCommand() {
	_ = s.srv.Register("ROLE", s.roleCommand)

	if s.master != nil {
		_ = s.master.Server().Register("ROLE", s.masterRoleCommand)
	}

	if s.replica != nil && (s.master != nil && s.replica != s.master) {
		_ = s.replica.Server().Register("ROLE", s.replicaRoleCommand)
	}
}

// cmdsSentinel - entry point for all commands that start with SENTINEL
func (s *Sentinel) cmdsSentinel(c *server.Peer, cmd string, args []string) {
	if !isSentinelCmd(cmd) {
		c.WriteError(fmt.Sprintf(msgInvalidSentinelCommand, cmd))
		return
	}
	if len(args) > 2 {
		c.WriteError(errWrongNumber(cmd))
		return
	}
	if !s.handleAuth(c) {
		return
	}

	subCmd := strings.ToUpper(args[0])
	switch subCmd {
	case "SENTINELS":
		err := s.sentinelsCommand(c, cmd, args)
		if err != nil {
			c.WriteError(err.Error())
		}
		return
	case "MASTERS":
		err := s.mastersCommand(c, cmd, args)
		if err != nil {
			c.WriteError(err.Error())
		}
		return
	case "GET-MASTER-ADDR-BY-NAME":
		err := s.getMasterAddrByNameCommand(c, cmd, args)
		if err != nil {
			c.WriteError(err.Error())
		}
		return
	case "SLAVES":
		err := s.slavesCommand(c, cmd, args)
		if err != nil {
			c.WriteError(err.Error())
		}
		return
	default:
		c.WriteError(fmt.Sprintf(msgInvalidSentinelCommand, subCmd))
		return
	}
}

func (s *Sentinel) getMasterAddrByNameCommand(c *server.Peer, cmd string, args []string) error {
	if !isSentinelCmd(cmd) {
		return fmt.Errorf(msgInvalidSentinelCommand, cmd)
	}
	subCmd := strings.ToUpper(args[0])
	if subCmd != "GET-MASTER-ADDR-BY-NAME" {
		return fmt.Errorf(msgInvalidSentinelCommand, subCmd)
	}
	if !strings.EqualFold(s.masterInfo.Name, args[1]) {
		c.WriteLen(-1)
		return nil
	}
	c.WriteLen(2)
	c.WriteBulk(s.master.Host())
	c.WriteBulk(s.master.Port())
	return nil
}

func (s *Sentinel) slavesCommand(c *server.Peer, cmd string, args []string) error {
	if !isSentinelCmd(cmd) {
		return fmt.Errorf(msgInvalidSentinelCommand, cmd)
	}
	subCmd := strings.ToUpper(args[0])
	if subCmd != "SLAVES" {
		return fmt.Errorf(msgInvalidSentinelCommand, subCmd)
	}
	c.WriteLen(1)
	c.WriteLen(40)
	t := reflect.TypeFor[ReplicaInfo]()
	v := reflect.ValueOf(s.replicaInfo)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		c.WriteBulk(tag)
		c.WriteBulk(v.Field(i).Interface().(string))
	}

	return nil
}

func (s *Sentinel) mastersCommand(c *server.Peer, cmd string, args []string) error {
	if !isSentinelCmd(cmd) {
		return fmt.Errorf(msgInvalidSentinelCommand, cmd)
	}
	subCmd := strings.ToUpper(args[0])
	if subCmd != "MASTERS" {
		return fmt.Errorf(msgInvalidSentinelCommand, subCmd)
	}
	c.WriteLen(1)
	c.WriteLen(40)
	t := reflect.TypeFor[MasterInfo]()
	v := reflect.ValueOf(s.masterInfo)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		c.WriteBulk(tag)
		c.WriteBulk(v.Field(i).Interface().(string))
	}

	return nil
}

func (s *Sentinel) sentinelsCommand(c *server.Peer, cmd string, args []string) error {
	if !isSentinelCmd(cmd) {
		return fmt.Errorf(msgInvalidSentinelCommand, cmd)
	}
	subCmd := strings.ToUpper(args[0])
	if subCmd != "SENTINELS" {
		return fmt.Errorf(msgInvalidSentinelCommand, subCmd)
	}

	// For now, we don't support sentinels, so we just return an empty array
	c.WriteLen(0)

	return nil
}

func (s *Sentinel) roleCommand(c *server.Peer, cmd string, args []string) {
	c.WriteLen(2)
	c.WriteBulk("sentinel")
	c.WriteLen(1)
	c.WriteBulk(s.masterInfo.Name)
}

func (s *Sentinel) masterRoleCommand(c *server.Peer, cmd string, args []string) {
	writeLen := 2
	if s.replica != nil {
		writeLen = 3
	}
	c.WriteLen(writeLen)
	c.WriteBulk("master")
	c.WriteInt(0)

	if s.replica != nil {
		c.WriteLen(3)
		c.WriteBulk(s.replica.Host())
		c.WriteBulk(s.replica.Port())
		c.WriteInt(0)
	}
}

func (s *Sentinel) replicaRoleCommand(c *server.Peer, cmd string, args []string) {
	portInt, err := strconv.Atoi(s.master.Port())
	if err != nil {
		c.WriteError(fmt.Sprintf("ERR invalid port: %s", s.master.Port()))
		return
	}

	c.WriteLen(5)
	c.WriteBulk("slave")
	c.WriteBulk(s.master.Host())
	c.WriteInt(portInt)
	c.WriteBulk("connected")
	c.WriteInt(0)
}

func isSentinelCmd(cmd string) bool {
	return strings.ToUpper(cmd) == "SENTINEL"
}
