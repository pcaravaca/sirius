<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header">
      <span class="logo">⭐</span>
      <span class="logo-text" v-if="!collapsed">Sirius</span>
    </div>

    <nav class="sidebar-nav">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="nav-item"
        :class="{ active: isActive(item.to) }"
        :title="item.label"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <span class="nav-label" v-if="!collapsed">{{ item.label }}</span>
        <span class="nav-badge" v-if="item.badge && !collapsed">{{ item.badge }}</span>
      </router-link>
    </nav>

    <div class="sidebar-footer">
      <button class="nav-item" @click="collapsed = !collapsed" :title="collapsed ? 'Expand' : 'Collapse'">
        <span class="nav-icon">{{ collapsed ? '→' : '←' }}</span>
        <span class="nav-label" v-if="!collapsed">Collapse</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const collapsed = ref(false)

const navItems = computed(() => [
  { to: '/', icon: '🏠', label: 'Dashboard' },
  { to: '/projects', icon: '📁', label: 'Projects' },
  { to: '/agents', icon: '🤖', label: 'Agents' },
  { to: '/analytics', icon: '📊', label: 'Analytics' },
  { to: '/system', icon: '🖥️', label: 'System' },
  { to: '/exchange', icon: '💱', label: 'Exchange' },
  { to: '/news', icon: '📰', label: 'News' },
  { to: '/settings', icon: '⚙️', label: 'Settings' },
])

function isActive(to) {
  if (to === '/') return route.path === '/'
  return route.path.startsWith(to)
}
</script>

<style scoped>
.sidebar {
  width: var(--sidebar-width);
  height: 100%;
  background: var(--color-panel);
  border-right: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-normal);
  flex-shrink: 0;
  overflow: hidden;
}

.sidebar.collapsed {
  width: var(--sidebar-collapsed);
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  height: var(--topbar-height);
  border-bottom: 1px solid var(--border-subtle);
}

.logo {
  font-size: 20px;
  flex-shrink: 0;
}

.logo-text {
  font-size: var(--font-small);
  font-weight: 590;
  letter-spacing: -0.165px;
  color: var(--color-text-primary);
  white-space: nowrap;
}

.sidebar-nav {
  flex: 1;
  padding: var(--space-2);
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  color: var(--color-text-muted);
  font-size: var(--font-small);
  font-weight: 510;
  letter-spacing: -0.165px;
  cursor: pointer;
  transition: all var(--transition-fast);
  text-decoration: none;
  white-space: nowrap;
}

.nav-item:hover {
  background: var(--color-surface);
  color: var(--color-text-secondary);
}

.nav-item.active {
  background: rgba(94, 106, 210, 0.15);
  color: var(--color-accent-bright);
}

.nav-icon {
  font-size: 18px;
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.nav-label {
  flex: 1;
}

.nav-badge {
  background: var(--color-accent);
  color: white;
  font-size: var(--font-micro);
  font-weight: 510;
  padding: 1px 6px;
  border-radius: var(--radius-full);
}

.sidebar-footer {
  padding: var(--space-2);
  border-top: 1px solid var(--border-subtle);
}

/* Mobile: sidebar becomes bottom nav */
@media (max-width: 768px) {
  .sidebar {
    width: 100% !important;
    height: auto;
    flex-direction: row;
    border-right: none;
    border-top: 1px solid var(--border-subtle);
    order: 999;
  }

  .sidebar-header,
  .sidebar-footer {
    display: none;
  }

  .sidebar-nav {
    display: flex;
    flex-direction: row;
    overflow-x: auto;
    padding: var(--space-1);
    gap: var(--space-1);
  }

  .nav-item {
    flex-direction: column;
    gap: 2px;
    padding: var(--space-1) var(--space-2);
    min-width: 56px;
    justify-content: center;
  }

  .nav-label {
    font-size: var(--font-micro);
  }

  .nav-icon {
    font-size: 20px;
  }
}
</style>