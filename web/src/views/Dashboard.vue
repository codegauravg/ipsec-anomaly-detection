<template>
  <div class="dashboard">
    <div class="header-actions">
      <h1>Active Network Threats</h1>
      <button @click="fetchAnomalies" class="btn-refresh">Refresh Data</button>
    </div>

    <div v-if="loading" class="loading-state">Loading real-time data...</div>
    <div v-else-if="error" class="error-state">{{ error }}</div>
    
    <AnomalyTable 
      v-else 
      :anomalies="anomalies" 
      @resolve="resolveAnomaly" 
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useAnomalies } from '@/composables/useAnomalies';
import AnomalyTable from '@/components/AnomalyTable.vue';

const { anomalies, loading, error, fetchAnomalies, resolveAnomaly } = useAnomalies();

onMounted(() => {
  fetchAnomalies();
});
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.btn-refresh {
  padding: 0.5rem 1rem;
  background-color: white;
  border: 1px solid var(--border);
  border-radius: 4px;
  cursor: pointer;
}

.btn-refresh:hover {
  background-color: #f8fafc;
}

.loading-state, .error-state {
  text-align: center;
  padding: 3rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.error-state {
  color: var(--danger);
}
</style>