package util

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"
)

func BufferWriteAppend(filename, content string) error {
	fileHandle, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	writer := bufio.NewWriter(fileHandle)
	if _, err := writer.WriteString(content + "\n"); err != nil {
		return err
	}

	return writer.Flush()
}
func OpenFolder(path string) error {
	var command string
	switch runtime.GOOS {
	case "linux":
		command = "xdg-open"
	case "windows":
		command = "explorer"
	case "darwin":
		command = "open"
	}

	cmd := exec.Command(command, path)
	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait() // 等待命令完成
}
