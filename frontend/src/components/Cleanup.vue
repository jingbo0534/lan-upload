<script setup>
import { onMounted, ref } from 'vue';
import { store, toast } from '../store.js';
import { CleanPreview, RunCleanup, TrashStats, EmptyTrash } from '../api.js';
import { fmtSize, iconFor, esc, cleanupTypes } from '../utils.js';

const types = cleanupTypes();
const keyword = ref('');
const typeVal = ref('');
const before = ref('');
const min = ref('');
const max = ref('');
const result = ref(null);    // CleanupPreview
const trash = ref({ count: 0, size: 0 });

function request() {
  return {
    paths: null,
    keyword: keyword.value.trim(),
    exts: typeVal.value ? typeVal.value.split(',') : [],
    before: before.value,
    minSize: Number(min.value || 0),
    maxSize: Number(max.value || 0),
  };
}

async function preview() {
  try {
    result.value = await CleanPreview(request());
    if (result.value.files.length === 0) toast('没有符合条件的文件', false);
  } catch (e) { toast(String(e), false); }
}

async function run() {
  if (!confirm('将把命中的 ' + (result.value?.files.length || 0) + ' 个文件移入回收站，确认执行？')) return;
  try {
    const r = await RunCleanup(request());
    toast(`已清理 ${r.deleted} 个，失败 ${r.failed} 个`);
    result.value = null;
    loadTrash();
  } catch (e) { toast(String(e), false); }
}

async function loadTrash() {
  try { trash.value = await TrashStats(); } catch (e) { /* 忽略 */ }
}

async function emptyTrash() {
  if (!confirm('清空回收站将彻底删除其中所有文件，确认？')) return;
  try { await EmptyTrash(); toast('回收站已清空'); loadTrash(); }
  catch (e) { toast(String(e), false); }
}

onMounted(() => loadTrash());
</script>

<template>
  <div>
    <div class="page-head">
      <h2>文件清理</h2>
      <span class="sub">筛选后预览确认，执行时移入回收站（可清空）</span>
    </div>

    <div class="card">
      <div class="filters">
        <label>名称关键字 <input v-model="keyword" placeholder="包含的文字" /></label>
        <label>文件类型
          <select v-model="typeVal">
            <option value="">全部</option>
            <option v-for="t in types" :key="t.label" :value="t.exts.join(',')">{{ t.label }}</option>
          </select>
        </label>
        <label>早于日期 <input v-model="before" type="date" /></label>
        <label>最小大小
          <select v-model="min">
            <option value="">不限</option>
            <option value="1048576">≥ 1 MB</option>
            <option value="10485760">≥ 10 MB</option>
            <option value="104857600">≥ 100 MB</option>
            <option value="1073741824">≥ 1 GB</option>
          </select>
        </label>
        <label>最大大小
          <select v-model="max">
            <option value="">不限</option>
            <option value="1048576">≤ 1 MB</option>
            <option value="10485760">≤ 10 MB</option>
          </select>
        </label>
      </div>
      <div class="btn-row">
        <button class="btn primary" @click="preview">预览命中文件</button>
        <button class="btn danger" :disabled="!result" @click="run">执行清理（移入回收站）</button>
      </div>
    </div>

    <div class="card">
      <div class="card-title">命中结果</div>
      <div class="trash-line">
        回收站：<b>{{ trash.count }}</b> 项，{{ fmtSize(trash.size) }}
        <button class="btn ghost sm" @click="emptyTrash">清空回收站</button>
      </div>
      <div v-if="result && result.files.length">
        <div class="preview-info">命中 <b>{{ result.total }}</b> 个文件，共 <b>{{ fmtSize(result.size) }}</b></div>
        <table class="table">
          <thead><tr><th>名称</th><th>大小</th><th>修改时间</th></tr></thead>
          <tbody>
            <tr v-for="f in result.files" :key="f.path">
              <td class="fname">{{ iconFor(f) }} {{ f.name }}</td>
              <td>{{ fmtSize(f.size) }}</td>
              <td class="muted">{{ f.modTime }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="empty">设置清理条件后点击「预览命中文件」</div>
    </div>
  </div>
</template>

<style scoped>
.filters{display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px;margin-bottom:8px}
.filters label{display:flex;flex-direction:column;gap:6px;font-size:12px;color:var(--mute)}
.filters input,.filters select{padding:8px 10px;border:1px solid var(--line);border-radius:8px;font-size:13px;color:var(--ink);background:#fff}
.preview-info{margin:10px 0 12px;color:#475569}
.trash-line{display:flex;align-items:center;gap:10px;background:#fefce8;border:1px solid #fde68a;border-radius:9px;padding:10px 12px;margin-bottom:14px;font-size:13px}
</style>