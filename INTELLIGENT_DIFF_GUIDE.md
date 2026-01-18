# 🧠 智能 Diff 分析指南

## 概述

GitAI 的智能 Diff 分析功能能够深入理解代码变更，即使在 diff 很大时也能生成准确的提交消息。

## 为什么需要智能 Diff 分析？

### 问题场景

当你进行以下类型的变更时，简单的 diff 截断会丢失重要信息：

1. **大型重构** - 重命名函数、重组代码结构
2. **多文件变更** - 跨多个文件的功能实现
3. **依赖更新** - 添加/删除/更新导入包
4. **混合变更** - 同时包含功能、测试、文档的变更

### 传统方式的局限

```
❌ 简单截断前 2000 字符
   → 可能只看到文件头，错过实际变更
   → 无法理解变更的整体结构
   → 生成的提交消息不准确
```

### 智能分析的优势

```
✅ 智能 Diff 分析
   → 提取每个文件的摘要
   → 识别关键代码变更（函数、类）
   → 检测依赖变化
   → 分析变更复杂度
   → 生成准确、有意义的提交消息
```

## 工作原理

### 1. 文件级分析

对每个修改的文件进行独立分析：

```
📄 src/auth/login.go [modified] +45/-12
   Key changes: function Login, function ValidateToken

📄 src/auth/user.go [added] +67/-0
   Key changes: type User, function NewUser

📄 go.mod [modified] +2/-1
   Import changes: + github.com/golang-jwt/jwt/v5
```

### 2. 代码模式识别

识别关键代码结构：

**Go:**
- `func NewAuthenticator()` → function NewAuthenticator
- `type User struct` → type User
- `type Handler interface` → type Handler

**JavaScript/TypeScript:**
- `function handleLogin()` → function handleLogin
- `class AuthService` → class AuthService
- `const validateUser = async` → function validateUser

**Python:**
- `def authenticate(user)` → function authenticate
- `class UserManager` → class UserManager
- `async def login()` → function login

### 3. 导入变更检测

自动识别依赖变化：

```go
+ import "github.com/golang-jwt/jwt/v5"
+ import "golang.org/x/crypto/bcrypt"
```

```javascript
+ import { useState, useEffect } from 'react'
+ const axios = require('axios')
```

### 4. 智能截断

当 diff 超过限制时，优先保留：

1. **文件头信息** - 文件名、修改类型
2. **关键代码变更** - 函数定义、类定义
3. **重要文件优先** - 非测试文件 > 测试文件
4. **变更较大的文件** - 更多改动的文件优先

### 5. 复杂度评估

```
simple    - < 100 行变更，< 3 个文件
moderate  - 100-500 行变更，3-10 个文件
complex   - > 500 行变更，> 10 个文件
```

## 配置选项

### 基础配置

```yaml
# .gitcommit.yaml
diff_analysis:
  enabled: true                    # 启用智能分析
  include_function_names: true     # 提取函数名
  include_imports: true            # 提取导入变更
  smart_truncate: true             # 智能截断
  context_lines: 3                 # 上下文行数
```

### 配置说明

| 选项 | 默认值 | 说明 |
|------|--------|------|
| `enabled` | `true` | 是否启用智能分析 |
| `include_function_names` | `true` | 是否提取函数/类名 |
| `include_imports` | `true` | 是否提取导入变更 |
| `smart_truncate` | `true` | 是否使用智能截断 |
| `context_lines` | `3` | diff 块中的上下文行数 |

### 性能调优

**对于大型项目**（slow performance）：

```yaml
diff_analysis:
  enabled: true
  include_function_names: false    # 关闭函数提取以提速
  include_imports: false           # 关闭导入分析
  smart_truncate: true
  context_lines: 1                 # 减少上下文行数
```

**对于小型项目**（want maximum accuracy）：

```yaml
diff_analysis:
  enabled: true
  include_function_names: true
  include_imports: true
  smart_truncate: true
  context_lines: 5                 # 增加上下文
```

**完全禁用**（use simple diff）：

```yaml
diff_analysis:
  enabled: false
```

## 实际示例

### 示例 1: 添加新功能

**变更内容：**
- 添加用户认证模块
- 新增 3 个文件
- 总计 +245/-0 行

**传统方式生成：**
```
feat: update code
```

**智能分析生成：**
```
feat(auth): add user authentication system

- Implement JWT-based authentication
- Add User model and Authenticator service
- Include bcrypt password hashing
- Add jwt dependency for token management
```

### 示例 2: 重构代码

**变更内容：**
- 重构认证逻辑
- 重命名多个函数
- 跨 5 个文件
- 总计 +180/-165 行

**传统方式生成：**
```
refactor: various changes
```

**智能分析生成：**
```
refactor(auth): restructure authentication flow

- Rename ValidateUser to AuthenticateUser
- Extract token validation to separate function
- Move session management to dedicated service
- Improve error handling consistency
```

### 示例 3: 依赖更新

**变更内容：**
- 更新 JWT 库
- 修改相关代码适配新 API
- 总计 +15/-18 行

**传统方式生成：**
```
chore: update dependencies
```

**智能分析生成：**
```
build: upgrade jwt library to v5

- Update github.com/golang-jwt/jwt/v4 to v5
- Adapt token validation to new API
- Update imports across auth package
```

