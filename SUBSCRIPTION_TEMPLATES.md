# 订阅链接配置模板功能

## 概述

此功能允许根据不同的规则集生成相应的 Surge 配置文件。系统提供三种预设的配置模板：

1. **最小规则 (minimal)** - 基础配置，仅包含必要的规则
2. **均衡规则 (balanced)** - 常用服务的规则配置
3. **全面规则 (comprehensive)** - 包含所有分类的完整规则配置

## 配置模板说明

### 1. 最小规则 (minimal.conf)

包含的规则分类：
- 私有网络 (private)
- 国内服务 (china-direct)  
- 非中国 (non-china)

**适用场景**：基础代理需求，配置文件体积小，加载快速

### 2. 均衡规则 (balanced.conf)

包含的规则分类：
- AI 服务 (ai-services)
- 油管视频 (youtube)
- 谷歌服务 (google)
- 私有网络 (private)
- 国内服务 (china-direct)
- 电报消息 (telegram)
- Github (github)
- 非中国 (non-china)

**适用场景**：日常使用，覆盖常用的国际服务

### 3. 全面规则 (comprehensive.conf)

包含的规则分类：
- 广告拦截 (ad-block)
- AI 服务 (ai-services)
- 哔哩哔哩 (bilibili)
- 油管视频 (youtube)
- 谷歌服务 (google)
- 私有网络 (private)
- 国内服务 (china-direct)
- 电报消息 (telegram)
- Github (github)
- 微软服务 (microsoft)
- 苹果服务 (apple)
- 社交媒体 (social-media)
- 流媒体 (streaming)
- 游戏平台 (gaming)
- 教育资源 (education)
- 金融服务 (finance)
- 云服务 (cloud)
- 非中国 (non-china)

**适用场景**：完整的分流需求，涵盖所有主流服务

## 使用方法

### 前端操作

1. 在订阅链接对话框中，开启"高级选项"开关
2. 在"规则选择"下拉框中选择所需的规则集：
   - 自定义 (custom)
   - 最小规则 (minimal)
   - 均衡规则 (balanced)
   - 全面规则 (comprehensive)
3. 选择规则后，相应的复选框会自动勾选
4. 也可以手动勾选/取消勾选规则分类进行自定义
5. 点击"重新生成"按钮生成订阅链接
6. 复制生成的订阅链接到 Surge 客户端使用

### API 调用

生成高级订阅链接的 API 格式：

```
GET /subscribe?token={token}&format=advanced&ruleSet={ruleSet}
```

参数说明：
- `token`: API 访问令牌（必需）
- `format`: 设置为 `advanced` 启用高级配置（必需）
- `ruleSet`: 规则集类型，可选值：
  - `minimal` - 最小规则
  - `balanced` - 均衡规则
  - `comprehensive` - 全面规则
  - `custom` - 自定义规则
- `customRules`: 自定义规则列表，逗号分隔（可选）

示例：
```
# 使用均衡规则
https://your-api.com/subscribe?token=xxx&format=advanced&ruleSet=balanced

# 使用全面规则
https://your-api.com/subscribe?token=xxx&format=advanced&ruleSet=comprehensive

# 自定义规则
https://your-api.com/subscribe?token=xxx&format=advanced&ruleSet=custom&customRules=youtube,google,telegram
```

## 配置文件结构

所有模板文件存放在 `templates/` 目录下：

```
templates/
├── minimal.conf        # 最小规则配置
├── balanced.conf       # 均衡规则配置
└── comprehensive.conf  # 全面规则配置
```

模板文件使用占位符 `{{NODES}}` 和 `{{NODE_NAMES}}`，在生成配置时会被替换为实际的节点信息。

## 模板自定义

如需自定义配置模板：

1. 在 `templates/` 目录下创建或修改 `.conf` 文件
2. 使用 `{{NODES}}` 占位符标记节点列表插入位置
3. 使用 `{{NODE_NAMES}}` 占位符标记节点名称列表插入位置
4. 保存文件后即可通过 ruleSet 参数调用

## 技术实现

### 后端（Go）

1. `utils/config_template.go` - 配置模板加载和生成工具
   - `LoadConfigTemplate()` - 从 templates 目录加载模板
   - `GenerateConfigFromTemplate()` - 替换占位符生成配置
   - `GetTemplateNameByRuleSet()` - 根据规则集确定模板名称

2. `handlers/handlers.go` - GetAdvancedSubscription 方法
   - 查询数据库获取所有节点
   - 根据 ruleSet 参数选择模板
   - 生成完整的 Surge 配置文件

### 前端（JavaScript）

1. `app.js` - 订阅链接生成逻辑
   - `updateSubscriptionUrl()` - 根据用户选择生成订阅 URL
   - `applyRulePreset()` - 应用规则预设，自动勾选复选框
   - `getSubscriptionUrl()` - 构建完整的订阅链接

2. `index.html` - UI 界面
   - 规则选择下拉框
   - 规则复选框网格
   - 高级选项开关

## 注意事项

1. 模板文件必须是有效的 Surge 配置格式
2. 占位符 `{{NODES}}` 和 `{{NODE_NAMES}}` 必须存在
3. 自定义规则会覆盖预设规则的自动选择
4. 更新模板后无需重启服务，下次生成时自动加载新模板
5. 如果模板加载失败，系统会返回错误信息

## 升级说明

从旧版本升级时：

1. 创建 `templates/` 目录
2. 将三个模板文件放入该目录
3. 后端代码已自动支持新功能
4. 前端需要清除浏览器缓存以加载新的 UI

## 故障排查

**问题：生成的配置为空或报错**
- 检查 templates 目录是否存在
- 检查模板文件格式是否正确
- 查看后端日志获取详细错误信息

**问题：规则与预期不符**
- 确认选择的规则集类型
- 检查是否有自定义规则覆盖
- 验证模板文件内容是否正确

**问题：订阅链接无法使用**
- 确认 API URL 和 Token 配置正确
- 检查 format=advanced 参数是否存在
- 验证 ruleSet 参数值是否有效
