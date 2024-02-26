![tls_check](https://github.com/gitlayzer/tls_check/assets/77761224/fc081129-5a1b-42f0-81db-f36c0306a29d)

TLS_Check 是一个用于检测 HTTPS 证书过期时间的工具，它支持如下的参数

```shell
  -a string
        This parameter is used to specify the address of the HTTP service
  -l string
        Please enter a valid HTTPS link, multiple domains are separated by commas
  -p string
        This parameter is used to specify the port of the HTTP service
  -t int
        Please enter a valid timeout, the unit is seconds, the default is 5 seconds (default 5)
  -w bool
        This parameter is used to enable web，default is false，default port is 8080
        
当你仅使用 -l 参数时，它就是一个 CLI 工具，示例如下：
❯ .\tlscheckctl -l www.baidu.com
{"domain":"www.baidu.com","subject":"baidu.com","expires_on":"2024-08-06","days_left":162}

当你使用 -t 参数时，它可以指定超时的时间（此参数可直接配置给 -w 参数，这样 web 接口的超时时间就是此参数），示例如下：
❯ .\tlscheckctl -t 10 -l www.google.com
2024/02/26 00:25:34 now Time: 2024-02-26 00:25:34, Dial www.google.com:443 error: dial tcp 108.160.161.20:443: i/o timeout

当你想使用 -w 时，你可以启用 -a，-p 参数分别配置 web 监听的地址与端口，示例如下：
Administrator in E:\codes\tls_check 10s 
❯ .\tlscheckctl.exe -w -a 0.0.0.0 -p 80

# 模拟请求接口
❯ curl.exe 127.0.0.1/check?domain=www.baidu.com
{"domain":"www.baidu.com","subject":"baidu.com","expires_on":"2024-08-06","days_left":162}
```

