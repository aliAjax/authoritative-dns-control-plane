package snapshot_release

import "example.com/authoritativedns/internal/zone_domain"

func cloneRecordSets(records []zone_domain.RecordSet) []zone_domain.RecordSet {
	cloned := make([]zone_domain.RecordSet, len(records))
	copy(cloned, records)
	return cloned
}
