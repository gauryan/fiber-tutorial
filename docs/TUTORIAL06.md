# 관리자 수정하기 : Fiber v2.x

### 1. `~/project/xyz/views/mgmt/admin/update_form.html` 을 작성한다.
```html
<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
  <h4 class="modal-title">관리자 수정</h4>
</div>
<div class="modal-body">
  <form name="update_form" action="/mgmt/admin/update" method="post">
    <div class="form-group">
      <label>아이디</label>
      <input type="text" name="userid" class="form-control" readonly required pattern="[a-zA-Z0-9]+" value="{{ .Admin.Userid }}"/>
      <input type="hidden" name="id" class="form-control" value="{{ .Admin.Sno }}" />
    </div>
    <div class="form-group">
      <label>별명 <small>(필수)</small></label>
      <input type="text" name="nick" class="form-control" required value="{{ .Admin.Nick }}"/>
    </div>
    <div class="form-group" style="text-align: right">
      <input class="btn btn-primary" type="submit" value="관리자 수정" />
    </div>
  </form>
</div>
```

### 2. `~/project/xyz/controllers/mgmt/admin.go` 에 다음을 추가한다.
```go
...
 
// 관리자 수정 폼
// /mgnt/admin/update_form/{id}
func UpdateForm (c *fiber.Ctx) error {
    type Admin struct {
        Sno    int
        Userid string
        Passwd string
        Nick   string
    }
    var admin Admin
 
    id := c.Params("id")
 
    db := database.DBConn
    db.Raw("CALL getAdmin(?)", id).First(&admin)
    data := fiber.Map{"Admin": admin}
    return c.Render("mgmt/admin/update_form", data)
}
 
...
```

### 3. `~/project/xyz/routes/web.go` 에 다음을 추가한다.
```go
    mgmtApp.Get("/admin/update_form/:id", mgmt.UpdateForm)
```

### 4. 이제, 수정 버튼을 클릭하면 수정할 수 있는 폼이 나타나게 될 것이다. 마지막으로 실제로 수정을 처리하는 루틴을 작성하고 라우터에 등록하자. 그전에 `updateAdmin` 이라는 저장 프로시저부터 만들어야겠지?

```bash
$ mysql -u xyz -pxyz123 xyz
```
```
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 262
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
CREATE PROCEDURE updateAdmin
    (i_sno INT,
     i_nick VARCHAR(255))
BEGIN
    UPDATE admins SET nick = i_nick WHERE sno = i_sno;
END $$
DELIMITER ;
```
```sql
mysql> CALL listAdmins();
+-----+---------+-----------+-------+
| sno | userid  | password  | nick  |
+-----+---------+-----------+-------+
|   1 | testid1 | passwd101 | nick1 |
|   2 | testid2 | passwd2   | nick2 |
|   3 | testid3 | passwd3   | nick3 |
|   4 | testid4 | passwd4   | nick4 |
|   5 | testid5 | passwd5   | nick5 |
|   9 | testid6 | passwd106 | nick6 |
+-----+---------+-----------+-------+
6 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> CALL updateAdmin(1, 'nick101');
Query OK, 1 ROW affected (0.00 sec)
 
mysql> CALL listAdmins();
+-----+---------+-----------+---------+
| sno | userid  | password  | nick    |
+-----+---------+-----------+---------+
|   1 | testid1 | passwd101 | nick101 |
|   2 | testid2 | passwd2   | nick2   |
|   3 | testid3 | passwd3   | nick3   |
|   4 | testid4 | passwd4   | nick4   |
|   5 | testid5 | passwd5   | nick5   |
|   9 | testid6 | passwd106 | nick6   |
+-----+---------+-----------+---------+
6 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 5. `~/project/xyz/controllers/mgmt/admin.go` 에 다음을 추가한다.
```
...
 
// 관리자 수정
// /mgmt/admin/update
func Update (c *fiber.Ctx) error {
    id   := c.FormValue("id")
    nick := c.FormValue("nick")
 
    db := database.DBConn
    db.Exec("CALL updateAdmin(?, ?)", id, nick)
 
    return c.Redirect("/mgmt/admin")
}
 
...
```

### 6. `~/project/xyz/routes/web.go` 에 다음을 추가하고, 수정 작업을 진행해보자.
```
    mgmtApp.Post("/admin/update", mgmt.Update)
```
