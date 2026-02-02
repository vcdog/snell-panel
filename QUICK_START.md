# 订阅链接配置模板 - 快速开始

## 🚀 功能简介

为 Snell Panel 添加了订阅链接配置模板功能，支持根据不同规则集生成完整的 Surge 配置文件。

## 📋 规则集说明

| 规则集 | 模板文件 | 规则数量 | 适用场景 |
|--------|----------|----------|----------|
| **最小规则** (minimal) | minimal.conf | 3 | 基础代理需求 |
| **均衡规则** (balanced) | balanced.conf | 8 | 日常使用 |
| **全面规则** (comprehensive) | comprehensive.conf | 18 | 完整分流 |
| **自定义** (custom) | 智能选择 | 可变 | 个性化配置 |

## 🎯 快速使用

### 方法一：通过 Web 界面（推荐）

1. 登录 Snell Panel 管理界面
2. 点击"生成订阅链接"按钮
3. 开启"高级选项"开关
4. 在"规则选择"下拉框中选择：
   - **全面规则** - 完整的分流配置（对应 surge1.conf）
   - **均衡规则** - 常用服务配置
   - **最小规则** - 基础配置
   - **自定义** - 手动勾选规则
5. 点击"重新生成"按钮
6. 复制生成的订阅链接到 Surge 使用

### 方法二：直接使用 API

```bash
# 全面规则（推荐）
https://your-domain.com/subscribe?token=YOUR_TOKEN&format=advanced&ruleSet=comprehensive

# 均衡规则
https://your-domain.com/subscribe?token=YOUR_TOKEN&format=advanced&ruleSet=balanced

# 最小规则
https://your-domain.com/subscribe?token=YOUR_TOKEN&format=advanced&ruleSet=minimal
```

## 📦 已包含的规则

### 最小规则包含：
- ✅ 私有网络直连
- ✅ 国内服务直连
- ✅ 其他流量走代理

### 均衡规则包含：
- ✅ AI 服务（ChatGPT、Claude 等）
- ✅ YouTube 视频
- ✅ Google 服务
- ✅ Telegram 消息
- ✅ Github 代码托管
- ✅ 私有网络和国内服务直连

### 全面规则包含（18个分类）：
- ✅ 广告拦截
- ✅ AI 服务
- ✅ 哔哩哔哩
- ✅ YouTube
- ✅ Google
- ✅ Telegram
- ✅ Github
- ✅ 微软服务
- ✅ 苹果服务
- ✅ 社交媒体（Twitter、Facebook、Instagram、TikTok）
- ✅ 流媒体（Netflix、Disney+、HBO 等）
- ✅ 游戏平台（Steam、Epic Games 等）
- ✅ 教育资源
- ✅ 金融服务
- ✅ 云服务（AWS、Azure 等）
- ✅ 私有网络和国内服务

## 💡 使用建议

### 新手用户
→ 推荐使用**均衡规则**，覆盖日常使用的主流服务

### 进阶用户  
→ 推荐使用**全面规则**，获得完整的分流体验

### 资深用户
→ 使用**自定义规则**，精确控制每个分类

## 🔧 配置示例

### Surge iOS/Mac 订阅配置

1. 打开 Surge
2. 选择"配置" → "从 URL 下载配置"
3. 粘贴生成的订阅链接
4. 设置更新间隔（推荐 24 小时）
5. 保存并启用配置

### 节点自动更新

订阅链接是动态的，每次访问都会：
- ✅ 自动获取最新的节点列表
- ✅ 应用选择的规则集配置
- ✅ 生成完整的 Surge 配置文件

## ❓ 常见问题

**Q: 全面规则和 surge1.conf 有什么关系？**

A: 全面规则的配置模板（comprehensive.conf）是参考 surge1.conf 创建的，包含了相同的规则集和代理组配置。

**Q: 可以修改模板吗？**

A: 可以！模板文件位于 `templates/` 目录，修改后无需重启服务即可生效。

**Q: 如何查看生成的配置内容？**

A: 将订阅链接粘贴到浏览器地址栏访问，即可查看生成的完整配置文件。

**Q: 为什么选择"自定义"后还要手动勾选规则？**

A: "自定义"模式让您完全控制需要的规则分类，可以精确配置，避免不必要的规则造成配置文件过大。

**Q: 不同规则集的配置文件大小差异大吗？**

A: 是的：
- 最小规则：~2KB
- 均衡规则：~5KB  
- 全面规则：~12KB

## 📝 更新日志

### Version 1.0.0 (2026-02-02)

- ✨ 新增配置模板系统
- ✨ 支持三种预设规则集（最小、均衡、全面）
- ✨ 全面规则参考 surge1.conf 实现
- ✨ 支持自定义规则组合
- ✨ 优化订阅链接生成逻辑
- 📚 完善功能文档

## 🔗 相关文档

- [SUBSCRIPTION_TEMPLATES.md](./SUBSCRIPTION_TEMPLATES.md) - 详细功能文档
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - 实现总结
- [surge1.conf](./surge1.conf) - 全面规则参考配置

## 💬 反馈与支持

如有问题或建议，欢迎通过以下方式反馈：
- 提交 GitHub Issue
- Telegram: @missuo
- Email: support@example.com

---

**立即开始使用，享受更智能的代理配置体验！** 🎉
