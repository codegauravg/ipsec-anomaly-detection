import express from "express";
import dotenv from "dotenv";
import { securityHeaders, corsConfig } from "./middlewares/security.js";
import { apiLimiter } from "./middlewares/rateLimiter.js";
import { errorHandler } from "./middlewares/errorHandler.js";
import anomalyRoutes from "./routes/anomalyRoutes.js";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// Security & Parsing Middlewares
app.use(securityHeaders);
app.use(corsConfig);
app.use(express.json());

// Apply Rate Limiting to all /api/ endpoints
app.use("/api", apiLimiter);

// Health check endpoint for Kubernetes probes (liveness & readiness)
app.get("/healthz", (req, res) => {
  res.status(200).json({ status: "healthy", timestamp: new Date().toISOString() });
});

// Domain Routes
app.use("/api", anomalyRoutes);

// Global Centralized Error Handler
app.use(errorHandler);

const server = app.listen(PORT, () => {
  console.log(`[BFF] IPSec Anomaly API Server running on port ${PORT}`);
});

// Graceful Shutdown for Kubernetes SIGTERM/SIGINT
process.on("SIGTERM", () => {
  console.log("SIGTERM signal received. Closing HTTP server...");
  server.close(() => {
    console.log("HTTP server closed.");
    process.exit(0);
  });
});