package socketfile

import (
	"fmt"
	"io/ioutil"
	"miego/klog"
	"miego/xerror"
	"net"
)

func Read(path string) ([]byte, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, xerror.Newf(err, "path=%v", path)
	}
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, xerror.Newf(err, "path=%v", path)
	}

	return data, nil
}

// 通过这个导出信息: Serv隐含信息是这个是一个监听的状态
func Serv(path string, data []byte) error {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return xerror.Newf(err, "path: %s", path)
	}

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				klog.E("listener.Accept: socketPath:%s, err:%v", path, err)
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()
				if _, err := conn.Write([]byte(data)); err != nil {
					fmt.Printf("Failed to send data to client: %v\n", err)
				}
			}(conn)
		}
	}()

	return nil
}
