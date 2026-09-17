// 纯工具函数
export function esc(s) {
  return String(s == null ? '' : s).replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));
}

export function fmtSize(n) {
  if (n == null || !isFinite(n)) return '0 B';
  const u = ['B', 'KB', 'MB', 'GB', 'TB'];
  let i = 0;
  while (n >= 1024 && i < u.length - 1) { n /= 1024; i++; }
  return (i ? n.toFixed(1) : n) + ' ' + u[i];
}

const ICONS = {
  png: '🖼️', jpg: '🖼️', jpeg: '🖼️', gif: '🖼️', webp: '🖼️', bmp: '🖼️',
  mp4: '🎬', mov: '🎬', mkv: '🎬', webm: '🎬',
  mp3: '🎵', wav: '🎵',
  pdf: '📄', txt: '📄', md: '📄',
  zip: '🗜️', rar: '🗜️', '7z': '🗜️', tar: '🗜️',
  doc: '📝', docx: '📝', xls: '📊', xlsx: '📊',
};

export function iconFor(e) {
  if (e.isDir) return '📁';
  return ICONS[e.ext] || '📎';
}

export function isPreviewable(f) {
  if (f.isDir) return false;
  return ['.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp', '.txt', '.md', '.log', '.json', '.xml', '.csv'].includes(f.ext);
}

export function cleanupTypes() {
  return [
    { label: '图片', exts: ['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp'] },
    { label: '视频', exts: ['mp4', 'mov', 'mkv', 'webm', 'avi'] },
    { label: '音频', exts: ['mp3', 'wav', 'ogg', 'flac'] },
    { label: '文档', exts: ['doc', 'docx', 'xls', 'xlsx', 'ppt', 'pptx', 'pdf', 'txt', 'md'] },
    { label: '压缩包', exts: ['zip', 'rar', '7z', 'tar', 'gz'] },
  ];
}