package snapshot_release

import "example.com/authoritativedns/internal/zone_domain"

func rollbackRecords(records []zone_domain.RecordSet) []zone_domain.RecordSet {
	rolledBack := cloneRecordSets(records)
	applyRollbackStatus(rolledBack)
	cloneRecordValues(rolledBack)
	return rolledBack
}
