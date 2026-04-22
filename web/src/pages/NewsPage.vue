<template>
  <div class="news-page">
    <div class="page-header">
      <h2>📰 News</h2>
      <button class="btn-ghost" @click="fetch" :disabled="loading">Refresh</button>
    </div>

    <div class="news-loading" v-if="loading && !news.length">
      <span>Loading news...</span>
    </div>

    <div class="news-list" v-else-if="news.length">
      <a
        v-for="(item, i) in news"
        :key="i"
        :href="item.url"
        target="_blank"
        rel="noopener"
        class="news-card"
      >
        <div class="news-meta">
          <span class="news-source">{{ item.source || 'News' }}</span>
          <span class="news-time" v-if="item.publishedAt">{{ formatDate(item.publishedAt) }}</span>
        </div>
        <h3 class="news-title">{{ item.title }}</h3>
        <p class="news-desc" v-if="item.description">{{ item.description }}</p>
        <div class="news-tags" v-if="item.category">
          <span class="tag">{{ item.category }}</span>
        </div>
      </a>
    </div>

    <div class="news-empty" v-else>
      <span class="empty-icon">📰</span>
      <p>No news available</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useNews } from '@/composables/useApi'

const { news, loading, fetch } = useNews()

function formatDate(ts) {
  if (!ts) return ''
  return new Date(ts).toLocaleDateString('es-CR', { day: 'numeric', month: 'short' })
}

onMounted(fetch)
</script>

<style scoped>
.news-page {
  max-width: 800px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
}

.page-header h2 {
  font-size: var(--font-h3);
  font-weight: 590;
}

.btn-ghost {
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-4);
  color: var(--color-text-secondary);
  cursor: pointer;
  font-size: var(--font-caption);
}

.btn-ghost:hover { background: var(--color-elevated); }
.btn-ghost:disabled { opacity: 0.5; }

.news-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.news-card {
  display: block;
  background: var(--color-surface);
  border: 1px solid var(--border-standard);
  border-radius: var(--radius-lg);
  padding: var(--space-4) var(--space-5);
  text-decoration: none;
  transition: all var(--transition-fast);
}

.news-card:hover {
  border-color: var(--border-strong);
  background: var(--color-elevated);
}

.news-meta {
  display: flex;
  gap: var(--space-3);
  margin-bottom: var(--space-2);
  font-size: var(--font-micro);
}

.news-source {
  color: var(--color-accent-bright);
  font-weight: 510;
}

.news-time {
  color: var(--color-text-subtle);
}

.news-title {
  font-size: var(--font-body);
  font-weight: 510;
  color: var(--color-text-primary);
  line-height: 1.4;
  margin-bottom: var(--space-2);
}

.news-desc {
  font-size: var(--font-caption);
  color: var(--color-text-muted);
  line-height: 1.5;
}

.news-tags {
  margin-top: var(--space-2);
  display: flex;
  gap: var(--space-2);
}

.tag {
  font-size: var(--font-micro);
  padding: 2px 8px;
  background: rgba(94, 106, 210, 0.1);
  color: var(--color-accent-bright);
  border-radius: var(--radius-full);
}

.news-loading, .news-empty {
  text-align: center;
  padding: var(--space-12);
  color: var(--color-text-muted);
}

.empty-icon {
  font-size: 32px;
  display: block;
  margin-bottom: var(--space-3);
}
</style>