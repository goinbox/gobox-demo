#!/bin/bash

urlPrefix="http://127.0.0.1/demo"
headerHost="Host: ${USER}.gdemo.com"

echo "====== test index ======"

echo curl "$urlPrefix/index?status=1&tid=abc" -H "$headerHost" -H "TRACE-ID: 12345"

echo ""
echo "====== test add ======"

curl "$urlPrefix/add?name=aaa&status=1" -H "$headerHost"

echo "===== enter id ======"
read id

echo ""
echo "====== test edit ======"

curl "$urlPrefix/edit?id=$id&name=bbb&status=1" -H "$headerHost"

echo ""
echo "====== test get ======"

curl "$urlPrefix/get?id=$id" -H "$headerHost"

echo ""
echo "====== test del ======"

curl "$urlPrefix/del?id=$id" -H "$headerHost"
