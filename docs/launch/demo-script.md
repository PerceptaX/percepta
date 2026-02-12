# Percepta Demo Script (5 minutes)

## Scene 1: The Problem (30 seconds)
"AI code generation tools like Embedder are great—they generate code from
datasheets in seconds. But there's a problem: does it actually work?"

[Show Embedder generating code]
"Embedder says this code should work. Let's flash it..."
[Show LED not blinking correctly]
"2Hz instead of 1Hz. Now I'm debugging AI-generated code."

## Scene 2: The Percepta Way (2 minutes)
"Percepta generates code AND validates it works on real hardware."

```bash
percepta generate "Blink LED at 1Hz" --board esp32 --output blink.c
```

[Show generation report with style validation]
"✓ BARR-C compliant, auto-fixes applied"

[Flash to hardware]

```bash
percepta observe my-esp32
```

[Show vision capturing LED blinking at 1.02 Hz]
"Percepta observes actual hardware behavior with computer vision."

## Scene 3: The Knowledge Graph (1 minute)
"Every validated pattern gets stored. The more you use Percepta, the smarter it gets."

```bash
percepta generate "Toggle LED on button press" --board esp32
```

[Show prompt retrieving similar validated patterns]
"It found 3 similar patterns that work on ESP32."

## Scene 4: Professional Code (1 minute)
"Code that looks like a senior engineer wrote it."

[Show generated code side-by-side with BARR-C checklist]
- Function names: Module_Function() ✓
- Types: uint8_t, uint16_t ✓
- No magic numbers ✓
- Const correctness ✓
- Doxygen comments ✓

## Scene 5: The Pitch (30 seconds)
"Percepta: Firmware that works. Code that professionals write."

"Get started free: github.com/Perceptax/percepta"
