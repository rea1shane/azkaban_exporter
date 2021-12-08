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

出现位置为 getExecInfos 函数, 应该是 json.UnMarshal 解析返回数据的时候解析错误.

因为 ids 中的元素出现错误, 有过大的值如 `3543824036068086856` 出现, 返回值出现 html.

可能是由于协程写数组导致.

```html

<html>

<head>
    <meta http-equiv="Content-Type" content="text/html; charset=ISO-8859-1"/>
    <title>Error 500 For input string: "3543824036068086856"</title>
</head>

<body>
<h2>HTTP ERROR 500</h2>
<p>Problem accessing /executor. Reason:
<pre>    For input string: "3543824036068086856"</pre>
</p>
<h3>Caused by:</h3>
<pre>java.lang.NumberFormatException: For input string: "3543824036068086856"
	at java.lang.NumberFormatException.forInputString(NumberFormatException.java:65)
	at java.lang.Integer.parseInt(Integer.java:583)
	at java.lang.Integer.parseInt(Integer.java:615)
	at azkaban.server.HttpRequestUtils.getIntParam(HttpRequestUtils.java:217)
	at azkaban.webapp.servlet.AbstractAzkabanServlet.getIntParam(AbstractAzkabanServlet.java:149)
	at azkaban.webapp.servlet.ExecutorServlet.handleAJAXAction(ExecutorServlet.java:118)
	at azkaban.webapp.servlet.ExecutorServlet.handleGet(ExecutorServlet.java:97)
	at azkaban.webapp.servlet.LoginAbstractAzkabanServlet.doGet(LoginAbstractAzkabanServlet.java:122)
	at javax.servlet.http.HttpServlet.service(HttpServlet.java:668)
	at javax.servlet.http.HttpServlet.service(HttpServlet.java:770)
	at org.mortbay.jetty.servlet.ServletHolder.handle(ServletHolder.java:511)
	at org.mortbay.jetty.servlet.ServletHandler.handle(ServletHandler.java:401)
	at org.mortbay.jetty.servlet.SessionHandler.handle(SessionHandler.java:182)
	at org.mortbay.jetty.handler.ContextHandler.handle(ContextHandler.java:766)
	at org.mortbay.jetty.handler.HandlerWrapper.handle(HandlerWrapper.java:152)
	at org.mortbay.jetty.Server.handle(Server.java:326)
	at org.mortbay.jetty.HttpConnection.handleRequest(HttpConnection.java:542)
	at org.mortbay.jetty.HttpConnection$RequestHandler.headerComplete(HttpConnection.java:928)
	at org.mortbay.jetty.HttpParser.parseNext(HttpParser.java:549)
	at org.mortbay.jetty.HttpParser.parseAvailable(HttpParser.java:212)
	at org.mortbay.jetty.HttpConnection.handle(HttpConnection.java:404)
	at org.mortbay.jetty.bio.SocketConnector$Connection.run(SocketConnector.java:228)
	at org.mortbay.thread.QueuedThreadPool$PoolThread.run(QueuedThreadPool.java:582)
</pre>
<hr/>
<i><small>Powered by Jetty://</small></i><br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>
<br/>

</body>

</html>
```

### request failure when call fetch-a-flow-execution api, reason: Cannot find execution '0'

出现位置为 getExecInfos 函数, 此为 azkaban 返回的错误.

检查发现 exec ids 数组中有编号为 `0` 的 exec id.

不清楚出现原因, 可能是由于协程写数组导致.