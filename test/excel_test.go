package test

import (
	"encoding/csv"
	"encoding/json"
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

type recommendConfig struct {
	GoodNum  string `json:"goodNum"`
	PopText  string `json:"popText"`
	Priority string `json:"priority"`
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
		m := make(map[string]recommendConfig, 1)
		for rows.Next() {
			columnIndex++
			columns, err := rows.Columns()
			if err != nil {
				config.Log.Errorw("ERROR", err)
				return false
			}
			if columnIndex < 3 {
				continue
			}
			m[columns[2]] = recommendConfig{columns[15], columns[14], columns[6]}
		}
		result, err := json.Marshal(m)
		if err != nil {
			return false
		}
		err = newFile.SetSheetRow(sheetName, "A1", &[]interface{}{string(result)})
		newFile.SetActiveSheet(index)
		break
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
	copyExcelData("C:\\Users\\DELL\\Desktop\\宽带推荐优先级v.2.xlsx")
}
