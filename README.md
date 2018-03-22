# gobox-demo

这个项目是[gobox](https://github.com/goinbox)在实际项目中的运用示例。

## 需要的第三方工具

1. 安装[dep](https://golang.github.io/dep/)
1. 安装[rigger](https://github.com/ligang1109/rigger)

```
ligang@vm-centos7 ~ $ which dep rigger
/usr/local/bin/dep
/usr/local/bin/rigger
```

## 项目初始化

获取项目：

```
ligang@vm-centos7 ~/devspace $ git clone git@github.com:goinbox/gobox-demo.git
Cloning into 'gobox-demo'...
remote: Counting objects: 258, done.
remote: Compressing objects: 100% (21/21), done.
remote: Total 258 (delta 6), reused 22 (delta 6), pack-reused 230
Receiving objects: 100% (258/258), 51.02 KiB | 0 bytes/s, done.
Resolving deltas: 100% (105/105), done.
Checking connectivity... done.
```

使用rigger初始化项目：

```
ligang@vm-centos7 ~/devspace $ cd gobox-demo/
ligang@vm-centos7 ~/devspace/gobox-demo $ ./init.sh 
[sudo] password for ligang:
```

初始化mysql：

```
mysql> source /home/ligang/devspace/gobox-demo/tools/gobox-demo.sql;
Query OK, 0 rows affected (0.00 sec)

Query OK, 0 rows affected (0.00 sec)

Query OK, 0 rows affected (0.00 sec)

Query OK, 0 rows affected (0.00 sec)

Query OK, 0 rows affected (0.00 sec)

Query OK, 2 rows affected (0.02 sec)

Query OK, 1 row affected (0.00 sec)

Database changed
Query OK, 0 rows affected, 1 warning (0.00 sec)

Query OK, 0 rows affected (0.03 sec)

Query OK, 10 rows affected (0.00 sec)
Records: 10  Duplicates: 0  Warnings: 0

Query OK, 0 rows affected, 1 warning (0.00 sec)

Query OK, 0 rows affected (0.01 sec)

Query OK, 1 row affected (0.01 sec)
```

## 开发调试

运行demo-api：

```
ligang@vm-centos7 ~/devspace/gobox-demo/src $ ./go.sh run gdemo/main/api/main.go -prjHome=/home/ligang/devspace/gobox-demo/
```

这会启动示例的demo-api进程：

```
ligang@vm-centos7 ~ $ ps auxww | grep gobox-demo
ligang     4980  0.0  0.0 113120  1460 pts/1    S+   10:21   0:00 /bin/bash ./go.sh run gdemo/main/api/main.go -prjHome=/home/ligang/devspace/gobox-demo/
ligang     4984  0.1  0.2 202644  9148 pts/1    Sl+  10:21   0:00 go run gdemo/main/api/main.go -prjHome=/home/ligang/devspace/gobox-demo/
ligang     5121  0.0  0.1  69108  4452 pts/1    Sl+  10:21   0:00 /tmp/go-build030739617/command-line-arguments/_obj/exe/main -prjHome=/home/ligang/devspace/gobox-demo/
```

可以简单测试下是否运行成功：

```
ligang@vm-centos7 ~ $ curl 'http://$USER.gdemo.com/demo/get?id=1' -x 127.0.0.1:80

返回json：

{
    "errno": 0, 
    "msg": "", 
    "v": "", 
    "data": {
        "id": 1, 
        "add_time": "2018-03-20 16:24:15", 
        "edit_time": "2018-03-20 16:24:15", 
        "name": "aa", 
        "status": 0
    }
}
```

## demo-api说明

本示例项目实现了对demo这个mysql表的增删改查操作，并配合使用redis作为缓存的简单使用, 

表结构请见项目中的`tools/gobox-demo.sql`

api入口main文件放在`src/gdemo/main/api/main.go`

对应增删改查操作的controller：

```
ligang@vm-centos7 ~/devspace/gobox-demo $ tree src/gdemo/controller/
src/gdemo/controller/
├── api
│   ├── base.go
│   └── demo
│       ├── add.go
│       ├── base.go
│       ├── del.go
│       ├── edit.go
│       ├── get.go
│       └── index.go
└── base.go
```

