package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"os"
	"strconv"
	"time"
)

func MoveToImg(filename string) {
	bitmap, ok := bitmapCache[filename]
	filePath := ResolvePath(filename)
	if !ok {
		f, err := os.Open(filePath)
		if err != nil {
			Quit(fmt.Sprintf("读取文件%s出错,%s", filename, err.Error()))
		}
		_ = f.Close()
		cb := robotgo.OpenBitmap(filePath)
		bitmap = robotgo.ToBitmap(cb)
		bitmapCache[filename] = bitmap
	}
	cb := robotgo.ToCBitmap(bitmap)
	var x, y int
	for {
		x, y = robotgo.FindBitmap(cb, nil, 0.1)
		if x != -1 && y != -1 {
			break
		}
		// 感觉暴力轮询不太好，设个延迟
		time.Sleep(200 * time.Millisecond)
	}
	x, y = GetCenter(x, y, bitmap)
	// retina屏幕分辨率处理
	if isRetina {
		x, y = x/2, y/2
	}
	robotgo.Move(x, y)
	time.Sleep(50 * time.Millisecond)
}

type Operation interface {
	Perform()
}

type Step struct {
	Row       int
	Operation string
}

type LeftClick struct {
	Step
}

func (l *LeftClick) Perform() {
	MoveToImg(l.Operation)
	robotgo.Click("left")
}

func NewLeftClick(step Step) Operation {
	return &LeftClick{Step: step}
}

type DoubleClick struct {
	Step
}

func (d DoubleClick) Perform() {
	MoveToImg(d.Operation)
	robotgo.Click("left", true)
}

func NewDoubleClick(step Step) Operation {
	return &DoubleClick{Step: step}
}

type RightClick struct {
	Step
}

func (r RightClick) Perform() {
	MoveToImg(r.Operation)
	robotgo.Click("right")
}

func NewRightClick(step Step) Operation {
	return &RightClick{Step: step}
}

type Input struct {
	Step
}

func (i Input) Perform() {
	robotgo.TypeStr(i.Operation)
}

func NewInput(step Step) Operation {
	return &Input{Step: step}
}

type Wait struct {
	Step
}

func (w Wait) Perform() {
	delay, err := strconv.ParseFloat(w.Operation, 64)
	if err != nil {
		Quit(fmt.Sprintf("第%d行操作不是数字", w.Row))
	}
	time.Sleep(time.Duration(delay * float64(time.Second)))
}

func NewWait(step Step) Operation {
	return &Wait{Step: step}
}

type Scroll struct {
	Step
}

func (s Scroll) Perform() {
	toY, err := strconv.Atoi(s.Operation)
	if err != nil {
		Quit(fmt.Sprintf("第%d行操作内容不是整数 ", s.Row))
	}

	// TODO 有待测试,可能换成relative
	robotgo.ScrollSmooth(toY, 1, 50)
}

func NewScroll(step Step) Operation {
	return &Scroll{Step: step}
}
