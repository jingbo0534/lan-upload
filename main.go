package main

import (
	"embed"

	"lan-upload/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 后端逻辑全部在 app 包中，main 仅作为入口
	application := app.NewApp()

	// 系统托盘独立 goroutine 运行
	go app.NewTray(application).Run()

	err := wails.Run(&options.App{
		Title:     "局域网文件接收",
		Width:     1180,
		Height:    840,
		MinWidth:  980,
		MinHeight: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 244, G: 246, B: 251, A: 1},
		OnStartup:        application.Startup,
		HideWindowOnClose: true, // 关闭窗口时最小化到托盘而非退出
		Bind: []interface{}{
			application,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}