这个版本的MSB没有把数据保存到Redis，因为服务不断的注册。保存这些没有意义。

现在只启动了两个进程，Nginx和msb。

`makesrc.sh`是用来编译生成msb的，原来msb是一个包，现在是一个独立的程序。

`makedocker.sh`生成docker，镜像的名字是msb。

`daorun.sh` 这个是从原来的msb中拷贝过来的。就是把镜像msb启动，启动前要删除老的msb，同时生成的容器的名字也叫做msb。

结论就是，依次运行上边三个脚本就可以把msb跑起来了。