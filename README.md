# Azkaban Exporter

## TODO

- 打包为 docker 镜像

## Develop Env

- azkaban: `3.90.0`
- go: `1.17.2`

## Usage

### Setp1

Add a user in azkaban, example: `metrics`

And edit `azkaban.yml`:

```yaml
server:
  protocol: http
  host: 127.0.0.1
  port: 20000
user:
  username: metrics
  password: password
```

### Setp2

Add `read` permissions to the project for this user.

Azkaban exporter monitors all projects that `metrics` have read permission.

### Setp3

Run it, basic usage:

```shell
$ azkaban_exporter --web.listen-address=:9900 --azkaban.conf=azkaban.yml
```


To see all available configuration flags:

```shell
$ azkaban_exporter -h
```

Then you can access `http://127.0.0.1:9900` to view azkaban metrics.

## Why add a user

More security, reduce the risk of major account leakage.

## Grafana Dashboard

grafana version: `8`

dashboard id: `15429`

https://grafana.com/grafana/dashboards/15429

![image](https://raw.githubusercontent.com/rea1shane/azkaban_exporter/feature-http-retry/img/1.png)

![image](https://raw.githubusercontent.com/rea1shane/azkaban_exporter/feature-http-retry/img/2.png)

**NEED config variables azkaban_address, example: http(s)://host:port** 

![image](https://raw.githubusercontent.com/rea1shane/azkaban_exporter/feature-http-retry/img/3.png)

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
