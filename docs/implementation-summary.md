# Clash Verge 支持 - 实现总结

## 已完成的修改

### 1. 新增常量
在 `main.go` 中添加了两个新的文件路径常量：
- `clashOutputFilePath`: 手动 Clash 配置文件路径
- `autoClashOutputFilePath`: 自动生成的 Clash 配置文件路径

### 2. 新增数据结构
```go
// ClashProxy - 单个 Clash 代理配置
type ClashProxy struct {
    Name      string
    Type      string
    Server    string
    Port      int
    UUID      string
    AlterID   int
    Cipher    string
    TLS       bool
    Network   string
    WSPath    string
    WSHeaders map[string]string
}

// ClashConfig - 完整的 Clash 配置
type ClashConfig struct {
    Proxies []ClashProxy
}
```

### 3. 新增函数

#### `generateClashProxies(instances []xui.V2rayInstance, ips string) []ClashProxy`
- 功能：将 V2Ray 实例转换为 Clash 代理格式
- 输入：V2Ray 实例列表和服务器 IP
- 输出：Clash 代理配置列表

#### `generateAndWriteClashConfigs(configs []config.V2rayInstance)`
- 功能：生成完整的 Clash 配置并写入文件
- 处理流程：
  1. 遍历所有配置
  2. 登录获取实例列表
  3. 过滤有效实例（未过期）
  4. 转换为 Clash 格式
  5. 反转顺序（与 V2Ray 保持一致）
  6. 写入文件

#### `writeClashFile(config ClashConfig, filepath string)`
- 功能：将 Clash 配置写入 YAML 文件
- 使用 `yaml.Marshal` 进行序列化

#### `clashHandler(w http.ResponseWriter, r *http.Request)`
- 功能：HTTP 处理器，提供 Clash 订阅服务
- 特点：
  - 设置 Content-Type 为 `application/yaml`
  - 自动合并手动和自动生成的配置
  - 返回 YAML 格式的配置

### 4. 修改现有代码

#### `task()` 函数
在原有任务的基础上，添加了 Clash 配置生成：
```go
// Generate Clash configs
generateAndWriteClashConfigs(configs)
```

#### `startWebServer()` 函数
新增了 Clash 订阅路由：
```go
http.HandleFunc("/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash", clashHandler)
```

#### Import 语句
添加了 YAML 库的导入：
```go
import "gopkg.in/yaml.v3"
```

### 5. 配置示例

生成的 Clash 配置格式：
```yaml
proxies:
  - name: 节点名称
    type: vmess
    server: 1.2.3.4
    port: 10000
    uuid: c3700e59-b55b-4db7-c836-c7c9b4c7d607
    alterId: 0
    cipher: auto
    tls: false
    network: ws
    ws-path: /
    ws-headers:
      Host: 1.2.3.4
```

## API 端点

### 新增端点
```
GET /me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
```
- 返回：YAML 格式的 Clash 配置
- Content-Type: `application/yaml; charset=utf-8`

### 原有端点（保持不变）
```
GET /me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ
```
- 返回：Base64 编码的 vmess 链接

## 兼容性

- ✅ Clash
- ✅ Clash for Windows
- ✅ Clash Verge
- ✅ ClashX (macOS)
- ✅ ClashA (Android)

## 文件生成

定时任务（每 24 小时）会自动生成以下文件：
- `autoClashConfig.yaml` - 自动生成的 Clash 配置
- `clashConfig.yaml` - 手动维护的 Clash 配置（需手动创建）

## 功能特点

1. **自动转换**：V2Ray 实例自动转换为 Clash 格式
2. **配置合并**：自动合并手动和自动生成的配置
3. **节点过滤**：自动过滤过期节点（超过 7 天）
4. **顺序一致**：节点顺序与 V2Ray 订阅保持一致
5. **错误处理**：完善的错误处理和日志记录

## 测试建议

1. 启动服务后，访问 Clash 端点验证响应格式
2. 在 Clash Verge 中添加订阅链接
3. 检查节点是否正确显示
4. 测试节点连接功能
5. 验证定时任务是否正常更新配置

## 注意事项

1. 确保 `gopkg.in/yaml.v3` 依赖已安装（已在 go.mod 中）
2. 如果文件不存在，会返回 500 错误（首次运行需等待定时任务执行）
3. 可以手动创建 `clashConfig.yaml` 添加额外的静态节点
4. 订阅链接应保密，避免泄露
