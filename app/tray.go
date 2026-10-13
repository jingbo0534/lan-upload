package app

import (
	_ "embed"
	"runtime"

	"github.com/getlantern/systray"

	wailsrt "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed tray_icon.ico
var trayIcon []byte

// Tray 系统托盘
type Tray struct {
	app *App
}

func NewTray(a *App) *Tray { return &Tray{app: a} }

// Run 阻塞运行系统托盘（放入 goroutine 调用）
func (t *Tray) Run() {
	// Windows 下托盘隐藏窗口与消息循环必须固定在同一 OS 线程，
	// 否则 goroutine 被调度到其它线程后，点击事件将永远收不到（表现为点击无反应）
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	systray.Run(t.onReady, t.onExit)
}

func (t *Tray) onReady() {
	systray.SetIcon(trayIcon)
	systray.SetTitle("局域网文件接收")
	systray.SetTooltip("局域网文件接收工具")

	mShow := systray.AddMenuItem("打开管理界面", "显示主窗口")
	mUrl := systray.AddMenuItem("复制上传地址", "复制扫码上传网址")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出并关闭接收服务")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				t.show()
			case <-mUrl.ClickedCh:
				t.copyURL()
			case <-mQuit.ClickedCh:
				_ = t.app.StopServer()
				if t.app.ctx != nil {
					wailsrt.Quit(t.app.ctx)
				}
				systray.Quit()
				return
			}
		}
	}()
}

func (t *Tray) onExit() {}

func (t *Tray) show() {
	if t.app.ctx == nil {
		return
	}
	wailsrt.WindowShow(t.app.ctx)
	wailsrt.WindowSetAlwaysOnTop(t.app.ctx, true)
	wailsrt.WindowSetAlwaysOnTop(t.app.ctx, false)
}

func (t *Tray) copyURL() {
	if t.app.ctx == nil {
		return
	}
	st := t.app.GetStatus()
	if st.URL != "" {
		_ = wailsrt.ClipboardSetText(t.app.ctx, st.URL)
	}
}