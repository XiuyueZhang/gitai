# GitAI Quick Start Guide

Get started with GitAI in 5 minutes!

## Prerequisites Check

Before you begin, ensure you have:

- [ ] Go 1.21 or higher installed (`go version`)
- [ ] Git installed (`git --version`)
- [ ] Ollama installed and running

## Step 1: Install Ollama

### macOS/Linux
```bash
curl -fsSL https://ollama.com/install.sh | sh
```
```brew
brew install ollama
```

### Windows
Download from [ollama.com/download](https://ollama.com/download)

### Start Ollama
```bash
ollama serve
```

Keep this terminal running!

## Step 2: Pull AI Model

In a new terminal:

```bash
# Recommended: Fast and accurate for code
ollama pull qwen2.5-coder:7b

# Alternative: Smaller, faster
ollama pull mistral:7b
```

This will download ~4-5GB. Wait for it to complete.

## Step 3: Build GitAI

```bash
# Navigate to the gitai directory
cd /path/to/gitai

# Download dependencies
go mod download

# Build
go build -o gitai

# Optional: Install to system PATH
sudo mv gitai /usr/local/bin/
```

## Step 4: Test It Out

### Create a test repository
```bash
mkdir test-repo
cd test-repo
git init
```

### Create some changes
```bash
echo "# Test Project" > README.md
echo "console.log('Hello');" > app.js
git add .
```

### Run GitAI
```bash
gitai commit
```

### Follow the prompts:
1. Select commit type (e.g., "feat")
2. Enter scope (or leave empty)
3. Wait for AI to generate message
4. Choose "Use this message"

Done! Check your commit:
```bash
git log
```

## Step 5: Use in Real Projects

### Navigate to your project
```bash
cd ~/my-project
```

### Stage changes
```bash
git add src/feature.js src/styles.css
```

### Generate commit
```bash
gitai commit
```

## Common Usage Patterns

### Quick commit with type specified
```bash
gitai commit --type fix --scope api
```

### Just generate, don't commit
```bash
gitai generate
```

### Preview without committing
```bash
gitai commit --dry-run
```

### Use different model
```bash
gitai commit --model mistral:7b
```

## Customize for Your Project

### Create config file
```bash
cd ~/my-project
gitai config --init
```

### Edit `.gitcommit.yaml`
```yaml
model: "qwen2.5-coder:7b"
language: "en"

scopes:
  - "frontend"
  - "backend"
  - "api"
  - "database"
```

### Now GitAI will suggest these scopes!
```bash
gitai commit
```

## Troubleshooting

### "Cannot connect to Ollama"
- Make sure `ollama serve` is running
- Check: `curl http://localhost:11434/api/tags`

### "Model not found"
- Pull the model: `ollama pull qwen2.5-coder:7b`
- Check available models: `ollama list`

### "No staged changes"
- Stage files first: `git add <files>`
- Check: `git status`

### Slow generation
- Use a smaller model: `--model mistral:7b`
- Check Ollama logs for issues

## Next Steps

1. **Read the full README**: Check [README.md](README.md) for advanced features
2. **Customize config**: Add project-specific scopes and types
3. **Create aliases**: Add to your shell:
   ```bash
   alias gai="gitai commit"
   ```

## Tips for Better Results

1. **Write clear code changes**: The AI analyzes your diff
2. **Stage related changes together**: Don't mix unrelated changes
3. **Check recent commits**: AI learns from your commit history
4. **Add a good README**: Helps AI understand your project

## Example Workflow

```bash
# Daily workflow
cd ~/my-project

# Make changes
vim src/auth.js

# Stage specific files
git add src/auth.js

# Let AI generate commit
gitai commit

# Select type: fix
# Select scope: auth
# Review and use message

# Push
git push
```

## Getting Help

```bash
# Show all commands
gitai --help

# Show commit options
gitai commit --help

# Show config options
gitai config --help
```

## Success!

You're now ready to use GitAI for all your commits!

Pro tip: The more you use it, the better it gets at matching your project's commit style.

Happy committing! 🚀
