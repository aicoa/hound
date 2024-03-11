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
type Unzip struct {}

func NewUnzip() *Unzip {
	return &Unzip{}
}
