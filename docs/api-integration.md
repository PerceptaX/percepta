# API Integration Guide

Integrate Percepta into your applications and workflows.

---

## Overview

Percepta can be used as:
1. **CLI tool** - Command-line interface (current release)
2. **Go library** - Import into Go applications
3. **MCP server** - Model Context Protocol server (planned for v2.1)
4. **REST API** - HTTP API (planned for v3.0)

---

## CLI Integration

Current release focuses on CLI integration with shell scripts and CI/CD pipelines.

### Shell Scripts

**Basic observation script:**

```bash
#!/bin/bash
# observe.sh - Capture and validate hardware behavior

DEVICE=$1

if [ -z "$DEVICE" ]; then
  echo "Usage: $0 <device-id>"
  exit 1
fi

# Observe hardware
percepta observe $DEVICE

# Validate expected behavior
percepta assert $DEVICE "led power is ON" || {
  echo "ERROR: Power LED not on"
  exit 1
}

echo "Hardware validation passed"
```

**Firmware comparison script:**

```bash
#!/bin/bash
# compare-firmware.sh - Compare two firmware versions

DEVICE=$1
FROM=$2
TO=$3

if [ -z "$DEVICE" ] || [ -z "$FROM" ] || [ -z "$TO" ]; then
  echo "Usage: $0 <device> <from-version> <to-version>"
  exit 1
fi

# Compare firmware versions
percepta diff $DEVICE --from $FROM --to $TO

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "✓ No behavioral changes"
elif [ $EXIT_CODE -eq 1 ]; then
  echo "⚠️  Behavioral changes detected"
  echo "Review diff output above"
else
  echo "✗ Error comparing firmware"
  exit 1
fi
```

### CI/CD Integration

**GitHub Actions:**

```yaml
name: Hardware Validation

on: [push, pull_request]

jobs:
  hardware-test:
    runs-on: self-hosted  # Must have hardware access

    steps:
      - uses: actions/checkout@v3

      - name: Install Percepta
        run: |
          curl -fsSL https://github.com/Perceptax/percepta/releases/latest/download/percepta-linux-amd64 -o /usr/local/bin/percepta
          chmod +x /usr/local/bin/percepta

      - name: Build firmware
        run: make

      - name: Flash firmware
        run: make flash

      - name: Hardware validation
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          percepta device set-firmware test-board ${{ github.sha }}
          percepta observe test-board
          percepta assert test-board "led power is ON"
          percepta assert test-board "boot time < 3000ms"

      - name: Compare to baseline
        if: github.ref != 'refs/heads/main'
        run: |
          percepta diff test-board --from main --to ${{ github.sha }}
```

**GitLab CI:**

```yaml
hardware-test:
  stage: test
  tags:
    - hardware  # Runner with hardware access

  variables:
    ANTHROPIC_API_KEY: $ANTHROPIC_API_KEY

  script:
    - make build
    - make flash
    - percepta device set-firmware test-board $CI_COMMIT_SHA
    - percepta observe test-board
    - percepta assert test-board "led power is ON"
    - percepta diff test-board --from main --to $CI_COMMIT_SHA

  only:
    - merge_requests
    - main
```

**Jenkins:**

```groovy
pipeline {
    agent {
        label 'hardware-test-rig'
    }

    environment {
        ANTHROPIC_API_KEY = credentials('anthropic-api-key')
    }

    stages {
        stage('Build') {
            steps {
                sh 'make clean && make'
            }
        }

        stage('Flash') {
            steps {
                sh 'make flash'
            }
        }

        stage('Hardware Validation') {
            steps {
                sh '''
                    percepta device set-firmware test-board ${GIT_COMMIT}
                    percepta observe test-board
                    percepta assert test-board "led power is ON"
                '''
            }
        }

        stage('Compare to Main') {
            when {
                not { branch 'main' }
            }
            steps {
                sh 'percepta diff test-board --from main --to ${GIT_COMMIT}'
            }
        }
    }
}
```

---

## Go Library Usage

Use Percepta as a Go library in your applications.

### Installation

