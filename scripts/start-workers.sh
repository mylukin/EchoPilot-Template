#!/bin/bash

# Workers 启动脚本
# 并行运行所有队列处理器

set -e

echo "🚀 Starting Workers..."

# 设置信号处理
cleanup() {
    echo "🛑 Received shutdown signal, stopping all workers..."
    # 杀死所有子进程
    jobs -p | xargs -r kill
    wait
    echo "✅ All workers stopped"
    exit 0
}

trap cleanup SIGTERM SIGINT

# 并行启动所有常驻进程
# echo "Starting worker..."
# /go/app/app-bin worker &

echo "✅ All workers started successfully"

# 等待所有后台进程
wait 