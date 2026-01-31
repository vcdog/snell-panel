# Snell Panel (Surge)

[English](README.md)

## 概述

Snell Panel 是一个全面的 Snell 代理节点管理系统，提供统一的节点管理和自动订阅链接生成。该系统由后端 API 服务器和多个前端界面（Web UI 和 iOS App）组成，旨在提供无缝的节点管理体验。

**核心组件：**
- **后端服务器**：管理节点并生成订阅链接的 RESTful API 服务器
- **Web UI**：基于浏览器的管理界面，访问地址：[snell-panel.owo.nz](http://snell-panel.owo.nz)
- **iOS App**：通过 TestFlight 提供的原生移动应用，提供增强的移动体验

## 功能特性

### 节点管理
- **多节点支持**：统一管理多个 Snell 代理节点
- **节点操作**：添加、删除和修改节点配置
- **节点重命名**：自定义节点名称以便更好地组织
- **中转节点**：支持添加子节点（中转/转发节点）以实现高级路由
- **实时监控**：跟踪节点状态和性能

### 订阅管理
- **自动生成**：生成兼容 Surge 和其他代理客户端的订阅链接
- **动态更新**：订阅链接自动反映节点变更
- **多种格式**：支持多种订阅格式

### API 与集成
- **RESTful API**：用于程序化节点管理的完整 API 端点
- **跨平台**：API 同时服务于 Web UI 和 iOS App，确保功能一致性
- **安全访问**：基于令牌（Token）的 API 安全认证

### 用户界面
- **Web UI**：功能丰富的浏览器管理界面
- **iOS App**：针对触摸操作优化的原生移动应用
- **响应式设计**：在所有设备上提供一致的体验

## 使用说明

### 方法 1：直接运行二进制文件

1. **配置环境变量**

   创建 `.env` 文件或设置环境变量：

   ```bash
   export API_TOKEN=your_token_here
   export DATABASE_URL=your_database_url_here  # 例如：postgres://user:pass@host:port/dbname
   export PORT=8080  # 可选，默认为 8080
   ```

2. **启动 Snell Panel 服务器**

   ```bash
   ./snell-panel
   ```

### 方法 2：使用 Docker Compose（推荐）

1. **配置 compose.yaml 文件**

   编辑 `compose.yaml` 文件并更新环境变量：

   ```yaml
   environment:
     - API_TOKEN=your_token_here
     - DATABASE_URL=your_database_url_here
   ```

2. **启动服务**

   ```bash
   docker-compose up -d
   ```

   服务将在端口 9997 上可用。

### 方法 3：使用 Docker

1. **构建 Docker 镜像**

   ```bash
   docker build -t snell-panel .
   ```

2. **运行容器**

   ```bash
   docker run -d \
     --name snell-panel \
     -p 8080:8080 \
     -e API_TOKEN=your_token_here \
     -e DATABASE_URL=your_database_url_here \
     snell-panel
   ```

### 方法 4：Vercel Serverless 部署

1. **设置 Supabase 数据库**

   - 前往 [supabase.com](https://supabase.com) 创建新项目
   - 导航至 Settings > Database 并复制连接字符串 (Connection String)
   - 连接字符串格式应为：`postgresql://postgres:[YOUR-PASSWORD]@[YOUR-HOST]:[YOUR-PORT]/postgres`

   **详细步骤：**
   - 创建 Supabase 账户和新项目
   - 转到 Project Settings → Database
   - 在 "Connection string" 下，选择 "URI"
   - 复制连接字符串（格式类似于 `postgresql://postgres:[YOUR-PASSWORD]@db.xxxxx.supabase.co:5432/postgres`）
   - 将 `[YOUR-PASSWORD]` 替换为您的实际数据库密码

2. **部署到 Vercel**

   点击下方按钮直接部署到 Vercel：

   [![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/missuo/snell-panel)

   或者手动部署：

   ```bash
   # 克隆仓库
   git clone https://github.com/missuo/snell-panel.git
   cd snell-panel

   # 安装 Vercel CLI
   npm i -g vercel

   # 部署到 Vercel
   vercel --prod
   ```

3. **在 Vercel 中配置环境变量**

   在您的 Vercel 仪表板中，进入项目设置并添加以下环境变量：

   ```
   API_TOKEN=your_token_here
   DATABASE_URL=your_supabase_database_url
   ```

   示例 Supabase DATABASE_URL：
   ```
   postgresql://postgres:your_password@db.abcdefghijklmnop.supabase.co:5432/postgres
   ```

4. **访问您的部署**

   您的 Snell Panel 将可通过 `https://your-project-name.vercel.app` 访问

   **Vercel 部署的优势：**
   - Serverless 架构，自动扩缩容
   - 全球边缘网络，响应速度更快
   - 自动配置 HTTPS 和自定义域名
   - 零服务器维护
   - 个人项目可使用免费层级

## 安装 Snell 服务端

使用以下命令**安装** Snell 服务端：

```bash
bash <(curl -Ls https://ssa.sx/sn) install your_panel_url your_token custom_node_name
```

使用以下命令**卸载** Snell 服务端：

```bash
bash <(curl -Ls https://ssa.sx/sn) uninstall your_panel_url your_token custom_node_name
```

`custom_node_name` 是可选的。如果您的节点名称包含空格，请使用引号。例如：

```bash
bash <(curl -Ls https://ssa.sx/sn) install your_panel_url your_token "My Node Name"
```

使用以下命令**更新** Snell 服务端：
```bash
bash <(curl -Ls https://ssa.sx/sn) update
```

## 访问 Web UI

使用以下链接访问 Web 管理界面：

[https://snell-panel.owo.nz](https://snell-panel.owo.nz)

**Web UI 功能：**
- **仪表板**：概览所有受管节点及其状态
- **节点管理**：添加、删除和配置 Snell 代理节点
- **节点自定义**：重命名节点以便更好地组织和识别
- **中转配置**：设置子节点（中转/转发节点）以应对高级路由场景
- **订阅链接**：生成并复制适用于各种代理客户端的订阅 URL
- **实时更新**：实时监控状态和配置变更

您可以从 Web UI 获取订阅链接。

### 替代方案：iOS App

您也可以使用 iOS App 代替 Web UI 以获得更好的移动体验。由于该应用尚未上架 App Store，您必须使用 TestFlight 进行安装：

**TestFlight Beta 下载：**
[https://testflight.apple.com/join/wKvw64P6](https://testflight.apple.com/join/wKvw64P6)

**iOS App 功能：**
- **原生界面**：针对 iOS 设备的触摸优化界面
- **全功能**：与 Web UI 匹配的完整节点管理能力
- **节点操作**：随时随地添加、删除、重命名和配置节点
- **中转管理**：配置子节点和路由策略
- **订阅分享**：轻松复制和分享订阅链接
- **离线访问**：即使离线也能查看节点配置

**安装说明：**
- TestFlight Beta 测试名额可能有限
- 您需要先在 iOS 设备上安装 [TestFlight](https://apps.apple.com/app/testflight/id899247664)
- iOS 应用提供与 Web 界面相同的功能，并进行了原生移动优化

## API 文档

Snell Panel 提供了全面的 RESTful API，用于程序化节点管理和集成。

**基本 URL：** `https://your-panel-domain.com`

**认证：** 所有 API 请求都需要 `token` 查询参数：
```
GET /entries?token=your_api_token_here
```

### 端点 (Endpoints)

#### 1. 欢迎消息
```
GET /
```
返回包含基本 API 信息的欢迎消息。

**响应：**
```json
{
  "status": "success",
  "message": "Welcome to Snell Panel. Please use the API to manage the entries.\n https://github.com/missuo/snell-panel"
}
```

#### 2. 创建节点条目
```
POST /entry?token=your_token
```

**请求体：**
```json
{
  "ip": "example.com",
  "port": 443,
  "psk": "your_psk_here",
  "node_name": "Custom Node Name",
  "version": "4"
}
```

**响应：**
```json
{
  "status": "success",
  "message": "Entry created successfully",
  "data": {
    "id": 1,
    "ip": "example.com",
    "port": 443,
    "psk": "your_psk_here",
    "country_code": "US",
    "isp": "Example ISP",
    "asn": 12345,
    "node_id": "uuid-generated-string",
    "node_name": "Custom Node Name",
    "version": "4"
  }
}
```

#### 3. 列出所有节点
```
GET /entries?token=your_token
```

**响应：**
```json
{
  "status": "success",
  "message": "Entries retrieved successfully",
  "data": [
    {
      "id": 1,
      "ip": "example.com",
      "port": 443,
      "psk": "your_psk_here",
      "country_code": "US",
      "isp": "Example ISP",
      "asn": 12345,
      "node_id": "uuid-string",
      "node_name": "Custom Node Name",
      "version": "4"
    }
  ]
}
```

#### 4. 通过 IP 删除节点
```
DELETE /entry/:ip?token=your_token
```

**示例：** `DELETE /entry/192.168.1.1?token=your_token`

**响应：**
```json
{
  "status": "success",
  "message": "Entry deleted successfully"
}
```

#### 5. 通过 Node ID 删除节点
```
DELETE /entry/node/:node_id?token=your_token
```

**示例：** `DELETE /entry/node/uuid-string?token=your_token`

**响应：**
```json
{
  "status": "success",
  "message": "Entry deleted successfully"
}
```

#### 6. 生成订阅链接
```
GET /subscribe?token=your_token
```

**响应：** 兼容 Surge 的纯文本订阅内容：
```
🇺🇸 Custom Node Name = snell, example.com, 443, psk = your_psk_here, version = 4
🇯🇵 JP Node = snell, jp.example.com, 443, psk = another_psk, version = 4
```

#### 7. 修改节点
```
PUT /modify/:node_id?token=your_token
```

**请求体：**
```json
{
  "node_name": "New Node Name",
  "ip": "new.example.com"
}
```

**响应：**
```json
{
  "status": "success",
  "message": "Node updated successfully"
}
```

### 数据模型

#### 条目模型
```json
{
  "id": 1,
  "ip": "string",
  "port": 443,
  "psk": "string",
  "country_code": "string",
  "isp": "string", 
  "asn": 12345,
  "node_id": "string",
  "node_name": "string",
  "version": "string"
}
```

#### API 响应模型
```json
{
  "status": "success|error|warning",
  "message": "string",
  "data": "object|array (optional)"
}
```

### 注意事项
- 创建条目时会自动生成 `node_id`
- IP 地址可以是域名或直接 IP - 地理位置信息会自动解析
- 如果未指定，默认版本为 "4"
- 如果令牌无效，所有需认证的端点将返回 401
- 对不存在的资源返回 404 响应

## 待办事项 (TODO)

- [x] **后端 API 服务器**：带认证的 RESTful API
- [x] **Web UI**：功能完备的浏览器界面
- [x] **iOS App**：原生移动应用
- [x] **节点管理**：节点的添加、删除、重命名和配置
- [x] **中转节点**：支持子节点和中转路由
- [x] **订阅生成**：自动链接生成和更新
- [ ] **节点健康监控**：实时健康检查和警报
- [ ] **高级分析**：使用统计和性能指标
- [ ] **Android App**：原生 Android 应用
- [ ] **多用户支持**：用户账户和权限管理

## Web UI 源代码

Snell Panel 项目由多个组件组成：

- **后端 API 服务器**：开源（本仓库）- 基于 Go 的 RESTful API 服务器
- **Web UI 前端**：闭源 - 基于浏览器的管理界面
- **iOS App**：闭源 - 原生 iOS 应用

**我们目前暂不考虑开源前端代码（Web UI 和 iOS App）。** 后端 API 服务器保持完全开源，并提供对所有功能的完整程序化访问。

## 贡献指南

欢迎为 Snell Panel 后端服务器做出贡献！您可以通过以下方式提供帮助：

### 开发设置

**先决条件：**
- 安装 Go 1.24+
- PostgreSQL 数据库（本地或远程）
- Git

**本地开发：**

1. **克隆仓库**
   ```bash
   git clone https://github.com/missuo/snell-panel.git
   cd snell-panel
   ```

2. **设置环境**
   ```bash
   cp .env.example .env
   # 编辑 .env 进行配置
   ```

3. **安装依赖**
   ```bash
   go mod tidy
   ```

4. **以开发模式运行**
   ```bash
   go run .
   ```

### 测试

```bash
go mod tidy
go test ./...
```

### 构建

**本地平台：**
```bash
go mod tidy
go build .
```

**跨平台构建：**
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o snell-panel-linux .

# macOS
GOOS=darwin GOARCH=amd64 go build -o snell-panel-macos .

# Windows  
GOOS=windows GOARCH=amd64 go build -o snell-panel-windows.exe .
```

### 提交指南

- Fork 仓库并创建功能分支
- 编写清晰的提交信息
- 为新功能添加测试
- 提交前确保所有测试通过
- 需要时更新文档
- 提交带有清晰描述的 Pull Request

## 许可证

本项目基于 GPL-3.0 许可证授权。
