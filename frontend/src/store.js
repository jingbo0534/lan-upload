import { reactive } from 'vue';

// 全局响应式状态
export const store = reactive({
  status: null,     // ServerStatus
  nav: 'dash',      // 当前页面: dash | files | clean | set
  // WiFi 卡片（自动持久化）
  wifi: { ssid: '', password: '', sec: 'wpa', qr: '' },
  // 文件管理
  filesPath: '',
  files: [],
  // 清理
  cleanupList: [],
  ready: false,
});

// 顶栏 Toast 状态
const toastState = reactive({ vis: false, msg: '', ok: true });
let toastTimer = null;

export function toast(msg, ok = true) {
  toastState.msg = msg;
  toastState.ok = ok;
  toastState.vis = true;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => { toastState.vis = false; }, 2600);
}

export { toastState };

// 右下角"收到新文件"弹窗
const uploadState = reactive({ vis: false, msg: '', files: [] });
export function notifyUpload(files) {
  uploadState.files = (files || []).slice(0, 5);
  uploadState.msg = files.length > 1 ? `${files.length} 个新文件已接收` : '收到新文件';
  uploadState.vis = true;
}
export function closeUpload() { uploadState.vis = false; }
export { uploadState };