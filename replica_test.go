package minisentinel

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/matryer/is"
)

func TestInitReplicaInfo(t *testing.T) {
	is := is.New(t)

	m, err := miniredis.Run()
	is.NoErr(err)
	defer m.Close()

	s, err := Run(m)
	is.NoErr(err)
	defer s.Close()

	replicaInfo := initReplicaInfo(s)
	is.Equal(replicaInfo.Name, "mymaster")
	is.Equal(replicaInfo.Port, m.Port())
	is.Equal(replicaInfo.IP, m.Host())
	is.True(replicaInfo.RunID != "")
	is.Equal(replicaInfo.Flags, "master")
	is.Equal(replicaInfo.LinkPendingCommands, "0")
	is.Equal(replicaInfo.LinkRefCount, "1")
	is.Equal(replicaInfo.LastPingSent, "0")
	is.Equal(replicaInfo.LastOkPingReply, "0")
	is.Equal(replicaInfo.LastPingReply, "0")
	is.Equal(replicaInfo.DownAfterMilliseconds, "5000")
	is.Equal(replicaInfo.InfoRefresh, "6295")
	is.Equal(replicaInfo.RoleReported, "replica")
	is.True(replicaInfo.RoleReportedTime != "")
	is.Equal(replicaInfo.MasterLinkDownTime, "0")
	is.Equal(replicaInfo.MasterLinkStatus, "ok")
	is.Equal(replicaInfo.MasterHost, m.Host())
	is.Equal(replicaInfo.MasterPort, m.Port())
	is.Equal(replicaInfo.SlavePriority, "100")
	is.Equal(replicaInfo.SlaveReplOffset, "1")
}

func TestNewReplicaInfoFromStrings(t *testing.T) {
	is := is.New(t)
	invalidsInfo := []string{"name", "mymaster", "ip"}
	_, err := NewReplicaInfoFromStrings(invalidsInfo)

	if err == nil {
		t.Fatal("Expected error when input is not a modulus of 2, but got none")
	}

	validInfo := []string{
		"name", "mymaster",
		"ip", "127.0.0.1",
		"port", "6379",
	}

	ri, err := NewReplicaInfoFromStrings(validInfo)
	if err != nil {
		t.Fatal("Expected no error when input is valid, but got:", err)
	}

	is.Equal(ri.Name, "mymaster")
	is.Equal(ri.IP, "127.0.0.1")
	is.Equal(ri.Port, "6379")
}
