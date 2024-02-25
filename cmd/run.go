package cmd

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// init 函数，初始化参数
func init() {
	// 传递参数 -l 指定域名
	flag.StringVar(&Domain, "l", "", "Please enter a valid HTTPS link, multiple domains are separated by commas")
	// 支持传递参数 -t 指定超时时间,单位为秒
	flag.IntVar(&TimeOut, "t", 5, "Please enter a valid timeout, the unit is seconds, the default is 5 seconds")
	// 支持传递参数 -w 指定是否开启 HTTP 服务
	flag.BoolVar(&Weh, "w", false, "This parameter is used to enable web，default is false，default port is 8080")
	// 支持传递参数 -a 指定 HTTP 服务的地址
	flag.StringVar(&Address, "a", "0.0.0.0", "This parameter is used to specify the address of the HTTP service")
	// 支持传递参数 -p 指定 HTTP 服务的端口
	flag.StringVar(&Port, "p", "8080", "This parameter is used to specify the port of the HTTP service")
	// 解析传递的参数到 Domain 这个变量
	flag.Parse()
}

// checkCertExpiration 函数，传递一个域名和超时时间，返回一个切片和错误信息
func checkCertExpiration(d string, t time.Duration) ([]byte, error) {
	// 创建 TCP 连接
	conn, err := net.DialTimeout("tcp", d+":443", t*time.Second)
	if err != nil {
		return nil, err
	}
	// 函数执行完后关闭连接
	defer conn.Close()

	// 配置 TLS 的参数，ServerName 为域名，也就是我们调用函数时传递的参数
	config := &tls.Config{
		ServerName: d,
	}

	// 创建一个 TLS 的连接
	tlsConn := tls.Client(conn, config)
	// 函数执行完后关闭连接
	defer tlsConn.Close()

	// 创建一个 TLS 的握手
	err = tlsConn.Handshake()
	if err != nil {
		return nil, err
	}

	// 获取证书信息，返回的是一个切片
	certs := tlsConn.ConnectionState().PeerCertificates
	for _, cert := range certs {
		info := CertInfo{
			Domain:    d,
			Subject:   cert.Subject.CommonName,
			ExpiresOn: cert.NotAfter.Format("2006-01-02"),
			DaysLeft:  int(cert.NotAfter.Sub(time.Now()).Hours() / 24),
		}
		// 将结构体转换为 JSON 格式
		return json.Marshal(info)
	}

	return nil, nil
}

// handleCheckCertExpiration 函数，处理 HTTP 请求
func handleCheckCertExpiration(w http.ResponseWriter, r *http.Request) {
	// 获取 Query 的 domain 参数
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Please provide a domain", http.StatusBadRequest)
		return
	}

	// 执行函数并获取返回值
	data, err := checkCertExpiration(domain, time.Duration(TimeOut))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	// 写入数据
	w.Write(data)
}

// Run 函数，执行函数
func Run() {
	// 判断是否开启 HTTP 服务
	if Weh {
		http.HandleFunc("/check", handleCheckCertExpiration)
		err := http.ListenAndServe(Address+":"+Port, nil)
		if err != nil {
			return
		}
		return
	}

	// 如果传递的参数为空，打印提示信息
	if Domain == "" {
		fmt.Println("Please enter a valid HTTPS link")
		return
	}

	// 执行函数并获取返回值
	data, err := checkCertExpiration(Domain, time.Duration(TimeOut))
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(data))
	}
}
