package menu

import (
	"bufio"
	"errors"
	"fmt"
	"hound/common/logger"
	"hound/gui/basic"
	"hound/gui/global"
	"hound/gui/mytheme"
	"hound/lib/qqwry"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	update "github.com/fynelabs/selfupdate"
)

const Qqwrypath = "./config/qqwry.dat"
const lastestClinetUrl = "https://github.com/aicoa/hound/releases/download/"
const remoteClientVersion = "https://gihub.com/aicoa/hound/main/version"
const updateClientContent = "https://github.com/aicoa/hound/main/update"

func DowdloadQqwry() {
	if err := qqwry.Download(Qqwrypath); err == nil {
		dialog.ShowInformation("", "update success", global.Widgets.Win)
	} else {
		dialog.ShowError(fmt.Errorf("update failed %v", err), global.Widgets.Win)
	}
}

func doUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			fmt.Printf("Failed to rollback from bad update: %v", rerr)
		}
	}
	return err
}

func CheckUpdate(removeTarget, localVersion string) (string, error) {
	r, err := http.Get(removeTarget)
	if err != nil {
		return "", err
	}
	b, err2 := io.ReadAll(r.Body)
	if err2 != nil {
		return "", err2
	}
	remoteVersion := string(b)
	if remoteVersion <= localVersion {
		return "", errors.New("当前已是最新版本 " + localVersion)
	}
	return remoteVersion, nil
}

// numOfTimes 为0时不显示dialog提示
func ConfrimUpdateClient(numOfTimes int) {
	if numOfTimes == 0 {
		return
	}
	if version, err := CheckUpdate(remoteClientVersion, fyne.CurrentApp().Metadata().Version); err != nil {
		dialog.ShowInformation("提示", "客户端"+err.Error(), global.Widgets.Win)
	} else {
		r, err := http.Get(updateClientContent)
		if err != nil {
			return
		}
		b, err2 := io.ReadAll(r.Body)
		if err2 != nil {
			return
		}
		dp := widget.NewProgressBarInfinite()
		dp.Hide()
		l := &widget.Label{Text: string(b), Wrapping: fyne.TextWrapBreak}
		content := container.NewBorder(nil, dp, nil, nil, container.NewVScroll(l))
		basic.ShowCustomDialog(mytheme.UpdateIcon(), "更新提醒,最新版本:"+version, "立即更新", content, func() {
			dp.Show()
			content.Refresh()
			if err3 := UpdateClinet(version); err3 != nil {
				l.SetText(fmt.Sprintf("更新失败: %v", err3))
			}
		}, fyne.NewSize(400, 300))
	}
}

func UpdateClinet(latestVersion string) error {
	var binaryFileName string
	switch runtime.GOOS {
	case "windows":
		binaryFileName = "hound.exe"
	case "linux":
		binaryFileName = "hound_linux_amd64"
	case "darwin":
		binaryFileName = "hound_darwin_amd64"
	}
	if err := doUpdate(lastestClinetUrl + "v" + latestVersion + "/" + binaryFileName); err != nil {
		return err
	}
	basic.ShowCustomDialog(theme.InfoIcon(), "提示", "立即重启客户端", basic.NewCenterLabel("更新成功!!"), func() {
		go func() {
			cmd := exec.Command(os.Args[0])
			err := cmd.Start()
			if err != nil {
				logger.Error(err)
			}
			// 退出当前的进程
			os.Exit(0)
		}()
	}, fyne.NewSize(100, 50))
	return nil
}

func download(target, dest string) (string, error) {
	fileName := path.Base(target)
	res, err := http.Get(target)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)
	file, err := os.Create(dest + fileName)
	if err != nil {
		return "", err
	}
	//获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)
	return fileName, nil
}
