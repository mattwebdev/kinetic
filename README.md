# 🚀 Kinetic

> Build Avalanche dApps at light speed.

Kinetic is a powerful development environment for Avalanche and Subnet development, available as both a native GUI application and a CLI tool for power users and automation.

## ✨ Features

- 🖥️ **Native GUI** – Streamlined desktop application for managing your Avalanche development
- ⌨️ **Power User CLI** – Full featured command-line interface for automation and scripting
- 🐳 **Local Validator** – Preconfigured Avalanche node via Docker
- 📝 **Smart Contract Templates** – Ready-to-use templates with customizable features:
  - ERC20 tokens with optional cap, minting, burning, and pause functionality
  - ERC721 NFTs with configurable supply, URI handling, and minting controls
  - Basic contracts with storage, events, whitelist, and emergency features
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

# Check node status
kinetic node status

# Create contracts from templates
kinetic contract create ERC20 MyToken --is-mintable --is-burnable
kinetic contract create ERC721 MyNFT --has-max-supply --has-base-uri
kinetic contract create Basic MyContract --has-storage --has-whitelist

# Deploy a contract
kinetic contract deploy MyToken --network local

# Get help for any command
kinetic --help
kinetic <command> --help
```

## 📝 Contract Templates

Kinetic provides flexible smart contract templates that can be customized via command-line flags:

### ERC20 Token
```bash
# Basic ERC20 token
kinetic contract create ERC20 MyToken

# Token with supply cap and minting
kinetic contract create ERC20 MyToken --has-cap --is-mintable

# Full featured token
kinetic contract create ERC20 MyToken \
  --has-cap \
  --is-mintable \
  --is-burnable \
  --is-pausable
```

### ERC721 NFT
```bash
# Basic NFT collection
kinetic contract create ERC721 MyNFT

# NFT with max supply and custom URIs
kinetic contract create ERC721 MyNFT \
  --has-max-supply \
  --has-custom-uri

# Full featured NFT
kinetic contract create ERC721 MyNFT \
  --has-max-supply \
  --has-base-uri \
  --has-custom-uri \
  --is-mintable \
  --is-burnable \
  --is-pausable
```

### Basic Contract
```bash
# Simple storage contract
kinetic contract create Basic MyContract --has-storage

# Contract with whitelist and emergency features
kinetic contract create Basic MyContract \
  --has-storage \
  --has-whitelist \
  --has-emergency-stop \
  --has-emergency-withdraw
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
- Go 1.21 or later (for building from source)

## 🔧 Development Setup

1. Clone the repository:
```bash
git clone https://github.com/kinetic-dev/kinetic
cd kinetic
```

2. Build from source:
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

# Contract Management
kinetic contract list          # List available templates
kinetic contract create        # Create from template
  --output-dir                 # Specify output directory (default: current directory)
  --has-*                      # Template-specific features
  --is-*                      # Token capabilities
kinetic contract deploy        # Deploy to network
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details. 