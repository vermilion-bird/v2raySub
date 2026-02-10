# 快速开始 - Clash Verge 支持

## 🚀 快速上手

### 1. 编译程序
```bash
go build -o v2raySub
```

### 2. 运行程序
```bash
./v2raySub
```

### 3. 获取订阅链接

#### Clash Verge 订阅（新功能）
```
http://localhost:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
```

#### V2Ray 订阅（原有功能）
```
http://localhost:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ
```

### 4. 在 Clash Verge 中使用

**步骤：**
1. 打开 Clash Verge 客户端
2. 点击左侧菜单的 "配置" 或 "Profiles"
3. 点击右上角的 "新建" 或 "+" 按钮
4. 选择 "URL" 类型
5. 输入订阅链接（替换 localhost 为您的服务器地址）：
   ```
   http://your-server:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
   ```
6. 点击 "确定" 或 "下载"
7. 选择刚添加的配置
8. 在 "代理" 页面选择节点
9. 开启系统代理即可使用

## 📝 配置说明

### 初次运行
程序启动后会立即开启 Web 服务，但 Clash 配置文件需要等待定时任务首次执行后才会生成。

**两种方式获取配置：**

#### 方式 1：等待自动生成（推荐）
- 程序会在启动后的 24 小时内首次执行
- 查看日志确认执行：`定时任务执行: 2024-xx-xx xx:xx:xx`
- 生成文件：`autoClashConfig.yaml`

#### 方式 2：手动创建配置（立即可用）
1. 复制示例文件：
   ```bash
   cp clashConfig.yaml.example clashConfig.yaml
   ```
2. 编辑 `clashConfig.yaml` 添加节点
3. 立即可通过订阅链接访问

### 配置文件说明

**自动生成的文件：**
- `autoClashConfig.yaml` - 由程序自动创建和更新
- 每 24 小时更新一次
- 包含所有有效节点（未过期）

**手动维护的文件：**
- `clashConfig.yaml` - 可选，手动创建
- 用于添加静态节点
- 不会被程序覆盖

**最终订阅内容 = 手动配置 + 自动配置**

## 🔧 测试验证

### 测试 1：验证 API 响应
```bash
# 测试 Clash 端点
curl http://localhost:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash

# 应该返回 YAML 格式的配置
# 如果返回空配置，说明还没有生成节点
```

### 测试 2：检查文件生成
```bash
# 查看自动生成的配置文件
cat autoClashConfig.yaml

# 查看手动配置文件（如果创建了）
cat clashConfig.yaml
```

### 测试 3：查看日志
程序运行时会输出日志信息：
- 服务器启动信息
- 定时任务执行时间
- 文件生成状态
- 错误信息（如有）

## ⚙️ 常见配置

### 修改服务器端口
编辑 `main.go`，修改常量：
```go
const serverPort = ":8889"  // 改为您想要的端口
```
重新编译运行即可。

### 修改定时任务间隔
编辑 `main.go`，修改常量：
```go
const tickerInterval = 1 * 24 * 60 * time.Minute  // 默认 24 小时
```
可以改为：
```go
const tickerInterval = 1 * 60 * time.Minute  // 1 小时
const tickerInterval = 30 * time.Minute      // 30 分钟
```

### 修改节点保留天数
编辑 `main.go`，修改常量：
```go
const keepDays = 7  // 默认保留 7 天
```

## 🌐 生产环境部署

### 使用 systemd（推荐）
创建服务文件 `/etc/systemd/system/v2raysub.service`：
```ini
[Unit]
Description=V2Ray Subscription Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/path/to/v2raySub
ExecStart=/path/to/v2raySub/v2raySub
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable v2raysub
sudo systemctl start v2raysub
sudo systemctl status v2raysub
```

### 使用 Nginx 反向代理（推荐）
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ {
        proxy_pass http://127.0.0.1:8889;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

添加 HTTPS（强烈推荐）：
```bash
sudo certbot --nginx -d your-domain.com
```

## 📊 监控和维护

### 查看日志
```bash
# 如果使用 systemd
sudo journalctl -u v2raysub -f

# 或者重定向输出
./v2raySub > app.log 2>&1 &
tail -f app.log
```

### 检查服务状态
```bash
# 检查端口是否监听
netstat -tlnp | grep 8889

# 或使用 ss
ss -tlnp | grep 8889
```

### 手动触发更新
修改代码让定时任务立即执行一次，或者重启服务。

## 🔒 安全建议

1. **修改默认路径**：将 `/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ` 改为随机字符串
2. **使用 HTTPS**：通过 Nginx 反向代理提供 SSL/TLS 加密
3. **限制访问**：使用防火墙或 Nginx 限制访问 IP
4. **定期更新**：保持程序和依赖项最新
5. **备份配置**：定期备份 `config/config.yaml` 文件

## ❓ 常见问题

### Q: 访问订阅链接返回 404
A: 检查路径是否正确，确保程序正在运行

### Q: 返回空的配置
A: 等待定时任务首次执行，或手动创建 `clashConfig.yaml`

### Q: Clash Verge 无法解析配置
A: 检查 YAML 格式是否正确，查看程序日志

### Q: 节点不更新
A: 检查定时任务是否正常执行，查看日志输出

### Q: 如何立即触发更新？
A: 重启程序，或等待下一个定时周期

## 📚 更多文档

- [文档索引](README.md) - 所有文档列表
- [项目说明](../README.md) - 项目总览
- [Clash 使用说明](clash-usage.md) - Clash 详细使用
- [实现总结](implementation-summary.md) - 技术实现细节
- [更新日志](../CHANGELOG.md) - 功能变更记录

## 💡 提示

- 首次使用建议先手动创建配置文件进行测试
- 生产环境建议使用 systemd 管理服务
- 使用域名和 HTTPS 提供更好的安全性
- 定期检查日志确保服务正常运行
