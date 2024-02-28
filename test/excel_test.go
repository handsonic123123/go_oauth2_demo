package test

import (
	"encoding/csv"
	"flag"
	"github.com/xuri/excelize/v2"
	"oauth_demo/config"
	"os"
	"testing"
)

func getExcelData(dir string) {
	f, err := excelize.OpenFile(dir)
	if err != nil {
		config.Log.Errorw("ERROR", err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			config.Log.Errorw("ERROR", err)
		}
	}()
	for _, name := range f.GetSheetMap() {
		config.Log.Info("find sheet:{}", name)
		WriterCSV := csv.NewWriter(os.Stdout)

		rows, err := f.Rows(name)
		if err != nil {
			config.Log.Errorw("ERROR", err)
			return
		}
		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				config.Log.Errorw("ERROR", err)
			}
			csverr := WriterCSV.Write(row)
			if csverr != nil {
				return
			}
		}
		WriterCSV.Flush()
	}
}

func copyExcelData(dir string) bool {
	originFile, err := excelize.OpenFile(dir)
	if err != nil {
		return false
	}
	defer func() {
		if err := originFile.Close(); err != nil {
			println(err)
		}
	}()
	newFile := excelize.NewFile()
	defer func() {
		if err = newFile.Close(); err != nil {
			println(err)
		}
	}()
	for _, sheetName := range originFile.GetSheetMap() {
		config.Log.Infof("find sheet:%v", sheetName)
		index, _ := newFile.NewSheet(sheetName)
		rows, err := originFile.Rows(sheetName)
		if err != nil {
			config.Log.Errorw("ERROR", err)
			return false
		}
		columnIndex := 0
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				config.Log.Errorw("ERROR", err)
				return false
			}
			if err != nil {
				config.Log.Errorw("ERROR", err)
				return false
			}
			columnIndex++
			name, err := excelize.CoordinatesToCellName(1, columnIndex)
			if err != nil {
				config.Log.Errorw("ERROR", err)
				return false
			}
			err = newFile.SetSheetRow(sheetName, name, &columns)
			if err != nil {
				config.Log.Errorw("ERROR", "errMsg", err, "key", "行处理出错")
				return false
			}
		}
		newFile.SetActiveSheet(index)
	}
	err = newFile.DeleteSheet("sheet1")
	if err != nil {
		config.Log.Errorw("ERROR", "删除出错", err)
		return false
	}
	err = os.Mkdir("C:\\Users\\DELL\\Desktop\\temp", os.ModePerm)
	if err != nil {
		config.Log.Errorw("ERROR", "保存出错", err)
		return false
	}
	filePath := "C:\\Users\\DELL\\Desktop\\temp\\test.xlsx"
	err = newFile.SaveAs(filePath)
	if err != nil {
		config.Log.Errorw("ERROR", "保存出错", err)
		return false
	}
	err = config.OpenFile(filePath)
	if err != nil {
		config.Log.Errorw("ERROR", "文件打开失败", err)
		return false
	}
	return false
}

func TestExcel(t *testing.T) {
	flag.Parse()
	dir := flag.Args()
	if len(dir) == 0 {
		config.Log.Errorw("DirRequired", "DirRequired", "")
		return
	}
	copyExcelData(dir[0])
}
