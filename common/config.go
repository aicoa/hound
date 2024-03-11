package common

import (
	"encoding/json"
	"fmt"
	"hound/common/logger"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	defaultConfigFile      = "./config/config.json"
	defaultConfigDirectory = "./config"
	DefaultWebTimeout      = 9
)

type Profiles struct {
	WebScan struct {
		Thread int
	}
	Subdomain struct {
		DNS1 string
		DNS2 string
	}
	PortScan struct {
		Thread  int
		Timeout int
	}
	Proxy struct {
		Enable   bool
		Mode     string
		Address  string
		Port     int
		Username string
		Password string
	}
	Hunter struct {
		Api string
	}
	Fofa struct {
		Email string
		Api   string
	}
	Quake struct {
		Api string
	}
}

var Profile Profiles

func init() {
	go FreeMemory() //回收内存
	if _, err := os.Stat(defaultConfigFile); err != nil {
		profile := Profiles{
			WebScan: struct{ Thread int }{100},
			Subdomain: struct {
				DNS1 string
				DNS2 string
			}{"223.5.5.5", "223.6.6.6"},
			PortScan: struct {
				Thread  int
				Timeout int
			}{2100, 9},
			Proxy: struct {
				Enable   bool
				Mode     string
				Address  string
				Port     int
				Username string
				Password string
			}{false, "HTTP", "127.0.0.1", 8080, "", ""},
			Hunter: struct{ Api string }{""},
			Fofa: struct {
				Email string
				Api   string
			}{"", ""},
			Quake: struct{ Api string }{""},
		}
		b, err := json.MarshalIndent(profile, "", "    ")
		if err != nil {
			logger.Error(err)
		}
		if _, err := os.Stat(defaultConfigDirectory); err != nil {
			os.Mkdir(defaultConfigDirectory, 0755)
		}
		f, err := os.Create(defaultConfigFile)
		if err != nil {
			logger.Error(err)
			return
		}
		defer f.Close()
		_, err = f.WriteString(string(b))
		if err != nil {
			logger.Error(err)
			return
		}
	}
	file, err := os.ReadFile(defaultConfigFile)
	if err != nil {
		logger.Error(err)
		return
	}
	err = json.Unmarshal(file, &Profile)
	if err != nil {
		logger.Error(err)
		return
	}

}

// 10s一次
func FreeMemory() {
	for {
		runtime.GC() //垃圾回收
		debug.FreeOSMemory()
		time.Sleep(10 * time.Second) //
	}
}

// 如果.old文件存在则删除
func init() {
	const oldPocZip = "./config/afrog-pocs.zip"
	currentMain := strings.Split(os.Args[0], "\\")
	dir, _ := os.Getwd()
	oldFile := fmt.Sprintf("%v\\.%v.old", dir, currentMain[len(currentMain)-1:][0])
	if _, err := os.Stat(oldPocZip); err == nil {
		if err2 := os.Remove(oldFile); err2 != nil {
			logger.Error(err)
		}
	}
	if _, err := os.Stat(oldPocZip); err == nil {
		if err2 := os.Remove(oldPocZip); err2 != nil {
			logger.Error(err)
		}
	}
}

func init() {
	app.SetMetadata(fyne.AppMetadata{
		Name:    "hound",
		Version: "0.0.1",
		Custom: map[string]string{
			"Issues": "https://github.com/aicoa/hound/issues/new",
		},
	})
}
