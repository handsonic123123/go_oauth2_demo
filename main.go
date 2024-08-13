package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/handsonic123123/go_oauth2_demo/config"
	pbLog "github.com/handsonic123123/go_oauth2_demo/proto/log"
	"github.com/handsonic123123/go_oauth2_demo/server"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"net/http"
	"time"
)

type log struct{}

func (*log) LogQuery(_ context.Context, in *pbLog.LogQueryReq) (*pbLog.LogQueryResp, error) {
	return &pbLog.LogQueryResp{Id: in.Id, Name: "1", Content: "2"}, nil
}

func main() {

	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		config.Log.Fatalf("Failed to listen:%s", err)
	}
	// Create a gRPC server object
	s := grpc.NewServer()
	pbLog.RegisterLogEventServer(s, &log{})
	config.Log.Infoln("Serving gRPC on 0.0.0.0:8080")
	go func() {
		config.Log.Fatalln(s.Serve(listen))
	}()

	// Create a client connection to the gRPC server we just started
	client, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		config.Log.Fatalf("Failed to dial server:%s", err)
	}
	// Create a new ServeMux for the gRPC-Gateway
	mux := runtime.NewServeMux()
	err = pbLog.RegisterLogEventHandler(context.Background(), mux, client)
	if err != nil {
		config.Log.Fatalf("Failed to register gateway:%s", err)
	}

	// Create a new HTTP server for the gRPC-Gateway
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}

	config.Log.Info("Serving gRPC-Gateway on http://0.0.0.0:8090")
	config.Log.Fatalln(gwServer.ListenAndServe())

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
		config.Log.Info("server start ", "duration：", time.Now().Sub(startTime).Microseconds())
		errChan <- http.ListenAndServe(":9000", nil)
	}()
	err = <-errChan
	if err != nil {
		config.Log.Error("server stop")
	}

}
