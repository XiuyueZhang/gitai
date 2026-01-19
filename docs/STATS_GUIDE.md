# 📊 GitAI Stats - Commit History Statistics

## 概述

`gitai stats` 命令提供强大的提交历史分析功能，帮助你理解团队的提交模式、改进提交质量，并保持一致性。

## 快速开始

### 基本用法

```bash
# 分析最近 100 条提交（默认）
gitai stats

# 分析最近 500 条提交
gitai stats --limit 500

# 分析最近 50 条提交（简短版）
gitai stats -n 50

# 导出统计数据到 JSON 文件
gitai stats --export stats.json
```

## 功能特性

### 1. 📈 Overview - 总览

提供项目提交的整体概况：

```
📈 Overview
  Total Commits: 297
  Average Subject Length: 52 characters
  With Scope: 120 (40.4%)
  With Body: 89 (30.0%)
  With Ticket: 45 (15.2%)
```

**指标说明：**
- **Total Commits**: 分析的提交总数
- **Average Subject Length**: 提交主题行的平均长度（建议 50-72 字符）
- **With Scope**: 包含 scope 的提交比例（如 `feat(auth):` 中的 `auth`）
- **With Body**: 包含详细说明的提交比例
- **With Ticket**: 包含工单号的提交比例（如 `[JIRA-123]`）

### 2. 🏷️ Commit Types - 提交类型分布

展示最常用的提交类型及其使用频率：

```
🏷️  Commit Types
  feat         [████████████████████] 120 (40.4%)
  fix          [████████████░░░░░░░░]  89 (30.0%)
  docs         [████░░░░░░░░░░░░░░░░]  25 (8.4%)
  refactor     [███░░░░░░░░░░░░░░░░░]  20 (6.7%)
  test         [██░░░░░░░░░░░░░░░░░░]  15 (5.1%)
  chore        [██░░░░░░░░░░░░░░░░░░]  12 (4.0%)
```

**提交类型说明：**
- **feat**: 新功能
- **fix**: Bug 修复
- **docs**: 文档更新
- **refactor**: 代码重构
- **test**: 测试相关
- **chore**: 构建/工具相关
- **perf**: 性能优化
- **style**: 代码格式
- **ci**: CI/CD 相关

### 3. 📦 Top Scopes - 常用作用域

显示最常修改的模块/组件：

```
📦 Top Scopes
  auth            [████████████████░░░]  45
  api             [██████████████░░░░░]  38
  ui              [██████████░░░░░░░░░]  28
  database        [████████░░░░░░░░░░░]  22
  config          [██████░░░░░░░░░░░░░]  15
```

**用途：**
- 识别项目的热点区域
- 了解团队工作重点
- 发现需要重点关注的模块

### 4. 🔤 Common Action Verbs - 常用动词

分析提交消息中最常使用的动词：

```
🔤 Common Action Verbs
  add            85
  update         62
  fix            54
  implement      42
  refactor       28
  improve        25
  remove         18
```

**最佳实践：**
- **add**: 添加新内容（文件、功能、测试）
- **update**: 更新现有内容
- **implement**: 实现功能
- **refactor**: 重构代码
- **improve**: 改进现有功能
- **remove**: 删除内容
- **fix**: 修复问题

### 5. 🌍 Language Usage - 语言使用

显示提交消息使用的语言分布：

```
🌍 Language Usage
  English    [████████████████████] 85.5%
  中文        [████░░░░░░░░░░░░░░░░] 12.1%
  日本語      [█░░░░░░░░░░░░░░░░░░░]  2.4%
```

**用途：**
- 了解团队的语言偏好
- 评估多语言支持的必要性
- 保持提交消息语言的一致性

### 6. ⏰ Commit Time Distribution - 提交时间分布

分析一天中的提交活跃时段：

```
⏰ Commit Time Distribution
  14:00  [████████████████████]  45
  15:00  [██████████████████░░]  38
  11:00  [████████████████░░░░]  35
  16:00  [██████████████░░░░░░]  30
  10:00  [████████████░░░░░░░░]  25
```

**洞察：**
- 识别团队最活跃的时段
- 了解工作习惯
- 优化协作时间安排

### 7. 📅 Commit Day Distribution - 提交星期分布

展示一周中每天的提交频率：

