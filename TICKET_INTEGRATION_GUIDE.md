# Ticket/Issue Number Integration Guide

GitAI支持多种方式提供和管理票号(Jira, GitHub Issues等),让你的commit自动包含工单信息。

## 🚀 快速开始

### 方式1：命令行参数(最简单)

```bash
# 直接指定ticket号
gitai commit --ticket PROJ-123

# 简写
gitai commit -k JIRA-456
```

生成的commit会自动包含:
```
feat(api): [PROJ-123] add user authentication

- Implement JWT-based authentication
- Add login and logout endpoints
```

---

### 方式2：自动从分支名提取(推荐)

如果你的分支命名包含ticket号,GitAI会自动识别:

```bash
# 分支名示例
git checkout -b feature/PROJ-123-add-login
git checkout -b bugfix/JIRA-456-fix-auth
git checkout -b hotfix/GH-789

# GitAI会自动提取 PROJ-123, JIRA-456, GH-789
gitai commit
```

**支持的分支命名格式**:
- `feature/PROJ-123-description`
- `bugfix/JIRA-456-fix-something`
- `PROJ-789` (直接用ticket号做分支名)
- `#123` (GitHub Issues)
- `GH-456` (GitHub格式)

GitAI会询问确认:
```
📝 Git Commit AI Assistant

Found ticket number in branch: PROJ-123
? Use this ticket number?
▸ Yes
  No
```

---

### 方式3：配置为必填项

在 `.gitcommit.yaml` 中配置:

```yaml
# 强制要求提供ticket号
require_ticket: true
ticket_prefix: "PROJ"                    # 默认前缀
ticket_pattern: "[A-Z]+-\\d+"            # 提取规则
```

**效果**:
- 如果分支名有ticket,自动提取
- 如果没有,会提示输入:

```
? Enter ticket number (e.g., PROJ-123): ▌
```

输入 `123`,会自动格式化为 `PROJ-123`

---

## 📝 配置详解

### 基础配置

```yaml
# .gitcommit.yaml

# 是否强制要求ticket号
require_ticket: false              # true=必须提供, false=可选

# 默认ticket前缀(当用户只输入数字时自动添加)
ticket_prefix: "PROJ"              # 如: 输入"123" → "PROJ-123"

# 从分支名提取ticket的正则表达式
ticket_pattern: "[A-Z]+-\\d+"      # 匹配 ABC-123, PROJ-456
```

### Jira项目配置

```yaml
require_ticket: true
ticket_prefix: "JIRA"
ticket_pattern: "[A-Z]+-\\d+"

custom_prompt: |
  格式要求: <type>(<scope>): [JIRA-XXX] <description>

  必须包含:
  Jira: JIRA-XXX
  Reviewer: @username
```

使用:
```bash
# 方式1: 从分支名自动提取
git checkout -b feature/JIRA-123-new-feature
gitai commit

# 方式2: 手动指定
gitai commit --ticket JIRA-123

# 方式3: 只输入数字,自动加前缀
? Enter ticket number (e.g., JIRA-123): 123
# 自动变成 JIRA-123
```

---

### GitHub Issues配置

```yaml
require_ticket: true
ticket_prefix: "GH"
ticket_pattern: "#\\d+|GH-\\d+"    # 匹配 #123 或 GH-123
```

使用:
```bash
# 分支名示例
git checkout -b fix/#123-bug-fix
gitai commit
# 自动提取 #123
```

---

### 中国企业配置

```yaml
require_ticket: true
ticket_prefix: "WO"                # Work Order
ticket_pattern: "WO-\\d{8}-\\d+"   # WO-20250106-001

custom_prompt: |
  提交格式: <类型>(<模块>): [工单号] <描述>

  工单号格式: WO-YYYYMMDD-XXX

  必须包含:
  工单号: WO-YYYYMMDD-XXX
  测试人: @姓名
```

---

## 🎯 使用场景

### 场景1: Jira工作流

**团队规范**: 所有commit必须关联Jira ticket

**配置**:
```yaml
require_ticket: true
ticket_prefix: "PROJ"
ticket_pattern: "PROJ-\\d+"
```

**日常使用**:
```bash
# 1. 从Jira创建分支
git checkout -b feature/PROJ-1234-add-payment

# 2. 开发代码
vim src/payment.js

# 3. 提交时自动提取ticket
git add .
gitai commit

# ✅ 生成: feat(payment): [PROJ-1234] add payment gateway
```

---

### 场景2: GitHub Flow

**团队规范**: Issue驱动开发

**配置**:
```yaml
require_ticket: false              # 可选
ticket_prefix: "GH"
ticket_pattern: "#\\d+|GH-\\d+"
```

**使用**:
```bash
# 1. 从Issue创建分支
git checkout -b fix/#456-memory-leak

# 2. GitAI自动提取 #456
gitai commit
```

---

### 场景3: 多项目不同规范

**项目A (Jira)**:
```yaml
# ~/projects/project-a/.gitcommit.yaml
require_ticket: true
ticket_prefix: "PROJA"
```

**项目B (GitHub)**:
```yaml
# ~/projects/project-b/.gitcommit.yaml
require_ticket: true
ticket_prefix: "GH"
ticket_pattern: "#\\d+"
```

