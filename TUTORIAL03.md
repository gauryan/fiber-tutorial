지난 글에서 admins 테이블을 생성하고, 기초 데이터 5개를 넣어두었습니다. 이제 http://xyz.test.com/mgmt/admin 를 접속하면 관리자 목록을 출력하는 페이지를 만들 것입니다. http://xyz.test.com/mgmt 에서 mgmt는 Management 를 줄인말입니다. xyz 사이트의 백엔드 프로그램이라고 생각하시면 될 듯 합니다.

### 1. 저장 프로시저 (`listAdmins`) 생성 : 관리자 목록을 뽑아오는 쿼리입니다.
```
$ mysql -u xyz -pxyz123 xyz
```
```
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 21
Server version: 8.0.26-0ubuntu0.20.04.2 (Ubuntu)

Copyright (c) 2000, 2021, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```
```sql
DELIMITER $$
CREATE PROCEDURE listAdmins()
BEGIN
  SELECT sno, userid, password, nick FROM admins;
END $$
DELIMITER ;
```
```sql
mysql> CALL listAdmins();
+-----+---------+----------+-------+
| sno | userid  | password | nick  |
+-----+---------+----------+-------+
|   1 | testid1 | passwd1  | nick1 |
|   2 | testid2 | passwd2  | nick2 |
|   3 | testid3 | passwd3  | nick3 |
|   4 | testid4 | passwd4  | nick4 |
|   5 | testid5 | passwd5  | nick5 |
+-----+---------+----------+-------+
5 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 2. `database` 디렉토리 아래 `database.go` 를 만든다.
```go
package database
 
import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)
 
var (
    DBConn *gorm.DB
)
 
func Init() {
    var err error
    // dsn := "xyz:xyz123@tcp(127.0.0.1:3306)/xyz?charset=utf8mb4&parseTime=True&loc=Local"
    dsn := "xyz:xyz123@tcp(127.0.0.1:3306)/xyz"
    DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    fmt.Println("Connection Opened to Database")
}
```
GROM을 사용하고, MySQL에 접속하는 함수를 작성합니다.

### 3. `~/project/xyz/main.go` 에 DB관련 코드 추가합니다.
```go
package main
 
import (
    "github.com/gauryan/xyz/routes"
    "github.com/gauryan/xyz/database"
)
 
 
func main() {
    app := routes.Router()
    database.Init()
    app.Listen(":3000")
}
```
기본적인 DB설정은 끝났네요. 이제 본격적으로 코드를 작성해봅시다.

### 4. `~/project/xyz/controllers` 밑에 `mgmt` 디렉토리를 만들고, `admin.go` 파일을 생성합니다.
```go
package mgmt
// controllers/mgmt
 
import (
    "github.com/gauryan/xyz/database"
    "github.com/gofiber/fiber/v2"
)
 
 
// Admin 목록
func ListAdmin(c *fiber.Ctx) error {
    type Admin struct {
        Sno    int
        Userid string
        Nick   string
    }
    var admins []Admin
 
    db := database.DBConn
    // db.Raw("SELECT sno, userid, nick FROM admins").Scan(&admins)
    db.Raw("CALL listAdmins()").Scan(&admins)
 
    data := fiber.Map{"Admins": admins}
    return c.Render("mgmt/admin/index", data, "mgmt/base")
}
```
`ListAdmin`함수가 관리자 목록을 보여주기 위한 컨트롤러가 되겠습니다. 먼저, DB에 접속해서 위에서 만든 `listAdmins` 저장프로시저를 호출해서 `admins` 배열 변수에 담아줍니다. 그리고, `Render` 함수에서 템플릿(View) 파일인 `view/mgmt/admin/index.html`을 호출합니다. 이때, `data` 변수를 통해서 관리자목록도 함께 전달합니다. `mgmt/base`는 레이아웃 파일입니다. ORM 기능을 사용하지 않더라도 저장프로시저를 사용하므로 DB작업이 간단합니다. 재사용 계획이 있는 DB작업의 경우에는 Model 작성하는 것이 좋겠지만, 이 번 과정에서는 그럴 일이 없으므로 컨트롤러에서 직접 DB작업을 하는 것이 좋을 듯 하다.

### 5. `~/project/xyz/views/mgmt/base.html` 파일을 생성한다.
```html
<html>
<head>
<meta charset="utf-8">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
<title>Adonis Tutorial</title>
</head>
<body>
  {{embed}}
