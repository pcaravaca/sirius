<template>
  <div class="app-layout">
    <Sidebar />
    <div class="app-main">
      <TopBar />
      <main class="app-content">
        <router-view />
      </main>
      <StatusBar />
    </div>
  </div>
</template>

<script setup>
import Sidebar from '@/components/Sidebar.vue'
import TopBar from '@/components/TopBar.vue'
import StatusBar from '@/components/StatusBar.vue'
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  background: var(--color-canvas);
}

.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.app-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-6);
  background: var(--color-canvas);
}

@media (max-width: 768px) {
  .app-layout {
    flex-direction: column;
  }

  .app-main {
    order: 1;
  }

  .app-content {
    padding: var(--space-3);
  }
}
</style>