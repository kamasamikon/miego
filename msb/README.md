这个版本的MSB没有把数据保存到Redis，因为服务不断的注册。保存这些没有意义。

现在只启动了两个进程，Nginx和msb。

`makesrc.sh`是用来编译生成msb的，原来msb是一个包，现在是一个独立的程序。

`makedocker.sh`生成docker，镜像的名字是msb。

`daorun.sh` 这个是从原来的msb中拷贝过来的。就是把镜像msb启动，启动前要删除老的msb，同时生成的容器的名字也叫做msb。

结论就是，依次运行上边三个脚本就可以把msb跑起来了。

nginx.conf.full.tmpl 是完整的nginx.conf文件，相当于/etc/nginx/nginx.conf。
nginx.conf.server.tmpl 相当于/etc/nginx/conf.d/目录下的文件。

### 容器模式

- 修改msb.conf
- 检查Dockerfile、msb.conf、nginx.conf.temp，保证这几个是一致的。
- **./makedocker.sh**生成容器。
- **./daorun.sh -d** 生成并运行容器。

### 独立模式

- 复制需要的文件
  - `cp msb.pem /etc/nginx/`
  - `cp msb.key /etc/nginx/`
  - `cp nginx.conf.server.tmpl /etc/nginx/conf.d/msb.conf`
  - `cp nginx.conf.server.tmpl /etc/nginx/conf.d/msb.conf.tmpl`
- 修改msb.conf
- `s:/msb/nginx/conf=/etc/nginx/conf.d/msb.conf`
- `s:/msb/nginx/tmpl=/etc/nginx/conf.d/msb.conf.tmpl`
- `s:/msb/nginx/exec=/usr/sbin/nginx`
- **./makesrc.sh**生成MSB程序。
- 直接运行。

http --verify=no https://127.0.0.1/msb/service
验证服务列表
