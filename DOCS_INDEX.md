# GitAI Documentation Index

Welcome to GitAI! This guide helps you find the right documentation for your needs.

## 🚀 Quick Links

| I want to... | Read this |
|-------------|-----------|
| **Get started in 5 minutes** | [QUICKSTART.md](QUICKSTART.md) |
| **Install GitAI** | [INSTALLATION.md](INSTALLATION.md) |
| **Understand what GitAI does** | [README.md](README.md) |
| **See usage examples** | [EXAMPLES.md](EXAMPLES.md) |
| **Learn technical details** | [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) |
| **Check if everything is implemented** | [CHECKLIST.md](CHECKLIST.md) |

---

## 📚 Documentation Overview

### For First-Time Users

**Start here** → [QUICKSTART.md](QUICKSTART.md) (3.7 KB)
- 5-minute setup guide
- Your first commit with AI
- Common usage patterns
- Troubleshooting basics

**Then read** → [README.md](README.md) (8.5 KB)
- Complete feature overview
- Installation options
- Configuration guide
- FAQ

---

### For Daily Users

**Bookmark** → [EXAMPLES.md](EXAMPLES.md) (9.0 KB)
- 20+ real-world scenarios
- Configuration examples
- Team workflows
- Integration tips
- Git aliases

**Configure** → [.gitcommit.example.yaml](.gitcommit.example.yaml) (1.9 KB)
- Example configuration file
- All available options
- Copy and customize for your project

---

### For Developers

**Architecture** → [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) (8.2 KB)
- Technical overview
- Project structure
- Implementation details
- Code statistics

**Status** → [CHECKLIST.md](CHECKLIST.md) (7.2 KB)
- Complete feature list
- Implementation status
- Testing coverage
- Quality metrics

---

### For Installation

**Complete Guide** → [INSTALLATION.md](INSTALLATION.md) (7.5 KB)
- All platforms (macOS, Linux, Windows)
- Prerequisites setup
- Build from source
- Troubleshooting
- System requirements

**Quick Install**:
```bash
cd /Users/yue/go-source/gitai
make build
sudo make install
```

---

## 📖 Document Details

### [README.md](README.md) - Main Documentation
**Size**: 8.5 KB | **Reading time**: 15 min

**Contents**:
- ✨ Features overview
- 📦 Installation methods
- 🎯 Quick start
- ⚙️  Configuration
- 🔧 Troubleshooting
- ❓ FAQ

**Best for**: Understanding the full capabilities

---

### [QUICKSTART.md](QUICKSTART.md) - Getting Started
**Size**: 3.7 KB | **Reading time**: 5 min

**Contents**:
- Step-by-step setup
- First test run
- Real project usage
- Customization basics

**Best for**: Getting up and running fast

---

### [INSTALLATION.md](INSTALLATION.md) - Setup Guide
**Size**: 7.5 KB | **Reading time**: 10 min

**Contents**:
- Prerequisites (Go, Ollama)
- Platform-specific instructions
- Multiple installation methods
- Verification steps
- Troubleshooting

**Best for**: Detailed installation help

---

### [EXAMPLES.md](EXAMPLES.md) - Usage Examples
**Size**: 9.0 KB | **Reading time**: 15 min

**Contents**:
- Basic usage (10 examples)
- Advanced usage (5 examples)
- Project configs (3 examples)
- Team workflows (2 examples)
- Integration examples

**Best for**: Learning by example

---

### [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Technical Docs
**Size**: 8.2 KB | **Reading time**: 15 min

**Contents**:
- Architecture overview
- Project structure
- Implementation details
- Dependencies
- Performance metrics
- Code statistics

**Best for**: Developers and contributors

---

### [CHECKLIST.md](CHECKLIST.md) - Implementation Status
**Size**: 7.2 KB | **Reading time**: 10 min

**Contents**:
- Feature completion (100%)
- Requirements checklist
- Code quality metrics
- Testing status
- Deliverables list

**Best for**: Verifying completeness

---

### [.gitcommit.example.yaml](.gitcommit.example.yaml) - Config Template
**Size**: 1.9 KB | **Type**: YAML

**Contents**:
- All configuration options
- Commented examples
- Default values
- Custom prompt template

**Best for**: Setting up your project

---

## 🎯 Documentation by Task

### "I'm brand new to GitAI"
1. Read: [QUICKSTART.md](QUICKSTART.md)
2. Try: Follow the 5-minute guide
3. Read: [README.md](README.md) - Features section

---

### "I want to install GitAI"
1. Read: [INSTALLATION.md](INSTALLATION.md)
2. Install: Ollama and model
3. Build: `make build && sudo make install`
4. Verify: `gitai --version`

---

### "I want to use GitAI daily"
1. Read: [EXAMPLES.md](EXAMPLES.md)
2. Configure: Copy `.gitcommit.example.yaml` to `.gitcommit.yaml`
3. Customize: Edit scopes and types for your project
4. Use: `gitai commit`

---

