import { ref } from 'vue';

export interface Anomaly {
  id: string;
  timestamp: string;
  source_ip: string;
  threat_type: string;
  severity: number;
  resolved: boolean;
}

export function useAnomalies() {
  const anomalies = ref<Anomaly[]>([]);
  const loading = ref(true);
  const error = ref<string | null>(null);

  const fetchAnomalies = async () => {
    loading.value = true;
    try {
      const response = await fetch('/api/anomalies');
      if (!response.ok) throw new Error('Failed to fetch anomalies');
      
      const json = await response.json();
      anomalies.value = json.data;
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred';
    } finally {
      loading.value = false;
    }
  };

  const resolveAnomaly = async (id: string) => {
    try {
      const response = await fetch(`/api/anomalies/${id}/resolve`, {
        method: 'PATCH',
      });
      
      if (!response.ok) throw new Error('Failed to resolve anomaly');
      
      // Optimistically update the local state
      const index = anomalies.value.findIndex(a => a.id === id);
      if (index !== -1) {
        anomalies.value[index].resolved = true;
      }
    } catch (err) {
      console.error(err);
      alert('Error resolving anomaly');
    }
  };

  return {
    anomalies,
    loading,
    error,
    fetchAnomalies,
    resolveAnomaly
  };
}