# 🍺 BrewPy

Homebrew Python Version Manager

## Why?

I was previously using `pyenv`, but it would break everytime there was a new version of Python from Homebrew. I would have to run `brew pyenv-sync` and `pyenv global` each time to get it working again. I wanted a solution that does one thing and does it well: a Python version manager designed **only** for Homebrew on macOS. 

## ✨ Features

- 🚀 **Fast & Lightweight** - Built in Go for maximum performance
- 🔗 **Symlink-based** - Clean PATH management without shell aliases
- 🛠️ **Simple Setup** - One command is all you need to get started

## 🚀 Installation

### Via Homebrew (Recommended)

```bash
brew install landoncrabtree/brewpy/brewpy
```

### From Source

```bash
# Clone the repository
git clone https://github.com/landoncrabtree/brewpy.git
cd brewpy

# Build the binary
make build

# Move to PATH
sudo mv brewpy /usr/local/bin/
```

## 📋 Usage

```bash
# List available Python versions
brewpy versions

# Switch to a specific version
brewpy use Python3.11

# Interactive version selection
brewpy use

# Show current active version
brewpy current
```

## 🔧 How it Works

BrewPy manages Python versions by:

1. **Detecting Architecture** - Automatically finds the correct Homebrew path
   - Apple Silicon: `/opt/homebrew/bin`
   - Intel: `/usr/local/bin`

2. **Creating Symlinks** - Links executables in `~/.brewpy/shims/`
   - `python` → `python3.11`
   - `python3` → `python3.11`
   - `pip` → `pip3.11`
   - `pip3` → `pip3.11`

3. **PATH Management** - Prepends shims directory to PATH

## 🛠️ Requirements

- macOS (Intel or Apple Silicon)
- Homebrew
- Python versions installed via Homebrew

## 🚀 Install Python Versions

```bash
# Install multiple Python versions
brew install python@3.9 python@3.10 python@3.11 python@3.12

# Then use brewpy to switch between them
brewpy use
```

## 📝 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📞 Support

If you encounter issues:
- Check that Python versions are installed via Homebrew
- Ensure your shell profile sources the brewpy init
- Restart your terminal after making changes
    - `rehash -f` to force symlink reload
    - `source ~/.zshrc` to reload your shell profile
