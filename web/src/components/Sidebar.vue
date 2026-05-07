<template>
  <aside class="sidebar" :class="{ collapsed }">
    <div class="sidebar-header">
      <div class="logo" @click="$emit('toggle')">
        <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/><circle cx="12" cy="12" r="4"/><circle cx="12" cy="12" r="1.5"/>
        </svg>
        <span v-show="!collapsed" class="logo-text">Sirius</span>
      </div>
    </div>
    <nav class="sidebar-nav">
      <router-link v-for="item in navItems" :key="item.path" :to="item.path" class="nav-item" :title="item.label">
        <span class="nav-icon" v-html="item.icon"></span>
        <span v-show="!collapsed" class="nav-label">{{ item.label }}</span>
      </router-link>
    </nav>
    <div class="sidebar-footer">
      <div class="version-badge" v-show="!collapsed">v3.1</div>
    </div>
  </aside>
</template>

<script setup>
defineProps({ collapsed: Boolean })
defineEmits(['toggle'])

const navItems = [
  { path: '/', label: 'Dashboard', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/></svg>' },
  { path: '/agents', label: 'Agentes', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>' },
  { path: '/projects', label: 'Proyectos', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>' },
  { path: '/analytics', label: 'Analytics', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/></svg>' },
  { path: '/system', label: 'Sistema', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>' },
  { path: '/exchange', label: 'Exchange', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>' },
  { path: '/news', label: 'Noticias', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="20" height="20"><path d="M4 22h16a2 2 0 0 0 2-2V4a2 2 0 0 0-2-2H8a2 2 0 0 0-2 2v16a2 2 0 0 1-4 0v-9"/><line x1="10" y1="8" x2="18" y2="8"/><line x1="10" y1="12" x2="18" y2="12"/></svg>' },
]
</script>

<style scoped>
.sidebar {
  flex-shrink: 0;
  width: var(--sidebar-width);
  height: 100vh;
  background: #0a0e14;
  border-right: 1px solid rgba(0, 153, 255, 0.08);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-normal);
  overflow: hidden;
}
.sidebar.collapsed { width: var(--sidebar-collapsed); }

.sidebar-header {
  padding: 14px 12px;
  border-bottom: 1px solid var(--border);
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 4px;
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}
.logo:hover { background: rgba(0,153,255,0.06); }

.logo-icon {
  width: 26px; height: 26px;
  color: var(--brand);
  flex-shrink: 0;
}
.logo-text {
  font-size: 18px; font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
  letter-spacing: -0.3px;
}

.sidebar-nav {
  flex: 1;
  padding: 8px 6px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 10px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  text-decoration: none;
  transition: all var(--transition-fast);
  white-space: nowrap;
  font-size: var(--font-size-sm);
  font-weight: 500;
  border-left: 2px solid transparent;
}
.nav-item:hover {
  background: rgba(0,153,255,0.06);
  color: var(--text-primary);
  border-left-color: rgba(0,153,255,0.2);
}
.nav-item.router-link-exact-active {
  background: rgba(0,153,255,0.1);
  color: var(--brand);
  border-left-color: var(--brand);
}

.nav-icon {
  display: flex; align-items: center; justify-content: center;
  width: 20px; height: 20px;
  flex-shrink: 0;
  opacity: 0.7;
}
.nav-item:hover .nav-icon,
.nav-item.router-link-exact-active .nav-icon { opacity: 1; }

.sidebar-footer {
  padding: 10px 12px;
  border-top: 1px solid var(--border);
}
.version-badge {
  font-size: var(--font-size-xs);
  color: var(--text-muted);
  font-family: var(--font-mono);
}
</style>
