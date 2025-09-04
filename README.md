# tago

Smart Git tag version management tool built with Golang.

---

## CHINESE README

[中文说明](README.zh.md)

## Key Features

🏷️ **Smart Tag Management**: Auto create and bump semantic version tags for Git repositories  
⚡ **Version Base System**: Configurable version carry-over rules (1/10/100)  
🎯 **Interactive Confirmation**: User confirmation for low version base, auto mode for higher bases  
🌍 **Submodule Support**: Independent tag management for main project and submodules  
📋 **Semantic Versioning**: Follows v{major}.{minor}.{patch} format standards

## Installation

```bash
go install github.com/go-mate/tago/cmd/tago@latest
```

# Usage

## Show tags
```bash
tago
```

output:
```
refs/tags/v0.0.0 Wed Feb 12 16:18:18 2025 +0700
refs/tags/v0.0.1 Thu Feb 13 16:43:08 2025 +0700
refs/tags/v0.0.2 Thu Feb 13 18:43:40 2025 +0700
refs/tags/v0.0.3 Wed Apr 30 15:18:56 2025 +0700
refs/tags/v0.0.4 Wed May 7 18:38:38 2025 +0700
```

### Bump Tag Version (Interactive Mode)

Bump from va.b.c to va.b.c+1 and push new tag with user confirmation:

```bash
tago bump
```

Output:
```
cd xxx && git push origin v0.0.5
```

### Bump Tag Version (Auto Mode)

Bump from va.b.c to va.b.c+1 and push new tag without user confirmation:

```bash
tago bump -b=100
```

Output:
```
cd xxx && git push origin v0.0.5
```

### Main Project Tag Management

For main project root DIR tag operations:

```bash
tago bump main
tago bump main -b=10
```

### Submodule Tag Management

For submodule DIR tag operations (with path prefix):

```bash
cd submodule-dir
tago bump sub-module
tago bump sub-module -b=100
```

## Version Base System

The version base (-b parameter) controls version carry-over rules:

- **0 or 1**: Interactive mode, requires user confirmation for each operation
- **≥ 2**: Auto mode, supports automatic version carry-over

### Version Carry-over Examples

With version base 10:
- `v1.2.9` → `v1.3.0` (patch reaches base, carries to minor)
- `v1.9.8` → `v1.9.9` (normal increment)
- `v1.9.9` → `v2.0.0` (minor reaches base, carries to major)

## Advanced Usage

### Command Examples

```bash
# Show all current tags
tago

# Bump tag (with confirmation)
tago bump

# Quick bump without confirmation
tago bump -b=100

# Main project tag bump
tago bump main -b=10

# Submodule tag bump (run from submodule DIR)
cd my-submodule
tago bump sub-module -b=10
```

### Version Control Workflow

1. **After development**: Run `tago` to view current tags
2. **Create new version**: Run `tago bump` to upgrade version
3. **Automation scenarios**: Use `tago bump -b=100` to skip confirmation
4. **Multi-module projects**: Use corresponding subcommands in different directories

## Technical Features

### Smart Version Management
- Auto parse existing tag formats
- Support semantic versioning standards
- Handle version number carry-over logic
- Validate Git repository status

### Flexible Confirmation System
- Low version base: Interactive confirmation for each operation
- High version base: Auto execution, suitable for scripting
- User-friendly prompt messages
- Operation cancellation support

### Multi-project Architecture Support
- Main project tags: `v{major}.{minor}.{patch}`
- Submodule tags: `{path}/v{major}.{minor}.{patch}`
- Path-aware tag management
- Git submodule compatibility

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Welcome contributions! Report bugs, suggest features, contribute code:

- 🐛 **Found a bug?** Submit an issue on GitHub with reproduction steps
- 💡 **Feature idea?** Create an issue to discuss your thoughts
- 📖 **Documentation unclear?** Report issues to help us improve docs
- 🚀 **Need a feature?** Share your use case to help us understand the need
- ⚡ **Performance issue?** Report slow operations to help us optimize
- 🔧 **Configuration trouble?** Ask questions about complex setups
- 📢 **Stay updated?** Watch the repository for new releases and features
- 🌟 **Success story?** Share how this package improved your workflow
- 💬 **General feedback?** All suggestions and feedback are welcome

---

## Code Contributions

For new code contributions, please follow this process:

1. **Fork**: Fork the repository on GitHub (use the web interface)
2. **Clone**: Clone your fork (`git clone https://github.com/yourname/repo-name.git`)
3. **Navigate**: Enter the cloned directory (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`)
5. **Code**: Implement your changes and write comprehensive tests
6. **Test**: (Golang projects) Ensure tests pass (`go test ./...`) and follow Go style conventions
7. **Document**: Update documentation for user-facing changes and use meaningful commit messages
8. **Stage**: Stage your changes (`git add .`)
9. **Commit**: Commit your changes (`git commit -m "Add feature xxx"`) with backward-compatible code
10. **Push**: Push to your branch (`git push origin feature/xxx`)
11. **PR**: Open a Pull Request on GitHub (on GitHub web interface) with detailed description

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

**Project Support:**

- ⭐ **Give a star** if this project helps you
- 🤝 **Share the project** with team members and (golang) programming friends
- 📝 **Write blogs** about development tools and workflows - we provide writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and (golang) development scenarios

**Happy coding with this package!** 🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/tago.svg?variant=adaptive)](https://starchart.cc/go-mate/tago)
