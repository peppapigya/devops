package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var buf = make([]byte, 512*1024) // 512KB缓冲区大小

// 文件工具类

// CreateDir 判断目录是否存在
func CreateDir(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return fmt.Errorf("create dir error: %v", err)
	}
	return nil
}

// Copy 文件拷贝
func Copy(dest io.Writer, src io.Reader, buffer []byte) (written int64, err error) {
	if buffer == nil {
		buffer = buf
	}
	return io.CopyBuffer(dest, src, buffer)
}
