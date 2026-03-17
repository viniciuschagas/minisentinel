package minisentinel

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/matryer/is"
)

func TestInitMasterInfo(t *testing.T) {
	is := is.New(t)
	m, err := miniredis.Run()
	is.NoErr(err)
	defer m.Close()

	s, err := Run(m)
	is.NoErr(err)
	defer s.Close()

	mInfo := initMasterInfo(s)

	is.Equal(mInfo.Name, "mymaster")
	is.Equal(mInfo.IP, m.Host())
	is.Equal(mInfo.Port, m.Port())
	is.True(mInfo.RunID != "")
	is.Equal(mInfo.Flags, "master")
	is.Equal(mInfo.LinkPendingCommands, "0")
	is.Equal(mInfo.LinkRefCount, "1")
	is.Equal(mInfo.LastPingSent, "0")
	is.Equal(mInfo.LastOkPingReply, "0")
	is.Equal(mInfo.LastPingReply, "0")
	is.Equal(mInfo.DownAfterMilliseconds, "5000")
	is.Equal(mInfo.InfoRefresh, "6295")
	is.Equal(mInfo.RoleReported, "master")
	is.True(mInfo.RoleReportedTime != "")
	is.Equal(mInfo.ConfigEpoch, "1")
	is.Equal(mInfo.NumSlaves, "1")
	is.Equal(mInfo.NumOtherSentinels, "0")
	is.Equal(mInfo.Quorum, "1")
	is.Equal(mInfo.FailoverTimeout, "60000")
	is.Equal(mInfo.ParallelSync, "1")
}

func TestNewMasterInfoFromStrings(t *testing.T) {
	invalidsInfo := []string{"name", "mymaster", "ip"}
	_, err := NewMasterInfoFromStrings(invalidsInfo)

	if err == nil {
		t.Fatal("Expected error when input is not a modulus of 2, but got none")
	}

	validInfo := []string{
		"name", "mymaster",
		"ip", "127.0.0.1",
		"port", "6379",
	}

	mi, err := NewMasterInfoFromStrings(validInfo)
	if err != nil {
		t.Fatal("Expected no error when input is valid, but got:", err)
	}

	is := is.New(t)
	is.Equal(mi.Name, "mymaster")
	is.Equal(mi.IP, "127.0.0.1")
	is.Equal(mi.Port, "6379")
}
