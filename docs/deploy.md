# 部署指南

## 快速开始

### 1. 配置部署脚本

编辑 `deploy.sh` 中的配置变量（或使用 `deploy.conf`）：

```bash
REMOTE_USER="ubuntu"           # 远程服务器用户名
REMOTE_HOST="192.0.2.1"     # 远程服务器IP（示例，请改为实际 IP）
REMOTE_PORT="22"               # SSH 端口
REMOTE_PATH="/home/ubuntu/v2raySub"  # 远程部署路径
SERVICE_PORT="8888"            # 服务端口
```

### 2. 配置 SSH 免密登录

```bash
# 生成 SSH 密钥（如果还没有）
ssh-keygen -t rsa

# 复制公钥到远程服务器
ssh-copy-id -p 22 ubuntu@192.0.2.1
```

### 3. 配置 Go 环境

```bash
# 临时设置（本次会话有效）
export GOROOT=$HOME/go
export PATH=$GOROOT/bin:$PATH

# 或者永久设置（添加到 ~/.bashrc）
echo 'export GOROOT=$HOME/go' >> ~/.bashrc
echo 'export PATH=$GOROOT/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

### 4. 执行部署

```bash
# 完整部署（推荐）
./deploy.sh deploy

# 或者分步执行
./deploy.sh build    # 仅编译
./deploy.sh upload   # 上传并重启
```

## 使用说明

### 命令列表

```bash
./deploy.sh deploy      # 完整部署（编译、上传、重启）
./deploy.sh build       # 仅编译
./deploy.sh upload      # 仅上传（不编译）
./deploy.sh restart     # 仅重启远程服务
./deploy.sh stop        # 停止远程服务
./deploy.sh status      # 查看远程服务状态
./deploy.sh logs        # 查看远程日志
./deploy.sh help        # 显示帮助信息
```

### 部署流程

完整部署流程包括以下步骤：

1. ✅ **检查本地环境** - 验证 Go 是否安装
2. ✅ **测试 SSH 连接** - 确保可以连接到远程服务器
3. ✅ **编译程序** - 使用 go build 编译
4. ✅ **停止远程服务** - 停止旧的进程
5. ✅ **备份远程文件** - 备份旧版本（保留最近3个）
6. ✅ **上传文件** - 上传新的可执行文件
7. ✅ **启动远程服务** - 后台启动服务
8. ✅ **健康检查** - 验证服务是否正常运行

### 示例输出

```bash
$ ./deploy.sh deploy

[INFO] 检查本地环境...
[SUCCESS] Go 环境: go version go1.25.7 linux/amd64
[INFO] 测试 SSH 连接...
[SUCCESS] SSH 连接正常
[INFO] 开始编译...
[SUCCESS] 编译完成，文件大小: 9.7M
[INFO] 停止远程服务...
[SUCCESS] 远程服务已停止
[INFO] 备份远程文件...
[SUCCESS] 备份完成
[INFO] 上传文件到远程服务器...
[SUCCESS] 文件上传完成
[INFO] 启动远程服务...
[SUCCESS] 远程服务已启动
[INFO] 执行健康检查...
✓ 进程运行中
✓ 端口 8888 监听中
✓ Clash API 响应正常
[SUCCESS] 健康检查通过

=========================================
  部署完成！
=========================================

访问地址:
  Clash: http://192.0.2.1:8888/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
  V2Ray: http://192.0.2.1:8888/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ

查看状态: ./deploy.sh status
查看日志: ./deploy.sh logs
```

## 常见问题

### 1. SSH 连接失败

**错误信息**：
```
[ERROR] 无法连接到远程服务器
```

**解决方法**：
```bash
# 检查 SSH 配置
ssh -p 22 ubuntu@192.0.2.1

# 配置免密登录
ssh-copy-id -p 22 ubuntu@192.0.2.1
```

### 2. Go 未安装

**错误信息**：
```
[ERROR] Go 未安装或不在 PATH 中
```

**解决方法**：
```bash
# 设置 Go 环境变量
export GOROOT=$HOME/go
export PATH=$GOROOT/bin:$PATH

