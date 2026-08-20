CREATE TABLE zones (id text PRIMARY KEY, name text NOT NULL UNIQUE, serial bigint NOT NULL, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE TABLE record_sets (id text PRIMARY KEY, zone_id text NOT NULL REFERENCES zones(id), name text NOT NULL, type text NOT NULL, ttl integer NOT NULL, values_json jsonb NOT NULL, status text NOT NULL, version bigint NOT NULL);
CREATE TABLE snapshots (id text PRIMARY KEY, zone_id text NOT NULL REFERENCES zones(id), serial bigint NOT NULL, digest text NOT NULL, records_json jsonb NOT NULL, created_at timestamptz NOT NULL);
CREATE TABLE audit_events (id text PRIMARY KEY, resource text NOT NULL, action text NOT NULL, payload jsonb NOT NULL, previous_hash text, hash text NOT NULL, created_at timestamptz NOT NULL);
CREATE INDEX record_sets_zone_name ON record_sets(zone_id,name,type);

