package repository

import (
	"testing"

	"example.com/authoritativedns/internal/zone_domain"
)

func TestPutRecordOwnsNestedValues(t *testing.T) {
	s := NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	s.PutZone(z)
	values := []zone_domain.RecordValue{{Value: "192.0.2.10"}}
	start := make(chan struct{})
	putDone := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() {
		<-start
		s.PutRecord(z.ID, zone_domain.RecordSet{ID: "www", Name: "www.example.com.", Type: zone_domain.A, TTL: 60, Values: values})
		close(putDone)
		done <- struct{}{}
	}()
	go func() {
		<-start
		<-putDone
		values[0].Value = "198.51.100.10"
		done <- struct{}{}
	}()
	close(start)
	<-done
	<-done
	got := s.ListRecords(z.ID)
	if len(got) != 1 || got[0].Values[0].Value != "192.0.2.10" {
		t.Fatalf("stored record changed through caller buffer: %#v", got)
	}
	got[0].Values[0].Value = "203.0.113.55"
	again := s.ListRecords(z.ID)
	if again[0].Values[0].Value != "192.0.2.10" {
		t.Fatalf("returned record changed stored state: %#v", again)
	}
}

func TestGetZoneOwnsNameservers(t *testing.T) {
	s := NewMemoryStore()
	z := zone_domain.NewZone("example.com.")
	z.NS = []string{"ns1.example.com."}
	start := make(chan struct{})
	putDone := make(chan struct{})
	done := make(chan struct{}, 2)
	go func() {
		<-start
		s.PutZone(z)
		close(putDone)
		done <- struct{}{}
	}()
	go func() {
		<-start
		<-putDone
		z.NS[0] = "poisoned.example.com."
		done <- struct{}{}
	}()
	close(start)
	<-done
	<-done
	got, ok := s.GetZone(z.ID)
	if !ok || got.NS[0] != "ns1.example.com." {
		t.Fatalf("stored zone changed through caller slice: %#v", got)
	}
	got.NS[0] = "returned.example.com."
	again, ok := s.GetZone(z.ID)
	if !ok || again.NS[0] != "ns1.example.com." {
		t.Fatalf("returned zone changed stored state: %#v", again)
	}
}
