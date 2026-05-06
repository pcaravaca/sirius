<template>
  <div class="ex-page">
    <div class="rate-hero">
      <div class="rate-card"><div class="rate-lbl">Compra</div><div class="rate-val">₡{{ exchange.usdToCrc?.toFixed(0) }}</div></div>
      <div class="rate-card"><div class="rate-lbl">Venta</div><div class="rate-val">₡{{ exchange.sellRate?.toFixed(0) || '—' }}</div></div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
const exchange = ref({})
onMounted(async () => { try { exchange.value = await fetch('/api/exchange').then(r => r.json()) } catch {} })
</script>
<style scoped>
.rate-hero { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; max-width: 400px; }
.rate-card { background: var(--panel); border: 1px solid var(--border); border-radius: var(--radius-lg); padding: 24px; text-align: center; }
.rate-lbl { font-size: 12px; font-weight: 600; color: var(--text-muted); text-transform: uppercase; margin-bottom: 8px; }
.rate-val { font-size: 32px; font-weight: 700; font-family: var(--font-mono); color: var(--success); }
</style>