## 支持的语言

### 完全支持（函数提取 + 导入检测）

- ✅ Go
- ✅ JavaScript / TypeScript
- ✅ Python
- ✅ Java
- ✅ Rust
- ✅ C / C++
- ✅ C#
- ✅ Ruby
- ✅ PHP
- ✅ Swift
- ✅ Kotlin

### 部分支持（仅导入检测）

- ✅ Shell / Bash
- ✅ SQL
- ✅ HTML / CSS

### 配置文件识别

自动识别并标记配置文件：
- `package.json`, `go.mod`, `Cargo.toml`
- `requirements.txt`, `Gemfile`, `pom.xml`
- `.yaml`, `.json`, `.toml`, `.env`
- `Dockerfile`, `Makefile`

## 最佳实践

### 1. 保持默认配置

对于大多数项目，默认配置已经足够好：

```yaml
diff_analysis:
  enabled: true
  include_function_names: true
  include_imports: true
  smart_truncate: true
  context_lines: 3
```

### 2. 合理拆分提交

即使有智能分析，也建议：

```bash
# 好 - 逻辑独立的小提交
git add src/auth/*.go
gitai commit  # "feat(auth): add authentication"

git add src/user/*.go
gitai commit  # "feat(user): add user management"

# 避免 - 混合不相关的变更
git add .
gitai commit  # 太多不相关的变更
```

### 3. 检查生成的消息

智能分析提高了准确性，但仍建议：

1. 查看生成的提交消息
2. 确认是否准确反映了变更
3. 必要时使用"编辑"功能调整
4. 使用"重新生成"获取不同的表述

### 4. 针对项目调整

根据项目特点调整配置：

**微服务项目**（many small files）：
```yaml
diff_analysis:
  enabled: true
  context_lines: 2  # 减少上下文，关注变更本身
```

**框架开发**（complex changes）：
```yaml
diff_analysis:
  enabled: true
  include_function_names: true
  context_lines: 5  # 更多上下文帮助理解
```

## 故障排查

### 问题：生成的消息仍然不准确

**可能原因：**
1. Diff 太大，即使智能截断也丢失了关键信息
2. 代码变更过于复杂
3. AI 模型对特定领域不熟悉

**解决方案：**
```bash
# 1. 拆分提交
git reset HEAD~1
git add src/feature1/
gitai commit

git add src/feature2/
gitai commit

# 2. 增加 max_diff_length
# .gitcommit.yaml
max_diff_length: 4000  # 增加到 4000

# 3. 使用更好的模型
# .gitcommit.yaml
model: "qwen2.5-coder:32b"  # 更大的模型
```

### 问题：分析速度慢

**可能原因：**
1. 变更文件太多
2. 函数提取过程耗时

**解决方案：**
```yaml
# .gitcommit.yaml
diff_analysis:
  enabled: true
  include_function_names: false  # 关闭函数提取
  include_imports: false         # 关闭导入分析
  context_lines: 1               # 减少上下文
```

### 问题：提示信息太长

**可能原因：**
分析结果包含太多详细信息

**解决方案：**
```yaml
# .gitcommit.yaml
max_diff_length: 1500           # 减少 diff 长度
diff_analysis:
  enabled: true
  smart_truncate: true
  context_lines: 2              # 减少上下文
```

## 技术细节

### 分析流程

```
1. 获取完整 diff
   ↓
2. 按文件分割
   ↓
3. 对每个文件进行分析
   - 检测文件类型
   - 统计增删行数
   - 提取关键代码变更
   - 提取导入变更
   ↓
4. 汇总分析结果
   - 计算总体复杂度
   - 排序文件（重要性）
   ↓
5. 生成智能 diff
   - 添加摘要头
   - 添加文件摘要
   - 选择重要 diff 块
   - 智能截断到限制长度
   ↓
6. 传递给 AI
```

### 内存使用

- **小型变更** (< 100 行): ~50KB
- **中型变更** (100-500 行): ~200KB
- **大型变更** (> 500 行): ~500KB (截断后)

### 性能基准

| 变更规模 | 文件数 | 分析时间 |
|---------|--------|---------|
| 小型 | 1-3 | < 10ms |
| 中型 | 3-10 | 10-50ms |
| 大型 | 10-50 | 50-200ms |
| 超大型 | > 50 | 200-500ms |

## 未来增强

计划中的功能：

- [ ] **AST 深度分析** - 使用抽象语法树理解语义变更
- [ ] **变更模式识别** - 自动识别常见重构模式
- [ ] **相关文件分析** - 分析变更的影响范围
- [ ] **智能分组** - 自动将相关变更分组
- [ ] **自定义规则** - 允许用户定义提取规则

## 总结

智能 Diff 分析显著提高了 GitAI 对代码变更的理解能力，特别是在以下场景：

✅ **大型重构** - 准确识别重构意图
✅ **多文件变更** - 理解跨文件的功能实现
✅ **依赖更新** - 检测和说明依赖变化
✅ **复杂变更** - 提供准确的复杂度评估

通过合理配置和使用，可以让 AI 生成更加准确、有意义的提交消息！

---

**相关文档：**
- [配置示例](.gitcommit.example.yaml)
- [主要文档](README.md)
- [多语言支持](MULTILINGUAL_GUIDE.md)
