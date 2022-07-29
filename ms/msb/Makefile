clean:
	rm -f msb

msb: *.go
	CGO_ENABLED=0 go build -ldflags "-w -s" msb.go
	upx msb

docker: msb
	date > now
	-sudo docker rm -f msb
	-sudo docker rmi -f msb
	sudo docker build -t msb .

dockerrun: docker
	sudo docker run --restart=always -it -d --name msb -p 9080:80 -p 1443:443 msb

msbtar: msb.tar.gz
	REMOTE=root@117.50.95.113
	sudo docker save -o msb.tar msb
	sudo chmod a+rwx msb.tar
	sudo gzip --best -c msb.tar >> msb.tar.gz

deploy: msbtar
	sudo scp msb.tar.gz $REMOTE:/tmp/msb.tar.gz
	echo ""
	echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker rm -f msb"
	ssh $REMOTE "docker rm -f msb"
	echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker load -i /tmp/msb.tar"
	ssh $REMOTE "gunzip /tmp/msb.tar.gz"
	ssh $REMOTE "docker load -i /tmp/msb.tar"
	echo ">>>>>>>>>>>>>>" ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 msb"
	ssh $REMOTE "docker run -d -it --restart=always --name msb -p 9080:80 -p 9443:443 msb"

