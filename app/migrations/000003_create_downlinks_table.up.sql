CREATE TABLE downlinks (
    id          TEXT PRIMARY KEY,           -- dl-<uuid>
    device_id   TEXT NOT NULL REFERENCES devices(id),  -- which device
    command     TEXT NOT NULL,              -- e.g. "reset", "calibrate", "set_interval"
    payload     JSONB,                      -- optional structured params (nullable)
    status      TEXT NOT NULL,              -- "queued" | "sent" | "failed"
    created_at  TIMESTAMPTZ NOT NULL,
    confirmed_by TEXT                       -- who authorized it (for HIL record)
);
