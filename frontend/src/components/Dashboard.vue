<script setup>
import { computed } from 'vue';
import { store } from '../store.js';
import { esc } from '../utils.js';
import WifiCard from './WifiCard.vue';

const s = computed(() => store.status || { running: false, port: 0, url: '', qrDataUri: '' });
</script>

<template>
  <div>
    <div class="page-head">
      <h2>概览</h2><span class="sub">接收服务器实时状态</span>
    </div>

    <div class="grid2">
      <div class="card">
        <div class="card-title">状态</div>
        <div class="status-big" :class="s.running ? 'on' : 'off'">{{ s.running ? '● 接收中' : '○ 已停止' }}</div>
        <div class="kv"><span>访问地址</span><code>{{ s.url || '—' }}</code></div>
        <div class="kv"><span>服务端口</span><code>{{ s.port }}</code></div>
      </div>

      <div class="card center">
        <div class="card-title">扫码上传</div>
        <img v-if="s.qrDataUri" class="qr" :src="s.qrDataUri" alt="扫码上传" />
        <div v-else class="qr-empty">服务未启动</div>
        <div class="qr-tip">手机相机扫码即可上传文件</div>
      </div>
    </div>

    <WifiCard />
  </div>
</template>

<style scoped>
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:14px;margin-bottom:14px}
.center{text-align:center}
.status-big{font-size:26px;font-weight:700;margin:6px 0 16px}
.status-big.on{color:var(--ok)}
.status-big.off{color:var(--mute)}
.kv{display:flex;justify-content:space-between;padding:9px 0;border-bottom:1px dashed var(--line);font-size:13px}
.kv code{background:#f1f5f9;padding:3px 8px;border-radius:6px;font-size:13px;color:var(--brand-dark)}
</style>