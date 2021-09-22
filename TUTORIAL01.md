시작 (설치 및 첫페이지 만들어보기) : Fiber v2.x Tutorial
========================================================

본 튜토리얼에서는 대부분의 사이트에 기본적으로 필요한 기능을 익히는데 목적이 있습니다. 예제에서는 관리자로 로그인/로그아웃하고, 관리자를 등록/수정/삭제 등을 할 수 있는 기본적인 기능을 포함하는 예제를 만들어봅니다. 그리고, 비밀번호의 단방향 암호화(SHA-256)를 해보는 기능도 포함됩니다. DB 연결시 GORM을 사용하지만, 제공되는 Model 메소드(ORM)를 사용하지 않고, 직접 쿼리(스토어드 프로시저)를 사용하여 처리할 것입니다.

OS(Ubuntu Linux 20.04)계정은 기본계정인 ubuntu 를 사용하는 것으로 가정한다. 사용할 프로젝트 디렉토리는 ~/project/xyz 로 될 것이다. Go 언어도 설치되어 있다고 가정하겠습니다. 참고로 저는 1.17.1 버전이 설치되어 있습니다.


## 0. 프로젝트 디렉토리 구조

```
xyz
├-- controllers
├-- database
├-- routes
└-- views
```
Fiber는 디렉토리 구조를 강제하지 않기 때문에 마음대로 만드셔도 됩니다. 위의 구조는 제 마음대로 그냥 정한 것입니다.


## 1. 서비스에 사용할 웹서버(nginx)를 설치하고 설정 해보자.

```
$ sudo apt-get install nginx
```
```
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
