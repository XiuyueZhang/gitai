# GitAI Usage Examples

Real-world examples of using GitAI in different scenarios.

## Basic Usage

### Example 1: Simple Feature Addition

```bash
# You added a new login feature
echo "function login(user, pass) { ... }" >> src/auth.js
git add src/auth.js

# Run GitAI
gitai commit
```

**Interactive Flow:**
```
📝 Git Commit AI Assistant

Changed files (1):
  ✓ src/auth.js (+15, -0)

? Select commit type:
▸ ✨ feat - A new feature
  🐛 fix - A bug fix
  📝 docs - Documentation changes

? Select scope: auth

🤖 Generating commit message...

Generated message:
┌──────────────────────────────────────────┐
│ feat(auth): add login function           │
└──────────────────────────────────────────┘

? What do you want to do?
▸ ✅ Use this message
  🔄 Regenerate
  ✏️  Edit manually
  ❌ Cancel

✨ Commit created successfully!
```

---

### Example 2: Bug Fix with Dry Run

```bash
# Fixed a bug in validation
vim src/validation.js
git add src/validation.js

# Preview without committing
gitai commit --dry-run
```

**Output:**
```
🔍 Dry-run mode - no commit will be created

Changed files (1):
  ✓ src/validation.js (+3, -1)

[... interactive selection ...]

Generated message:
┌──────────────────────────────────────────┐
│ fix(validation): handle null values      │
└──────────────────────────────────────────┘

Would commit with message:
fix(validation): handle null values
```

---

### Example 3: Skip Interactive Selection

```bash
# Update documentation
vim README.md
git add README.md

# Specify type and scope directly
gitai commit --type docs --scope readme
```

Only prompts for confirmation, skips type/scope selection.

---

## Advanced Usage

### Example 4: Multiple Files Refactoring

```bash
# Refactored auth system
git add src/auth/*.js
git add src/middleware/auth.js

gitai commit
```

**Generated Message:**
```
refactor(auth): restructure authentication module
```

AI understands the context from:
- Multiple related files
- File paths (auth directory)
- Actual code changes in diff

---

### Example 5: Using Different Models

```bash
# Use faster model for quick commits
gitai commit --model mistral:7b

# Use more powerful model for complex changes
gitai commit --model codellama:13b
```

---

### Example 6: Chinese Language

```bash
# Generate Chinese commit message
git add .
gitai commit --language zh
```

**Example Output:**
```
Generated message:
┌──────────────────────────────────────────┐
│ feat(api): 添加用户认证接口              │
└──────────────────────────────────────────┘
```

---

## Project-Specific Configuration

### Example 7: Frontend Project

Create `.gitcommit.yaml`:

```yaml
model: "qwen2.5-coder:7b"
language: "en"

types:
  - name: "feat"
    desc: "New feature"
    emoji: "✨"
  - name: "fix"
    desc: "Bug fix"
    emoji: "🐛"
  - name: "style"
    desc: "UI/UX changes"
    emoji: "💄"
  - name: "perf"
    desc: "Performance improvement"
    emoji: "⚡"

scopes:
  - "components"
  - "pages"
  - "hooks"
  - "styles"
  - "api"
  - "utils"
```

Now when you run `gitai commit`, it suggests these scopes!

---

### Example 8: Backend API Project

```yaml
model: "qwen2.5-coder:7b"
language: "en"

scopes:
  - "auth"
  - "database"
  - "api/users"
  - "api/posts"
  - "middleware"
  - "models"
  - "routes"
  - "services"

template: "{type}({scope}): {message}"
```

---

## Team Workflows

### Example 9: Consistent Team Messages

Team config in repo (`.gitcommit.yaml`):

```yaml
model: "qwen2.5-coder:7b"
language: "en"

# Team's commit convention
scopes:
  - "frontend"
  - "backend"
  - "devops"
  - "docs"
  - "tests"

# Custom prompt for team standards
custom_prompt: |
  Our team follows these rules:
  - Use imperative mood (e.g., "add" not "added")
  - Max 50 characters for subject line
  - Reference issue numbers when applicable
```

Everyone on the team gets consistent commit messages!

---

### Example 10: Monorepo Setup

```yaml
scopes:
  - "packages/ui"
  - "packages/api"
  - "packages/shared"
  - "apps/web"
  - "apps/mobile"
  - "tools"
  - "docs"
```

---

## Common Scenarios

### Example 11: First Commit

