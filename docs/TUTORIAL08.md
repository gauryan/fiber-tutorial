# 비밀번호 단방향암호화(SHA256) 하기 : Fiber v2.x

여기에서는 비밀번호의 단방향 암호화를 구현해보도록 하겠습니다. SHA256을 적용하겠습니다. Fiber에서는 아무런 작업을 하지 않고, DB(MySQL)에서만 작업을 하겠습니다.

### 1. 일단 DBMS에 접속해봅니다.
```bash
$ mysql -u xyz -pxyz123 xyz
mysql: [Warning] Using a password on the command line interface can be insecure.
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A
 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 282
Server version: 8.0.26-0ubuntu0.20.04.2 (Ubuntu)
 
Copyright (c) 2000, 2021, Oracle and/or its affiliates.
 
Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.
 
Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
 
mysql>
```

### 2. `Admins` 테이블에 있는 모든 항목을 삭제합니다.
```sql
mysql> CALL listAdmins();
+-----+---------+-----------+-------+
| sno | userid  | password  | nick  |
+-----+---------+-----------+-------+
|   1 | testid1 | passwd101 | nick1 |
|   2 | testid2 | passwd2   | nick2 |
|   3 | testid3 | passwd3   | nick3 |
+-----+---------+-----------+-------+
3 ROWS IN SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> DELETE FROM admins;
Query OK, 3 ROWS affected (0.01 sec)
 
mysql> CALL listAdmins();
Empty SET (0.00 sec)
 
Query OK, 0 ROWS affected (0.00 sec)
```

### 3. 저장프로시저(`insertAdmin`, `updateAdminPassword`)를 삭제합니다.
```sql
mysql> DROP PROCEDURE insertAdmin;
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> DROP PROCEDURE updateAdminPassword;
Query OK, 0 ROWS affected (0.01 sec)
```

### 4. 저장프로시저(`insertAdmin`, `updateAdminPassword`)에 `SHA256`을 적용하여 다시 생성합니다.
```sql
DELIMITER $$
CREATE PROCEDURE insertAdmin
  (i_userid VARCHAR(255),
   i_password VARCHAR(255),
   i_nick VARCHAR(255))
BEGIN
  INSERT INTO admins(userid, password, nick) VALUES(i_userid, SHA2(i_password, 256), i_nick);
END $$
DELIMITER ;
```
```sql
DELIMITER $$
CREATE PROCEDURE updateAdminPassword
    (i_sno INT,
     i_password VARCHAR(255))
BEGIN
    UPDATE admins SET password = SHA2(i_password, 256) WHERE sno = i_sno;
END $$
DELIMITER ;
```

### 5. 관리자 3명을 추가해봅니다.
```sql
mysql> CALL insertAdmin('userid1', 'passwd1', 'nick1');
Query OK, 1 row affected (0.00 sec)

mysql> CALL insertAdmin('userid2', 'passwd2', 'nick2');
Query OK, 1 row affected (0.01 sec)

mysql> CALL insertAdmin('userid3', 'passwd3', 'nick3');
Query OK, 1 row affected (0.00 sec)
```

### 6. 관리자 목록을 조회해보면… 비밀번호가 제대로 암호화되었다는 것을 확인할 수 있습니다.
```sql
mysql> CALL listAdmins();
+-----+---------+------------------------------------------------------------------+-------+
| sno | userid  | password                                                         | nick  |
+-----+---------+------------------------------------------------------------------+-------+
|  10 | userid1 | 3b1d7e9a7c37141350fb473fa099b8b18030cde1909f363e3758e52d4ea1a7b4 | nick1 |
|  11 | userid2 | 5a7d362627a891441ee34012b087915f03a6958c1062fe4cf01de24abecee053 | nick2 |
|  12 | userid3 | 44f1471b4598a6f5577221f7caf011743343b8b3b29c9675738cd225055563b7 | nick3 |
+-----+---------+------------------------------------------------------------------+-------+
3 rows in set (0.00 sec)

Query OK, 0 rows affected (0.00 sec)

mysql> exit
```
