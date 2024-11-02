这个程序运行在远程目标机器的主机上(docker的外边)，返回docker的相关信息。服务通过它获取msb相关的信息。



- 用sudo的方式运行，这样才能保证docker命令可以执行。
- 默认的端口号是11111，可以通过`-port xyz` 的方式修改为 xyz。



接口：

- 获取容器的ID
  - `/container?byPort=xxx`
  - `/container?byName=xxx`
- 获取容器的ID，名称，端口
  - `/info?byPort=xxx`
  - `/info?byName=xxx`



注册为服务：

vim /usr/lib/systemd/system/dockerhelper.service

```ini
[Unit]
Description=dockerhelp
Environment=
After=syslog.target
After=network.target

[Service]
Type=forking
Restart=always
RestartSec=2
StartLimitInterval=0
StartLimitBurst=5

EnvironmentFile=
ExecStart=/usr/local/bin/dockerhelper -p 11111
User=root
ExecStop=/bin/kill -s TERM $MAINPID

[Install]
WantedBy=multi-user.target
```

