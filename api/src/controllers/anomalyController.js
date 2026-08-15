import { z } from "zod";
import * as anomalyService from "../services/anomalyService.js";

const querySchema = z.object({
  limit: z.coerce.number().min(1).max(100).default(50),
  offset: z.coerce.number().min(0).default(0),
  sourceIp: z.string().ip().optional(),
  resolved: z.enum(["true", "false"]).transform((val) => val === "true").optional(),
});

export const listAnomalies = async (req, res, next) => {
  try {
    const validatedQuery = querySchema.parse(req.query);
    const anomalies = await anomalyService.getAnomalies(validatedQuery);
    res.json({ data: anomalies });
  } catch (error) {
    if (error instanceof z.ZodError) {
      return res.status(400).json({ error: error.errors });
    }
    next(error);
  }
};

export const getStats = async (req, res, next) => {
  try {
    const stats = await anomalyService.getAnomalyStats();
    res.json({ data: stats });
  } catch (error) {
    next(error);
  }
};

export const markResolved = async (req, res, next) => {
  try {
    const { id } = req.params;
    const updated = await anomalyService.resolveAnomaly(id);

    if (!updated) {
      return res.status(404).json({ error: "Anomaly incident not found." });
    }

    res.json({ data: updated });
  } catch (error) {
    next(error);
  }
};