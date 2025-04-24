# 🚀 Kinetic

> Build Avalanche dApps at light speed.

Kinetic is a powerful development environment for Avalanche and Subnet development, available as both a native GUI application and a CLI tool for power users and automation.

## ✨ Features

- 🖥️ **Native GUI** – Streamlined desktop application for managing your Avalanche development
- ⌨️ **Power User CLI** – Full featured command-line interface for automation and scripting
- 🐳 **Local Validator** – Preconfigured Avalanche node via Docker
- 📝 **Smart Contract Templates** – Ready-to-use ERC20, ERC721, and more
- 💧 **Built-in Faucet** – Test tokens for development
- 🧬 **Subnet Wizard** – Create and deploy custom Subnets with a visual interface
- 🧪 **Testing Suite** – Comprehensive testing utilities and guides

## 🚦 Quick Start

### GUI Application
```bash
# Download the latest release for your platform:
# - Windows: kinetic-windows.exe
# - macOS: kinetic-macos.dmg
# - Linux: kinetic-linux.AppImage

# Run the application and follow the setup wizard
```

### CLI Usage
```bash
# Install Kinetic CLI
go install github.com/kinetic-dev/kinetic/cmd/kinetic@latest

# Start a local node
kinetic node start

# Create a new subnet
kinetic subnet create mysubnet --vm subnet-evm

# Deploy a contract
kinetic contract deploy --template erc20 --name "Test Token"

# Get help for any command
kinetic --help
kinetic <command> --help
```

## 🏗 Architecture

Kinetic consists of:

- **Native GUI**: Cross-platform desktop application for intuitive interaction
- **CLI Interface**: Comprehensive command-line tool for automation
- **Docker Integration**: Manages local Avalanche nodes and development tools
- **Smart Contract Management**: Template system and deployment tools
- **Subnet Tools**: Visual and CLI-based subnet creation and management

## 📦 Prerequisites

- Docker and Docker Compose
- 4GB RAM minimum (8GB recommended)
- 20GB free disk space

## 🔧 Development Setup

1. Install Go 1.21 or later
2. Clone the repository:
```bash
git clone https://github.com/kinetic-dev/kinetic
cd kinetic
```

3. Build from source:
```bash
# Build both GUI and CLI
go build ./cmd/kinetic

# Build CLI only
go build ./cmd/kinetic-cli
```

## 📚 Documentation

For detailed documentation, visit [docs link coming soon]

### CLI Reference
```bash
# Node Management
kinetic node start              # Start local node
kinetic node stop               # Stop local node
kinetic node status            # Check node status

# Subnet Management
kinetic subnet list            # List all subnets
kinetic subnet create          # Create new subnet
kinetic subnet deploy          # Deploy subnet to network

# Contract Operations
kinetic contract list          # List available templates
kinetic contract create        # Create from template
kinetic contract deploy        # Deploy to network
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details. 