package util

/*
 * @Author: aicoa
 * @Date: 2024-03-10 22:54:00
 * @Last Modified by:   aicoa
 * @Last Modified time: 2024-03-10 22:54:00
 */
// reference: https://raw.githubusercontent.com/artdarek/go-unzip/master/pkg/unzip/unzip.go
import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Unzip struct{}

func NewUnzip() *Unzip {
	return &Unzip{}
}

func (uz Unzip) Extract(source, des string) ([]string, error) {
	r, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(des, 0755)
	if err != nil {
		return nil, err
	}

	var extractedFiles []string

	for _, f := range r.File {
		err := uz.extractAndWriteFile(des, f)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, f.Name)
	}
	return extractedFiles, nil
}

func (Unzip) extractAndWriteFile(des string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	path := filepath.Join(des, f.Name)
	if !strings.HasPrefix(path, filepath.Clean(des)+string(os.PathListSeparator)) {
		return fmt.Errorf("%s: 非法文件路径", path)
	}
	if f.FileInfo().IsDir() {
		err = os.MkdirAll(path, f.Mode())
		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(path), f.Mode())
		if err != nil {
			return err
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
