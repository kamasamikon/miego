package nationality

var nationalities []string = []string{
	"东乡",
	"乌孜别克",
	"京",
	"仡佬",
	"仫佬",
	"佤",
	"侗",
	"俄罗斯",
	"保安",
	"傣",
	"僳僳",
	"哈尼",
	"哈萨克",
	"回",
	"土家",
	"土",
	"基诺",
	"塔吉克",
	"塔塔尔",
	"壮",
	"布依",
	"布朗",
	"彝",
	"德昂",
	"怒",
	"拉祜",
	"撒拉",
	"普米",
	"景颇",
	"朝鲜",
	"柯尔克孜",
	"毛南",
	"水",
	"汉",
	"满",
	"独龙",
	"珞巴",
	"瑶",
	"畲",
	"白",
	"纳西",
	"维吾尔",
	"羌",
	"苗",
	"蒙古",
	"藏",
	"裕固",
	"赫哲",
	"达斡尔",
	"鄂伦春",
	"鄂温克",
	"锡伯",
	"门巴",
	"阿昌",
	"高山",
	"黎",
}

func List() []string {
	return nationalities
}

func Include(n string) bool {
	for _, x := range nationalities {
		if x == n {
			return true
		}
	}
	return false
}
