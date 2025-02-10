CREATE TABLE telemetry_metadata (
                                    id SERIAL PRIMARY KEY,
                                    satellite_id TEXT NOT NULL,
                                    version TEXT NOT NULL,
                                    uploaded_at TIMESTAMP DEFAULT NOW(),
                                    UNIQUE (satellite_id, version)
);
