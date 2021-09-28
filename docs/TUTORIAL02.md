# MySQL 설치와 연결 : Fiber v2.x Tutorial


### 1. 설치하면서 root 비밀번호를 물어보면 적당한 것으로 입력한다.
```bash
$ sudo apt-get install mysql-client-8.0 mysql-server-8.0 libmysqlclient-dev
$ sudo /etc/init.d/mysql start
```

### 2. 외부접속 허용 설정
```bash
$ sudo vi /etc/mysql/mysql.conf.d/mysqld.cnf
bind-address = 0.0.0.0
```

### 3. MySQL root 계정 비밀번호 변경 (mysql 을 설치하면 root 계정의 비밀번호가 설정되어 있지 않으므로 반드시 설정하여야 한다.)
```bash
$ mysqladmin -u root password 새로운비밀번호 -p
```

### 4. root 계정 외부에서 접속 설정 (필요한 경우에만 설정)
```bash
$ sudo mysql -u root
```
```sql
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'root의 비밀번호';
FLUSH privileges;
```

### 5. Character set 을 UTF-8로 DB를 생성하고, 사용자 계정을 생성하는 방법은 다음과 같다.
```
DB Name   : xyz
User Name : xyz
Password  : xyz123
```
```bash
$ sudo mysql -u root
```
```sql
Welcome TO the MySQL monitor.  Commands END WITH ; OR \g.
Your MySQL connection id IS 13
Server version: 8.0.26-0ubuntu0.20.04.2 (Ubuntu)
 
Copyright (c) 2000, 2021, Oracle AND/OR its affiliates.
 
Oracle IS a registered trademark OF Oracle Corporation AND/OR its
affiliates. Other names may be trademarks OF their respective
owners.
 
TYPE 'help;' OR '\h' FOR help. TYPE '\c' TO clear the CURRENT INPUT statement.
 
mysql> CREATE DATABASE xyz DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
Query OK, 1 ROW affected, 2 warnings (0.06 sec)
 
mysql> CREATE USER 'xyz'@'%' IDENTIFIED BY 'xyz123';
Query OK, 0 ROWS affected (0.03 sec)
 
mysql> GRANT ALL PRIVILEGES ON xyz.* TO 'xyz'@'%' WITH GRANT OPTION;
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> CREATE USER 'xyz'@'localhost' IDENTIFIED BY 'xyz123';
Query OK, 0 ROWS affected (0.01 sec)
 
mysql> GRANT ALL PRIVILEGES ON xyz.* TO 'xyz'@'localhost' WITH GRANT OPTION;
Query OK, 0 ROWS affected (0.01 sec)
 
mysql> FLUSH privileges;
Query OK, 0 ROWS affected (0.00 sec)
 
mysql> exit
```

### 6. 테이블 생성 및 기초 데이터 입력
```bash
$ mysql -u xyz -pxyz123 xyz
```
```sql
mysql: [Warning] USING a password ON the command line interface can be insecure.
Welcome TO the MySQL monitor.  Commands END WITH ; OR \g.
Your MySQL connection id IS 14
Server version: 8.0.26-0ubuntu0.20.04.2 (Ubuntu)
 
Copyright (c) 2000, 2021, Oracle AND/OR its affiliates.
 
Oracle IS a registered trademark OF Oracle Corporation AND/OR its
affiliates. Other names may be trademarks OF their respective
owners.
 
TYPE 'help;' OR '\h' FOR help. TYPE '\c' TO clear the CURRENT INPUT statement.
 
mysql> CREATE TABLE `users` (
    ->     `sno`      INT NOT NULL AUTO_INCREMENT,
    ->     `username` VARCHAR(80) NOT NULL,
    ->     `password` VARCHAR(60) NOT NULL,
    ->     `email` VARCHAR(254) NOT NULL,
    ->     PRIMARY KEY (`sno`)
    -> ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
Query OK, 0 ROWS affected, 1 warning (0.02 sec)
 
mysql>
mysql> CREATE TABLE `admins` (
    ->     `sno`      INT NOT NULL AUTO_INCREMENT,
    ->     `userid`   VARCHAR(255) NOT NULL,
    ->     `password` VARCHAR(255) NOT NULL,
    ->     `nick`     VARCHAR(255),
    ->     PRIMARY KEY (`sno`)
    -> ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
Query OK, 0 ROWS affected, 1 warning (0.02 sec)
 
mysql> SHOW TABLES;
+---------------+
| Tables_in_xyz |
+---------------+
| admins        |
| users         |
+---------------+
2 ROWS IN SET (0.00 sec)
 
mysql> INSERT INTO admins(userid, password, nick) VALUES('testid1', 'passwd1', 'nick1');
Query OK, 1 ROW affected (0.01 sec)
 
mysql> INSERT INTO admins(userid, password, nick) VALUES('testid2', 'passwd2', 'nick2');
Query OK, 1 ROW affected (0.00 sec)
 
mysql> INSERT INTO admins(userid, password, nick) VALUES('testid3', 'passwd3', 'nick3');
Query OK, 1 ROW affected (0.00 sec)
 
mysql> INSERT INTO admins(userid, password, nick) VALUES('testid4', 'passwd4', 'nick4');
Query OK, 1 ROW affected (0.00 sec)
 
mysql> INSERT INTO admins(userid, password, nick) VALUES('testid5', 'passwd5', 'nick5');
Query OK, 1 ROW affected (0.01 sec)
 
mysql> SELECT * FROM admins;
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
 
mysql> exit
```
