package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/xuri/excelize/v2"
	"strconv"
)

var (
	operations = []func(Step) Operation{
		NewLeftClick,
		NewDoubleClick,
		NewRightClick,
		NewInput,
		NewWait,
		NewScroll,
	}
	bitmapCache = make(map[string]robotgo.Bitmap)
	isRetina    = false
)

func init() {
	robotgo.MouseSleep = 10
}

func main() {
	f, err := excelize.OpenFile(ResolvePath("cmd.xlsx"))
	if err != nil {
		if errors.Is(err, zip.ErrFormat) {
			Quit("文件格式错误，请使用xlsx文件（不支持xls）")
		}
		Quit("未找到cmd.xls", err)
	}
	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		Quit("sheet不存在")
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		Quit(err)
	}

	tableHead := rows[0]
	if len(tableHead) >= 4 && tableHead[3] == "retina" {
		isRetina = true
		fmt.Println("正在使用Retina配置")
	}
	for rowNumber, row := range rows[1:] {
		operation := parseRow(rowNumber, row)

		repeatTimes := 1
		if len(row) == 3 {
			repeatTimes, err = strconv.Atoi(row[2])
			if err != nil {
				Quit(DataErr(rowNumber+2, 3))
			}
		}

		// TODO 可用性存疑
		if repeatTimes == -1 {
			for {
				operation.Perform()
			}
		}
		for i := 0; i < repeatTimes; i++ {
			operation.Perform()
		}
	}
}

func parseRow(rowNumber int, row []string) Operation {
	if len(row) != 2 && len(row) != 3 {
		Quit(fmt.Sprintf("第%d行数据出错:有%d列", rowNumber+2, len(row)))
	}
	operationType, err := strconv.Atoi(row[0])
	if err != nil {
		Quit(DataErr(rowNumber+2, 1))
	}

	constructor := operations[operationType-1]

	return constructor(Step{
		Row:       rowNumber,
		Operation: row[1],
	})
}
