package repository

import "example.com/authoritativedns/internal/zone_domain"

func cloneRecordSet(in zone_domain.RecordSet) zone_domain.RecordSet {
	out := in
	if in.Values != nil {
		out.Values = make([]zone_domain.RecordValue, len(in.Values))
		copy(out.Values, in.Values)
	}
	return out
}

func cloneZone(in zone_domain.Zone) zone_domain.Zone {
	out := in
	out.NS = append([]string(nil), in.NS...)
	return out
}

func cloneSnapshot(in zone_domain.Snapshot) zone_domain.Snapshot {
	out := in
	out.Records = make([]zone_domain.RecordSet, 0, len(in.Records))
	for _, record := range in.Records {
		out.Records = append(out.Records, cloneRecordSet(record))
	}
	return out
}
