package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/sub", handler)         // 设置路由处理函数
	err := http.ListenAndServe(":8888", nil) // 启动服务器并监听端口
	if err != nil {
		fmt.Println("服务器启动失败:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("v2rayConfig.txt") // 读取文本文件内容
	if err != nil {
		fmt.Println("读取文件出错:", err)
		return
	}

	text := string(content) // 将字节切片转换为字符串
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	fmt.Fprint(w, encoded) // 在响应中写入内容
}
