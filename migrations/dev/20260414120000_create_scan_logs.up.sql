CREATE TABLE IF NOT EXISTS scan_logs (
    id              SERIAL PRIMARY KEY,
    active_pond_id  INTEGER NOT NULL REFERENCES active_ponds(id),
    month           VARCHAR(7) NOT NULL,
    image_paths     JSONB NOT NULL DEFAULT '[]',
    raw_response    TEXT,
    extracted_data  JSONB,
    confidence_scores JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending_review',
    reviewed_by     VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      VARCHAR(255),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by      VARCHAR(255),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_scan_logs_active_pond_id ON scan_logs(active_pond_id);
CREATE INDEX idx_scan_logs_deleted_at ON scan_logs(deleted_at);
