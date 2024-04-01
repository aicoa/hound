package menu

import (
	"fmt"
	"hound/common/logger"
	"hound/gui/basic"
	"hound/gui/global"
	"hound/gui/mytheme"
	"hound/lib/util"
	"net/url"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

type GlobalShortcut struct {
	dcs *desktop.CustomShortcut
	fun func(shortcut fyne.Shortcut)
}

func AddGlobalShortcutsToWindow(shortcuts []GlobalShortcut) {
	for _, sc := range shortcuts {
		global.Widgets.Win.Canvas().AddShortcut(sc.dcs, sc.fun)
	}
}

func MyMenu() *fyne.MainMenu {
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	logShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierShortcutDefault}
	AddGlobalShortcutsToWindow([]GlobalShortcut{
		{dcs: settingsShortcut, fun: func(shortcut fyne.Shortcut) { ChangeTheme() }},
		{dcs: logShortcut, fun: func(shortcut fyne.Shortcut) { basic.ConsoleWindow() }},
	})

	return fyne.NewMainMenu(
		createFileMenu(),
		createConfigMenu(),
		createUpdateMenu(),
		createThemeMenu(settingsShortcut),
		createLogMenu(logShortcut),
	)
}

func createFileMenu() *fyne.Menu {
	return fyne.NewMenu("文件", createOpenDirectoryMenuItem())
}

func createOpenDirectoryMenuItem() *fyne.MenuItem {
	return &fyne.MenuItem{
		Label: "打开目录",
		Icon:  theme.FolderOpenIcon(),
		Action: func() {
			if dir, err := os.Getwd(); err == nil {
				util.OpenFolder(dir)
			} else {
				logger.Info(fmt.Sprintf("Error getting working directory: %v", err))
			}
		},
	}
}
func createConfigMenu() *fyne.Menu {
	return fyne.NewMenu("配置", createConfigMenuItem()...)
}
func createConfigMenuItem() []*fyne.MenuItem {
	return []*fyne.MenuItem{
		{
			Label:  "修改配置",
			Icon:   theme.SettingsIcon(),
			Action: ConfigCenter,
		},
		{
			Label:  "内置字典",
			Icon:   mytheme.DictIcon(),
			Action: BuiltIn_Dict,
		},
	}
}

func createUpdateMenu() *fyne.Menu {
	return fyne.NewMenu("更新", createUpdateMenuItem()...)

}
func createUpdateMenuItem() []*fyne.MenuItem {
	return []*fyne.MenuItem{
		{
			Label:  "IP纯真库更新",
			Icon:   theme.DownloadIcon(),
			Action: DowdloadQqwry,
		},
		{
			Label: "客户端更新",
			Icon:  theme.ViewRefreshIcon(),
			Action: func() {
				ConfrimUpdateClient(1)
			},
		},
		{
			Label:  "暂不支持远程更新poc,请在本地文件中更新",
			Icon:   mytheme.UpdateIcon(),
			Action: nil,
		},
		{
			Label: "功能建议",
			Icon:  mytheme.GithubIcon(),
			Action: func() {
				u, _ := url.ParseRequestURI(fyne.CurrentApp().Metadata().Custom["Issue"])
				fyne.CurrentApp().OpenURL(u)
			},
		},
	}
}

func createLogMenu(shortcut fyne.Shortcut) *fyne.Menu {
	return fyne.NewMenu("日志", &fyne.MenuItem{
		Label:    "查看日志",
		Icon:     theme.FileTextIcon(),
		Shortcut: shortcut,
	})
}

func createThemeMenu(shortcut fyne.Shortcut) *fyne.Menu {
	return fyne.NewMenu("主题", &fyne.MenuItem{
		Label:    "修改主题",
		Icon:     mytheme.ThemeIcon(),
		Shortcut: shortcut,
		Action:   ChangeTheme,
	})
}
