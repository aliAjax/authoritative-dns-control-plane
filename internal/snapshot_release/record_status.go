package snapshot_release

import "example.com/authoritativedns/internal/zone_domain"

func applyRollbackStatus(records []zone_domain.RecordSet) {
	for i := range records {
		records[i].Status = zone_domain.RolledBack
	}
}
