# 一键部署脚本使用示例

## 最简手动部署

若只需快速上传并重启，可在项目根目录执行：

```bash
scp ./v2raysub ubuntu@192.0.2.1:~/v2raySub
ssh ubuntu@192.0.2.1 "sudo fuser -k 8889/tcp; cd ~/v2raySub && nohup ./v2raysub &"
```

建议日常使用 [deploy.sh](../deploy.sh) 做完整部署，参见 [部署指南](deploy.md)。

---

## 📦 部署脚本功能

### 核心功能
- ✅ 自动编译 Go 程序
- ✅ SSH 连接测试
- ✅ 远程服务停止
- ✅ 自动备份旧版本
- ✅ 文件上传
- ✅ 远程服务启动
- ✅ 健康检查
- ✅ 日志查看
- ✅ 状态监控

### 安全特性
- ✅ 部署前备份（保留最近3个版本）
- ✅ 颜色化输出，清晰易读
- ✅ 错误处理和回滚支持
- ✅ 连接超时检测

---

## 🚀 快速开始

### 1. 首次配置

```bash
# 进入项目目录
cd /path/to/v2raySub

# 查看帮助
./deploy.sh help

# 配置 Go 环境（如果需要）
export GOROOT=$HOME/go
export PATH=$GOROOT/bin:$PATH
```

### 2. 配置 SSH 免密登录

```bash
# 如果还没有 SSH 密钥，先生成
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

# 复制公钥到服务器
ssh-copy-id ubuntu@192.0.2.1

# 测试连接
ssh ubuntu@192.0.2.1 "echo 'SSH 连接成功'"
```

### 3. 第一次部署

```bash
# 完整部署
./deploy.sh deploy
```

---

## 📝 使用示例

### 场景 1: 完整部署（开发完成后）

```bash
# 一键部署：编译 + 上传 + 重启 + 检查
./deploy.sh deploy
```

**输出示例：**
```
[INFO] 检查本地环境...
[SUCCESS] Go 环境: go version go1.25.7 linux/amd64
[INFO] 测试 SSH 连接...
[SUCCESS] SSH 连接正常
[INFO] 开始编译...
已删除旧的可执行文件
[SUCCESS] 编译完成，文件大小: 9.7M
[INFO] 停止远程服务...
已停止旧进程
已释放端口 8888
[SUCCESS] 远程服务已停止
[INFO] 备份远程文件...
已备份为: v2raySub.backup.20260208_221700
已清理旧备份
[SUCCESS] 备份完成
[INFO] 上传文件到远程服务器...
[SUCCESS] 文件上传完成
[INFO] 启动远程服务...
服务已启动，PID: 12345
进程运行正常
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

---

### 场景 2: 只修改代码，快速部署

```bash
# 编译并部署
./deploy.sh deploy
```

---

### 场景 3: 只重启服务（配置文件修改）

```bash
# 仅重启，不上传新文件
./deploy.sh restart
```

**输出示例：**
```
[INFO] 测试 SSH 连接...
[SUCCESS] SSH 连接正常
[INFO] 停止远程服务...
[SUCCESS] 远程服务已停止
[INFO] 启动远程服务...
[SUCCESS] 远程服务已启动
[INFO] 执行健康检查...
✓ 进程运行中
✓ 端口 8888 监听中
✓ Clash API 响应正常
[SUCCESS] 健康检查通过
```

---

### 场景 4: 查看服务状态

```bash
# 查看远程服务状态
./deploy.sh status
```

**输出示例：**
```
[INFO] 获取远程服务状态...
=========================================
服务状态
=========================================
进程状态: 运行中
ubuntu    12345  0.0  0.2 1900312 15732 ?  Sl  22:17  0:00 ./v2raySub

端口监听:
LISTEN 0  4096  *:8888  *:*  users:(("v2raySub",pid=12345,fd=7))

最近日志 (最后10行):
2026/02/08 22:18:54 定时任务执行: 2026-02-08 22:18:54
2026/02/08 22:18:54 User: Breakingbad
2026/02/08 22:19:03 File written successfully.
2026/02/08 22:19:05 Clash config file written successfully.

=========================================
```

---

### 场景 5: 查看实时日志

```bash
# 查看最近50行日志
./deploy.sh logs
```

---

### 场景 6: 只编译不部署（本地测试）

```bash
# 仅编译
./deploy.sh build

# 本地测试
./v2raySub
```

---

### 场景 7: 紧急停止服务

```bash
# 停止远程服务
./deploy.sh stop
```

---

## 🔧 高级用法

### 自定义配置

创建 `deploy.conf` 文件：

```bash
# 复制配置模板
cp deploy.conf.example deploy.conf

