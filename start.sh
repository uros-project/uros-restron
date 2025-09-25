#!/bin/bash

echo "启动 Uros Restron 数字孪生平台..."

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go，请先安装 Go"
    exit 1
fi

# 下载依赖
echo "下载依赖包..."
go mod tidy

# 编译项目
echo "编译项目..."
go build -o uros-restron .

if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

echo "编译成功！"

# 启动服务
echo "启动服务..."
./uros-restron
