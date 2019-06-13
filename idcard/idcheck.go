package idcard

import (
	"fmt"
	"strconv"
	"time"
)

func byte2int(x byte) byte {
	if x == 88 || x == 120 {
		return 'X'
	}
	return (x - 48)
}

func idCheck(id []byte) int {
	arry := make([]int, 17)

	for index, value := range id {
		arry[index], _ = strconv.Atoi(string(value))
	}

	var wi [17]int = [...]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	var res int
	for i := 0; i < 17; i++ {
		res += arry[i] * wi[i]
	}

	return (res % 11)
}

func idVerify(verify int, id_v byte) bool {
	var temp byte
	a18 := [11]byte{1, 0, 'X', 9, 8, 7, 6, 5, 4, 3, 2}

	for i := 0; i < 11; i++ {
		if i == verify {
			temp = a18[i]
			break
		}
	}

	return temp == id_v
}

type KIDCardInfo struct {
	Province     int
	City         int
	Distict      int
	Year         int
	Month        int
	Day          int
	SerialNumber int
	CheckNumber  string
}

type KIDCardString string

func (idstring KIDCardString) Parse() (*KIDCardInfo, bool) {
	provines := map[int]string{
		11: "北京",
		12: "天津",
		13: "河北",
		14: "山西",
		15: "内蒙古",
		21: "辽宁",
		22: "吉林",
		23: "黑龙江",
		31: "上海",
		32: "江苏",
		33: "浙江",
		34: "安徽",
		35: "福建",
		36: "江西",
		37: "山东",
		41: "河南",
		42: "湖北",
		43: "湖南",
		44: "广东",
		45: "广西",
		46: "海南",
		50: "重庆",
		51: "四川",
		52: "贵州",
		53: "云南",
		54: "西藏",
		61: "陕西",
		62: "甘肃",
		63: "青海",
		64: "宁夏",
		65: "新疆",
		71: "台湾",
		81: "香港",
		82: "澳门",
		91: "国外",
	}

	// Check Length
	if len(idstring) != 18 {
		return nil, false
	}

	// Check date
	date := string(idstring[6:14])
	fmt.Println(date)
	if _, err := time.Parse("20060102", date); err != nil {
		fmt.Println(err.Error())
		return nil, false
	}

	// Check whole
	var idCardByte [18]byte
	for k, v := range []byte(idstring) {
		idCardByte[k] = byte(v)
	}
	valid := idVerify(idCheck(idCardByte[0:17]), byte2int(idCardByte[17]))
	if !valid {
		return nil, false
	}

	// Check Province
	ok := false
	Province, _ := strconv.Atoi(string(idstring[0:2]))
	for p, _ := range provines {
		if p == Province {
			ok = true
			break
		}
	}
	if !ok {
		return nil, false
	}

	City, _ := strconv.Atoi(string(idstring[2:4]))
	Distict, _ := strconv.Atoi(string(idstring[4:6]))
	Year, _ := strconv.Atoi(string(idstring[6:10]))
	Month, _ := strconv.Atoi(string(idstring[10:12]))
	Day, _ := strconv.Atoi(string(idstring[12:14]))
	SerialNumber, _ := strconv.Atoi(string(idstring[14:17]))
	CheckNumber := string(idstring[17:])

	ci := KIDCardInfo{}

	ci.Province = Province
	ci.City = City
	ci.Distict = Distict
	ci.Year = Year
	ci.Month = Month
	ci.Day = Day
	ci.SerialNumber = SerialNumber

	if CheckNumber == "x" {
		CheckNumber = "X"
	}
	ci.CheckNumber = CheckNumber

	return &ci, true
}
