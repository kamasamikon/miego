dockerhelper: 
	GO111MODULE=auto CGO_ENABLED=0 go build -o dockerhelper -tags "nomsgpack" -ldflags "-w -s" 
	upx dockerhelper

clean:
	rm -frv dockerhelper
