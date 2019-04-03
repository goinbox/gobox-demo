#!/bin/bash

urlPrefix="http://127.0.0.1/demo"
headerHost="Host: ${USER}.gdemo.com"

echo "====== test index ======"

echo curl "$urlPrefix/index?status=1" -H "$headerHost"

echo ""
echo "====== test add ======"

curl "$urlPrefix/add?name=aaa&status=1" -H "$headerHost"

echo ""
echo "====== test edit ======"

curl "$urlPrefix/edit?id=101&name=bbb&status=1" -H "$headerHost"

echo ""
echo "====== test get ======"

curl "$urlPrefix/get?id=61" -H "$headerHost"

echo ""
echo "====== test del ======"

curl "$urlPrefix/del?id=101" -H "$headerHost"
