## Cookie

### 环境变量设置

#### cookie.HttpOnly

浏览器不要在除HTTP（和 HTTPS)请求之外暴露Cookie。默认为true，可选false关闭

#### cookie.Secure

仅在https层面上安全传输。默认为true，采用http可以设置false

#### cookie.Path

生效路径。默认为/，即当前域名全部路径。可以设置指定路径，如/user

#### cookie.Domain

生效域名，默认为空即当前域名，不包含子域名。可以设置.domain.com允许子域名

#### cookie.MaxAge

生命周期秒数。默认一天即86400秒，可以设置其他数值