# 编辑配置
vim deploy.conf
```

修改配置内容：
```bash
REMOTE_USER="your_user"
REMOTE_HOST="your_server_ip"
REMOTE_PORT="22"
REMOTE_PATH="/path/to/deploy"
SERVICE_PORT="8888"
```

### 使用别名简化命令

添加到 `~/.bashrc`：

```bash
alias v2deploy='cd /path/to/v2raySub && ./deploy.sh'
alias v2status='cd /path/to/v2raySub && ./deploy.sh status'
alias v2logs='cd /path/to/v2raySub && ./deploy.sh logs'
alias v2restart='cd /path/to/v2raySub && ./deploy.sh restart'
```

然后可以：
```bash
v2deploy          # 快速部署
v2status          # 查看状态
v2logs            # 查看日志
v2restart         # 重启服务
```

---

## 🐛 故障排查

### 问题 1: SSH 连接失败

```bash
[ERROR] 无法连接到远程服务器
```

**解决方法：**
```bash
# 测试 SSH 连接
ssh -v ubuntu@192.0.2.1

# 检查 SSH 密钥
ls -la ~/.ssh/

# 重新配置免密登录
ssh-copy-id ubuntu@192.0.2.1
```

---

### 问题 2: Go 环境未找到

```bash
[ERROR] Go 未安装或不在 PATH 中
```

**解决方法：**
```bash
# 设置 Go 环境
export GOROOT=$HOME/go
export PATH=$GOROOT/bin:$PATH

# 验证
go version

# 永久设置（添加到 ~/.bashrc）
echo 'export GOROOT=$HOME/go' >> ~/.bashrc
echo 'export PATH=$GOROOT/bin:$PATH' >> ~/.bashrc
source ~/.bashrc
```

---

### 问题 3: 端口已占用

```bash
服务器启动失败:listen tcp :8888: bind: address already in use
```

**解决方法：**
```bash
# 使用脚本停止服务
./deploy.sh stop

# 或手动释放端口
ssh ubuntu@192.0.2.1 "sudo fuser -k 8888/tcp"
```

---

### 问题 4: 编译失败

```bash
[ERROR] 编译失败
```

**解决方法：**
```bash
# 检查依赖
go mod tidy

# 清理缓存
go clean -cache

# 重新编译
./deploy.sh build
```

---

## 📊 监控和维护

### 设置定时健康检查

在远程服务器上：

```bash
# 上传健康检查脚本
scp health_check.sh ubuntu@192.0.2.1:/home/ubuntu/v2raySub/
ssh ubuntu@192.0.2.1 "chmod +x /home/ubuntu/v2raySub/health_check.sh"

# 添加到 crontab（每5分钟检查一次）
ssh ubuntu@192.0.2.1 "crontab -l | { cat; echo '*/5 * * * * /home/ubuntu/v2raySub/health_check.sh'; } | crontab -"
```

### 查看健康检查日志

```bash
ssh ubuntu@192.0.2.1 "tail -50 /home/ubuntu/v2raySub/health_check.log"
```

---

## 💡 最佳实践

1. **部署前测试**
   ```bash
   # 本地测试
   ./deploy.sh build
   ./v2raySub  # 前台运行测试
   ```

2. **小步快跑**
   - 每次只修改少量代码
   - 频繁部署和测试
   - 出现问题容易定位

3. **查看日志**
   ```bash
   # 部署后立即查看日志
   ./deploy.sh logs
   ```

4. **定期备份配置**
   ```bash
   # 备份配置文件
   scp ubuntu@192.0.2.1:/home/ubuntu/v2raySub/config/config.yaml ./config.yaml.backup
   ```

5. **使用版本标签**
   ```bash
   # 打标签
   git tag v1.0.0
   git push origin v1.0.0
   
   # 部署时记录版本
   echo "v1.0.0" > VERSION
   ```

---

## 🎯 工作流程建议

### 日常开发流程

```bash
# 1. 修改代码
vim main.go

# 2. 本地测试（可选）
./deploy.sh build
./v2raySub

# 3. 部署到服务器
./deploy.sh deploy

# 4. 查看状态
./deploy.sh status

# 5. 查看日志确认
./deploy.sh logs

# 6. 测试 API
curl http://192.0.2.1:8888/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash | head -20
```

### 紧急回滚流程

```bash
# SSH 到服务器
ssh ubuntu@192.0.2.1

cd /home/ubuntu/v2raySub

# 查看备份
ls -lht v2raySub.backup.*

# 停止服务
pkill v2raySub

# 恢复备份（选择最近的备份）
cp v2raySub.backup.20260208_221700 v2raySub

# 重启
nohup ./v2raySub > v2raySub.log 2>&1 &

# 验证
ps aux | grep v2raySub
```

---

## 📞 获取帮助

```bash
# 显示帮助信息
./deploy.sh help

# 查看脚本版本
head -10 deploy.sh
```

---

## ✨ 总结

部署脚本提供了：
- 🚀 **一键部署** - 从编译到运行全自动
- 🔒 **安全可靠** - 自动备份，支持回滚
- 📊 **状态监控** - 实时查看服务状态
- 🐛 **错误处理** - 详细的错误提示
- 📝 **日志管理** - 方便的日志查看

**现在就开始使用吧！** 🎉
