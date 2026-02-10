# Clash Verge 订阅使用说明

## 新增功能

本次更新新增了对 Clash Verge 的支持，现在可以生成符合 Clash 配置格式的订阅链接。

## 订阅链接

### Clash 订阅端点
```
http://your-server:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
```

## 在 Clash Verge 中使用

1. 打开 Clash Verge 客户端
2. 进入 "订阅" 或 "Profiles" 页面
3. 点击 "新建" 或 "Add"
4. 输入订阅链接：`http://your-server:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash`
5. 点击 "确定" 或 "Submit"
6. 更新订阅即可获取最新节点

## 配置格式说明

生成的 Clash 配置包含以下字段：

```yaml
proxies:
  - name: 节点名称                    # 节点显示名称
    type: vmess                       # 协议类型
    server: 1.2.3.4                   # 服务器地址
    port: 443                         # 端口号
    uuid: xxxxxxxx-xxxx-xxxx...       # UUID
    alterId: 0                        # AlterID
    cipher: auto                      # 加密方式
    tls: false                        # 是否启用 TLS
    network: ws                       # 传输协议
    ws-path: /                        # WebSocket 路径
    ws-headers:                       # WebSocket 请求头
      Host: 1.2.3.4
```

## 支持的客户端

- Clash
- Clash for Windows
- Clash Verge
- ClashX (macOS)
- ClashA (Android)

## 文件说明

程序会自动生成以下文件：

- `clashConfig.yaml` - 手动维护的 Clash 配置
- `autoClashConfig.yaml` - 自动生成的 Clash 配置

订阅端点会自动合并这两个文件的内容。

## 注意事项

1. 确保服务器防火墙允许端口 8889 的访问
2. 订阅链接应保密，避免泄露
3. 定时任务会每 24 小时自动更新节点配置
4. 旧节点（超过 7 天）会被自动清理

## 技术细节

- 配置格式：YAML
- 响应类型：`application/yaml; charset=utf-8`
- 自动合并手动和自动生成的配置
- 节点顺序与 V2Ray 订阅保持一致（反向排序）
