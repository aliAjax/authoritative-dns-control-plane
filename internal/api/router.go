package api

import (
	"context"
	"encoding/json"
	"example.com/authoritativedns/internal/policy_engine"
	"example.com/authoritativedns/internal/resolver_runtime"
	"example.com/authoritativedns/internal/service"
	"example.com/authoritativedns/internal/zone_domain"
	"net/http"
	"strings"
)

type Router struct {
	app *service.Application
	mux *http.ServeMux
	dns *resolver_runtime.Server
}

func NewRouter(a *service.Application) *Router {
	r := &Router{app: a, mux: http.NewServeMux(), dns: resolver_runtime.NewServer(a)}
	r.routes()
	return r
}
func (r *Router) ServeHTTP(w http.ResponseWriter, q *http.Request) { r.mux.ServeHTTP(w, q) }
func (r *Router) routes() {
	r.mux.HandleFunc("/healthz", r.health)
	r.mux.HandleFunc("/readyz", r.health)
	r.mux.HandleFunc("/dns-query", r.dns.DoH)
	r.mux.HandleFunc("/api/v1/zones", r.zones)
	r.mux.HandleFunc("/api/v1/zones/", r.zoneSubresource)
	r.mux.HandleFunc("/api/v1/policies", r.policies)
	r.mux.HandleFunc("/api/v1/anycast-nodes", r.nodes)
	r.mux.HandleFunc("/api/v1/routing-intents/reconcile", r.reconcile)
}
func (r *Router) health(w http.ResponseWriter, q *http.Request) {
	write(w, 200, map[string]string{"status": "ok"})
}
func (r *Router) zones(w http.ResponseWriter, q *http.Request) {
	if q.Method == "GET" {
		write(w, 200, r.app.ListZones(contextFor(q)))
		return
	}
	if q.Method != "POST" {
		writeErr(w, 405, "method not allowed")
		return
	}
	var in struct {
		Name string `json:"name"`
	}
	if json.NewDecoder(q.Body).Decode(&in) != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	z, e := r.app.CreateZone(contextFor(q), in.Name)
	if e != nil {
		writeErr(w, 400, e.Error())
		return
	}
	write(w, 201, z)
}
func (r *Router) zoneSubresource(w http.ResponseWriter, q *http.Request) {
	parts := strings.Split(strings.Trim(q.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeErr(w, 404, "not found")
		return
	}
	id := parts[3]
	if len(parts) == 5 && parts[4] == "record-sets" {
		if q.Method == "GET" {
			write(w, 200, r.app.ListRecords(contextFor(q), id))
			return
		}
		var in zone_domain.RecordSet
		if json.NewDecoder(q.Body).Decode(&in) != nil {
			writeErr(w, 400, "invalid json")
			return
		}
		v, e := r.app.AddRecord(contextFor(q), id, in)
		if e != nil {
			writeErr(w, 400, e.Error())
			return
		}
		write(w, 201, v)
		return
	}
	if len(parts) == 5 && parts[4] == "snapshots" {
		write(w, 200, r.app.Store.ListSnapshots(id))
		return
	}
	if len(parts) == 5 && parts[4] == "validate" {
		v, e := r.app.Validate(contextFor(q), id)
		if e != nil {
			writeErr(w, 404, e.Error())
			return
		}
		write(w, 200, v)
		return
	}
	if len(parts) == 5 && parts[4] == "publish" {
		if q.Method != "POST" {
			writeErr(w, 405, "post required")
			return
		}
		v, e := r.app.Publish(contextFor(q), id)
		if e != nil {
			writeErr(w, 409, e.Error())
			return
		}
		write(w, 200, v)
		return
	}
	if len(parts) == 5 && parts[4] == "rollback" {
		var in struct {
			SnapshotID string `json:"snapshot_id"`
		}
		json.NewDecoder(q.Body).Decode(&in)
		v, e := r.app.Rollback(contextFor(q), id, in.SnapshotID)
		if e != nil {
			writeErr(w, 409, e.Error())
			return
		}
		write(w, 200, v)
		return
	}
	writeErr(w, 404, "not found")
}
func (r *Router) policies(w http.ResponseWriter, q *http.Request) {
	if q.Method != "POST" {
		write(w, 200, r.app.Policies())
		return
	}
	var p policyInput
	if json.NewDecoder(q.Body).Decode(&p) != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	r.app.AddPolicy(p.Policy)
	write(w, 201, p.Policy)
}
func (r *Router) nodes(w http.ResponseWriter, q *http.Request) {
	if q.Method == "GET" {
		write(w, 200, []any{})
		return
	}
	write(w, 202, map[string]string{"status": "accepted"})
}
func (r *Router) reconcile(w http.ResponseWriter, q *http.Request) {
	write(w, 200, map[string]string{"status": "reconciled"})
}

type policyInput struct {
	Policy policy_engine.Policy `json:"policy"`
}

func contextFor(*http.Request) context.Context { return context.Background() }
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	write(w, status, map[string]string{"error": msg})
}
