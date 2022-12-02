package xlssave

import (
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"

	"github.com/tealeg/xlsx"
)

// Save *records to *path with *columns define
// columns: xmapKey:showName, 列标题和列数据对应的Key
func Save(path string, records []xmap.Map, columns []string) error {
	addRow := func(sheet *xlsx.Sheet, cells []string) {
		row := sheet.AddRow()
		row.SetHeightCM(1)
		for _, c := range cells {
			row.AddCell().Value = c
		}
	}

	file := xlsx.NewFile()

	var showNames []string
	var xmapKeys []string

	if columns != nil {
		for i := 0; i < len(columns)/2; i++ {
			xmapKey := columns[2*i]
			showName := columns[2*i+1]
			if showName == "" {
				showName = xmapKey
			}

			showNames = append(showNames, showName)
			xmapKeys = append(xmapKeys, xmapKey)
		}
	} else if records != nil {
		record := records[0]
		for k := range record {
			showNames = append(showNames, k)
			xmapKeys = append(xmapKeys, k)
		}
	}

	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		klog.E("%s", err.Error())
		return err
	}

	// 写入列标题
	addRow(sheet, showNames)

	// 写入数据
	for _, record := range records {
		var x []string
		for _, key := range xmapKeys {
			x = append(x, record.S(key))
		}
		addRow(sheet, x)
	}

	return file.Save(path)
}
