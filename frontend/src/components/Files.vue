<script setup>
import { onMounted, ref, computed } from 'vue';
import { store, toast } from '../store.js';
import { ListFiles, DeleteFile, NewFolder, OpenInExplorer, ReadFilePreview } from '../api.js';
import { esc, fmtSize, iconFor, isPreviewable } from '../utils.js';

const files = ref([]);
const crumbParts = computed(() => (store.filesPath || '').split('/').filter(Boolean));

const preview = ref(null); // {name, kind, content}

function sum(list) { return list.reduce((a, b) => a + b.size, 0); }

async function load(path) {
  store.filesPath = path || '';
  try { store.files = await ListFiles(store.filesPath); }
  catch (e) { toast(String(e), false); store.files = []; }
  files.value = store.files;
}

function go(relPath) { load(relPath); }

async function doOpen(f) {
  if (f.isDir) load(f.path);
  else if (isPreviewable(f)) openPreview(f);
}

async function openPreview(f) {
  try {
    const p = await ReadFilePreview(f.path);
    if (!p.previewable) { toast('该类型不支持在线预览'); return; }
    preview.value = { name: f.name, kind: p.kind, content: p.content };
  } catch (e) { toast(String(e), false); }
}
function closePreview() { preview.value = null; }

function openContainingDir(f) {
  OpenInExplorer(f.path).catch((e) => toast(String(e), false));
}

function openFolder() {
  OpenInExplorer(store.filesPath).catch((e) => toast(String(e), false));
}

function doDelete(f) {
  if (!confirm('确定删除「' + f.name + '」吗？\n（会移入回收站，可后续彻底删除）')) return;
  DeleteFile(f.path).then(() => load(store.filesPath)).catch((e) => toast(String(e), false));
}

function newFolder() {
  const name = prompt('新建文件夹名称：');
  if (!name) return;
  NewFolder(store.filesPath, name).then(() => load(store.filesPath)).catch((e) => toast(String(e), false));
}

function crumbClick(i) {
  const parts = crumbParts.value.slice(0, i + 1);
  load(parts.join('/'));
}

onMounted(() => load(store.filesPath));
</script>

<template>
  <div>
    <div class="page-head">
      <h2>文件管理</h2><span class="sub">已接收 {{ fmtSize(sum(files)) }}</span>
      <div class="btn-row right">
        <button class="btn ghost" @click="load(store.filesPath)">🔄 刷新</button>
        <button class="btn ghost" @click="newFolder">新建文件夹</button>
        <button class="btn ghost" @click="openFolder">在资源管理器打开</button>
      </div>
    </div>

    <div class="breadcrumb">
      <span class="crumb" @click="load('')">根目录</span>
      <template v-for="(p, i) in crumbParts" :key="i">
        <span class="sep">/</span>
        <span class="crumb" @click="crumbClick(i)">{{ p }}</span>
      </template>
    </div>

    <table class="table">
      <thead>
        <tr><th>名称</th><th>大小</th><th>修改时间</th><th style="width:240px">操作</th></tr>
      </thead>
      <tbody>
        <tr v-for="f in files" :key="f.path">
          <td class="fname" @dblclick="doOpen(f)">
            {{ iconFor(f) }} <span>{{ f.name }}</span>
          </td>
          <td>{{ fmtSize(f.size) }}</td>
          <td class="muted">{{ f.modTime }}</td>
          <td class="ops">
            <button v-if="f.isDir" @click="doOpen(f)">打开</button>
            <button v-else-if="isPreviewable(f)" @click="openPreview(f)">预览</button>
            <button @click="openContainingDir(f)">打开目录</button>
            <button class="del" @click="doDelete(f)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="!files.length" class="empty">该目录暂无文件</div>

    <!-- 预览弹窗 -->
    <div v-if="preview" class="modal">
      <div class="modal-mask" @click="closePreview"></div>
      <div class="modal-card">
        <div class="modal-head">
          <b>{{ preview.name }}</b>
          <button class="close" @click="closePreview">✕</button>
        </div>
        <div class="modal-body">
          <img v-if="preview.kind === 'image'" :src="preview.content" class="pv-img" />
          <pre v-else>{{ preview.content }}</pre>
        </div>
        <div class="modal-foot"><button class="btn ghost" @click="closePreview">关闭</button></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.breadcrumb{display:flex;align-items:center;flex-wrap:wrap;gap:2px;margin-bottom:14px;font-size:13px}
.crumb{color:var(--brand);cursor:pointer;padding:3px 6px;border-radius:6px}
.crumb:hover{background:#eef4ff}
.sep{color:var(--mute)}
.modal{position:fixed;inset:0;z-index:50}
.modal-mask{position:absolute;inset:0;background:rgba(15,23,42,.45)}
.modal-card{position:absolute;inset:0;margin:1.5vh auto;background:#fff;border-radius:14px;width:96vw;max-height:97vh;display:flex;flex-direction:column;overflow:hidden;box-shadow:0 20px 50px rgba(0,0,0,.3)}
.modal-head{display:flex;justify-content:space-between;align-items:center;padding:14px 18px;border-bottom:1px solid var(--line)}
.modal-head .close{border:none;background:none;font-size:18px;cursor:pointer;color:var(--mute)}
.modal-body{padding:16px;overflow:auto;background:#f8fafc;display:flex;align-items:center;justify-content:center}
.pv-img{display:block;margin:auto;max-width:100%;max-height:calc(97vh - 150px);object-fit:contain}
.modal-body pre{white-space:pre-wrap;word-break:break-all;background:#0f172a;color:#e2e8f0;padding:16px;border-radius:10px;max-height:60vh;overflow:auto;font-size:13px;margin:0}
.modal-body img{display:block;margin:0 auto;border-radius:8px}
.modal-foot{padding:12px 18px;border-top:1px solid var(--line);display:flex;justify-content:flex-end}
</style>