```bash
# Initialize new project
git init
echo "# My Project" > README.md
git add README.md

gitai commit --type chore --scope init
```

**Generated:**
```
chore(init): initialize project with README
```

---

### Example 12: Merge Conflict Resolution

```bash
# After resolving merge conflicts
git add .
gitai commit --type fix --scope merge
```

---

### Example 13: Version Bump

```bash
# Updated package.json version
vim package.json
git add package.json

gitai commit --type chore
```

**Generated:**
```
chore: bump version to 1.2.0
```

---

### Example 14: Dependencies Update

```bash
# Updated dependencies
npm update
git add package.json package-lock.json

gitai commit --type chore --scope deps
```

**Generated:**
```
chore(deps): update dependencies
```

---

## Generate-Only Mode

### Example 15: Get Ideas Without Committing

```bash
# Just want to see what AI suggests
git add src/feature.js

gitai generate
```

**Output:**
```
Generated message:
┌──────────────────────────────────────────┐
│ feat(feature): add new feature X         │
└──────────────────────────────────────────┘

Copy this message and use it with: git commit -m "<message>"
```

---

### Example 16: Multiple Message Options

```bash
# Generate, copy message, then commit manually
gitai generate --type feat --scope api

# If you like it:
git commit -m "feat(api): add new endpoint"

# Or modify it:
git commit -m "feat(api): add user listing endpoint with pagination"
```

---

## Integration with Git Aliases

### Example 17: Create Git Alias

Add to `~/.gitconfig`:

```ini
[alias]
    ai = !gitai commit
    aigen = !gitai generate
    aidry = !gitai commit --dry-run
```

Usage:
```bash
git add .
git ai              # Run gitai commit
git aigen           # Just generate
git aidry           # Dry run
```

---

### Example 18: Shell Alias

Add to `~/.bashrc` or `~/.zshrc`:

```bash
alias gai='gitai commit'
alias gaigen='gitai generate'
alias gaidry='gitai commit --dry-run'
alias gaifix='gitai commit --type fix'
alias gaifeat='gitai commit --type feat'
```

Usage:
```bash
git add .
gai                 # Quick commit
gaifix              # Quick fix commit
gaifeat --scope api # Quick feature commit
```

---

## Troubleshooting Examples

### Example 19: Model Not Found

```bash
$ gitai commit
Error: model 'qwen2.5-coder:7b' not found
Install it with:
  $ ollama pull qwen2.5-coder:7b
```

**Solution:**
```bash
ollama pull qwen2.5-coder:7b
gitai commit  # Try again
```

---

### Example 20: Ollama Not Running

```bash
$ gitai commit
Error: Cannot connect to Ollama
Please make sure Ollama is running:
  $ ollama serve
```

**Solution:**
```bash
# Terminal 1
ollama serve

# Terminal 2
gitai commit
```

---

## Real-World Examples

### Example 21: E-commerce Project

```bash
# Added shopping cart feature
git add src/components/Cart.tsx
git add src/hooks/useCart.ts
git add src/api/cart.ts

gitai commit
```

**Generated:**
```
feat(cart): implement shopping cart functionality
```

---

### Example 22: Database Migration

```bash
# Created new migration
git add migrations/20250106_add_users_table.sql

gitai commit --type feat --scope database
```

**Generated:**
```
feat(database): add users table migration
```

---

### Example 23: Security Patch

```bash
# Fixed security vulnerability
git add src/auth/password.js

gitai commit --type fix --scope security
```

**Generated:**
```
fix(security): patch password hashing vulnerability
```

---

## Tips for Best Results

### ✅ Do This

```bash
# Stage related changes together
git add src/auth/login.js src/auth/logout.js
gitai commit --type feat --scope auth
```

### ❌ Avoid This

```bash
# Don't mix unrelated changes
git add src/auth/login.js src/ui/button.css src/api/posts.js
gitai commit  # Will generate confusing message
```

---

### ✅ Incremental Commits

```bash
# Commit 1: Feature
git add src/feature.js
gitai commit --type feat

# Commit 2: Tests
git add tests/feature.test.js
gitai commit --type test

# Commit 3: Docs
git add README.md
gitai commit --type docs
```

Better than one large commit!

---

## Summary

GitAI works best when:
- ✅ Changes are focused and related
- ✅ Project has a README
- ✅ Commit history is consistent
- ✅ Configuration matches project structure
- ✅ Staged changes are meaningful

Happy committing! 🚀
