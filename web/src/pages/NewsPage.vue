<template>
  <div class="news-page">
    <div class="news-grid">
      <div v-for="(item,i) in news" :key="i" class="news-card">
        <span class="news-cat">{{ item.category || 'general' }}</span>
        <div class="news-title">{{ item.title }}</div>
        <div class="news-desc" v-if="item.description">{{ item.description.substring(0,120) }}</div>
        <div class="news-footer">
          <span>{{ item.source }}</span>
          <a v-if="item.url" :href="item.url" target="_blank">Leer</a>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const news = ref([])
onMounted(async () => { try { const n = await fetch('/api/news').then(r => r.json()); news.value = Array.isArray(n) ? n : [] } catch {} })
</script>
<style scoped>
.news-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 12px; }
.news-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-md); padding: 16px; }
.news-cat { font-size: 10px; font-weight: 700; text-transform: uppercase; color: var(--brand); }
.news-title { font-size: 15px; font-weight: 700; margin: 6px 0; }
.news-desc { font-size: 13px; color: var(--text-muted); }
.news-footer { display: flex; justify-content: space-between; font-size: 12px; color: var(--text-muted); margin-top: 8px; }
</style>
