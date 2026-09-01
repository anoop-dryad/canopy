INSERT INTO devices (id, name, status, battery_pct, last_seen, created_at, updated_at) VALUES
    ('dev-001', 'North Ridge Sensor',   'online',   87, now() - interval '2 minutes',  now(), now()),
    ('dev-002', 'East Valley Sensor',   'online',   12, now() - interval '5 minutes',  now(), now()),
    ('dev-003', 'South Gate Sensor',    'offline',  0,  now() - interval '3 hours',    now(), now()),
    ('dev-004', 'West Ridge Sensor',    'degraded', 45, now() - interval '90 minutes', now(), now()),
    ('dev-005', 'Central Hub Sensor',   'online',   64, now() - interval '1 minute',   now(), now())
ON CONFLICT (id) DO NOTHING;
