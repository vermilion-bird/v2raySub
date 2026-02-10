#!/bin/bash

#=============================================================================
# v2raySub 健康检查脚本
# 用于定时检查服务状态，如果服务停止则自动重启
#=============================================================================

WORK_DIR="/home/ubuntu/v2raySub"
LOG_FILE="${WORK_DIR}/v2raySub.log"
HEALTH_LOG="${WORK_DIR}/health_check.log"
SERVICE_PORT="8888"

# 记录日志
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$HEALTH_LOG"
}

# 检查进程
check_process() {
    if pgrep -x "v2raySub" > /dev/null; then
        return 0
    else
        return 1
    fi
}

# 检查端口
check_port() {
    if ss -tlnp 2>/dev/null | grep ":${SERVICE_PORT}" > /dev/null; then
        return 0
    else
        return 1
    fi
}

# 检查 API
check_api() {
    if curl -s --max-time 5 http://localhost:${SERVICE_PORT}/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash | head -1 | grep -q "proxies"; then
        return 0
    else
        return 1
    fi
}

# 重启服务
restart_service() {
    log "尝试重启服务..."
    
    # 停止旧进程
    pkill -9 v2raySub 2>/dev/null
    sleep 2
    
    # 启动新进程
    cd "$WORK_DIR"
    nohup ./v2raySub > "$LOG_FILE" 2>&1 &
    
    sleep 3
    
    if check_process; then
        log "服务重启成功"
        return 0
    else
        log "服务重启失败"
        return 1
    fi
}

# 主逻辑
main() {
    # 检查进程
    if ! check_process; then
        log "警告: 进程未运行"
        restart_service
        exit 0
    fi
    
    # 检查端口
    if ! check_port; then
        log "警告: 端口未监听"
        restart_service
        exit 0
    fi
    
    # 检查 API
    if ! check_api; then
        log "警告: API 响应异常"
        restart_service
        exit 0
    fi
    
    # 一切正常，静默退出
    # log "服务运行正常"
}

main
