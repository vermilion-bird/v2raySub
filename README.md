# v2raySub
节点订阅

## 功能

- 自动生成 V2Ray vmess 订阅链接（Base64 编码格式）
- 自动生成 Clash Verge 订阅配置（YAML 格式）
- 定时任务自动管理节点

## API 端点

### V2Ray 订阅（原有功能）
```
GET http://your-server:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ
```
返回 Base64 编码的 vmess 链接，适用于 V2Ray 客户端

### Clash Verge 订阅（新增功能）
```
GET http://your-server:8889/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash
```
返回 YAML 格式的 Clash 配置，适用于 Clash、Clash Verge 等客户端

## Clash 配置格式示例

生成的 Clash 配置格式如下：

```yaml
proxies:
  - name: 节点名称
    type: vmess
    server: 1.2.3.4
    port: 443
    uuid: c3700e59-b55b-4db7-c836-c7c9b4c7d607
    alterId: 0
    cipher: auto
    tls: false
    network: ws
    ws-path: /
    ws-headers:
      Host: 1.2.3.4
```

## 使用方法

1. 配置 `config/config.yaml` 文件
2. 运行程序：`./v2raySub`
3. 在客户端中添加订阅链接：
   - V2Ray 客户端：使用 `/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ`
   - Clash Verge 客户端：使用 `/me/rVMhVnCboe75XPMxVw9aVAN1u6wHZ/clash`

## 生成的文件

- `v2rayConfig.txt` - 手动 V2Ray 配置
- `autoV2rayConfig.txt` - 自动生成的 V2Ray 配置
- `clashConfig.yaml` - 手动 Clash 配置
- `autoClashConfig.yaml` - 自动生成的 Clash 配置

## 文档

详细说明见 [docs/](docs/README.md) 目录：

- [快速开始](docs/quickstart.md) - 编译、运行与 Clash 订阅
- [Clash 使用说明](docs/clash-usage.md) - 订阅链接与配置
- [部署指南](docs/deploy.md) - 部署脚本与运维
- [部署示例](docs/deploy-examples.md) - 一键部署使用示例
- [更新日志](CHANGELOG.md) - 功能变更记录
