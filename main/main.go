package main

import (
	"net/http"
	"oauth_demo/config"
	"oauth_demo/server"
	"time"
)

func main() {

	startTime := time.Now()

	server.Init()

	// auth_server 授权入口
	http.HandleFunc("/oauth2/authorize", server.AuthorizeHandler)

	// auth_server 发现未登录状态, 跳转到的登录handler
	http.HandleFunc("/oauth2/login", server.LoginHandler)

	// auth_server拿到 client以后重定向到的地址, 也就是 auth_client 获取到了code, 准备用code换取accesstoken
	//http.HandleFunc("/oauth2/code_to_token", server.CodeToToken)

	// auth_server 处理由code 换取access token 的handler
	http.HandleFunc("/oauth2/token", server.TokenHandler)

	// 登录完成, 同意授权的页面
	http.HandleFunc("/oauth2/agree-auth", server.AgreeAuthHandler)

	// access token 换取用户信息的handler
	http.HandleFunc("/oauth2/getuserinfo", server.GetUserInfoHandler)

	http.Handle("/", http.FileServer(http.Dir("./static"))) //http://localhost:9000

	errChan := make(chan error)
	go func() {
		config.Log.Info("server start", "duration", time.Now().Sub(startTime))
		errChan <- http.ListenAndServe(":9000", nil)
	}()
	err := <-errChan
	if err != nil {
		config.Log.Error("server stop")
	}

}
