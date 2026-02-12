# Percepta Demo Video Storyboard (5 minutes)

## Video Overview
- **Duration:** 5 minutes
- **Format:** Screen recording + webcam insert + voiceover
- **Target audience:** Embedded firmware developers, hardware engineers
- **Call to action:** Visit github.com/Perceptax/percepta

## Scene Breakdown

### Scene 1: Introduction - The Problem (30 seconds)
**Duration:** 0:00 - 0:30

**Visual:**
- Split screen: Code editor (left) + Hardware on desk (right)
- Show Embedder generating code from prompt "Blink LED at 1Hz"
- Code appears, looks reasonable
- Flash to ESP32 board
- LED blinks, but clearly faster than 1Hz
- Overlay timer showing "2.1 Hz" vs expected "1.0 Hz"

**Voiceover:**
"AI code generation tools like Embedder are impressive—they generate firmware from datasheets in seconds. But there's a critical gap: they can't tell you if the code actually works on hardware. This LED should blink at 1Hz, but it's running at 2Hz. Now I'm debugging AI-generated code."

**Key frame:**
- LED blinking wrong frequency with red X overlay

### Scene 2: The Percepta Way - Generation (1 minute 30 seconds)
**Duration:** 0:30 - 2:00

**Visual:**
- Terminal window, clean prompt
- Type command: `percepta generate "Blink LED at 1Hz on GPIO2" --board esp32 --output blink.c`
- Show progress spinner: "Generating firmware..."
- Generation report appears:
  ```
  ✓ Code generated (12.1s)
  ✓ BARR-C style check: 0 violations
  ✓ Auto-fixes applied:
    - Converted int to uint8_t (3 locations)
    - Added #include <stdint.h>
    - Renamed led_pin → Led_Pin (BARR-C module naming)
  ✓ Pattern stored in knowledge graph
  ```
- Show generated code side-by-side with BARR-C checklist
- Highlight: uint8_t types, named constants, Doxygen comments

**Voiceover:**
"Percepta takes a different approach. Start with a natural language prompt, just like Embedder. But Percepta enforces BARR-C coding standards automatically. It converts generic types to stdint.h types, eliminates magic numbers, and enforces professional naming conventions. This isn't just working code—it's code review-ready code."

**Key frame:**
- Code with BARR-C checkmarks overlay

### Scene 3: Hardware Validation (1 minute)
**Duration:** 2:00 - 3:00

**Visual:**
- Terminal: `platformio run -t upload` (flash firmware)
- Board boots, LED starts blinking
- Terminal: `percepta observe my-esp32`
- Webcam feed appears in terminal as ASCII preview
- Real-time vision analysis overlay:
  ```
  Detected signals:
  - LED1 (GPIO2): ON, Green, BlinkHz: 1.02Hz ✓

  Confidence: 98%
  ```
- Show 5-second observation window with temporal smoothing
- Final report: "✓ LED blinking at 1.02Hz (expected 1.0Hz, within tolerance)"

**Voiceover:**
"Now comes the magic. Flash the firmware to your board, then run 'percepta observe.' Point your webcam at the hardware. Percepta uses computer vision to watch actual LED behavior—not simulation, not hoping it works—actual hardware observation. It detects LED states, measures blink frequency, and validates against your specification. This LED is blinking at 1.02Hz. Close enough."

**Key frame:**
- Split screen: webcam feed + analysis report with green checkmark

### Scene 4: The Knowledge Graph (1 minute)
**Duration:** 3:00 - 4:00

**Visual:**
- Terminal: `percepta generate "Toggle LED on button press with debouncing" --board esp32`
- Show semantic search finding similar patterns:
  ```
  Found 3 similar validated patterns:
  1. "Button with interrupt-driven debounce" (similarity: 0.89)
     Board: ESP32-DevKitC | Validated: 2026-02-10
  2. "LED toggle on GPIO interrupt" (similarity: 0.84)
     Board: ESP32-S3 | Validated: 2026-02-09
  3. "Debounced button state machine" (similarity: 0.81)
     Board: ESP32-C3 | Validated: 2026-02-08

  Using patterns to guide generation...
  ```
- Show generated code incorporating debounce logic
- Highlight knowledge graph diagram: Spec → Pattern → Board → Observation