# 验证
go version
```

### 3. 端口占用

**错误信息**：
```
服务器启动失败:listen tcp :8888: bind: address already in use
```

**解决方法**：
```bash
# 在远程服务器上执行
sudo fuser -k 8888/tcp

# 或者使用部署脚本停止
./deploy.sh stop
```

### 4. 权限问题

**错误信息**：
```
Permission denied
```

**解决方法**：
```bash
# 确保脚本有执行权限
chmod +x deploy.sh

# 确保远程目录有写权限
ssh ubuntu@192.0.2.1 "chmod 755 /home/ubuntu/v2raySub"
```

## 高级配置

### 使用配置文件

创建 `deploy.conf`：

```bash
cp deploy.conf.example deploy.conf
vim deploy.conf
```

修改脚本以支持配置文件：

```bash
# 在脚本开头添加
[ -f deploy.conf ] && source deploy.conf
```

### 配置 systemd 服务

在远程服务器上创建 systemd 服务：

```bash
sudo tee /etc/systemd/system/v2raysub.service > /dev/null <<EOF
[Unit]
Description=V2Ray Subscription Service
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/v2raySub
ExecStart=/home/ubuntu/v2raySub/v2raySub
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable v2raysub
sudo systemctl start v2raysub
sudo systemctl status v2raysub
```

### 配置日志轮转

在远程服务器上创建日志轮转配置：

```bash
sudo tee /etc/logrotate.d/v2raysub > /dev/null <<EOF
/home/ubuntu/v2raySub/v2raySub.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 ubuntu ubuntu
}
EOF
```

### 监控告警

使用 cron 定时检查服务状态：

```bash
# 添加到 crontab
*/5 * * * * /home/ubuntu/v2raySub/health_check.sh
```

创建健康检查脚本：

```bash
#!/bin/bash
if ! pgrep -x "v2raySub" > /dev/null; then
    echo "v2raySub service is down! Restarting..."
    cd /home/ubuntu/v2raySub && nohup ./v2raySub > v2raySub.log 2>&1 &
    # 发送告警邮件或通知
fi
```

## 回滚操作

如果新版本有问题，可以快速回滚：

```bash
# SSH 到远程服务器
ssh ubuntu@192.0.2.1

cd /home/ubuntu/v2raySub

# 停止当前服务
pkill v2raySub

# 恢复备份
cp v2raySub.backup.20260208_221700 v2raySub

# 重启服务
nohup ./v2raySub > v2raySub.log 2>&1 &
```

## 安全建议

1. **使用 SSH 密钥认证**，禁用密码登录
2. **修改默认 SSH 端口**
3. **配置防火墙**，只开放必要端口
4. **定期备份配置文件**
5. **使用 HTTPS**（配置 Nginx 反向代理）
6. **隐藏订阅路径**，使用随机字符串

## 监控与维护

### 查看服务状态

```bash
./deploy.sh status
```

### 查看实时日志

```bash
# 本地执行
./deploy.sh logs

# 或远程执行
ssh ubuntu@192.0.2.1 "tail -f /home/ubuntu/v2raySub/v2raySub.log"
```

### 性能监控

```bash
# CPU 和内存使用
ssh ubuntu@192.0.2.1 "ps aux | grep v2raySub"

# 网络连接
ssh ubuntu@192.0.2.1 "ss -tnp | grep v2raySub"
```

## 故障排查

### 检查清单

- [ ] Go 环境是否正确配置
- [ ] SSH 连接是否正常
- [ ] 远程服务器磁盘空间是否充足
- [ ] 配置文件 config.yaml 是否存在
- [ ] 端口 8888 是否被占用
- [ ] 防火墙是否允许端口 8888
- [ ] 日志文件中是否有错误信息

### 调试模式

```bash
# 在远程服务器上前台运行，查看实时输出
ssh ubuntu@192.0.2.1
cd /home/ubuntu/v2raySub
./v2raySub
```

## 更新日志

查看最近的更新记录：

```bash
ssh ubuntu@192.0.2.1 "ls -lht /home/ubuntu/v2raySub/v2raySub.backup.*"
```

## 技术支持

如遇问题，请查看：

- [项目说明](../README.md)
- [快速开始](quickstart.md)
- 日志文件：`v2raySub.log`
