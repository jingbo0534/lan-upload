<script setup>
import { onMounted, ref } from 'vue';
import { store, toast } from '../store.js';
import { GetConfig, SetSaveDir, SetAutoFolder, SetServer, ChooseSaveDir } from '../api.js';

const port = ref(3210);
const dir = ref('');
const auto = ref(true);
const ip = ref('');

onMounted(async () => {
  try {
    const c = await GetConfig();
    port.value = c.port;
    dir.value = c.saveDir;
    auto.value = c.autoFolder;
    ip.value = c.ip || (store.status?.ips?.[0] || '');
  } catch (e) { toast(String(e), false); }
});

async function pickDir() {
  try { dir.value = await ChooseSaveDir(); } catch (e) { toast(String(e), false); }
}

async function save() {
  try {
    if (dir.value) await SetSaveDir(dir.value);
    await SetAutoFolder(auto.value);
    store.status = await SetServer(ip.value, Number(port.value));
    toast('设置已保存');
  } catch (e) { toast(String(e), false); }
}
</script>

<template>
  <div>
    <div class="page-head"><h2>设置</h2></div>
    <div class="card">
      <div class="set-row">
        <label>对外广播地址</label>
        <select v-model="ip" class="ipselect">
          <option v-for="a in (store.status?.ips || [])" :key="a" :value="a">{{ a }}</option>
          <option v-if="!(store.status?.ips?.length)" value="">未检测到局域网地址</option>
        </select>
        <span class="hint-inline">默认自动取第一个，VPN 多地址时可手动指定</span>
      </div>
      <div class="set-row">
        <label>服务端口</label>
        <input v-model.number="port" type="number" min="1" max="65535" />
      </div>
      <div class="set-row">
        <label>保存位置</label>
        <input :value="dir" type="text" readonly />
        <button class="btn ghost" @click="pickDir">浏览…</button>
      </div>
      <div class="set-row">
        <label>按日期自动分文件夹</label>
        <input v-model="auto" type="checkbox" />
      </div>
      <div class="btn-row"><button class="btn primary" @click="save">保存设置</button></div>
    </div>
  </div>
</template>

<style scoped>
.set-row{display:flex;align-items:center;gap:12px;margin-bottom:16px}
.set-row label{width:160px;flex:0 0 auto;color:#475569;font-size:13px}
.set-row input[type=text],.set-row input[type=number],.set-row select{flex:1;padding:8px 10px;border:1px solid var(--line);border-radius:8px;font-size:13px}
.set-row select.ipselect{flex:0 1 220px;background:#fff}
.set-row input[readonly]{background:#f8fafc;color:var(--mute)}
.hint-inline{font-size:12px;color:var(--mute)}
</style>