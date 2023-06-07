package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"net/http"
	"oauth_demo/config"
	"os"
	"time"
)

var manager *manage.Manager

var srv *server.Server

// 用户信息结构体
type UserInfo struct {
	Username string `json:"username"`
	Gender   string `json:"gender"`
}

// 用一个 map 存储用户信息
var user_info_map = make(map[string]UserInfo)

func Init() {

	// 设置 client 信息
	client_store := store.NewClientStore()
	client_store.Set("demo", &models.Client{ID: "demo", Secret: "xxxxxx", Domain: "http://localhost:9000"})

	// 设置 manager, manager 参与校验 code/access token 请求
	manager = manage.NewDefaultManager()

	//校验 redirect_uri 和 client 的 Domain, 简单起见, 不做校验
	//manager.SetValidateURIHandler(func(baseURI, redirectURI string) error {
	//	config.Info("ValidateURI", "baseURI", baseURI, "redirectURI", redirectURI)
	//	return nil
	//})

	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// manger 包含 client 信息
	manager.MapClientStorage(client_store)

	// server 也包含 manger, client 信息
	srv = server.NewServer(server.NewConfig(), manager)

	// 根据 client id 从 manager 中获取 client info, 在获取 access token 校验过程中会被用到
	srv.SetClientInfoHandler(func(r *http.Request) (clientID, clientSecret string, err error) {
		client_info, err := srv.Manager.GetClient(r.Context(), r.URL.Query().Get("client_id")) //r.URL.Query().Get("client_id")
		if err != nil {
			config.Info("get client error", "err", err)
			return "", "", err
		}
		return client_info.GetID(), client_info.GetSecret(), nil
	})

	// 设置为 authorization code 模式
	srv.SetAllowedGrantType(oauth2.AuthorizationCode)

	// authorization code 模式,  第一步获取code,然后再用code换取 access token, 而不是直接获取 access token
	srv.SetAllowedResponseType(oauth2.Code)

	// 校验授权请求用户的handler, 会重定向到 登陆页面, 返回"", nil
	srv.SetUserAuthorizationHandler(userAuthorizationHandler)

	// 校验授权请求的用户的账号密码, 给 LoginHandler 使用, 简单起见, 只允许一个用户授权
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (userID string, err error) {
		if username == "admin" && password == "admin" {
			return "0001", nil
		}
		return "", errors.New("username or password error")
	})

	// 允许使用 get 方法请求授权
	srv.SetAllowGetAccessRequest(true)

	// 储存用户信息的一个 map
	user_info_map["0001"] = UserInfo{
		"admin", "Male",
	}

}

// 授权入口, demo.html 和 agree-auth.html 按下 button 后
func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleAuthorizeRequest(w, r)
	if err != nil {
		config.Info("authorize fail", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// AuthorizeHandler 内部使用, 用于查看是否有登陆状态
func userAuthorizationHandler(w http.ResponseWriter, r *http.Request) (user_id string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		config.Info("userAuthorization fail", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uid, ok := store.Get("LoggedInUserId")
	// 如果没有查询到登陆状态, 则跳转到 登陆页面
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return "", nil
	}
	// 若有登录状态, 返回 user id
	user_id = uid.(string)
	return user_id, nil
}

// 登录页面的handler
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		config.Info("session start fail", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		user_id, err := srv.PasswordAuthorizationHandler(r.Context(), "demo", r.Form.Get("username"), r.Form.Get("password"))
		if err != nil {
			config.Info("password authorization fail", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		store.Set("LoggedInUserId", user_id) // 保存登录状态
		store.Save()

		// 跳转到 同意授权页面
		w.Header().Set("Location", "/oauth2/agree-auth")
		w.WriteHeader(http.StatusFound)
		return
	}

	// 若请求方法错误, 提供login.html页面
	outputHTML(w, r, "static/login.html")
}

// 若发现登录状态则提供 agree-auth.html, 否则跳转到 登陆页面
func AgreeAuthHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		config.Info("agree auth fail", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 如果没有查询到登陆状态, 则跳转到 登陆页面
	if _, ok := store.Get("LoggedInUserId"); !ok {
		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	// 如果有登陆状态, 会跳转到 确认授权页面
	outputHTML(w, r, "static/agree-auth.html")
}

// code 换取 access token
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	err := srv.HandleTokenRequest(w, r)
	if err != nil {
		config.Info("get access token fail", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// access token 换取用户信息
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 获取 access token
	access_token, ok := srv.BearerAuth(r)
	if !ok {
		config.Info("Failed to get access token from request")
		return
	}

	root_ctx := context.Background()
	ctx, cancle_func := context.WithTimeout(root_ctx, time.Second)
	defer cancle_func()

	// 从 access token 中获取 信息
	token_info, err := srv.Manager.LoadAccessToken(ctx, access_token)
	if err != nil {
		config.Info("load access token fail", "error", err)
		return
	}

	// 获取 user id
	user_id := token_info.GetUserID()
	grant_scope := token_info.GetScope()

	user_info := UserInfo{}

	// 根据 grant scope 决定获取哪些用户信息
	if grant_scope != "read_user_info" {
		config.Info(`invalid grant scope`)
		w.Write([]byte("invalid grant scope"))
		return
	}

	user_info = user_info_map[user_id]
	resp, err := json.Marshal(user_info)
	w.Write(resp)
	return
}

// 提供 HTML 文件显示
func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		config.Info("out put error", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