### "I need help with configuration"
1. Read: [README.md](README.md) - Configuration section
2. Copy: `.gitcommit.example.yaml` to `.gitcommit.yaml`
3. Examples: See [EXAMPLES.md](EXAMPLES.md) - Project configs
4. Initialize: `gitai config --init`

---

### "Something's not working"
1. Read: [README.md](README.md) - Troubleshooting section
2. Read: [INSTALLATION.md](INSTALLATION.md) - Troubleshooting
3. Check: Ollama is running (`ollama serve`)
4. Check: Model is installed (`ollama list`)

---

### "I want to understand the code"
1. Read: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
2. Read: [CHECKLIST.md](CHECKLIST.md)
3. Browse: Source code in `cmd/` and `internal/`
4. Check: Comments in code files

---

### "I want to contribute"
1. Read: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
2. Read: [CHECKLIST.md](CHECKLIST.md)
3. Run: `make test`
4. Build: `make build`

---

## 📊 Documentation Statistics

| Document | Size | Purpose | Reading Time |
|----------|------|---------|--------------|
| README.md | 8.5 KB | Main docs | 15 min |
| QUICKSTART.md | 3.7 KB | Getting started | 5 min |
| INSTALLATION.md | 7.5 KB | Setup guide | 10 min |
| EXAMPLES.md | 9.0 KB | Usage examples | 15 min |
| PROJECT_SUMMARY.md | 8.2 KB | Technical overview | 15 min |
| CHECKLIST.md | 7.2 KB | Implementation status | 10 min |
| .gitcommit.example.yaml | 1.9 KB | Config template | 2 min |
| **Total** | **46 KB** | **Complete docs** | **~1 hour** |

---

## 🔍 Find Specific Topics

### Configuration
- Main guide: [README.md](README.md#configuration)
- Examples: [EXAMPLES.md](EXAMPLES.md#project-specific-configuration)
- Template: [.gitcommit.example.yaml](.gitcommit.example.yaml)

### Installation
- Complete guide: [INSTALLATION.md](INSTALLATION.md)
- Quick install: [README.md](README.md#installation)
- Prerequisites: [QUICKSTART.md](QUICKSTART.md#prerequisites-check)

### Usage Examples
- Basic: [EXAMPLES.md](EXAMPLES.md#basic-usage)
- Advanced: [EXAMPLES.md](EXAMPLES.md#advanced-usage)
- Daily workflow: [QUICKSTART.md](QUICKSTART.md#example-workflow)

### Troubleshooting
- Common issues: [README.md](README.md#troubleshooting)
- Install issues: [INSTALLATION.md](INSTALLATION.md#troubleshooting-installation)
- Quick fixes: [QUICKSTART.md](QUICKSTART.md#troubleshooting)

### Technical Details
- Architecture: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md#technical-implementation)
- Code structure: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md#project-structure)
- Features: [CHECKLIST.md](CHECKLIST.md#mvp-requirements-all-complete)

---

## 📁 Other Important Files

| File | Description |
|------|-------------|
| [LICENSE](LICENSE) | MIT License (1.0 KB) |
| [Makefile](Makefile) | Build automation (1.9 KB) |
| [go.mod](go.mod) | Go dependencies |
| [.gitignore](.gitignore) | Git ignore rules |
| FILES.txt | File listing |

---

## 🎓 Learning Path

### Beginner
1. ☑️ [QUICKSTART.md](QUICKSTART.md) - 5 min setup
2. ☑️ [README.md](README.md) - Overview (15 min)
3. ☑️ [EXAMPLES.md](EXAMPLES.md) - Basic examples (5 min)

**Time**: ~25 minutes

---

### Intermediate
1. ☑️ [README.md](README.md) - Full read (15 min)
2. ☑️ [EXAMPLES.md](EXAMPLES.md) - All examples (15 min)
3. ☑️ Configure your project (10 min)

**Time**: ~40 minutes

---

### Advanced
1. ☑️ [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) (15 min)
2. ☑️ [CHECKLIST.md](CHECKLIST.md) (10 min)
3. ☑️ Browse source code (20 min)
4. ☑️ Run tests (5 min)

**Time**: ~50 minutes

---

## 💡 Tips

- **New to Git commits?** Start with [README.md](README.md) - Conventional Commits section
- **Team setup?** See [EXAMPLES.md](EXAMPLES.md) - Team Workflows
- **Installation issues?** Check [INSTALLATION.md](INSTALLATION.md) - Troubleshooting
- **Want examples?** Jump to [EXAMPLES.md](EXAMPLES.md)
- **Quick reference?** Keep [README.md](README.md) bookmarked

---

## 🚀 Next Steps

Choose your path:

**Just want to use it?**
→ [QUICKSTART.md](QUICKSTART.md)

**Need detailed install help?**
→ [INSTALLATION.md](INSTALLATION.md)

**Want to learn everything?**
→ [README.md](README.md)

**Looking for examples?**
→ [EXAMPLES.md](EXAMPLES.md)

**Technical deep dive?**
→ [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

---

Happy reading! 📚
