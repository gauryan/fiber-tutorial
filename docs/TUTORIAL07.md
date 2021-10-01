# 관리자 삭제하기 : Fiber v2.x

### 1. 관리자 삭제를 위한 저장 프로시저(`deleteAdmin`)을 만들어봅시다.
```
$ mysql -u xyz -pxyz123 xyz
```
```
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 271
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
CREATE PROCEDURE deleteAdmin (i_sno INT)
BEGIN
    DELETE FROM admins WHERE sno = i_sno;
END $$
DELIMITER ;
```
```
mysql> CALL listAdmins();
+-----+---------+-----------+---------+
| sno | userid  | password  | nick    |
+-----+---------+-----------+---------+
|   1 | testid1 | passwd101 | nick101 |
|   2 | testid2 | passwd2   | nick2   |
|   3 | testid3 | passwd3   | nick3   |
|   4 | testid4 | passwd4   | nick4   |
|   5 | testid5 | passwd5   | nick502 |
|   9 | testid6 | passwd106 | nick601 |
+-----+---------+-----------+---------+
6 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> CALL deleteAdmin(9);
Query OK, 1 ROW affected (0.01 sec)
 
mysql> CALL listAdmins();
+-----+---------+-----------+---------+
| sno | userid  | password  | nick    |
+-----+---------+-----------+---------+
|   1 | testid1 | passwd101 | nick101 |
|   2 | testid2 | passwd2   | nick2   |
|   3 | testid3 | passwd3   | nick3   |
|   4 | testid4 | passwd4   | nick4   |
|   5 | testid5 | passwd5   | nick502 |
+-----+---------+-----------+---------+
5 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 2. `~/project/xyz/controllers/mgmt/admin.go` 에 다음을 추가한다.
```go
...
 
// 관리자 삭제
// /mgmt/admin/delete/{id}
func Delete (c *fiber.Ctx) error {
    id := c.Params("id")
 
    db := database.DBConn
    db.Exec("CALL deleteAdmin(?)", id)
    return c.Redirect("/mgmt/admin")
}
 
...
```

### 3. `~/project/xyz/routes/web.go` 에 다음을 추가하고, 관리자 삭제를 해봅니다.
```go
    mgmtApp.Get("/admin/delete/:id", mgmt.Delete)
```

### 4. 여기까지 하면, 기본적인 CRUD 기능을 모두 작성할 수 있게 됩니다. *^^*
