# Domain Redirect Service

一个轻量级的域名重定向服务，支持多目标轮询转发。该服务可以根据配置的域名映射规则，将请求转发到不同的目标地址，并支持轮询负载均衡。

## 功能特点

- 支持多域名配置
- 支持多目标轮询转发
- 自动添加缓存预防参数
- 轻量级设计，资源占用低（Docker 镜像仅 4.2MiB，运行时内存占用约 1.277MiB）
- 支持 Docker 部署

## 配置说明

服务通过环境变量进行配置，支持多个域名映射规则。每个映射规则使用 `DOMAIN_MAPPING_*` 格式的环境变量进行配置。

### 环境变量格式

```
DOMAIN_MAPPING_<任意名称>=<域名>-><目标地址1>,<目标地址2>,...
```

例如：

```
DOMAIN_MAPPING_1=example.com->https://target1.com,https://target2.com
```

### 其他环境变量

- `PORT`: 服务监听端口，默认为 8080
- `PRESERVE_PATH`: 是否保持原始路径，默认为 false。设置为 "true" 时，重定向时会保留原始请求路径。

## 部署方式

### Docker 部署

1. 拉取镜像：

```bash
docker pull nexmoe/domain-redirect
```

2. 构建镜像（可选）：

```bash
docker build -t domain-redirect .
```

3. 运行容器：

```bash
docker run -d \
  -p 8080:8080 \
  -e DOMAIN_MAPPING_1=example.com->https://target1.com,https://target2.com \
  nexmoe/domain-redirect
```

### Docker Compose 部署

1. 创建 `docker-compose.yml` 文件：

```yaml
version: '3'
services:
  domain-redirect:
    image: nexmoe/domain-redirect
    container_name: domain-redirect
    ports:
      - "8080:8080"
    environment:
      - DOMAIN_MAPPING_1=example.com->https://target1.com,https://target2.com
    restart: unless-stopped
```

2. 启动服务：

```bash
docker-compose up -d
```

3. 停止服务：

```bash
docker-compose down
```

### 直接运行

1. 编译：

```bash
go build -o domain-redirect main.go
```

2. 运行：

```bash
export DOMAIN_MAPPING_1=example.com->https://target1.com,https://target2.com
./domain-redirect
```

## 使用示例

假设配置了以下映射：

```
DOMAIN_MAPPING_1=example.com->https://target1.com,https://target2.com
```

当访问 `http://example.com/any/path` 时，服务会：

1. 在目标地址之间轮询选择
2. 将请求重定向到选中的目标地址
3. 保持原始路径（如果配置了 PRESERVE_PATH 为 true）
4. 添加时间戳参数防止缓存

## 注意事项

- 确保配置的域名映射规则格式正确
- 目标地址必须是有效的 URL
- 服务默认监听 8080 端口，可以通过 PORT 环境变量修改