GitAI会自动使用当前项目的配置！

---

## 🔧 高级功能

### 自定义提取规则

支持自定义正则表达式:

```yaml
# 匹配复杂格式
ticket_pattern: "(PROJ|TASK|BUG)-\\d+"    # PROJ-123 或 TASK-456

# 匹配多种格式
ticket_pattern: "[A-Z]{2,10}-\\d+|#\\d+"  # ABC-123 或 #456
```

### 智能格式化

如果用户输入不完整,自动补全:

```yaml
ticket_prefix: "JIRA"
```

用户输入:
- `123` → 自动格式化为 `JIRA-123`
- `JIRA-123` → 保持不变
- `PROJECT-456` → 保持不变

### 分支名模式

支持以下分支命名模式:

```
✅ feature/PROJ-123-description
✅ bugfix/JIRA-456-fix-bug
✅ hotfix/PROJ-789
✅ PROJ-123
✅ fix/#123
✅ feature/GH-456-new-feature
✅ 123-add-feature (需要配置ticket_prefix)
```

---

## 📋 最佳实践

### 1. 统一分支命名规范

```bash
# 推荐格式
<type>/<ticket>-<description>

# 示例
feature/PROJ-123-add-login
bugfix/PROJ-456-fix-crash
hotfix/PROJ-789-security-patch
```

### 2. 配置Git Branch Template

在 `~/.gitconfig` 或项目 `.git/config`:

```ini
[alias]
    # 创建分支时自动提示ticket
    nb = "!f() { \
        read -p 'Ticket number: ' ticket; \
        read -p 'Description: ' desc; \
        git checkout -b \"feature/$ticket-$desc\"; \
    }; f"
```

使用:
```bash
git nb
Ticket number: PROJ-123
Description: add-payment
# 创建分支: feature/PROJ-123-add-payment
```

### 3. 团队配置模板

创建团队共享的配置模板:

```bash
# .gitcommit.team.yaml (提交到仓库)
require_ticket: true
ticket_prefix: "PROJ"
ticket_pattern: "PROJ-\\d+"

custom_prompt: |
  团队commit规范:
  - 必须包含Jira ticket: [PROJ-XXX]
  - 必须包含Reviewer: @username
```

团队成员拉取后:
```bash
git pull
cp .gitcommit.team.yaml .gitcommit.yaml
```

---

## ❓ 常见问题

### Q: 可以强制要求ticket吗?

**A**: 可以!设置 `require_ticket: true`

```yaml
require_ticket: true
```

如果用户不提供ticket,GitAI会报错:
```
❌ Error: ticket number required but not provided
```

### Q: 分支名不包含ticket怎么办?

**A**: GitAI会提示输入:

```
? Enter ticket number (e.g., PROJ-123): ▌
```

### Q: 可以跳过ticket吗?

**A**: 如果 `require_ticket: false`,可以跳过:

```bash
# 不提供ticket也能正常使用
gitai commit
```

### Q: 支持哪些ticket系统?

**A**: 支持所有系统,只需配置正确的pattern:

- ✅ Jira
- ✅ GitHub Issues
- ✅ GitLab Issues
- ✅ Azure DevOps
- ✅ 自定义工单系统

### Q: ticket号出现在哪里?

**A**: 在commit subject line中:

```
feat(api): [PROJ-123] add new endpoint

详细说明...

Jira: PROJ-123
```

### Q: 可以自定义ticket格式吗?

**A**: 可以!在 `custom_prompt` 中指定:

```yaml
custom_prompt: |
  格式要求:
  - Ticket号必须放在最前面: [PROJ-123]
  - 或者放在footer: Ticket: PROJ-123
```

---

## 🎨 示例commit输出

### 带Jira ticket

```
feat(auth): [JIRA-456] add OAuth2 login support

Implemented OAuth2 authentication flow for enterprise SSO.
This allows users to login using their company credentials.

Business Impact:
- Enables enterprise customer onboarding
- Improves security compliance

Technical Details:
- Added OAuth2 library integration
- Implemented callback endpoint

Jira: JIRA-456
Reviewer: @tech-lead
```

### 带GitHub Issue

```
fix(api): [#123] resolve memory leak in connection pool

Fixed connection pool not releasing connections properly.
This was causing server crashes under high load.

- Implement proper connection cleanup
- Add connection timeout handling
- Update connection pool configuration

Closes #123
```

### 带工单号(中文)

```
feat(支付): [WO-20250106-001] 新增支付宝支付功能

实现支付宝扫码支付接口集成。

改动内容:
- 新增支付宝SDK集成
- 实现支付回调处理
- 添加支付状态同步

业务价值:
- 支持更多支付方式
- 提升用户体验

工单号: WO-20250106-001
测试人: @测试工程师
```

---

## 🚀 总结

GitAI提供4种方式管理ticket号:

1. **命令行参数** - `--ticket PROJ-123`
2. **分支名自动提取** - 从 `feature/PROJ-123-xxx` 提取
3. **交互式输入** - 提示用户输入
4. **配置默认值** - 使用 `ticket_prefix` 自动补全

选择适合你团队的方式,让commit规范更自动化！