```bash
go get github.com/Perceptax/percepta
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/Perceptax/percepta/pkg/percepta"
    "github.com/Perceptax/percepta/internal/config"
    "github.com/Perceptax/percepta/internal/storage"
)

func main() {
    // Load config
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Get device config
    deviceCfg, ok := cfg.Devices["my-esp32"]
    if !ok {
        log.Fatal("Device not found")
    }

    // Initialize storage
    store, err := storage.NewSQLiteStorage()
    if err != nil {
        log.Fatal(err)
    }
    defer store.Close()

    // Create Percepta core
    core, err := percepta.NewCore(deviceCfg.CameraID, store)
    if err != nil {
        log.Fatal(err)
    }

    // Observe hardware
    obs, err := core.Observe("my-esp32")
    if err != nil {
        log.Fatal(err)
    }

    // Print signals
    fmt.Printf("Captured %d signals\n", len(obs.Signals))
    for _, signal := range obs.Signals {
        fmt.Printf("Signal: %v\n", signal)
    }
}
```

### Assertions

```go
package main

import (
    "fmt"
    "log"

    "github.com/Perceptax/percepta/internal/assertions"
    "github.com/Perceptax/percepta/pkg/percepta"
    "github.com/Perceptax/percepta/internal/storage"
)

func main() {
    // Initialize
    store, _ := storage.NewSQLiteStorage()
    defer store.Close()

    core, _ := percepta.NewCore("/dev/video0", store)

    // Observe
    obs, _ := core.Observe("my-board")

    // Parse assertion
    assertion, err := assertions.Parse("led power is ON")
    if err != nil {
        log.Fatal(err)
    }

    // Evaluate
    result := assertion.Evaluate(obs)

    if result.Passed {
        fmt.Println("✓ Assertion passed")
    } else {
        fmt.Printf("✗ Assertion failed: %s\n", result.Message)
    }
}
```

### Code Generation

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/Perceptax/percepta/internal/codegen"
    "github.com/Perceptax/percepta/internal/knowledge"
    "github.com/Perceptax/percepta/internal/style"
)

