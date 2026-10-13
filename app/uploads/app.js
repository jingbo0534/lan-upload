(function () {
  var dropZone = document.getElementById('dropZone');
  var fileInput = document.getElementById('fileInput');
  var camInput = document.getElementById('camInput');
  var pickFiles = document.getElementById('pickFiles');
  var pickCam = document.getElementById('pickCam');
  var list = document.getElementById('fileList');
  var counter = document.getElementById('counter');
  var uploadBtn = document.getElementById('uploadBtn');
  var clearBtn = document.getElementById('clearBtn');
  var doneBox = document.getElementById('doneBox');

  var files = []; // {file, id, status:'pending'|'uploading'|'ok'|'err', pct, msg}

  function fmt(bs) {
    if (bs < 1024) return bs + ' B';
    if (bs < 1048576) return (bs / 1024).toFixed(1) + ' KB';
    if (bs < 1073741824) return (bs / 1048576).toFixed(1) + ' MB';
    return (bs / 1073741824).toFixed(2) + ' GB';
  }

  function refresh() {
    counter.textContent = '已选 ' + files.length + ' 个文件';
    if (files.length === 0) {
      uploadBtn.disabled = true;
      clearBtn.hidden = true;
    } else {
      uploadBtn.disabled = false;
      clearBtn.hidden = false;
    }
    list.innerHTML = '';
    files.forEach(function (f) {
      var li = document.createElement('li');
      var name = document.createElement('div');
      name.className = 'name'; name.textContent = f.file.name;
      var size = document.createElement('div');
      size.className = 'size'; size.textContent = fmt(f.file.size);
      var state = document.createElement('div');
      state.className = 'state ' + f.status;
      if (f.status === 'pending') state.textContent = '待上传';
      else if (f.status === 'uploading') state.textContent = '上传中 ' + (f.pct || 0) + '%';
      else if (f.status === 'ok') state.textContent = '已完成';
      else state.textContent = '失败：' + (f.msg || '');
      li.appendChild(name); li.appendChild(size); li.appendChild(state);
      list.appendChild(li);
    });
  }

  // 进入选择器：选文件（相册/文件，任意类型）；拍照（Camera 摄像头）
  pickFiles.addEventListener('click', function (e) { e.stopPropagation(); fileInput.click(); });
  pickCam.addEventListener('click', function (e) { e.stopPropagation(); camInput.click(); });

  dropZone.addEventListener('click', function () { fileInput.click(); });
  ['dragenter', 'dragover'].forEach(function (ev) {
    dropZone.addEventListener(ev, function (e) { e.preventDefault(); dropZone.classList.add('dragover'); });
  });
  ['dragleave', 'drop'].forEach(function (ev) {
    dropZone.addEventListener(ev, function (e) { e.preventDefault(); dropZone.classList.remove('dragover'); });
  });
  dropZone.addEventListener('drop', function (e) {
    var dt = e.dataTransfer;
    if (dt && dt.files) addFiles(dt.files);
  });
  fileInput.addEventListener('change', function () {
    if (fileInput.files) addFiles(fileInput.files);
    fileInput.value = '';
  });
  camInput.addEventListener('change', function () {
    if (camInput.files) addFiles(camInput.files);
    camInput.value = '';
  });

  function addFiles(fileList) {
    Array.prototype.forEach.call(fileList, function (file) {
      files.push({ file: file, status: 'pending', pct: 0, msg: '' });
    });
    refresh();
  }

  clearBtn.addEventListener('click', function () {
    files = files.filter(function (f) { return f.status === 'uploading'; });
    refresh();
  });

  uploadBtn.addEventListener('click', function () {
    var pending = files.filter(function (f) { return f.status === 'pending'; });
    if (pending.length === 0) return;
    doneBox.hidden = true;

    // 一次请求提交本次选中的所有文件（服务端按顺序返回每条结果）
    var fd = new FormData();
    pending.forEach(function (f) {
      fd.append('files', f.file, f.file.name);
      f.status = 'uploading';
      f.pct = 0;
    });
    refresh();

    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/upload');
    if ('upload' in xhr) {
      xhr.upload.onprogress = function (e) {
        if (!e.lengthComputable) return;
        var pct = Math.round(e.loaded / e.total * 100);
        pending.forEach(function (f) { f.pct = pct; });
        refresh();
      };
    }

    function markAll(status, msg) {
      pending.forEach(function (f) { f.status = status; f.msg = msg || ''; });
    }

    xhr.onload = function () {
      var okN = 0;
      if (xhr.status === 200) {
        try {
          var r = JSON.parse(xhr.responseText);
          var res = r.results || [];
          pending.forEach(function (f, i) {
            var item = res[i];
            if (item && item.ok) { f.status = 'ok'; f.msg = ''; okN++; }
            else { f.status = 'err'; f.msg = (item && item.err) || '上传失败'; }
          });
        } catch (e) { markAll('err', '解析失败'); }
      } else {
        markAll('err', 'HTTP ' + xhr.status);
      }
      refresh();
      showDone(okN, pending.length);
    };
    xhr.onerror = function () {
      markAll('err', '连接失败');
      refresh();
      showDone(0, pending.length);
    };
    xhr.send(fd);
  });

  function showDone(ok, total) {
    doneBox.hidden = false;
    doneBox.textContent = ok === total
      ? '全部 ' + total + ' 个文件上传成功 ✅'
      : '上传完成：成功 ' + ok + ' / ' + total + '，失败的可在列表底部查看';
  }
})();