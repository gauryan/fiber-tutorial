# 로그인/로그아웃 처리하기 : Fiber v2.x

명색이 관리자 화면인데, 아무나 들어와서 조작하면 안되겠지요? 이제는 마지막으로 로그인/로그아웃 처리를 해보겠습니다. 로그인/로그아웃은 세션을 이용해서 구현합니다.

### 1. 로그인 화면을 구성합니다. `~/project/xyz/views/mgmt/index.html` 을 작성합니다.
```html
<!DOCTYPE html>
<html lang="ko">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Login</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
</head>
<body>
 
<div class="container" style="margin-top: 20px">
  <form action="/mgmt/login" method="post" class="form-horizontal" style="margin: 0 auto; max-width: 360px;">
    <div class="form-group">
      <label for="userid" class="col-sm-3 control-label">아이디</label>
      <div class="col-sm-9">
        <input type="text" id="userid" name="userid" class="form-control" placeholder="당신의 ID를 입력하세요..." required autofocus>
      </div>
    </div>
    <div class="form-group">
      <label for="passwd" class="col-sm-3 control-label">비밀번호</label>
      <div class="col-sm-9">
        <input type="password" id="passwd" name="passwd" class="form-control" placeholder="비밀번호를 입력하세요..." required>
      </div>
    </div>
    <input type="submit" class="btn btn-primary btn-block" value="로그인" />
  </form>
</div>
 
<!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
<!-- Include all compiled plugins (below), or include individual files as needed -->
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
</body>
</html>
```

### 2. `~/project/xyz/controllers/mgmt/main.go` 를 생성합니다.
```go
package mgmt
 
// controllers/mgmt
 
import (
    "github.com/gofiber/fiber/v2"
)
 
// MGMT Login 화면
func Index(c *fiber.Ctx) error {
    return c.Render("mgmt/index", fiber.Map{})
}
```

### 3. `~/project/xyz/routes/web.go` 다음을 추가합니다.
```go
...
    mgmtApp.Get("/", mgmt.Index)
...
```

### 4. http://xyz.test.com/mgmt 에 접속하면 로그인화면이 나올 것입니다.

### 5. MySQL Function 작성할 수 있도록 설정 변경
```
$ sudo mysql -u root
mysql> SET GLOBAL log_bin_trust_function_creators = 1;
mysql> exit
$
```

### 6. MySQL에서 `isMember` Function 생성
```bash
$ mysql -u xyz -pxyz123 xyz
```
```sql
DELIMITER $$
CREATE FUNCTION isMember (
  i_userid VARCHAR(255),
  i_password VARCHAR(255)
) RETURNS INT
BEGIN
  DECLARE CNT INT;
 
  SELECT COUNT(*) INTO CNT FROM admins 
  WHERE userid = i_userid AND password = SHA2(i_password, 256);
 
  RETURN CNT;
END $$
DELIMITER ;
```
현재 Admin 목록 확인
```sql
mysql> CALL listAdmins();
+-----+---------+------------------------------------------------------------------+-------+
| sno | userid  | password                                                         | nick  |
+-----+---------+------------------------------------------------------------------+-------+
|  10 | userid1 | 3b1d7e9a7c37141350fb473fa099b8b18030cde1909f363e3758e52d4ea1a7b4 | nick1 |
|  11 | userid2 | 5a7d362627a891441ee34012b087915f03a6958c1062fe4cf01de24abecee053 | nick2 |
|  12 | userid3 | 44f1471b4598a6f5577221f7caf011743343b8b3b29c9675738cd225055563b7 | nick3 |
|  13 | userid4 | 34344e4d60c2b6d639b7bd22e18f2b0b91bc34bf0ac5f9952744435093cfb4e6 | nick4 |
+-----+---------+------------------------------------------------------------------+-------+
```
isMember Function 테스트
```sql
mysql> SELECT isMember('userid1', 'passwd1');
+--------------------------------+
| isMember('userid1', 'passwd1') |
+--------------------------------+
|                              1 |
+--------------------------------+
1 ROW IN SET (0.00 sec)
 
mysql> SELECT isMember('userid1', 'password2');
+----------------------------------+
| isMember('userid1', 'password2') |
+----------------------------------+
|                                0 |
+----------------------------------+
1 ROW IN SET (0.00 sec)
 
mysql>
```

