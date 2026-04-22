import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'dashboard',
    component: () => import('@/pages/DashboardPage.vue'),
    meta: { icon: '🏠', label: 'Dashboard' },
  },
  {
    path: '/projects',
    name: 'projects',
    component: () => import('@/pages/ProjectsPage.vue'),
    meta: { icon: '📁', label: 'Projects' },
  },
  {
    path: '/projects/:id',
    name: 'project-detail',
    component: () => import('@/pages/ProjectDetailPage.vue'),
    meta: { icon: '📁', label: 'Project', hidden: true },
  },
  {
    path: '/agents',
    name: 'agents',
    component: () => import('@/pages/AgentsPage.vue'),
    meta: { icon: '🤖', label: 'Agents' },
  },
  {
    path: '/agents/:id',
    name: 'agent-detail',
    component: () => import('@/pages/AgentDetailPage.vue'),
    meta: { icon: '🤖', label: 'Agent', hidden: true },
  },
  {
    path: '/analytics',
    name: 'analytics',
    component: () => import('@/pages/AnalyticsPage.vue'),
    meta: { icon: '📊', label: 'Analytics' },
  },
  {
    path: '/system',
    name: 'system',
    component: () => import('@/pages/SystemPage.vue'),
    meta: { icon: '🖥️', label: 'System' },
  },
  {
    path: '/exchange',
    name: 'exchange',
    component: () => import('@/pages/ExchangePage.vue'),
    meta: { icon: '💱', label: 'Exchange' },
  },
  {
    path: '/news',
    name: 'news',
    component: () => import('@/pages/NewsPage.vue'),
    meta: { icon: '📰', label: 'News' },
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/pages/SettingsPage.vue'),
    meta: { icon: '⚙️', label: 'Settings' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router