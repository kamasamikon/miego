package common

// KService : Micro Service definition
type KService struct {
	ServiceName string `json:"serviceName"` // 服务基本信息：名字
	Version     string `json:"version"`     // 服务基本信息：版本号，一般是v1，但基本没有用
	Desc        string `json:"desc"`        // 服务基本信息：一句话说明
	ProjName    string `json:"projName"`    // 源码：的目录名
	ProjVersion string `json:"projVersion"` // 源码：的Git版本
	ProjTime    string `json:"projTime"`    // 源码：编译的时间
	IPAddr      string `json:"ipAddr"`      // 服务：所在机器的IP地址
	Port        int    `json:"port"`        // 服务：监听的地址
	HostName    string `json:"hostName"`    // 服务：所在容器的主机名
	CreatedAt   int64  `json:"createdAt"`   // 服务：服务启动的时间
	Kind        string `json:"kind"`        // Nginx：服务类型，http或者grpc。
	Upstream    string `json:"upstream"`    // Nginx：proxy_pass对应的upstream，若为空，则和服务名相同
	RegInterval int64  `json:"regInterval"` // 两次注册之间的时间，单位是秒
}
