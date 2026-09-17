<script setup>
import { onMounted } from 'vue';
import { store, toastState, toast, uploadState, notifyUpload } from './store.js';
import { GetStatus } from './api.js';
import { EventsOn, SendNotification, InitializeNotifications } from '../wailsjs/runtime/runtime';

import Dashboard from './components/Dashboard.vue';
import Files from './components/Files.vue';
import Cleanup from './components/Cleanup.vue';
import Settings from './components/Settings.vue';

const navs = [
  { key: 'dash', label: '概览', icon: '📶' },
  { key: 'files', label: '文件管理', icon: '📂' },
  { key: 'clean', label: '文件清理', icon: '🧹' },
  { key: 'set', label: '设置', icon: '⚙️' },
];

async function refreshStatus() {
  try { store.status = await GetStatus(); } catch (e) { toast(String(e), false); }
}

function switchTab(key) { store.nav = key; }

async function onUpload(files) {
  // 界面右下角弹窗
  notifyUpload(files);
  // 系统通知（尽力而为，便携版可能受限于应用标识）
  try {
    await SendNotification({
      title: '局域网文件接收',
      body: files.length > 1 ? `收到 ${files.length} 个新文件` : files[0],
    });
  } catch (e) { /* 忽略系统通知失败 */ }
}

onMounted(async () => {
  EventsOn('server:status', (st) => { store.status = st; });
  EventsOn('server:upload', onUpload);
  try { await InitializeNotifications(); } catch (e) { /* 忽略 */ }
  await refreshStatus();
  store.ready = true;
});
</script>

<template>
  <div class="shell">
    <aside class="sidebar">
      <div class="brand">LA</div>
      <div class="brand-sub">局域网文件接收</div>
      <nav>
        <button
          v-for="n in navs"
          :key="n.key"
          class="nav-item"
          :class="{ active: store.nav === n.key }"
          @click="switchTab(n.key)"
        >{{ n.icon }} {{ n.label }}</button>
      </nav>
      <div class="side-foot">关闭窗口不会退出，程序仍继续接收<br />可最小化到系统托盘</div>
    </aside>

    <main class="content">
      <Dashboard v-if="store.nav === 'dash'" />
      <Files v-else-if="store.nav === 'files'" />
      <Cleanup v-else-if="store.nav === 'clean'" />
      <Settings v-else />
    </main>

    <div id="toast" class="toast" :class="{ show: toastState.vis, err: !toastState.ok }">{{ toastState.msg }}</div>

    <!-- 右下角新文件弹窗 -->
    <div class="upload-pop" :class="{ show: uploadState.vis }" @click="uploadState.vis = false">
      <div class="up-title">📥 {{ uploadState.msg }}</div>
      <div v-for="f in uploadState.files" :key="f" class="up-file">{{ f }}</div>
    </div>
  </div>
</template>

<style scoped>
.shell{display:flex;height:100vh;overflow:hidden}
.sidebar{width:208px;flex:0 0 auto;background:#1e2a3a;color:#cbd5e1;display:flex;flex-direction:column;padding:18px 12px;user-select:none}
.sidebar .brand{font-size:22px;font-weight:800;color:#fff;background:linear-gradient(135deg,#3b82f6,#8b5cf6);-webkit-background-clip:text;background-clip:text;-webkit-text-fill-color:transparent;letter-spacing:1px}
.sidebar .brand-sub{font-size:12px;color:#64748b;margin-bottom:22px}
.sidebar nav{display:flex;flex-direction:column;gap:4px}
.nav-item{background:none;border:none;color:#cbd5e1;text-align:left;padding:11px 12px;border-radius:9px;cursor:pointer;font-size:14px;transition:.15s}
.nav-item:hover{background:#2b3a4d;color:#fff}
.nav-item.active{background:var(--brand);color:#fff;font-weight:600}
.side-foot{margin-top:auto;font-size:11px;color:#64748b;line-height:1.6}
.content{flex:1;overflow:auto;padding:18px 22px}
.toast{position:fixed;left:50%;bottom:28px;transform:translateX(-50%) translateY(20px);background:#1e2a3a;color:#fff;padding:10px 18px;border-radius:10px;font-size:14px;opacity:0;transition:.25s;z-index:99;pointer-events:none}
.toast.show{opacity:1;transform:translateX(-50%) translateY(0)}
.toast.err{background:var(--danger)}
.toast.ok{background:#0f172a}
.upload-pop{position:fixed;right:20px;bottom:20px;width:320px;background:#fff;border:1px solid var(--line);border-left:4px solid var(--ok);border-radius:10px;box-shadow:0 10px 30px rgba(0,0,0,.18);padding:12px 14px;z-index:98;opacity:0;transform:translateY(10px);pointer-events:none;transition:.25s}
.upload-pop.show{opacity:1;transform:translateY(0);pointer-events:auto}
.upload-pop .up-title{font-weight:700;font-size:14px;margin-bottom:6px}
.upload-pop .up-file{font-size:12px;color:#475569;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;line-height:1.7}
</style>