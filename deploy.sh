#!/bin/bash

#=============================================================================
# v2raySub 一键部署脚本
# 功能：编译、上传、部署、健康检查
#=============================================================================

set -e  # 遇到错误立即退出

# 配置变量
REMOTE_USER="ubuntu"
REMOTE_HOST="192.0.2.1"
REMOTE_PATH="/home/ubuntu/v2raySub"
REMOTE_PORT="22"
LOCAL_BINARY="./v2raySub"
SERVICE_PORT="8888"
PROCESS_NAME="v2raySub"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查本地环境
check_local_env() {
    log_info "检查本地环境..."
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        log_error "Go 未安装或不在 PATH 中"
        log_info "请设置 Go 环境变量，例如："
        log_info "export GOROOT=\$HOME/go"
        log_info "export PATH=\$GOROOT/bin:\$PATH"
        exit 1
    fi
    
    log_success "Go 环境: $(go version)"
}

# 编译程序
build() {
    log_info "开始编译..."
    
    # 清理旧的编译文件
    if [ -f "$LOCAL_BINARY" ]; then
        rm -f "$LOCAL_BINARY"
        log_info "已删除旧的可执行文件"
    fi
    
    # 编译
    go build -o v2raySub || {
        log_error "编译失败"
        exit 1
    }
    
    # 检查编译结果
    if [ ! -f "$LOCAL_BINARY" ]; then
        log_error "编译后未找到可执行文件"
        exit 1
    fi
    
    local size=$(du -h "$LOCAL_BINARY" | cut -f1)
    log_success "编译完成，文件大小: $size"
}

# 测试 SSH 连接
test_ssh() {
    log_info "测试 SSH 连接..."
    
    if ssh -p "$REMOTE_PORT" -o ConnectTimeout=5 "${REMOTE_USER}@${REMOTE_HOST}" "echo ok" &>/dev/null; then
        log_success "SSH 连接正常"
    else
        log_error "无法连接到远程服务器"
        log_info "请检查："
        log_info "1. 服务器地址: ${REMOTE_HOST}"
        log_info "2. SSH 端口: ${REMOTE_PORT}"
        log_info "3. 用户名: ${REMOTE_USER}"
        log_info "4. SSH 密钥配置"
        exit 1
    fi
}

# 停止远程服务
stop_remote_service() {
    log_info "停止远程服务..."
    
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << 'ENDSSH'
        # 查找并停止进程
        if pgrep -x "v2raySub" > /dev/null; then
            pkill -9 v2raySub
            echo "已停止旧进程"
            sleep 2
        else
            echo "没有运行中的进程"
        fi
        
        # 释放端口
        if command -v fuser &> /dev/null; then
            sudo fuser -k 8888/tcp 2>/dev/null || true
            echo "已释放端口 8888"
        fi
ENDSSH
    
    log_success "远程服务已停止"
}

# 备份远程文件
backup_remote() {
    log_info "备份远程文件..."
    
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << ENDSSH
        if [ -f "${REMOTE_PATH}/v2raySub" ]; then
            backup_name="v2raySub.backup.\$(date +%Y%m%d_%H%M%S)"
            cp "${REMOTE_PATH}/v2raySub" "${REMOTE_PATH}/\${backup_name}"
            echo "已备份为: \${backup_name}"
            
            # 只保留最近3个备份
            cd "${REMOTE_PATH}"
            ls -t v2raySub.backup.* 2>/dev/null | tail -n +4 | xargs -r rm
            echo "已清理旧备份"
        else
            echo "没有旧文件需要备份"
        fi
ENDSSH
    
    log_success "备份完成"
}

# 上传文件
upload() {
    log_info "上传文件到远程服务器..."
    
    # 确保远程目录存在
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_PATH}"
    
    # 上传可执行文件
    scp -P "$REMOTE_PORT" "$LOCAL_BINARY" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/" || {
        log_error "上传失败"
        exit 1
    }
    
    # 设置执行权限
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" "chmod +x ${REMOTE_PATH}/v2raySub"
    
    log_success "文件上传完成"
}

# 启动远程服务
start_remote_service() {
    log_info "启动远程服务..."
    
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << ENDSSH
        cd "${REMOTE_PATH}"
        
        # 检查配置文件
        if [ ! -f "config/config.yaml" ]; then
            echo "警告: 配置文件不存在，请手动创建 config/config.yaml"
        fi
        
        # 后台启动
        nohup ./v2raySub > v2raySub.log 2>&1 &
        echo "服务已启动，PID: \$!"
        
        # 等待启动
        sleep 3
        
        # 检查进程
        if pgrep -x "v2raySub" > /dev/null; then
            echo "进程运行正常"
        else
            echo "进程启动失败，请查看日志"
            tail -20 v2raySub.log
            exit 1
        fi
ENDSSH
    
    log_success "远程服务已启动"
}

