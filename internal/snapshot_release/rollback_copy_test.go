package snapshot_release

import (
	"example.com/authoritativedns/internal/zone_domain"
	"testing"
)

func TestRollbackCopiesRecordValues(t *testing.T) {
	source := []zone_domain.RecordSet{{Values: []zone_domain.RecordValue{{Value: "192.0.2.3"}}}}
	rolled := rollbackRecords(source)
	source[0].Values[0].Value = "203.0.113.3"
	if rolled[0].Values[0].Value != "192.0.2.3" {
		t.Fatalf("rollback record aliases source values: %#v", rolled)
	}
}
