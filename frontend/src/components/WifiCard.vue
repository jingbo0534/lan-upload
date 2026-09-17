<script setup>
import { onMounted, ref } from 'vue';
import { store, toast } from '../store.js';
import { GetConfig, MakeWifiQR, SaveWifi, CurrentWifiSSID } from '../api.js';
import { esc } from '../utils.js';

const ssid = ref(store.wifi.ssid);
const pass = ref(store.wifi.password);
const sec = ref(store.wifi.sec);
const qr = ref(store.wifi.qr);

let saveTimer = null;
function persist() {
  clearTimeout(saveTimer);
  saveTimer = setTimeout(() => {
    store.wifi.ssid = ssid.value.trim();
    store.wifi.password = pass.value;
    store.wifi.sec = sec.value;
    SaveWifi(ssid.value.trim(), pass.value, sec.value).catch((e) => toast(String(e), false));
  }, 600);
}
function setSsid(v) { store.wifi.ssid = v.trim(); }
function setPass(v) { store.wifi.password = v; }
function setSec(v) { store.wifi.sec = v; persist(); }

async function fillCurrent() {
  try {
    const name = await CurrentWifiSSID();
    if (name) { ssid.value = name; setSsid(name); toast('已填入：' + name); persist(); }
    else toast('未能获取当前 WiFi 名称', false);
  } catch (e) { toast(String(e), false); }
}

async function generate() {
  const s = ssid.value.trim();
  if (!s) { toast('请输入 WiFi 名称', false); return; }
  try {
    qr.value = await MakeWifiQR(s, pass.value, sec.value);
    store.wifi.ssid = s; store.wifi.password = pass.value; store.wifi.sec = sec.value;
    saveWifiNow(s, pass.value, sec.value);
  } catch (e) { toast(String(e), false); }
}

async function saveWifiNow(s, p, sc) {
  store.wifi.qr = qr.value;
  SaveWifi(s, p, sc).catch((e) => toast(String(e), false));
}

onMounted(async () => {
  // 若有已保存 WiFi，自动填充并生成二维码
  try {
    const cfg = await GetConfig();
    if (cfg?.wifiSsid) {
      ssid.value = cfg.wifiSsid;
      pass.value = cfg.wifiPass || '';
      sec.value = cfg.wifiSec || 'wpa';
      store.wifi.ssid = cfg.wifiSsid; store.wifi.password = pass.value; store.wifi.sec = sec.value;
      qr.value = await MakeWifiQR(cfg.wifiSsid, pass.value, sec.value).catch(() => '');
      store.wifi.qr = qr.value;
    }
  } catch (e) { /* 忽略 */ }
});
</script>

<template>
  <div class="card">
    <div class="card-title">🔐 扫码连 WiFi</div>
    <div class="wifi-grid">
      <div class="wifi-fields">
        <label>WiFi 名称
          <input v-model="ssid" type="text" placeholder="例如 MyWiFi" @input="setSsid($event.target.value); persist()" />
        </label>
        <label>WiFi 密码
          <input v-model="pass" type="password" placeholder="开放网络请选下方加密方式" @input="setPass($event.target.value); persist()" />
        </label>
        <label>加密方式
          <select v-model="sec" @change="setSec($event.target.value)">
            <option value="wpa">WPA / WPA2</option>
            <option value="wep">WEP</option>
            <option value="open">开放（无密码）</option>
          </select>
        </label>
        <div class="btn-row">
          <button class="btn ghost" @click="fillCurrent">填入当前 WiFi</button>
          <button class="btn primary" @click="generate">生成二维码</button>
        </div>
        <p class="hint2">请用手机<strong>系统相机</strong>扫码自动连接；微信扫码仅显示文本，无法直接加入 WiFi。</p>
      </div>
      <div class="wifi-qr">
        <img v-if="qr" class="qr" :src="qr" alt="WiFi二维码" />
        <div v-else class="qr-empty">输入后点“生成二维码”</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.wifi-grid{display:grid;grid-template-columns:1fr 1fr;gap:20px;align-items:center}
.wifi-fields{display:flex;flex-direction:column;gap:12px}
.wifi-fields label{display:flex;flex-direction:column;gap:6px;font-size:12px;color:var(--mute)}
.wifi-fields input,.wifi-fields select{padding:8px 10px;border:1px solid var(--line);border-radius:8px;font-size:13px;color:var(--ink);background:#fff;width:100%;max-width:100%}
.wifi-qr{display:flex;flex-direction:column;justify-content:center;align-items:center;text-align:center}
.hint2{color:var(--mute);font-size:12px;margin:6px 0 12px}
</style>