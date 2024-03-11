/*
 * @Author: aicoa
 * @Date: 2024-03-08 21:31:02
 * @Last Modified by: aicoa
 * @Last Modified time: 2024-03-08 22:05:48
 */
package common

import (
	"bufio"
	"hound/common/logger"
	"net/url"
	"os"
	"path"
	"strings"

	"fyne.io/fyne/v2/widget"
)

const (
	Mode_Other = iota
	Mode_Url
)

// 读取文件，读成每行的目标
func ParseFile(filePath string) (targets []string) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("[ERRO]" + filePath + " open failed")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" { //去除空行
			targets = append(targets, scanner.Text())
		}
	}
	return targets
}

// 这里将text文本转换为扫描目标对象，mode=1表示目标是url
func ParseTargets(text string, mode int) []string {
	var tmp, targets []string
	tmp = strings.Split(text, "\n") //一般字典格式
	if mode == 1 {
		for _, v := range tmp {
			v = strings.ReplaceAll(strings.ReplaceAll(v, "\r", ""), "", "")
			if v != "" {
				v = strings.TrimSuffix(v, "/") //如果末尾是/结尾则有必要删除/
				targets = append(targets, v)
			}
		}
	} else {
		for _, v := range tmp {
			v = strings.ReplaceAll(strings.ReplaceAll(v, "\r", ""), "", "")
			if v != "" {
				targets = append(targets, v)
			}
		}
	}
	return targets
}

func ParseURLWithoutSlash(text string) (string, error) {
	if _, err := url.ParseRequestURI(text); err != nil {
		return "", err
	}
	text = strings.ReplaceAll(text, " ", "")
	if strings.HasSuffix(text, "/") {
		return text, nil
	} else {
		return text + "/", nil
	}
}

// 分辨参数是否为.txt字典or单个用户名
func ParseDict(this *widget.Entry, defaulsts []string) []string {
	var dicts []string
	if this.Text != "" && !this.Disabled() { //使用自定义字典
		if path.Ext(this.Text) == ".txt" { //如果用户输入的是.txt字典
			dicts = append(dicts, ParseFile(this.Text)...)
		} else {
			dicts = append(dicts, this.Text)
		}
	} else {
		dicts = defaulsts //使用默认内置字典
	}
	return dicts

}

func ParseDirectoryDict(filepath, old string, new []string) (dict []string) {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Error(err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	for s.Scan() {
		if s.Text() != "" {
			if len(new) > 0 {
				if strings.Contains(s.Text(), old) { //如果新数组不为空，将old字段替换成new数组
					for _, v := range new {
						dict = append(dict, strings.ReplaceAll(s.Text(), old, v))
					}

				} else {
					dict = append(dict, s.Text())
				}

			} else {
				if !strings.Contains(s.Text(), old) {
					dict = append(dict, s.Text())
				}
			}
		}

	}
	return dict
}
