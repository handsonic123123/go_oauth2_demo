package test

import (
	"encoding/csv"
	"fmt"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/handsonic123123/go_oauth2_demo/config"
	"github.com/xuri/excelize/v2"
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

type Area struct {
	Name string  `json:"name"`
	Code string  `json:"code"`
	Next []*Area `json:"next"`
}

func copyExcelData(dir string) bool {
	originFile, err := excelize.OpenFile(dir)
	if err != nil {
		return false
	}
	defer func() {
		if err := originFile.Close(); err != nil {
			config.Log.Errorw("ERROR", err)
		}
	}()
	root := Area{"根节点", "0", make([]*Area, 0)}
	newFile := excelize.NewFile()
	defer func() {
		if err = newFile.Close(); err != nil {
			config.Log.Errorw("ERROR", err)
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
		//m := make(map[string]recommendConfig, 1)
		for rows.Next() {
			columnIndex++
			columns, err := rows.Columns()
			if err != nil {
				config.Log.Errorw("ERROR", err)
				return false
			}
			if columnIndex < 2 {
				continue
			}
			//m[columns[2]] = recommendConfig{columns[15], columns[14], columns[6]}
			switch len(columns[1]) {
			case len("002034"):
				{
					provinceSlice := append(root.Next, &Area{columns[0], columns[1], make([]*Area, 0)})
					root.Next = provinceSlice
					break
				}
			case len("002034001"):
				{
					for _, province := range root.Next {
						if stringUtils.StartsWith(columns[1], province.Code) {
							citySlice := append(province.Next, &Area{columns[0], columns[1], make([]*Area, 0)})
							province.Next = citySlice
							break
						}
					}
					break
				}
			case len("002034001003"):
				{
					for _, province := range root.Next {
						if stringUtils.StartsWith(columns[1], province.Code) {
							for _, city := range province.Next {
								if stringUtils.StartsWith(columns[1], city.Code) {
									areaSlice := append(city.Next, &Area{columns[0], columns[1], make([]*Area, 0)})
									city.Next = areaSlice
									break
								}
							}
						}
					}
					break
				}

			}
		}
		newSheetIndex := 0
		for _, province := range root.Next {
			newSheetIndex++
			err = newFile.SetSheetRow(sheetName, fmt.Sprintf("%s%d", "A", newSheetIndex), &[]interface{}{fmt.Sprintf(
				"INSERT INTO `t_area_sz_dx_bdwt_sd` ( `F_LEVEL`, `F_PROVINCE_NAME`, `F_PROVINCE_CODE`, `F_CITY_NAME`, `F_CITY_CODE`, `F_COUNTY_NAME`, `F_COUNTY_CODE`, `F_COUNTY_NAME_HIS`) VALUES ('0', '%s', '%s', NULL, NULL, NULL, NULL, NULL);",
				province.Name, province.Code)})
			for _, city := range province.Next {
				newSheetIndex++
				err = newFile.SetSheetRow(sheetName, fmt.Sprintf("%s%d", "A", newSheetIndex), &[]interface{}{fmt.Sprintf(
					"INSERT INTO `t_area_sz_dx_bdwt_sd` ( `F_LEVEL`, `F_PROVINCE_NAME`, `F_PROVINCE_CODE`, `F_CITY_NAME`, `F_CITY_CODE`, `F_COUNTY_NAME`, `F_COUNTY_CODE`, `F_COUNTY_NAME_HIS`) VALUES ( '1', '%s', '%s', '%s', '%s', NULL, NULL, NULL);", province.Name, province.Code, city.Name, city.Code)})
				for _, county := range city.Next {
					newSheetIndex++
					err = newFile.SetSheetRow(sheetName, fmt.Sprintf("%s%d", "A", newSheetIndex), &[]interface{}{fmt.Sprintf(
						"INSERT INTO `t_area_sz_dx_bdwt_sd` ( `F_LEVEL`, `F_PROVINCE_NAME`, `F_PROVINCE_CODE`, `F_CITY_NAME`, `F_CITY_CODE`, `F_COUNTY_NAME`, `F_COUNTY_CODE`, `F_COUNTY_NAME_HIS`) VALUES ( '2', '%s', '%s', '%s', '%s', '%s', '%s', NULL);", province.Name, province.Code, city.Name, city.Code, county.Name, county.Code)})
				}
			}
		}
		//err = newFile.SetSheetRow(sheetName, "A1", &[]interface{}{string(result)})
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
	copyExcelData("C:\\Users\\DELL\\Desktop\\深圳电信地址.xlsx")
}

type Node struct {
	height int
	pre    *Node
	next   *Node
}

func trap(height []int) int {
	length := len(height)
	if length == 0 {
		return 0
	}
	first, last := &Node{height[0], nil, nil}, &Node{height[length-1], nil, nil}
	pre := first
	for _, item := range height[1:] {
		now := &Node{item, pre, nil}
		pre.next = now
		pre = now
	}
	last.pre = pre

	leftMax, rightMax := 0, 0
	water := 0

	for first != last {
		if first.height < last.height {
			if first.height >= leftMax {
				leftMax = first.height
			} else {
				water += leftMax - first.height
			}
			first = first.next
		} else {
			if last.height >= rightMax {
				rightMax = last.height
			} else {
				water += rightMax - last.height
			}
			last = last.pre
		}
	}

	return water
}

func TestTrap(t *testing.T) {

	i := []int{4, 2, 0, 3, 2, 5}

	i2 := trap(i)

	config.Log.Info(i2)

}