```
📅 Commit Day Distribution
  Monday     [████████████████████]  52
  Tuesday    [██████████████████░░]  48
  Wednesday  [█████████████████░░░]  45
  Thursday   [████████████████░░░░]  42
  Friday     [████████████████░░░░]  40
  Saturday   [████░░░░░░░░░░░░░░░░]  12
  Sunday     [██░░░░░░░░░░░░░░░░░░]   8
```

**用途：**
- 了解团队工作节奏
- 识别加班情况
- 优化开发流程

### 8. 📊 Recent Activity - 近期活跃度

分析最近的提交趋势：

```
📊 Recent Activity
  Last 30 days: 156 commits (5.2/day avg)
  Last 7 days:  42 commits
  Most active:  2026-01-15
```

**指标：**
- **Last 30 days**: 最近 30 天的提交总数
- **Average per day**: 平均每天提交次数
- **Last 7 days**: 最近 7 天的提交数
- **Most active**: 最活跃的日期

### 9. 👥 Top Contributors - 主要贡献者

显示提交最多的贡献者：

```
👥 Top Contributors
  Alice Smith              [████████████████████] 42.3%
  Bob Johnson             [████████████░░░░░░░░] 28.7%
  Carol Williams          [████████░░░░░░░░░░░░] 18.5%
  David Brown             [████░░░░░░░░░░░░░░░░] 10.5%
```

### 10. 📏 Extremes - 极值

展示最长和最短的提交消息：

```
📏 Extremes
  Longest:  feat(auth): implement comprehensive JWT-based authentic...
  Shortest: fix: typo
```

**用途：**
- 识别过长或过短的提交消息
- 改进提交消息质量

### 11. 💡 Top Commit Patterns - 常用模式

识别最常用的提交模式：

```
💡 Your Top Commit Patterns:
  1. feat(auth) - used 45 times
  2. fix(api) - used 38 times
  3. docs - used 25 times
```

### 12. 💭 Insights & Recommendations - 洞察与建议

基于统计数据提供个性化建议：

```
💭 Insights & Recommendations:
  ⚠️  Your average subject line is quite long (>72 chars)
     Consider using shorter, more concise subjects

  💡 You rarely use scopes in commits (<20%)
     Scopes help organize changes by component/module

  💡 Most commits have no body (<10%)
     Consider adding details for non-trivial changes

  🔥 Very high commit frequency (>10/day avg)
     Great activity! Consider squashing related commits
```

## 高级用法

### 导出统计数据

将统计数据导出为 JSON 格式，便于进一步分析或生成报告：

```bash
gitai stats --export stats.json
```

导出的 JSON 包含：
- 总提交数
- 平均长度
- 类型分布
- Scope 分布
- 语言使用情况
- 等等...

### 与 CI/CD 集成

在 CI/CD 流程中使用 stats 命令来监控提交质量：

```yaml
# .github/workflows/commit-quality.yml
name: Commit Quality Check

on:
  schedule:
    - cron: '0 0 * * 0'  # 每周日运行

jobs:
  check-quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 100  # 获取最近 100 条提交

      - name: Install GitAI
        run: |
          curl -sSL https://raw.githubusercontent.com/xyue92/gitai/main/scripts/install.sh | bash

      - name: Generate Stats
        run: |
          gitai stats --limit 100 --export weekly-stats.json

      - name: Upload Stats
        uses: actions/upload-artifact@v4
        with:
          name: commit-stats
          path: weekly-stats.json
```

### 定期报告

创建一个脚本定期生成团队提交统计报告：

```bash
#!/bin/bash
# generate-monthly-report.sh

# 生成当月统计
gitai stats --limit 500 > report-$(date +%Y-%m).txt

# 导出 JSON
gitai stats --limit 500 --export report-$(date +%Y-%m).json

# 可选：发送到 Slack
# curl -X POST -H 'Content-type: application/json' \
#   --data "{\"text\":\"$(cat report-$(date +%Y-%m).txt)\"}" \
#   $SLACK_WEBHOOK_URL
```

## 实际应用场景

### 场景 1: 团队入职培训

新成员加入时，使用 `gitai stats` 展示团队的提交习惯：

```bash
gitai stats --limit 200
```

帮助新成员了解：
- 团队使用的提交类型
- 常见的 scope 命名
- 提交消息的平均长度
- 是否使用工单号

### 场景 2: 代码审查改进

