import { pool } from "./db.js";

export const getAnomalies = async ({ limit = 50, offset = 0, sourceIp, resolved }) => {
  const values = [];
  let query = `
    SELECT 
      id,
      timestamp,
      source_ip,
      threat_type,
      severity,
      raw_payload,
      resolved
    FROM ipsec_anomalies
    WHERE 1=1
  `;

  if (sourceIp) {
    values.push(sourceIp);
    query += ` AND source_ip = $${values.length}`;
  }

  if (typeof resolved === "boolean") {
    values.push(resolved);
    query += ` AND resolved = $${values.length}`;
  }

  query += ` ORDER BY timestamp DESC`;

  values.push(limit);
  query += ` LIMIT $${values.length}`;

  values.push(offset);
  query += ` OFFSET $${values.length}`;

  const { rows } = await pool.query(query, values);
  return rows;
};

export const getAnomalyStats = async () => {
  const query = `
    SELECT 
      COUNT(*) AS total_alerts,
      COUNT(*) FILTER (WHERE resolved = false) AS active_threats,
      COUNT(DISTINCT source_ip) AS unique_attacker_ips,
      AVG(severity)::numeric(10,2) AS average_severity
    FROM ipsec_anomalies
    WHERE timestamp >= NOW() - INTERVAL '24 HOURS';
  `;

  const { rows } = await pool.query(query);
  return rows[0];
};

export const resolveAnomaly = async (id) => {
  const query = `
    UPDATE ipsec_anomalies
    SET resolved = true
    WHERE id = $1
    RETURNING *;
  `;
  const { rows } = await pool.query(query, [id]);
  return rows[0] || null;
};