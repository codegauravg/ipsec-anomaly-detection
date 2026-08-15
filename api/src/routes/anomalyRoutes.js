import { Router } from "express";
import { listAnomalies, getStats, markResolved } from "../controllers/anomalyController.js";

const router = Router();

router.get("/anomalies", listAnomalies);
router.get("/anomalies/stats", getStats);
router.patch("/anomalies/:id/resolve", markResolved);

export default router;