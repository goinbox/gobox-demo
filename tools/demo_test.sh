#!/bin/bash

urlPrefix="http://127.0.0.1/demo"
headerHost="Host: ${USER}.gdemo.com"

echo "====== test index ======"

curl "$urlPrefix/index?status=1&traceId=abc" -H "$headerHost" -H "TRACE-ID: 12345"

echo ""
echo "====== test add ======"

curl "$urlPrefix/add?name=aaa&status=1" -H "$headerHost"

echo ""
echo "====== test edit ======"

curl "$urlPrefix/edit?id=63&name=bbb&status=1" -H "$headerHost"

echo ""
echo "====== test get ======"

curl "$urlPrefix/get?id=63" -H "$headerHost"

echo ""
echo "====== test del ======"

curl "$urlPrefix/del?id=101" -H "$headerHost"
