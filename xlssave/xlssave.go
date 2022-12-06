package xlssave

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xmap"
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

// 查询并返回为Excel
func Export(c *gin.Context, doQuery func(xmap.Map) []xmap.Map, columns []string) {
	mp := xmap.MapQuery(c, true)

	mp.Put("PageSize", "10000000")
	mp.Put("PageNumber", "1")

	// 获取数据
	records := doQuery(mp)

	// 生成临时文件
	temp, err := ioutil.TempFile("/tmp", "rbvision")
	if err != nil {
		klog.E(err.Error())
		c.Redirect(200, "/")
		return
	}
	tmpName := temp.Name()
	defer os.Remove(tmpName)

	// 把数据保存到临时文件
	Save(tmpName, records, columns)
	content, _ := ioutil.ReadFile(tmpName)

	// 生成文件名
	fileName := mp.Str("fileName", "导出结果")
	fileName += ".xlsx"

	// 写入Excel文件
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write([]byte(content))
}
