package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"os"
	"path/filepath"
)

func GetCenter(x int, y int, img robotgo.Bitmap) (int, int) {
	return x + int(1.0/2.0*float64(img.Width)), y + int(1.0/2.0*float64(img.Height))
}

func DataErr(row int, column int) error {
	return fmt.Errorf("第%d行，第%d列数据出错", row, column)
}

func Quit(msg ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, msg)
	fmt.Println("\n按任意键退出...")

	var key string
	_, _ = fmt.Scanln(&key)
	os.Exit(1)
}

func ResolvePath(src string) string {
	executable, err := os.Executable()
	if err != nil {
		Quit(err)
	}
	return filepath.Join(filepath.Dir(executable), src)
}