定期运行 stats 命令，识别提交质量问题：

```bash
# 每月初运行
gitai stats --limit 500 --export monthly-report.json
```

关注指标：
- Scope 使用率（建议 >40%）
- Body 使用率（建议 >30%）
- 平均长度（建议 50-72 字符）

### 场景 3: 项目健康监控

在项目 README 中添加提交统计徽章：

```bash
# 生成统计数据
gitai stats --limit 100

# 在 README 中展示关键指标
```

### 场景 4: 性能基准

使用 stats 来评估提交频率和质量：

```bash
# 季度回顾
gitai stats --limit 1000

# 关注：
# - 提交频率趋势
# - 类型分布变化
# - 质量指标改进
```

## 最佳实践

### 1. 定期运行

建议每周或每月运行一次 stats 命令：

```bash
# 每周
gitai stats --limit 50

# 每月
gitai stats --limit 200

# 季度
gitai stats --limit 1000
```

### 2. 团队讨论

在团队会议中分享统计结果：
- 讨论改进点
- 制定提交规范
- 表彰优秀实践

### 3. 持续改进

基于统计洞察制定行动计划：

| 问题 | 行动 |
|------|------|
| Scope 使用率低 | 创建 scope 列表文档 |
| 提交消息过短 | 鼓励添加 body |
| 缺少工单号 | 配置 git hooks |
| 提交频率不均 | 讨论工作流程 |

### 4. 文档化标准

根据 stats 结果更新团队规范：

```markdown
# 我们的提交标准

基于最近的统计分析（gitai stats），我们确定了以下标准：

- ✅ 使用 conventional commits 格式
- ✅ 60% 以上的提交应包含 scope
- ✅ 复杂变更应添加 body 说明
- ✅ 主题行控制在 50-72 字符
- ✅ 关联 JIRA 工单号（如适用）
```

## 故障排查

### 问题：统计数据不准确

**可能原因：**
- 提交格式不规范
- 分析数量不足

**解决方案：**
```bash
# 增加分析数量
gitai stats --limit 500

# 检查提交格式
git log --oneline -20
```

### 问题：某些类型未被识别

**可能原因：**
- 使用了非标准的提交类型

**解决方案：**
统计只识别符合 Conventional Commits 格式的提交：
```
type(scope): subject
```

确保使用标准类型：`feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `style`, `ci`

### 问题：导出失败

**可能原因：**
- 没有写入权限

**解决方案：**
```bash
# 使用完整路径
gitai stats --export ~/reports/stats.json

# 或指定其他目录
gitai stats --export /tmp/stats.json
```

## 与其他工具对比

| 功能 | gitai stats | git-stats | gitinspector |
|------|------------|-----------|--------------|
| 类型分布 | ✅ | ❌ | ❌ |
| Scope 分析 | ✅ | ❌ | ❌ |
| 可视化图表 | ✅ | ✅ | ✅ |
| 语言检测 | ✅ | ❌ | ❌ |
| 洞察建议 | ✅ | ❌ | ❌ |
| JSON 导出 | ✅ | ✅ | ✅ |
| 易于安装 | ✅ | ❌ | ❌ |

## 未来增强

计划中的功能：

- [ ] **交互式 Web 仪表板** - 在浏览器中查看统计
- [ ] **趋势图表** - 显示指标随时间的变化
- [ ] **团队对比** - 对比不同贡献者的风格
- [ ] **质量评分** - 自动评估提交质量
- [ ] **自定义规则** - 定义团队特定的检查规则
- [ ] **Slack/Teams 集成** - 自动发送报告
- [ ] **更多导出格式** - HTML, PDF, Markdown

## 总结

`gitai stats` 是一个强大的工具，帮助你：

✅ **了解团队习惯** - 通过数据洞察提交模式
✅ **改进提交质量** - 基于统计建议优化
✅ **保持一致性** - 识别并推广最佳实践
✅ **监控项目健康** - 跟踪关键指标趋势
✅ **优化工作流程** - 基于活跃度数据调整

开始使用 `gitai stats` 让你的提交更加专业和一致！

---

**相关文档：**
- [GitAI 主文档](../README.md)
- [智能 Diff 分析](../INTELLIGENT_DIFF_GUIDE.md)
- [多语言支持](../MULTILINGUAL_GUIDE.md)