func main() {
    // Get API key
    apiKey := os.Getenv("ANTHROPIC_API_KEY")
    if apiKey == "" {
        log.Fatal("ANTHROPIC_API_KEY not set")
    }

    // Initialize components
    patternStore, _ := knowledge.NewPatternStore()
    defer patternStore.Close()

    styleChecker := style.NewStyleChecker()
    styleFixer := style.NewStyleFixer()
    claudeClient := codegen.NewClaudeClient(apiKey)
    promptBuilder := codegen.NewPromptBuilder(patternStore)

    // Create pipeline
    pipeline := codegen.NewGenerationPipeline(
        claudeClient,
        promptBuilder,
        styleChecker,
        styleFixer,
        patternStore,
    )

    // Generate code
    result, err := pipeline.Generate(
        "Blink LED at 1Hz",
        "esp32",
        "my-esp32",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Print result
    fmt.Println("Generated code:")
    fmt.Println(result.Code)
    fmt.Printf("\nStyle compliant: %v\n", result.StyleCompliant)
}
```

### Custom Camera Driver

Implement custom camera backend:

```go
package main

import (
    "github.com/Perceptax/percepta/internal/camera"
)

type CustomCamera struct {
    devicePath string
}

func NewCustomCamera(path string) camera.CameraDriver {
    return &CustomCamera{devicePath: path}
}

func (c *CustomCamera) Capture() ([]byte, error) {
    // Your custom capture logic
    // Must return JPEG bytes

    // Example: read from custom camera API
    jpegData, err := readFromCamera(c.devicePath)
    if err != nil {
        return nil, err
    }

    return jpegData, nil
}

func (c *CustomCamera) Close() error {
    // Cleanup
    return nil
}
```

Use custom camera:

```go
customCam := NewCustomCamera("/dev/mydevice")
core := percepta.NewCoreWithCamera(customCam, storage)
```

---

## MCP Server Mode (Planned)

**Status:** Planned for v2.1 (Q2 2026)

Percepta will support Model Context Protocol (MCP), allowing Claude Desktop and other MCP clients to access hardware observation capabilities.

### Configuration (Preview)

**claude_desktop_config.json:**

```json
{
  "mcpServers": {
    "percepta": {
      "command": "percepta",
      "args": ["mcp"],
      "env": {
        "ANTHROPIC_API_KEY": "sk-ant-..."
      }
    }
  }
}
```

### Available Tools (Preview)

**observe_hardware:**
```json
{
  "name": "observe_hardware",
  "description": "Observe hardware behavior via camera",
  "inputSchema": {
    "type": "object",
    "properties": {
      "device_id": {
        "type": "string",
        "description": "Device identifier"
      }
    },
    "required": ["device_id"]
  }
}
```

**assert_behavior:**
```json
{
  "name": "assert_behavior",
  "description": "Validate hardware behavior",
  "inputSchema": {
    "type": "object",
    "properties": {
      "device_id": {"type": "string"},
      "assertion": {"type": "string"}
    },
    "required": ["device_id", "assertion"]
  }
}
```

**generate_firmware:**
```json
{
  "name": "generate_firmware",
  "description": "Generate BARR-C compliant firmware",
  "inputSchema": {
    "type": "object",
    "properties": {
      "spec": {"type": "string"},
      "board": {"type": "string"}
    },
    "required": ["spec", "board"]
  }
}
```

### Usage Example (Preview)

```
User: Observe my ESP32 and check if the status LED is blinking at 1Hz

Claude: I'll observe the ESP32 hardware and check the LED behavior.

[Uses observe_hardware tool]

The ESP32 status LED is blinking at 0.98 Hz, which is very close to the expected 1Hz (within normal tolerance).

Would you like me to:
1. Generate code to adjust it to exactly 1Hz?
2. Create an assertion to validate this behavior in CI?
3. Compare this behavior to a previous firmware version?
```

---

## REST API (Planned)

**Status:** Planned for v3.0 (Q3 2026)

Future releases will include an HTTP REST API for language-agnostic integration.

### Endpoints (Preview)

**POST /api/v1/observe**

Capture hardware observation.

Request:
```json
{
  "device_id": "my-esp32",
  "duration": "5s"
}
```

Response:
```json
{
  "observation_id": "obs-12345",
  "device_id": "my-esp32",
  "timestamp": "2026-02-13T10:30:00Z",
  "signals": [
    {
      "type": "led",
      "name": "LED1",
      "state": "on",
      "color": {"r": 0, "g": 0, "b": 255},
      "blink_hz": 0,
      "confidence": 0.92
    }
  ]
}
```

**POST /api/v1/assert**

Validate hardware behavior.

Request:
```json
{
  "device_id": "my-esp32",
  "assertion": "led power is ON"
}
```

Response:
```json
{
  "passed": true,
  "expected": "LED 'power' is ON",
  "actual": "LED 'power' is ON",
  "confidence": 0.89
}
```

**POST /api/v1/generate**

Generate firmware code.

Request:
```json
{
  "spec": "Blink LED at 1Hz",
  "board": "esp32"
}
```

Response:
```json
{
  "code": "/* Generated C code */\n...",
  "style_compliant": true,
  "violations": [],
  "pattern_stored": true
}
```

---

## Python Integration (Community)

While Percepta is written in Go, Python integration is possible via subprocess or HTTP API (when available).

### Subprocess Wrapper

```python
import subprocess
import json

class Percepta:
    def observe(self, device_id):
        result = subprocess.run(
            ['percepta', 'observe', device_id],
            capture_output=True,
            text=True
        )

        if result.returncode != 0:
            raise Exception(f"Observation failed: {result.stderr}")

        return result.stdout

    def assert_behavior(self, device_id, assertion):
        result = subprocess.run(
            ['percepta', 'assert', device_id, assertion],
            capture_output=True,
            text=True
        )

        return result.returncode == 0

    def diff(self, device_id, from_version, to_version):
        result = subprocess.run(
            ['percepta', 'diff', device_id,
             '--from', from_version, '--to', to_version],
            capture_output=True,
            text=True
        )

        return {
            'has_changes': result.returncode == 1,
            'output': result.stdout
        }

# Usage
percepta = Percepta()

# Observe
output = percepta.observe('my-esp32')
print(output)

# Assert
passed = percepta.assert_behavior('my-esp32', 'led power is ON')
print(f"Assertion passed: {passed}")

# Diff
result = percepta.diff('my-esp32', 'v1.0', 'v1.1')
print(f"Changes detected: {result['has_changes']}")
```

---

## JavaScript/TypeScript Integration

### Node.js Subprocess Wrapper

```javascript
const { exec } = require('child_process');
const util = require('util');
const execPromise = util.promisify(exec);

class Percepta {
  async observe(deviceId) {
    try {
      const { stdout, stderr } = await execPromise(
        `percepta observe ${deviceId}`
      );
      return stdout;
    } catch (error) {
      throw new Error(`Observation failed: ${error.message}`);
    }
  }

  async assert(deviceId, assertion) {
    try {
      await execPromise(
        `percepta assert ${deviceId} "${assertion}"`
      );
      return true;
    } catch (error) {
      return false;
    }
  }

  async diff(deviceId, fromVersion, toVersion) {
    try {
      const { stdout } = await execPromise(
        `percepta diff ${deviceId} --from ${fromVersion} --to ${toVersion}`
      );
      return { hasChanges: false, output: stdout };
    } catch (error) {
      if (error.code === 1) {
        return { hasChanges: true, output: error.stdout };
      }
      throw error;
    }
  }
}

// Usage
(async () => {
  const percepta = new Percepta();

  // Observe
  const output = await percepta.observe('my-esp32');
  console.log(output);

  // Assert
  const passed = await percepta.assert('my-esp32', 'led power is ON');
  console.log(`Assertion passed: ${passed}`);

  // Diff
  const result = await percepta.diff('my-esp32', 'v1.0', 'v1.1');
  console.log(`Changes detected: ${result.hasChanges}`);
})();
```

---

## Rust Integration

### Using FFI (Foreign Function Interface)

```rust
use std::process::Command;
use std::str;

pub struct Percepta;

impl Percepta {
    pub fn observe(device_id: &str) -> Result<String, String> {
        let output = Command::new("percepta")
            .arg("observe")
            .arg(device_id)
            .output()
            .map_err(|e| format!("Failed to execute: {}", e))?;

        if output.status.success() {
            Ok(String::from_utf8_lossy(&output.stdout).to_string())
        } else {
            Err(String::from_utf8_lossy(&output.stderr).to_string())
        }
    }

    pub fn assert_behavior(device_id: &str, assertion: &str) -> Result<bool, String> {
        let output = Command::new("percepta")
            .arg("assert")
            .arg(device_id)
            .arg(assertion)
            .output()
            .map_err(|e| format!("Failed to execute: {}", e))?;

        Ok(output.status.success())
    }
}

// Usage
fn main() {
    // Observe
    match Percepta::observe("my-esp32") {
        Ok(output) => println!("{}", output),
        Err(e) => eprintln!("Error: {}", e),
    }

    // Assert
    match Percepta::assert_behavior("my-esp32", "led power is ON") {
        Ok(passed) => println!("Assertion passed: {}", passed),
        Err(e) => eprintln!("Error: {}", e),
    }
}
```

---

## Docker Integration

Run Percepta in Docker with USB camera access.

**Dockerfile:**

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o percepta ./cmd/percepta

FROM alpine:latest

RUN apk --no-cache add ca-certificates v4l-utils

COPY --from=builder /app/percepta /usr/local/bin/percepta

# Default config directory
RUN mkdir -p /root/.config/percepta
VOLUME /root/.config/percepta

# Storage directory
RUN mkdir -p /root/.local/share/percepta
VOLUME /root/.local/share/percepta

ENTRYPOINT ["percepta"]
```

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  percepta:
    build: .
    devices:
      - /dev/video0:/dev/video0  # USB camera
    environment:
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - OPENAI_API_KEY=${OPENAI_API_KEY}
    volumes:
      - ./config:/root/.config/percepta
      - ./data:/root/.local/share/percepta
    command: observe my-board
```

**Usage:**

```bash
# Build
docker-compose build

# Run observation
docker-compose run percepta observe my-board

# Run assertion
docker-compose run percepta assert my-board "led power is ON"
```

---

## Webhook Integration

Trigger webhooks on observation events (requires wrapper script).

**webhook-observer.sh:**

```bash
#!/bin/bash
# Observe hardware and POST results to webhook

DEVICE=$1
WEBHOOK_URL=$2

if [ -z "$DEVICE" ] || [ -z "$WEBHOOK_URL" ]; then
  echo "Usage: $0 <device> <webhook-url>"
  exit 1
fi

# Capture observation
OUTPUT=$(percepta observe $DEVICE 2>&1)
EXIT_CODE=$?

# Build JSON payload
PAYLOAD=$(cat <<EOF
{
  "device": "$DEVICE",
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "success": $([ $EXIT_CODE -eq 0 ] && echo "true" || echo "false"),
  "output": $(echo "$OUTPUT" | jq -Rs .)
}
EOF
)

# POST to webhook
curl -X POST \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD" \
  $WEBHOOK_URL

echo "Observation sent to webhook"
```

**Usage:**

```bash
./webhook-observer.sh my-esp32 https://hooks.example.com/percepta
```

---

## See Also

- [Commands Reference](commands.md) - CLI commands
- [Examples](examples.md) - Usage examples
- [Configuration Guide](configuration.md) - Config options
- [GitHub Repository](https://github.com/Perceptax/percepta) - Source code
