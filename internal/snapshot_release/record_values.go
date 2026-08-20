package snapshot_release

import "example.com/authoritativedns/internal/zone_domain"

func cloneRecordValues(records []zone_domain.RecordSet) {
	for i := range records {
		values := make([]zone_domain.RecordValue, len(records[i].Values))
		copy(values, records[i].Values)
		records[i].Values = values
	}
}
