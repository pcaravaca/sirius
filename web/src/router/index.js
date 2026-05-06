import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', name: 'Dashboard', component: () => import('@/pages/DashboardPage.vue') },
  { path: '/agents', name: 'Agentes', component: () => import('@/pages/AgentsPage.vue') },
  { path: '/projects', name: 'Proyectos', component: () => import('@/pages/ProjectsPage.vue') },
  { path: '/analytics', name: 'Analytics', component: () => import('@/pages/AnalyticsPage.vue') },
  { path: '/system', name: 'Sistema', component: () => import('@/pages/SystemPage.vue') },
  { path: '/exchange', name: 'Exchange', component: () => import('@/pages/ExchangePage.vue') },
  { path: '/news', name: 'Noticias', component: () => import('@/pages/NewsPage.vue') },
]

const router = createRouter({ history: createWebHistory(), routes })
export default router
