# 관리자 비밀번호 변경하기 : Fiber v2.x

### 1. 저장 프로시저 (`getAdmin`) 생성
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
CREATE PROCEDURE getAdmin(i_sno INT)
BEGIN
  SELECT sno, userid, password, nick FROM admins WHERE sno = i_sno LIMIT 1;
END $$
DELIMITER ;
```
```sql
mysql> CALL getAdmin(2);
+------+---------+----------+-------+
| sno  | userid  | password | nick  |
+------+---------+----------+-------+
|    2 | testid2 | passwd2  | nick2 |
+------+---------+----------+-------+
1 ROW IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 2. 뷰 디렉토리에 `~/project/xyz/views/mgmt/admin/chg_passwd_form.html` 파일을 생성한다.
```html
<div class="modal-header">
  <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
  <h4 class="modal-title">관리자 비밀번호 변경</h4>
</div>
<div class="modal-body">
  <form name="chg_passwd_form" action="/mgmt/admin/chg_passwd" method="post">
    <div class="form-group">
      <label>아이디</label>
      <input type="text" name="userid" class="form-control" readonly required value="{{ .Admin.Userid }}"/>
      <input type="hidden" name="id" value="{{ .Admin.Sno }}" />
    </div>
    <div class="form-group">
      <label>비밀번호 <small>(필수)</small></label>
      <input type="password" name="passwd1" class="form-control" required />
    </div>
    <div class="form-group">
      <label>비밀번호 확인 <small>(필수)</small></label>
      <input type="password" name="passwd2" class="form-control" required />
    </div>
    <div class="form-group" style="text-align: right">
      <input class="btn btn-primary" type="submit" value="관리자 비밀번호 변경" />
    </div>
  </form>
</div>
```

### 3. 컨트롤러(`~/project/xyz/controllers/mgmt/admin.go`)에 다음을 추가한다.
```go
...
 
// 관리자 비밀번호변경 폼
// /mgnt/admin/chg_passwd_form/:id
func ChgPasswdForm (c *fiber.Ctx) error {
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
    return c.Render("mgmt/admin/chg_passwd_form", data)
}
 
...
```

### 4. 라우터(`~/project/xyz/routes/web.go`)에 다음을 추가한다.
```go
    mgmtApp.Get("/admin/chg_passwd_form/:id", mgmt.ChgPasswdForm)
```

### 5. 이제, `비밀번호변경` 버튼을 클릭하면 비밀번호변경을 위한 모달 다이얼로그박스가 나타날 것이다.

### 6. 실제로 비밀번호를 변경하는 작업을 해보자. 우선 저장 프로시저 (`updateAdminPassword`) 생성하자.
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
CREATE PROCEDURE updateAdminPassword
    (i_sno INT,
     i_password VARCHAR(255))
BEGIN
    UPDATE admins SET password = i_password WHERE sno = i_sno;
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
|   9 | testid6 | passwd6  | nick6 |
+-----+---------+----------+-------+
6 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> CALL updateAdminPassword(1, 'passwd101');
Query OK, 1 ROW affected (0.00 sec)
 
mysql> CALL listAdmins();
+-----+---------+-----------+-------+
| sno | userid  | password  | nick  |
+-----+---------+-----------+-------+
|   1 | testid1 | passwd101 | nick1 |
|   2 | testid2 | passwd2   | nick2 |
|   3 | testid3 | passwd3   | nick3 |
|   4 | testid4 | passwd4   | nick4 |
|   5 | testid5 | passwd5   | nick5 |
|   9 | testid6 | passwd6   | nick6 |
+-----+---------+-----------+-------+
6 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 7. 그리고, 컨트롤러(`~/project/xyz/controllers/mgmt/admin.go`)에 다음을 추가한다.
```go
...
 
// 관리자 비밀번호변경
// /mgmt/admin/chg_passwd
func ChgPasswd (c *fiber.Ctx) error {
    id      := c.FormValue("id")
    passwd1 := c.FormValue("passwd1")
    passwd2 := c.FormValue("passwd2")
 
    if passwd1 != passwd2 {
        return c.Redirect("/mgmt/admin")
    }
    db := database.DBConn
    db.Exec("CALL updateAdminPassword(?, ?)", id, passwd1)
 
    return c.Redirect("/mgmt/admin")
}
 
...
```

### 8. 라우터(`~/project/xyz/routes/web.go`)에는 다음을 추가한다.
```go
    mgmtApp.Post("/admin/chg_passwd", mgmt.ChgPasswd)
```

### 9. 비밀번호를 변경해보고, DB의 내용이 잘 반영되었는지 확인해보자.
```sql
$ mysql -u xyz -pxyz123 xyz
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 241
Server version: 8.0.26-0ubuntu0.20.04.2 (Ubuntu)
 
Copyright (c) 2000, 2021, Oracle and/or its affiliates.
 
Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.
 
Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
 
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
6 rows in set (0.00 sec)
 
Query OK, 0 rows affected (0.00 sec)
 
mysql>
```
