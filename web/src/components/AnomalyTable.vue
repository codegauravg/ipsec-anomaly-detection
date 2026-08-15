<template>
  <div class="table-container">
    <table>
      <thead>
        <tr>
          <th>Timestamp</th>
          <th>Source IP</th>
          <th>Threat Type</th>
          <th>Severity</th>
          <th>Status</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="anomalies.length === 0">
          <td colspan="6" style="text-align: center;">No anomalies detected.</td>
        </tr>
        <tr v-for="anomaly in anomalies" :key="anomaly.id">
          <td>{{ new Date(anomaly.timestamp).toLocaleString() }}</td>
          <td><strong>{{ anomaly.source_ip }}</strong></td>
          <td>{{ anomaly.threat_type }}</td>
          <td>
            <span :class="['badge', anomaly.severity > 5 ? 'severity-high' : 'severity-medium']">
              {{ anomaly.severity }}
            </span>
          </td>
          <td>{{ anomaly.resolved ? 'Resolved' : 'Active' }}</td>
          <td>
            <button 
              v-if="!anomaly.resolved" 
              @click="$emit('resolve', anomaly.id)" 
              class="btn-resolve"
            >
              Mark Resolved
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { Anomaly } from '@/composables/useAnomalies';

defineProps<{
  anomalies: Anomaly[]
}>();

defineEmits<{
  (e: 'resolve', id: string): void
}>();
</script>