**Voiceover:**
"Every pattern that works gets stored in a knowledge graph. Board-specific quirks, validated configurations, patterns that actually run on hardware. When you generate new firmware, Percepta searches for similar validated patterns. This isn't just code generation—it's learning from your hardware. The more you use Percepta, the smarter it gets."

**Key frame:**
- Knowledge graph visualization with highlighted path

### Scene 5: Professional Code Quality (30 seconds)
**Duration:** 4:00 - 4:30

**Visual:**
- Split screen comparison: Generic AI code (left) vs Percepta code (right)
- Highlight differences:
  - Left: `int led = 2;` → Right: `#define LED_PIN ((uint8_t)2)`
  - Left: No comments → Right: Doxygen function headers
  - Left: `while(1)` → Right: `while(true)` with explicit return
  - Left: No error checks → Right: `if(gpio_set_level(...) != ESP_OK)`
- Show BARR-C compliance report: 98% pass rate

**Voiceover:**
"Here's the difference. Generic AI code versus Percepta's BARR-C compliant output. Professional naming conventions. No magic numbers. Explicit error handling. Type safety. This is code that passes code review. Code that ships in production firmware."

**Key frame:**
- Side-by-side code with annotations

### Scene 6: Call to Action (30 seconds)
**Duration:** 4:30 - 5:00

**Visual:**
- Terminal: Quick start sequence (4 commands)
  ```bash
  # Install
  curl -fsSL https://github.com/Perceptax/percepta/releases/latest/download/percepta-linux-amd64 -o /usr/local/bin/percepta
  chmod +x /usr/local/bin/percepta

  # Add device
  percepta device add my-esp32 --camera /dev/video0

  # Observe hardware
  percepta observe my-esp32

  # Generate firmware
  percepta generate "Blink LED at 1Hz" --board esp32 --output blink.c
  ```
- Show GitHub repo: github.com/Perceptax/percepta
- Overlay key stats: "Open source | MIT license | 500+ GitHub stars"

**Voiceover:**
"Percepta is open source, MIT licensed, and free for unlimited local usage. Install in seconds, add your board, start generating production-ready firmware. Stop debugging AI-generated code. Generate firmware that works. Get started at github.com/Perceptax/percepta."

**Key frame:**
- GitHub repo page with "Get Started" button highlighted

## Technical Requirements

### Screen Recordings Needed
1. Terminal session showing full workflow (5 min continuous)
2. Code editor showing generated code (1 min)
3. Webcam feed with LED blinking (30 sec)
4. Browser showing GitHub repo (15 sec)

### Hardware Setup
- ESP32-DevKitC-32E board
- Single LED on GPIO2 (or built-in LED)
- Webcam on tripod pointing at board
- Clean desk background (minimal clutter)

### Editing
- Add text overlays for key points
- Add progress bars/spinners for long operations
- Add comparison table overlay (Scene 5)
- Add knowledge graph visualization (Scene 4)
- Background music: Subtle, tech-focused (optional)

### Voiceover Script
- Clear, confident tone
- Technical but accessible
- Emphasize "works on hardware" repeatedly
- Avoid attacking competitors (focus on "we validate" not "they don't")

### Export Settings
- Resolution: 1920x1080 (1080p)
- Frame rate: 30fps
- Format: MP4 (H.264)
- Audio: 192kbps AAC
- File size target: <100MB

## Production Checklist

- [ ] Write detailed voiceover script with timestamps
- [ ] Record clean terminal session (no typos, smooth execution)
- [ ] Capture webcam footage with good lighting
- [ ] Record voiceover with quality microphone
- [ ] Edit together screen recordings with voiceover sync
- [ ] Add text overlays and animations
- [ ] Add comparison table and knowledge graph graphics
- [ ] Add background music (if using)
- [ ] Review for pacing (5 minutes exactly)
- [ ] Export high-quality MP4
- [ ] Upload to YouTube
- [ ] Create thumbnail with "Firmware That Works" tagline
- [ ] Write video description with timestamps and links
- [ ] Add to GitHub README and docs/README.md

## Distribution

**Primary:**
- YouTube (main hosting)
- GitHub README (embedded player)
- Blog post (embedded)

**Secondary:**
- Twitter/X (1-minute clip)
- LinkedIn (2-minute professional version)
- Reddit r/embedded (link in launch post)

## Success Metrics

Track after launch:
- YouTube views (target: 1000 in first week)
- Average watch time (target: >3 minutes)
- Click-through rate to GitHub (target: >10%)
- Conversion to installation (track downloads after video views)
