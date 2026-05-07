<template>
  <div class="app-layout">
    <Sidebar :collapsed="sidebarCollapsed" @toggle="sidebarCollapsed = !sidebarCollapsed" />
    <div class="main-area">
      <TopBar @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed" :collapsed="sidebarCollapsed" />
      <div class="content-area">
        <router-view />
      </div>
      <StatusBar />
    </div>
  </div>
</template>

<script setup>
import { ref, provide } from 'vue'
import Sidebar from '@/components/Sidebar.vue'
import TopBar from '@/components/TopBar.vue'
import StatusBar from '@/components/StatusBar.vue'

const sidebarCollapsed = ref(false)
const wsConnected = ref(false)
provide('wsConnected', wsConnected)
</script>

<style scoped>
/* ── CRITICAL: Flexbox layout que NO shiftea ── */
.app-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.main-area {
  flex: 1 1 0;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  max-width: 100%;
}
</style>