# 健康检查
health_check() {
    log_info "执行健康检查..."
    
    sleep 5  # 等待服务完全启动
    
    # 检查进程
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << 'ENDSSH'
        if pgrep -x "v2raySub" > /dev/null; then
            echo "✓ 进程运行中"
        else
            echo "✗ 进程未运行"
            exit 1
        fi
        
        # 检查端口
        if command -v ss &> /dev/null; then
            if ss -tlnp | grep :8888 > /dev/null; then
                echo "✓ 端口 8888 监听中"
            else
                echo "✗ 端口 8888 未监听"
                exit 1
            fi
        fi
        
        # 测试 API
        if command -v curl &> /dev/null; then
            if curl -s http://localhost:8888/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash | head -1 | grep -q "proxies"; then
                echo "✓ Clash API 响应正常"
            else
                echo "✗ Clash API 响应异常"
            fi
        fi
ENDSSH
    
    log_success "健康检查通过"
}

# 显示远程状态
show_status() {
    log_info "获取远程服务状态..."
    
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << 'ENDSSH'
        echo "========================================="
        echo "服务状态"
        echo "========================================="
        
        # 进程信息
        if pgrep -x "v2raySub" > /dev/null; then
            echo "进程状态: 运行中"
            ps aux | grep v2raySub | grep -v grep
        else
            echo "进程状态: 未运行"
        fi
        
        echo ""
        
        # 端口信息
        echo "端口监听:"
        ss -tlnp 2>/dev/null | grep :8888 || echo "端口 8888 未监听"
        
        echo ""
        
        # 最近日志
        echo "最近日志 (最后10行):"
        if [ -f "/home/ubuntu/v2raySub/v2raySub.log" ]; then
            tail -10 /home/ubuntu/v2raySub/v2raySub.log
        else
            echo "日志文件不存在"
        fi
        
        echo ""
        echo "========================================="
ENDSSH
}

# 显示使用说明
show_usage() {
    cat << EOF
v2raySub 一键部署脚本

用法: $0 [选项]

选项:
    deploy      完整部署（编译、上传、重启）
    build       仅编译
    upload      仅上传（不编译）
    restart     仅重启远程服务
    stop        停止远程服务
    status      查看远程服务状态
    logs        查看远程日志
    help        显示此帮助信息

示例:
    $0 deploy           # 完整部署
    $0 status           # 查看状态
    $0 restart          # 重启服务

配置:
    远程服务器: ${REMOTE_HOST}
    远程用户:   ${REMOTE_USER}
    远程路径:   ${REMOTE_PATH}
    服务端口:   ${SERVICE_PORT}
EOF
}

# 查看日志
show_logs() {
    log_info "获取远程日志..."
    
    ssh -p "$REMOTE_PORT" "${REMOTE_USER}@${REMOTE_HOST}" << ENDSSH
        if [ -f "${REMOTE_PATH}/v2raySub.log" ]; then
            tail -50 "${REMOTE_PATH}/v2raySub.log"
        else
            echo "日志文件不存在"
        fi
ENDSSH
}

# 完整部署流程
deploy() {
    log_info "开始完整部署流程..."
    echo ""
    
    check_local_env
    test_ssh
    build
    stop_remote_service
    backup_remote
    upload
    start_remote_service
    health_check
    
    echo ""
    log_success "========================================="
    log_success "  部署完成！"
    log_success "========================================="
    echo ""
    log_info "访问地址:"
    log_info "  Clash: http://${REMOTE_HOST}:${SERVICE_PORT}/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash"
    log_info "  V2Ray: http://${REMOTE_HOST}:${SERVICE_PORT}/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ"
    echo ""
    log_info "查看状态: $0 status"
    log_info "查看日志: $0 logs"
    echo ""
}

# 主函数
main() {
    case "${1:-deploy}" in
        deploy)
            deploy
            ;;
        build)
            check_local_env
            build
            ;;
        upload)
            test_ssh
            stop_remote_service
            backup_remote
            upload
            start_remote_service
            health_check
            ;;
        restart)
            test_ssh
            stop_remote_service
            start_remote_service
            health_check
            ;;
        stop)
            test_ssh
            stop_remote_service
            ;;
        status)
            test_ssh
            show_status
            ;;
        logs)
            test_ssh
            show_logs
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            log_error "未知命令: $1"
            echo ""
            show_usage
            exit 1
            ;;
    esac
}

# 执行主函数
main "$@"
