# Azkaban Exporter

To see all available configuration flags:

```shell
./azkaban_exporter -h
```

**更改为单个用户**, example: `watcher / watcher`

- 减少请求次数
- 更加安全, 降低主要账号泄露可能
- 降低维护 session 成本
- 省略 project 去重操作