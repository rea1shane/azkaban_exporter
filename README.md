# Azkaban Exporter

## TODO

- 协程写数组改为写 ch
- 去掉调试用的 panic

To see all available configuration flags:

```shell
./azkaban_exporter -h
```

**更改为单个用户**, example: `watcher / watcher`

- 减少请求次数
- 更加安全, 降低主要账号泄露可能
- 降低维护 session 成本
- 省略 project 去重操作

## 可能出现的报错

### connect: can't assign requested address

请求接口过于频繁导致端口释放不及, 从而无法建立新连接. 可能是对 exporter 拉取频率太高, 也可能是有过多的 project 与 flow 导致, 修改 exporter 所在服务器配置:

```shell
# vim /etc/sysctl.conf

sysctl -w net.ipv4.tcp_fin_timeout=30  #修改系統默认的TIMEOUT时间，默认为60s 
sysctl -w net.ipv4.tcp_timestamps=1    #修改tcp/ip协议配置， 通过配置/proc/sys/net/ipv4/tcp_tw_resue, 默认为0，修改为1，释放TIME_WAIT端口给新连接使用
sysctl -w net.ipv4.tcp_tw_recycle=1     #修改tcp/ip协议配置，快速回收socket资源，默认为0，修改为1：
sysctl -w net.ipv4.tcp_tw_reuse = 1     #允许端口重用
```

### invalid character '<' looking for beginning of value