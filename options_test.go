package minisentinel

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestGetDefaultOptions(t *testing.T) {
	opts := getDefaultOptions()
	if opts.masterName != "mymaster" {
		t.Errorf("Expected default master name to be 'mymaster', got '%s'", opts.masterName)
	}
	if opts.master != nil {
		t.Error("Expected default master to be nil")
	}
	if opts.replica != nil {
		t.Error("Expected default replica to be nil")
	}
}

func TestWithMasterName(t *testing.T) {
	opts := GetOpts(WithMasterName("custommaster"))
	if opts.masterName != "custommaster" {
		t.Errorf("Expected master name to be 'custommaster', got '%s'", opts.masterName)
	}
}

func TestWithMaster(t *testing.T) {
	m := &miniredis.Miniredis{}
	opts := GetOpts(WithMaster(m))
	if opts.master != m {
		t.Error("Expected master to be set correctly")
	}
}

func TestWithReplica(t *testing.T) {
	r := &miniredis.Miniredis{}
	opts := GetOpts(WithReplica(r))
	if opts.replica != r {
		t.Error("Expected replica to be set correctly")
	}
}

func TestGetOpts(t *testing.T) {
	m := &miniredis.Miniredis{}
	r := &miniredis.Miniredis{}
	opts := GetOpts(WithMasterName("custommaster"), WithMaster(m), WithReplica(r))

	if opts.masterName != "custommaster" {
		t.Errorf("Expected master name to be 'custommaster', got '%s'", opts.masterName)
	}
	if opts.master != m {
		t.Error("Expected master to be set correctly")
	}
	if opts.replica != r {
		t.Error("Expected replica to be set correctly")
	}
}