</body>
</html>
```
`base.html` 파일은 레이아웃입니다.

### 6. `~/project/xyz/views/mgmt/admin/index.html` 파일을 생성한다.
```html
<div class="container">
  <a href="/mgmt/logout">[로그아웃]</a>
  <div class="page-header">
    <h1>Administrator (관리자)</h1>
  </div>
 
  <div style="text-align: right; padding-bottom: 10px">
    <a href="/mgmt/admin/insert_form" class="btn btn-default" data-toggle="modal" data-target="#myModal">관리자 추가</a>
  </div>
 
  <table class="table table-striped table-hover table-condensed">
  <tr>
    <th style="text-align: center">아이디</th>
    <th style="text-align: center">별명</th>
    <th style="text-align: center">수정/삭제</th>
  </tr>
  {{range .Admins}}
  <tr>
    <td style="text-align: center">{{.Userid}}</td>
    <td style="text-align: center">{{.Nick}}</td>
    <td style="text-align: center">
      <a href="/mgmt/admin/chg_passwd_form/{{.Sno}}" class="btn btn-default btn-xs" data-toggle="modal" data-target="#myModal">비밀번호변경</a>
      <a href="/mgmt/admin/update_form/{{.Sno}}" class="btn btn-default btn-xs" data-toggle="modal" data-target="#myModal">수정</a>
      <button onclick="delete_admin('/mgmt/admin/delete/{{.Sno}}')" class="btn btn-default btn-xs">삭제</button>
    </td>
  </tr>
  {{end}}
  </table>
</div>
 
<div id="myModal" class="modal fade" role="dialog" tabindex="-1" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
    </div>
  </div>
</div>
 
<script>
// Modal Remote Reload
$(document).on('hidden.bs.modal', function (e) {
    $(e.target).removeData('bs.modal');
})
</script>
 
<script>
function delete_admin(url) {
    var result = confirm("관리자를 정말로 삭제하시겠습니까?");
    if( result == false ) return;
    location.href = url;
}
</script>
```
`.Admins` 가 컨트롤러에서 넘겨받은 데이터이다. 이것을 `range` 로 루프를 돌려주면, 그 안에서 `.Userid`, `.Nick`, `.Sno` 등이 튀어나온다.

### 7. `~/project/xyz/routes/web.go` 에 다음처럼 수정한다.
```go
package routes
 
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
    // "github.com/gauryan/xyz/controllers"
    "github.com/gauryan/xyz/controllers/mgmt"
)
 
func Router() *fiber.App {
    // App 생성과 템플릿 설정
    app := fiber.New(fiber.Config{
        Views: html.New("./views", ".html"),
    })
 
    // Route 설정
    // app.Get("/", controllers.Index)
    mgmtApp := app.Group("/mgmt")
    mgmtApp.Get("/admin", mgmt.ListAdmin)
 
    return app
}
```
기존과 다른 점은 `/mgmt` 로 그룹을 만들고, 그 밑에 필요한 디렉토리들을 등록한다. 이렇게 그룹으로 관리하면 여러가지로 편한 면이 있다.

### 8. 모듈 의존성 업데이트
```bash
$ go mod tidy
```

### 9. 개발서버 실행
```
$ ./run.sh
 ┌---------------------------------------------------┐  
 │                   Fiber v2.19.0                   │ 
 │               http://127.0.0.1:3000               │ 
 │                                                   │ 
 │ Handlers ............. 2  Processes ........... 1 │ 
 │ Prefork ....... Disabled  PID ............. 35382 │ 
 └---------------------------------------------------┘ 
```

### 10. 이제, 웹브라우저에서 http://xyz.test.com/mgmt/admin 에 접속하면 관리자 목록을 볼 수 있을 것이다.

