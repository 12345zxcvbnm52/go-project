#!/bin/bash

# 检查是否传入目录路径
if [ -z "$1" ]; then
  echo "错误：请传入一个目录路径作为参数。"
  exit 1
fi

# 获取传入的目录路径
BASE_DIR="$1"

# 创建目录结构
mkdir -p "$BASE_DIR/etc"
mkdir -p "$BASE_DIR/internal/config"
mkdir -p "$BASE_DIR/internal/form"
mkdir -p "$BASE_DIR/internal/controller"
mkdir -p "$BASE_DIR/internal/data"
mkdir -p "$BASE_DIR/internal/resource"
mkdir -p "$BASE_DIR/internal/model"
mkdir -p "$BASE_DIR/log"

# 创建日志文件
touch "$BASE_DIR/log/error.log"
touch "$BASE_DIR/log/info.log"
touch "$BASE_DIR/log/db.log"

# 创建 main.go 文件
cat <<EOL > "$BASE_DIR/main.go"
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
EOL

# 输出成功信息
echo "文件结构已成功创建在 $BASE_DIR"