# tago

用 Golang 为 Git 仓库设置标签的智能版本管理工具。

---

## 英文文档

[ENGLISH README](README.md)

## 核心特性

🏷️ **智能标签管理**: 自动创建和升级 Git 仓库的语义版本标签  
⚡ **版本基数系统**: 支持可配置的版本进位规则（1/10/100）  
🎯 **交互式确认**: 低版本基数时提供用户确认，高基数时自动执行  
🌍 **子模块支持**: 支持主项目和子模块的独立标签管理  
📋 **语义版本控制**: 遵循 v{major}.{minor}.{patch} 格式规范

## 安装

```bash
go install github.com/go-mate/tago/cmd/tago@latest
```

## 使用方法

### 查看标签列表

```bash
tago
```

输出示例：
```
refs/tags/v0.0.0 Wed Feb 12 16:18:18 2025 +0700
refs/tags/v0.0.1 Thu Feb 13 16:43:08 2025 +0700
refs/tags/v0.0.2 Thu Feb 13 18:43:40 2025 +0700
refs/tags/v0.0.3 Wed Apr 30 15:18:56 2025 +0700
refs/tags/v0.0.4 Wed May 7 18:38:38 2025 +0700
```

### 升级标签版本（交互模式）

从 va.b.c 升级到 va.b.c+1 并推送新标签，会要求用户确认：

```bash
tago bump
```

输出：
```
cd xxx && git push origin v0.0.5
```

### 升级标签版本（自动模式）

从 va.b.c 升级到 va.b.c+1 并推送新标签，无需用户确认：

```bash
tago bump -b=100
```

输出：
```
cd xxx && git push origin v0.0.5
```

### 主项目标签管理

专门用于主项目根目录的标签操作：

```bash
tago bump main
tago bump main -b=10
```

### 子模块标签管理

用于子模块目录的标签操作（带路径前缀）：

```bash
cd submodule-dir
tago bump sub-module
tago bump sub-module -b=100
```

## 版本基数系统说明

版本基数（-b 参数）控制版本号的进位规则：

- **0 或 1**: 交互模式，每次操作都需要用户确认
- **≥ 2**: 自动模式，支持版本号自动进位

### 版本进位示例

假设版本基数为 10：
- `v1.2.9` → `v1.3.0`（patch 达到基数时进位到 minor）
- `v1.9.8` → `v1.9.9`（正常递增）  
- `v1.9.9` → `v2.0.0`（minor 达到基数时进位到 major）

## 高级用法

### 命令组合示例

```bash
# 查看当前所有标签
tago

# 升级标签（需确认）
tago bump

# 快速升级不需确认
tago bump -b=100

# 主项目标签升级
tago bump main -b=10

# 子模块标签升级（在子模块目录中运行）
cd my-submodule
tago bump sub-module -b=10
```

### 版本控制工作流

1. **开发完成后**：运行 `tago` 查看当前标签
2. **创建新版本**：运行 `tago bump` 升级版本
3. **自动化场景**：使用 `tago bump -b=100` 跳过确认
4. **多模块项目**：在不同目录使用对应的子命令

## 技术特点

### 智能版本管理
- 自动解析现有标签格式
- 支持语义版本控制规范
- 处理版本号进位逻辑
- 验证 Git 仓库状态

### 灵活的确认机制
- 低版本基数：交互式确认每个操作
- 高版本基数：自动执行，适合脚本化
- 用户友好的提示信息
- 操作可取消性

### 多项目架构支持
- 主项目标签：`v{major}.{minor}.{patch}`
- 子模块标签：`{path}/v{major}.{minor}.{patch}`
- 路径感知的标签管理
- Git 子模块兼容性

---

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **意见反馈？** 欢迎所有建议和宝贵意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Pull Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Pull Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**使用这个包快乐编程！** 🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-mate/tago.svg?variant=adaptive)](https://starchart.cc/go-mate/tago)