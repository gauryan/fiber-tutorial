# 관리자 추가하기 : Fiber v2.x

### 1. 저장 프로시저 (insertAdmin) 생성
```bash
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
CREATE PROCEDURE insertAdmin
  (userid VARCHAR(255),
   password VARCHAR(255),
   nick VARCHAR(255))
BEGIN
  INSERT INTO admins(userid, password, nick) VALUES(userid, password, nick);
END $$
DELIMITER ;

```
```sql
mysql> select * from admins;
+-----+---------+----------+-------+
| sno | userid  | password | nick  |
+-----+---------+----------+-------+
|   1 | testid1 | passwd1  | nick1 |
|   2 | testid2 | passwd2  | nick2 |
|   3 | testid3 | passwd3  | nick3 |
|   4 | testid4 | passwd4  | nick4 |
|   5 | testid5 | passwd5  | nick5 |
+-----+---------+----------+-------+
5 rows in set (0.00 sec)

mysql> CALL insertAdmin('testid6', 'passwd6', 'nick6');
Query OK, 1 row affected (0.01 sec)

mysql> select * from admins;
+-----+---------+----------+-------+
| sno | userid  | password | nick  |
+-----+---------+----------+-------+
|   1 | testid1 | passwd1  | nick1 |
|   2 | testid2 | passwd2  | nick2 |
|   3 | testid3 | passwd3  | nick3 |
|   4 | testid4 | passwd4  | nick4 |
|   5 | testid5 | passwd5  | nick5 |
|   6 | testid6 | passwd6  | nick6 |
+-----+---------+----------+-------+
6 rows in set (0.00 sec)

mysql> exit
```

### 2. 관리자 입력 양식을 만들자. `~/project/xyz/views/mgmt/admin/insert_form.html`
```html
<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
  <h4 class="modal-title">관리자 추가</h4>
</div>
<div class="modal-body">
  <form name="insert_form" action="/mgmt/admin/insert" method="post">
    <div class="form-group">
      <label>아이디 <small>(필수)</small></label>
      <input type="text" name="userid" class="form-control" required>
    </div>
    <div class="form-group">
      <label>비밀번호 <small>(필수)</small></label>
      <input type="password" id="password" name="passwd1" class="form-control" required>
    </div>
    <div class="form-group">
      <label>비밀번호 확인 <small>(필수)</small></label>
      <input type="password" name="passwd2" class="form-control" required>
    </div>
    <div class="form-group">
      <label>별명 <small>(필수)</small></label>
      <input type="text" name="nick" class="form-control" required>
    </div>
    <div class="form-group" style="text-align: right">
      <input class="btn btn-primary" type="submit" value="관리자 추가" />
    </div>
  </form>
</div>
```

### 3. `~/project/xyz/controllers/mgmt/admin.go` 에 다음 코드를 추가한다.
```go
...

// 관리자 추가 폼
func InsertForm(c *fiber.Ctx) error {
    return c.Render("mgmt/admin/insert_form")
}

...
```

### 4. `~/project/xyz/routes/web.go` 에서 `mgmt` 그룹안에 다음을 추가한다.
```go
    mgmtApp.Get("/admin/insert_form", mgmt.InsertForm)
```

### 5. `관리자 추가` 버튼을 클릭하면 모달 다이얼로그 박스 형식의 입력 양식이 나올 것이다.

### 6. 이제, 실제로 DB에 관리자를 추가해보자. `~/project/xyz/controllers/mgmt/admin.go` 에 다음을 추가한다.
```go

...

// 관리자 추가
func Insert (c *fiber.Ctx) error {
    userid  := c.FormValue("userid")
    passwd1 := c.FormValue("passwd1")
    passwd2 := c.FormValue("passwd2")
    nick    := c.FormValue("nick")

    if passwd1 != passwd2 {
        return c.Redirect("/mgmt/admin")
    }
    db := database.DBConn
    db.Exec("CALL insertAdmin(?, ?, ?)", userid, passwd1, nick)

    return c.Redirect("/mgmt/admin")
}

...
```

### 7. `~/project/xyz/routes/web.go` 에서 `mgmt` 그룹안에 다음을 추가한다.
```go
    mgmtApp.Post("/admin/insert", mgmt.Insert)
```

### 8. 코드 작성은 완료되었으니, 실제 화면에서 관리자를 등록해보면 목록에 표시되는 것을 볼 수 있을 것이다.

