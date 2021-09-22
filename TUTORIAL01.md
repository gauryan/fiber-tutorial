시작 (설치 및 첫페이지 만들어보기) : Fiber v2.x Tutorial
========================================================

본 튜토리얼에서는 대부분의 사이트에 기본적으로 필요한 기능을 익히는데 목적이 있습니다. 예제에서는 관리자로 로그인/로그아웃하고, 관리자를 등록/수정/삭제 등을 할 수 있는 기본적인 기능을 포함하는 예제를 만들어봅니다. 그리고, 비밀번호의 단방향 암호화(SHA-256)를 해보는 기능도 포함됩니다. DB 연결시 GORM을 사용하지만, 제공되는 Model 메소드(ORM)를 사용하지 않고, 직접 쿼리(스토어드 프로시저)를 사용하여 처리할 것입니다.

OS(Ubuntu Linux 20.04)계정은 기본계정인 ubuntu 를 사용하는 것으로 가정한다. 사용할 프로젝트 디렉토리는 **`~/project/xyz`** 로 될 것이다. Go 언어도 설치되어 있다고 가정하겠습니다. 참고로 저는 1.17.1 버전이 설치되어 있습니다.


### 1. 프로젝트 디렉토리 구조
```
xyz
├-- controllers
├-- database
├-- routes
└-- views
```
Fiber는 디렉토리 구조를 강제하지 않기 때문에 마음대로 만드셔도 됩니다. 위의 구조는 제 마음대로 그냥 정한 것입니다.


### 2. 서비스에 사용할 웹서버(nginx)를 설치하고 설정 해보자.
```bash
$ sudo apt-get install nginx
```
```bash
$ cd /etc/nginx/sites-available
$ sudo vi xyz
```
```
server { 
	listen 80; 
	server_name xyz.test.com; # 자신이 원하는 도메인주소 입력 

	location / { 
		proxy_pass http://localhost:3000; 
		proxy_http_version 1.1; 
		proxy_set_header Upgrade $http_upgrade; 
		proxy_set_header Connection 'upgrade'; 
		proxy_set_header Host $host; 
		proxy_set_header X-Real-IP $remote_addr; 
		proxy_set_header X-Forwarded-Proto $scheme; 
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for; 
		proxy_cache_bypass $http_upgrade;
	}
}
```
```bash
$ cd ../sites-enabled
$ sudo ln -s /etc/nginx/sites-available/xyz xyz
```

### 3. nginx 재시작
```bash
$ sudo /etc/init.d/nginx restart
```

### 4. 작업 디렉토리를 만들고, 모듈 초기화하기
```bash
$ cd ~
$ mkdir -p project/xyz
$ cd project/xyz
$ go mod init github.com/gauryan/xyz
go: creating new go.mod: module github.com/gauryan/xyz
$
```
모듈이름은 자신의 원하는 것으로 하세요~! 단, 이후 import 할 때, 주의하세요.

### 5. 하위 폴더 만들고, 시작 스크립트 만들기
```bash
$ mkdir controllers; mkdir database; mkdir routes; mkdir views
$ vi run.sh
go run main.go
 
$ chmod 755 run.sh
```

### 6. 이제, 첫화면을 위한 controller를 만들어볼까요? `~/project/xyz/controllers` 디렉토리에 `main.go` 를 생성합니다.
```go
package controllers
 
import (
    "github.com/gofiber/fiber/v2"
)
 
func Index(c *fiber.Ctx) error {
    data := fiber.Map{ "Title": "Hello, World!", }
    return c.Render("index", data)
}
```
Index 함수는 http://xyz.test.com 으로 접속할 때, 연결되는 것입니다. 이후에 `index.html` 을 열어서 보여주게 됩니다.

### 7. 그러면, `~/project/xyz/views/index.html` 을 만듭니다.
```html
<!DOCTYPE html>
<html>
<body>
  <h1>{{.Title}}</h1>
</body>
</html>
```
`.Title`은 컨트롤러에서 넘겨받은 데이터 `“Hello, World!”` 입니다.

### 8. URL과 컨트롤러를 연결해주는 Router를 작성합니다. 작성할 파일은 `~/project/xyz/routes/web.go` 입니다.
```go
package routes
 
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
    "github.com/gauryan/xyz/controllers"
)
 
 
func Router() *fiber.App {
    // App 생성과 템플릿 설정
    app := fiber.New(fiber.Config{
        Views: html.New("./views", ".html"),
    })
 
    // Route 설정
    app.Get("/", controllers.Index)
 
    return app
}
```

### 9. 이제, 마지막으로 `~/project/xyz/main.go` 를 작성해봅시다.
```go
package main
 
import (
    "github.com/gauryan/xyz/routes"
)
 
 
func main() {
    app := routes.Router()
    app.Listen(":3000")
}
```

### 10. 모듈 의존성 업데이트
```bash
$ go mod tidy
```

### 11. 개발서버 실행
```bash
$ ./run.sh
 ┌---------------------------------------------------┐  
 │                   Fiber v2.19.0                   │ 
 │               http://127.0.0.1:3000               │ 
 │                                                   │ 
 │ Handlers ............. 2  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 35382 │ 
 └---------------------------------------------------┘ 
```

### 12. PC의 `hosts` 파일에 xyz.test.com 을 설정한 후에, 브라우저에서 http://xyz.test.com 을 입력하면 출력화면을 볼 수 있을 것이다. 그러면, 설치완료~!