### 7. `~/project/xyz/store` 디렉토리를 생성하고, 그 안에 `store.go` 를 작성한다
```go
package store
 
import (
    "github.com/gofiber/fiber/v2/middleware/session"
)
 
var SessionStore *session.Store
 
 
func Init() {
    SessionStore = session.New()
}
```

### 8. `~/project/xyz/main.go` 를 다음과 같이 수정한다. 이렇게 하면, Session을 사용할 수 있도록 초기화한다.
```go
package main
 
import (
    "github.com/gauryan/xyz/routes"
    "github.com/gauryan/xyz/database"
    "github.com/gauryan/xyz/store"
)
 
 
func main() {
    app := routes.Router()
    database.Init()
    store.Init()
    app.Listen(":3000")
}
```

### 9. `~/project/xyz/controllers/mgmt/main.go` 에 다음을 추가한다.
```go
...
 
    "github.com/gauryan/xyz/store"
 
...
 
// 로그인
func Login(c *fiber.Ctx) error {
    type Result struct {
        IsMember int
    }
    var result Result
 
    session, err := store.SessionStore.Get(c)
    if err != nil {
        panic(err)
    }
 
    userid := c.FormValue("userid")
    passwd := c.FormValue("passwd")
 
    db := database.DBConn
    db.Raw("SELECT isMember(?, ?) as is_member", userid, passwd).First(&result)
 
    if result.IsMember == 1 {
        session.Set("mgmt-login", true)
        session.Save()
 
        return c.Redirect("/mgmt/admin")
    }
    return c.Redirect("/mgmt")
}
 
 
...
```

### 10. `~/project/xyz/routes/web.go` 에 다음을 추가한다.
```go
...
 
    mgmtApp.Post("/login", mgmt.Login)
 
...
```

### 11. `~/project/xyz/controllers/mgmt/main.go` 에 다음을 추가한다.
```go
...
 
// 로그아웃
func Logout (c *fiber.Ctx) error {
    session, err := store.SessionStore.Get(c)
    if err != nil {
        panic(err)
    }
    session.Destroy()
    return c.Redirect("/mgmt")
}
 
...
```

### 12. `~/project/xyz/routes/web.go` 에 다음을 추가한다.
```go
...
 
    mgmtApp.Get("/logout", mgmt.Logout)
 
...
```
이제, 로그인도 해보고, 로그아웃도 해보세요. ^^

### 13. `~/project/xyz/routes/web.go` 파일을 아래처럼 수정한다. 이 안에서 `authMgmt` 미들웨어 함수를 작성하였고, `mgmt` 그룹을 2개로 나누어서 한쪽에는 `authMgmt` 미들웨어를 적용하였다.
```go
package routes
 
import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/template/html"
    "github.com/gauryan/xyz/store"
    "github.com/gauryan/xyz/controllers/mgmt"
)
 
// authMgmt 미들웨어
func authMgmt(c *fiber.Ctx) error {
    session, err := store.SessionStore.Get(c)
    if err != nil {
        panic(err)
    }
 
    mgmt_login := session.Get("mgmt-login")
    if (mgmt_login != true) {
        return c.Redirect("/mgmt")
    }
 
    return c.Next()
}
 
func Router() *fiber.App {
    // App 생성과 템플릿 설정
    app := fiber.New(fiber.Config{
        Views: html.New("./views", ".html"),
    })
 
    // Route 설정
    // app.Get("/", controllers.Index)
    mgmtApp1 := app.Group("/mgmt")
    mgmtApp1.Get("/", mgmt.Index)
    mgmtApp1.Post("/login", mgmt.Login)
 
    mgmtApp2 := app.Group("/mgmt", authMgmt)
    mgmtApp2.Get("/logout", mgmt.Logout)
    mgmtApp2.Get("/admin", mgmt.ListAdmin)
    mgmtApp2.Get("/admin/insert_form", mgmt.InsertForm)
    mgmtApp2.Post("/admin/insert", mgmt.Insert)
    mgmtApp2.Get("/admin/chg_passwd_form/:id", mgmt.ChgPasswdForm)
    mgmtApp2.Post("/admin/chg_passwd", mgmt.ChgPasswd)
    mgmtApp2.Get("/admin/update_form/:id", mgmt.UpdateForm)
    mgmtApp2.Post("/admin/update", mgmt.Update)
    mgmtApp2.Get("/admin/delete/:id", mgmt.Delete)
 
    return app
}
```

### 14. 여기까지 작성하면, 로그인해야만 http://xyz.test.com/mgmt/admin 을 접근할 수 있게 된다.

