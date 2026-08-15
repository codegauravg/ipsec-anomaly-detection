CREATE TABLE ipsec_anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ DEFAULT NOW(),
    source_ip VARCHAR(45) NOT NULL,
    threat_type VARCHAR(50) NOT NULL,
    severity INT NOT NULL,
    raw_payload JSONB NOT NULL,
    resolved BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_anomalies_timestamp ON ipsec_anomalies(timestamp);
CREATE INDEX idx_anomalies_source_ip ON ipsec_anomalies(source_ip);