# go_oauth2_demo

inde.html 点下button后会请求auth server 进行授权, 也就是尝试获取用户信息。客户端携带client_id, response_type, scope, redirect_uri 请求第三方登录中心的 authorization server 授权

login.html 按button后,发现没有登陆状态, 则会跳转的 login.html页面。让用户输入用户名和密码,  然后将请求发送到/oauth2/login,请求登录

agree-auth.html 当 GitHub 的 /oauth2/login 对应的 handler 接收到传过来的 账户密码 并且校验成功后, 就会跳转到 auth.html 这个页面，等待用户去点击那个 同意授权 button。 这个button和 index.html中的 使用 GitHub 登录的 button 作用一样, 都是请求授权,如果没有发现登陆状态都会跳转回login.html。一般来说, 用户在没有登陆状态的情况下会先经过 index.html, 再经过 auth.html。如果有登录状态, 则在 inde.html不经过 auth.html,直接拿到用户信息用于注册/登录

code-to-user-info.html 接受返回的 code, 也就是请求中的redirect_uri 地址。其内部主要实现了两个功能, 接受code换取access token, 使用 access token 换取用户信息。 特别的是, 我在内部定义了一个 httpRequest 参数, 这个函数 使得 js 自带的 xmlHttpRequest 从一个异步函数强制变成了一个同步函数, 避免了还没有获取到 access token 就请求换取用户信息
