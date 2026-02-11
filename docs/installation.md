# Installation

Percepta is distributed as a standalone binary for Linux, macOS, and Windows. Choose your preferred installation method below.

## Binary Installation (Recommended)

**1. Download the latest release for your platform:**

Visit [GitHub Releases](https://github.com/Perceptax/percepta/releases) and download the binary for your OS:

- **Linux (x86_64):** `percepta-linux-amd64.tar.gz`
- **Linux (ARM64):** `percepta-linux-arm64.tar.gz`
- **macOS (Intel):** `percepta-darwin-amd64.tar.gz`
- **macOS (Apple Silicon):** `percepta-darwin-arm64.tar.gz`
- **Windows (x86_64):** `percepta-windows-amd64.zip`

**2. Extract the archive:**

```bash
# Linux/macOS
tar -xzf percepta-*.tar.gz

# Windows
# Extract the .zip file using Windows Explorer or:
unzip percepta-windows-amd64.zip
```

**3. Move binary to PATH:**

```bash
# Linux/macOS
sudo mv percepta-* /usr/local/bin/percepta
sudo chmod +x /usr/local/bin/percepta

# Windows
# Move percepta-windows-amd64.exe to a directory in your PATH
# For example: C:\Program Files\Percepta\
```

**4. Verify installation:**

```bash
percepta --help
```

You should see the Percepta CLI help output.

## Build from Source

**Requirements:**
- Go 1.20 or later
- Git

**1. Install Go:**

```bash
# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# macOS
brew install go

# Windows
# Download from https://go.dev/dl/
```

**2. Clone the repository:**

```bash
git clone https://github.com/Perceptax/percepta.git
cd percepta
```

**3. Build:**

```bash
go build -o percepta ./cmd/percepta
```

**4. Install (optional):**

```bash
# Install to GOPATH/bin
go install ./cmd/percepta

# Or use go install directly:
go install github.com/Perceptax/percepta/cmd/percepta@latest
```

**5. Verify:**

```bash
percepta --help
```

## Environment Setup

**Claude API Key (Required):**

Percepta uses the Claude Vision API for hardware observation. Set your Anthropic API key:

```bash
# Linux/macOS (add to ~/.bashrc or ~/.zshrc for persistence)
export ANTHROPIC_API_KEY="your-api-key-here"

# Windows (PowerShell)
$env:ANTHROPIC_API_KEY="your-api-key-here"

# Windows (Command Prompt)
set ANTHROPIC_API_KEY=your-api-key-here
```

**Get an API key:**
1. Sign up at [console.anthropic.com](https://console.anthropic.com)
2. Navigate to API Keys
3. Create a new key
4. Copy and set as environment variable

**Camera Access:**

**Linux:**
- Camera devices are typically at `/dev/video0`, `/dev/video1`, etc.
- Check available cameras: `ls /dev/video*`
- If permission denied, add your user to the `video` group:
  ```bash
  sudo usermod -a -G video $USER
  # Log out and log back in for changes to take effect
  ```

**macOS:**
- Built-in camera is usually device `0`
- macOS will prompt for camera permissions on first use

**Windows:**
- Camera devices are indexed numerically (0, 1, 2, etc.)
- Windows will prompt for camera permissions on first use

## Quick Verification

**1. Check Percepta is installed:**

```bash
percepta --version
```

**2. Check API key is set:**

```bash
# Linux/macOS
echo $ANTHROPIC_API_KEY

# Windows (PowerShell)
echo $env:ANTHROPIC_API_KEY
```

**3. Check camera access:**

```bash
# Linux
ls -l /dev/video0

# macOS/Windows
# Connect a USB camera or use built-in webcam
```

**4. Create a test device:**

```bash
percepta device add test-board
# Follow prompts
```

**5. Test observation (requires hardware):**

```bash
percepta observe test-board
```

## Troubleshooting

**"command not found: percepta"**
- Binary is not in PATH. Either move it to `/usr/local/bin/` or add its directory to PATH.

**"ANTHROPIC_API_KEY not set"**
- Set the environment variable as shown above.
- Make sure it's exported in your shell profile for persistence.

**"failed to open camera: /dev/video0: permission denied"**
- Add your user to the `video` group (Linux): `sudo usermod -a -G video $USER`
- Log out and log back in.

**"failed to open camera: device not found"**
- Check available cameras: `ls /dev/video*` (Linux)
- Try a different camera ID in device config (e.g., `/dev/video1`)
- Ensure USB camera is connected and recognized by OS.

**"API call failed: authentication error"**
- Verify API key is correct.
- Check that key has not expired or been revoked.
- Ensure API key is for Anthropic (not OpenAI or other provider).

**"observation timeout"**
- Claude Vision API can be slow on first call (~2-3 seconds).
- If consistent timeout, check internet connection.
- Verify API key has valid credits/quota.

## Next Steps

Installation complete! Continue to [Getting Started](getting-started.md) for a step-by-step walkthrough.
