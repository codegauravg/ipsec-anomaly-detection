<template>
  <div class="app-layout">
    <header class="app-header">
      <div class="logo">IPSec Security Monitor</div>
    </header>
    <main class="app-content">
      <Dashboard />
    </main>
  </div>
</template>

<script setup lang="ts">
import Dashboard from '@/views/Dashboard.vue';
</script>