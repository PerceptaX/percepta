# Percepta: Full Platform PRD
## Vision-Based Hardware Observation + AI Firmware Generation

**Version:** 2.0 (Full Platform)  
**Date:** February 12, 2026  
**Status:** Phase 2 Vision (Post-MVP)

---

## Executive Summary

**Percepta is the complete embedded AI development platform** that generates professional, hardware-validated firmware. It combines vision-based physical observation with AI code generation to deliver firmware that not only compiles — but actually works on real hardware and passes code review.

### The Complete Vision

**Phase 1 (Months 0-6): Perception Layer**
- Vision-based hardware observation
- Behavioral memory
- 500+ users validating hardware behavior

**Phase 2 (Months 7-12): Code Generation**
- AI firmware generation with hardware validation
- Style-compliant code (BARR-C, MISRA-C)
- Automatic validation loop
- "Better than Embedder" positioning

**Phase 3 (Months 13-24): Platform**
- Cloud HIL farm (50+ boards)
- Community code hub (hardware-verified patterns)
- Enterprise features
- "The firmware platform professionals trust"

### Why This Beats Everything

| Dimension | Embedder | Traditional HIL | Percepta |
|-----------|----------|-----------------|----------|
| **Code Generation** | ✅ 95% compiles | ❌ N/A | ✅ 100% works (validated) |
| **Style Compliance** | ❌ Generic AI | ❌ N/A | ✅ BARR-C/MISRA-C |
| **Hardware Validation** | ❌ Simulation only | ✅ $50K-500K | ✅ $0-49/month |
| **Physical Observation** | ❌ Blind | ⚠️ Electrical only | ✅ Vision-based |
| **Iteration Speed** | ⚠️ Manual | ⚠️ Slow setup | ✅ Automatic loop |
| **Behavioral Knowledge** | ❌ None | ❌ None | ✅ 10M+ observations |

**Unique position:** Only platform that generates professional code AND validates it works on real hardware.

### Success Metrics

| Timeline | Key Milestone | Metric |
|----------|---------------|--------|
| **Month 6** | Perception PMF | 500 weekly active users |
| **Month 12** | Code gen launch | 200 paying users, $10K MRR |
| **Month 18** | Platform features | HIL farm live, 1000 paying users |
| **Month 24** | Market leadership | $100K MRR, Embedder killer |
| **Month 36** | Exit opportunity | $10M ARR, $50-100M valuation |

---

## Part I: Foundation (Months 0-6)

*See Percepta_PRD_Final.md for complete Phase 1 specifications*

### Quick Summary: Perception MVP

**What we build:**
- `percepta observe` - Capture hardware state via vision
- `percepta assert` - Validate expected behavior
- `percepta diff` - Compare firmware versions
- SQLite storage + optional Mem0 integration
- CLI + optional MCP server mode

**Why this comes first:**
1. Builds credibility ("they actually shipped something")
2. Creates perception moat (vision expertise is hard)
3. Gathers 100K+ observations for training data
4. Proves market demand before investing in code gen

**Exit criteria:**
- 500 weekly active users
- 50%+ weekly retention
- Users say: "I can't develop without this anymore"

---

## Part II: Code Generation (Months 7-12)

### The Big Reveal: "We Generate Better Code Than Embedder"

After 6 months of building trust as "the vision layer," we reveal that Percepta can also **generate firmware** — with a critical difference:

**Embedder's promise:**
> "We generate code cited to datasheets with 95% accuracy"

**Percepta's promise:**
> "We generate professional firmware that we VALIDATE works on your actual hardware. Every line follows BARR-C standards. 100% accuracy because we test it with vision."

### 2.1 Architecture (Code Generation Layer)

```
┌─────────────────────────────────────────────────────────────────┐
│                      PERCEPTA PLATFORM                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              USER INTERFACE LAYER                          │ │
│  │                                                            │ │
│  │  CLI           Web App        VSCode Ext      API         │ │
│  │  ─────         ─────────      ────────────    ───         │ │
│  │  Terminal      Dashboard      IDE Native      REST/MCP    │ │
│  └────────────────────────────────────────────────────────────┘ │
│                              │                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              AGENT ORCHESTRATION LAYER                     │ │
│  │                                                            │ │
│  │  Code Gen Agent    Validation Agent    Style Agent        │ │
│  │  ───────────────   ─────────────────   ───────────────    │ │
│  │  Writes firmware   Tests on hardware   Enforces BARR-C    │ │
│  │  from spec         via perception      code standards     │ │
│  └────────────────────────────────────────────────────────────┘ │
│                              │                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              CORE CAPABILITIES                             │ │
│  │                                                            │ │
│  │  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐ │ │
│  │  │ PERCEPTION   │  │ GENERATION   │  │ KNOWLEDGE GRAPH │ │ │
│  │  │              │  │              │  │                 │ │ │
│  │  │ Vision       │  │ LLM fine-    │  │ Board behaviors │ │ │
│  │  │ Observation  │  │ tuned on     │  │ Silicon quirks  │ │ │
│  │  │ Assertion    │  │ validated    │  │ Style patterns  │ │ │
│  │  │ Diff         │  │ patterns     │  │ Timing data     │ │ │
│  │  └──────────────┘  └──────────────┘  └─────────────────┘ │ │
│  └────────────────────────────────────────────────────────────┘ │
│                              │                                   │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │              INFRASTRUCTURE LAYER                          │ │
│  │                                                            │ │
│  │  Local Dev       Cloud HIL       Style Checker            │ │
│  │  ─────────       ──────────      ──────────────           │ │
│  │  Your boards     50+ boards      BARR-C/MISRA             │ │
│  │  w/ camera       w/ cameras      enforcement              │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 Core Innovation: Validation Loop

**The key difference from Embedder:**

```python
# Embedder's workflow:
def embedder_generate(spec, board):
    code = llm.generate(spec, datasheet=get_datasheet(board))
    if compile(code):
        return code  # Hope it works! 🤞
    else:
        retry()

# Percepta's workflow:
def percepta_generate(spec, board, max_iterations=5):
    for i in range(max_iterations):
        # 1. Generate with context
        code = llm.generate(
            spec=spec,
            datasheet=get_datasheet(board),
            board_quirks=knowledge_graph.get_quirks(board),
            validated_patterns=knowledge_graph.find_similar(spec, board),
            style_template=style_graph.get_template(spec, board)
        )
        
        # 2. Enforce style
        code = style_checker.enforce_barrc(code)
        
        # 3. Flash to hardware
        flash(code, board)
        wait(2)  # Boot time
        
        # 4. Observe with vision
        observation = vision.observe(board)
        
        # 5. Check against spec
        result = vision.assert(spec, observation)
        
        if result.passed:
            # Store as validated pattern
            knowledge_graph.store(
                spec=spec,
                board=board,
                code=code,
                observation=observation,
                style_compliant=True
            )
            return ValidatedCode(code, iterations=i+1)
        
        # 6. Give LLM feedback for next iteration
        feedback = f"""
        Attempt {i+1} failed validation.
        Expected: {result.expected}
        Actual: {result.actual}
        
        Common issue on {board}: {knowledge_graph.get_failure_hints(result)}
        """
        llm.add_context(feedback)
    
    raise ValidationError(f"Could not generate working code after {max_iterations} attempts")
```

**Why this beats Embedder:**
- Embedder: 95% works on first try (no feedback loop)
- Percepta: 100% works after N iterations (automatic validation)

### 2.3 Code Generation Engine

```python
class PerceptaCodeGenerator:
    """
    Generates hardware-validated, style-compliant embedded firmware.
    
    Key innovations:
    1. Fine-tuned on hardware-validated code (not just GitHub)
    2. Behavioral knowledge graph (what actually works)
    3. Style enforcement (BARR-C/MISRA-C compliance)
    4. Automatic validation loop (vision feedback)
    """
    
    def __init__(self):
        # Base model fine-tuned on style-compliant code
        self.base_model = "claude-sonnet-4"
        self.lora_adapter = load_lora("percepta-embedded-barrc-v1")
        
        # Knowledge graphs
        self.behavioral_knowledge = BehavioralKnowledgeGraph()
        self.style_knowledge = StyleKnowledgeGraph()
        
        # Style enforcement
        self.style_checker = BARRCStyleChecker()
        
        # Hardware interface
        self.vision = PerceptionEngine()
        self.flasher = HardwareFlasher()
    
    def generate_and_validate(
        self, 
        spec: str, 
        board: str,
        max_iterations: int = 5
    ) -> ValidatedCode:
        """
        Generate firmware and validate on real hardware.
        
        Args:
            spec: Natural language specification (e.g., "Blink LED at 1Hz")
            board: Board type (e.g., "esp32-devkit-v1")
            max_iterations: Max validation attempts
            
        Returns:
            ValidatedCode with:
                - code: The firmware source
                - style_compliant: True
                - hardware_validated: True
                - iterations: Number of attempts needed
                - observation: Final hardware state
        """
        
        # 1. Get contextual information
        style_template = self.style_knowledge.get_template(spec, board)
        board_quirks = self.behavioral_knowledge.get_quirks(board)
        validated_patterns = self.behavioral_knowledge.find_patterns(spec, board)
        
        # System prompt with style requirements
        system_prompt = f"""
You are an expert embedded firmware engineer writing BARR-C compliant code.

Style requirements:
- Function names: Module_Function() format
- Variables: snake_case
- Constants: UPPER_SNAKE with meaningful names
- Types: Use stdint.h (uint8_t, uint16_t, etc.)
- No magic numbers: Define all constants
- Const correctness: const uint8_t* not uint8_t*
- Doxygen comments: /** ... */ for all functions
- Non-blocking: Use timers, not delays
- Error handling: Return codes, not assertions
- Static allocation: No malloc/free
- Explicit casts: (uint16_t)value

Board-specific context:
{board_quirks}

Similar validated patterns:
{validated_patterns}
"""
        
        # 2. Validation loop
        for iteration in range(max_iterations):
            # Generate code
            code = self.base_model.generate(
                spec=spec,
                system_prompt=system_prompt,
                template=style_template,
                lora_adapter=self.lora_adapter
            )
            
            # Enforce style
            code = self.style_checker.enforce(code)
            violations = self.style_checker.check(code)
            
            if violations:
                logger.warning(f"Style violations: {violations}")
            
            # Flash to hardware
            try:
                self.flasher.flash(board, code)
                time.sleep(2)  # Wait for boot
            except FlashError as e:
                # Code didn't even flash
                system_prompt += f"\n\nCompilation error: {e}\n"
                continue
            
            # Observe behavior
            observation = self.vision.observe(board)
            
            # Validate against spec
            result = self.validate_observation(spec, observation)
            
            if result.passed:
                # Success! Store for future reference
                self.behavioral_knowledge.store_validated_pattern(
                    spec=spec,
                    board=board,
                    code=code,
                    observation=observation,
                    style_compliant=True
                )
                
                self.style_knowledge.store_template(
                    spec=spec,
                    board=board,
                    template=extract_template(code)
                )
                
                return ValidatedCode(
                    code=code,
                    style_compliant=True,
                    hardware_validated=True,
                    iterations=iteration + 1,
                    observation=observation
                )
            
            # Failed validation - provide feedback
            feedback = f"""
Validation failed on iteration {iteration + 1}:

Expected: {result.expected}
Actual: {result.actual}
Confidence: {result.confidence}

Common causes on {board}:
{self.behavioral_knowledge.get_failure_hints(result)}

Please regenerate the code addressing this issue.
"""
            system_prompt += f"\n\n{feedback}\n"
        
        # Exhausted attempts
        raise ValidationError(
            f"Could not generate working code after {max_iterations} iterations"
        )
    
    def validate_observation(self, spec: str, observation: Observation) -> ValidationResult:
        """
        Validate if observed behavior matches specification.
        
        Uses LLM to interpret both spec and observation.
        """
        
        validation_prompt = f"""
Given specification: {spec}
Observed behavior: {observation.to_natural_language()}

Does the observed behavior match the specification?

Respond with:
- passed: true/false
- expected: what should have happened
- actual: what actually happened
- confidence: 0.0-1.0
"""
        
        response = self.base_model.generate(validation_prompt)
        return ValidationResult.from_llm_response(response)
```

### 2.4 Behavioral Knowledge Graph

```python
class BehavioralKnowledgeGraph:
    """
    Stores relationships between:
    - Board types
    - Code patterns
    - Observed behaviors
    - Success/failure modes
    - Style templates
    
    This is your moat. Embedder can't build this without vision.
    """
    
    def __init__(self):
        self.graph = Neo4j()  # Or similar graph DB
        self.vector_store = Qdrant()  # For semantic search
    
    def store_validated_pattern(
        self,
        spec: str,
        board: str,
        code: str,
        observation: Observation,
        style_compliant: bool
    ):
        """
        After successful validation, store all relationships.
        
        Creates graph:
        (Spec)-[:IMPLEMENTED_BY]->(Code)
        (Code)-[:RUNS_ON]->(Board)
        (Code)-[:PRODUCES]->(Observation)
        (Board)-[:HAS_QUIRK]->(SiliconBug)
        (Code)-[:FOLLOWS_STYLE]->(StyleTemplate)
        """
        
        # Extract patterns from code
        patterns = extract_code_patterns(code)
        
        # Store relationships
        self.graph.create_relationships([
            (f"Spec:{hash(spec)}", "IMPLEMENTED_BY", f"Code:{hash(code)}"),
            (f"Code:{hash(code)}", "RUNS_ON", f"Board:{board}"),
            (f"Code:{hash(code)}", "PRODUCES", f"Observation:{observation.id}"),
            (f"Board:{board}", "VALIDATED_PATTERN", f"Code:{hash(code)}"),
        ])
        
        # Store for semantic search
        self.vector_store.add(
            text=f"{spec} -> {observation.to_natural_language()}",
            metadata={
                "board": board,
                "code": code,
                "style_compliant": style_compliant,
                "validated_at": datetime.now().isoformat()
            }
        )
    
    def get_quirks(self, board: str) -> str:
        """
        Get known quirks/bugs for this board.
        
        Example return:
        "ESP32-S3 rev 1.1:
         - I2C has silicon bug, use clock stretching workaround
         - Timer prescaler must account for 80MHz crystal (not 40MHz)
         - GPIO2 controls onboard LED (datasheet says GPIO5)"
        """
        
        quirks = self.graph.query(f"""
            MATCH (b:Board {{type: '{board}'}})-[:HAS_QUIRK]->(q:Quirk)
            WHERE q.validation_count > 5
            RETURN q.description, q.workaround, q.validation_count
            ORDER BY q.validation_count DESC
        """)
        
        return format_quirks(quirks)
    
    def find_patterns(self, spec: str, board: str) -> list[CodePattern]:
        """
        Find validated code patterns similar to this spec.
        
        Uses semantic search + graph filtering.
        """
        
        # Semantic search
        similar_specs = self.vector_store.search(
            query=spec,
            filter={"board": board, "style_compliant": True},
            top_k=5
        )
        
        return [
            CodePattern(
                code=result.metadata["code"],
                observation=result.metadata["observation"],
                validation_count=result.metadata.get("validation_count", 1)
            )
            for result in similar_specs
        ]
    
    def get_failure_hints(self, result: ValidationResult) -> str:
        """
        Given a validation failure, suggest common causes.
        
        Example:
        "LED blinking at 2Hz instead of 1Hz on ESP32-DevKitC is commonly
         caused by incorrect timer prescaler calculation. The board uses
         an 80MHz crystal, not 40MHz as stated in some datasheets."
        """
        
        # Query for similar failures
        similar_failures = self.graph.query(f"""
            MATCH (b:Board)-[:HAD_FAILURE]->(f:Failure)
            WHERE f.expected = '{result.expected}'
              AND f.actual = '{result.actual}'
            RETURN f.cause, f.fix, COUNT(*) as occurrences
            ORDER BY occurrences DESC
            LIMIT 3
        """)
        
        return format_failure_hints(similar_failures)
```

### 2.5 Style Knowledge Graph

```python
class StyleKnowledgeGraph:
    """
    Stores relationships between:
    - Coding standards (BARR-C, MISRA-C)
    - Board types
    - Code templates
    - Naming conventions
    
    Ensures generated code looks like it was written by a senior engineer.
    """
    
    def get_template(self, spec: str, board: str) -> StyleTemplate:
        """
        Get style template for this type of code.
        
        Example for "Blink LED" on ESP32:
        - Non-blocking timer-based architecture
        - Proper function naming: StatusLED_Toggle()
        - Type-safe constants: const uint16_t LED_PERIOD_MS = 1000U;
        - Doxygen comments
        """
        
        templates = self.graph.query(f"""
            MATCH (s:Spec)-[:HAS_TEMPLATE]->(t:StyleTemplate)
            WHERE s.category = '{categorize_spec(spec)}'
              AND t.board_type = '{board}'
              AND t.style_standard = 'BARR-C'
            RETURN t.template, t.example_code
            ORDER BY t.usage_count DESC
            LIMIT 1
        """)
        
        if templates:
            return StyleTemplate.from_db(templates[0])
        else:
            # Return generic BARR-C template
            return StyleTemplate.default_barrc()
```

### 2.6 Style Checker (BARR-C Enforcement)

```python
class BARRCStyleChecker:
    """
    Enforces BARR-C Embedded C Coding Standard.
    https://barrgroup.com/embedded-systems/books/embedded-c-coding-standard
    """
    
    RULES = {
        # Naming conventions (Rule 2)
        'function_names': r'^[A-Z][a-zA-Z0-9]*_[A-Z][a-zA-Z0-9]*$',  # Module_Function
        'variable_names': r'^[a-z][a-z0-9_]*$',  # snake_case
        'constant_names': r'^[A-Z][A-Z0-9_]*$',  # UPPER_SNAKE
        'type_names': r'^[a-z][a-z0-9_]*_t$',  # name_t
        
        # Type safety (Rule 3)
        'use_stdint': True,  # uint8_t not unsigned char
        'no_implicit_conversions': True,
        
        # Magic numbers (Rule 4)
        'no_magic_numbers': True,  # #define not hardcoded
        
        # Const correctness (Rule 5)
        'const_pointers': True,  # const uint8_t*
        
        # Architecture (Rule 6-8)
        'no_dynamic_allocation': True,  # no malloc/free
        'no_recursion': True,  # stack safety
        'explicit_casts': True,  # (uint16_t)value
        
        # Comments (Rule 9)
        'doxygen_functions': True,  # /** ... */ for all functions
    }
    
    def check(self, code: str) -> list[StyleViolation]:
        """Check code against BARR-C rules."""
        violations = []
        ast = parse_c_code(code)
        
        # Check function names
        for func in ast.functions:
            if not re.match(self.RULES['function_names'], func.name):
                violations.append(StyleViolation(
                    rule='function_names',
                    line=func.line,
                    severity='error',
                    message=f"Function '{func.name}' should use Module_Function format",
                    suggestion=self.suggest_function_name(func.name)
                ))
            
            # Check for Doxygen comment
            if not func.has_doxygen_comment():
                violations.append(StyleViolation(
                    rule='doxygen_functions',
                    line=func.line,
                    severity='warning',
                    message=f"Function '{func.name}' missing Doxygen comment"
                ))
        
        # Check for magic numbers
        for number in ast.numeric_literals:
            if number.is_magic():  # Not 0, 1, or in #define
                violations.append(StyleViolation(
                    rule='no_magic_numbers',
                    line=number.line,
                    severity='warning',
                    message=f"Magic number {number.value} should be #define constant"
                ))
        
        # Check type usage
        for var in ast.variables:
            if not var.uses_stdint():
                violations.append(StyleViolation(
                    rule='use_stdint',
                    line=var.line,
                    severity='error',
                    message=f"Use uint8_t/uint16_t instead of {var.type}"
                ))
        
        return violations
    
    def enforce(self, code: str) -> str:
        """Auto-fix violations where possible."""
        violations = self.check(code)
        
        fixed_code = code
        for v in violations:
            if v.rule == 'function_names' and v.severity == 'error':
                fixed_code = self.fix_function_name(fixed_code, v)
            elif v.rule == 'doxygen_functions':
                fixed_code = self.add_doxygen_template(fixed_code, v)
            elif v.rule == 'use_stdint':
                fixed_code = self.replace_type(fixed_code, v)
        
        return fixed_code
    
    def add_doxygen_template(self, code: str, violation: StyleViolation) -> str:
        """Add Doxygen comment template before function."""
        template = """
/**
 * @brief [Brief description]
 * @param [parameter name] [description]
 * @return [return value description]
 * @note [additional notes if needed]
 */
"""
        # Insert template before function
        lines = code.split('\n')
        lines.insert(violation.line - 1, template)
        return '\n'.join(lines)
```

### 2.7 CLI Interface (Code Generation)

```bash
# Generate firmware with validation
$ percepta generate "Blink LED at 1Hz" \
    --board esp32-devkit-v1 \
    --validate

Output:
[1] Generating firmware from specification...
    ✅ Code generated (BARR-C compliant)
    
[2] Flashing to board esp32-devkit-v1...
    ✅ Flashed successfully
    
[3] Observing hardware behavior...
    📹 Capturing video...
    🔍 LED is blinking at 2.1 Hz (expected 1.0 Hz)
    ❌ Validation failed
    
[4] Regenerating with feedback...
    💡 Context: ESP32-DevKitC uses 80MHz crystal, recalculating prescaler
    ✅ Code updated
    
[5] Flashing updated firmware...
    ✅ Flashed successfully
    
[6] Re-observing hardware behavior...
    📹 Capturing video...
    🔍 LED is blinking at 1.02 Hz
    ✅ Validation passed!
    
📄 Firmware written to: ./src/main.c
📊 Stored as validated pattern for future use
⏱️  Total time: 45 seconds (2 iterations)

Would you like to:
  [1] View the generated code
  [2] Test additional scenarios
  [3] Deploy to production
  [q] Quit
```

```bash
# Compare Percepta vs Embedder
$ percepta benchmark "Blink LED at 1Hz" \
    --boards esp32-devkit-v1,stm32-nucleo \
    --compare-with embedder

Output:
Benchmark Results: "Blink LED at 1Hz"

┌────────────┬────────────┬─────────────────┬──────────────┬─────────────┐
│ Tool       │ Compiles   │ Works on HW     │ Style Check  │ Time        │
├────────────┼────────────┼─────────────────┼──────────────┼─────────────┤
│ Embedder   │ ✅ 100%    │ ⚠️  95% (est)   │ ❌ 45% pass  │ 10s         │
│ Percepta   │ ✅ 100%    │ ✅ 100% (valid) │ ✅ 98% pass  │ 45s (2 iter)│
└────────────┴────────────┴─────────────────┴──────────────┴─────────────┘

Winner: Percepta
  - 100% hardware validation (Embedder: manual testing required)
  - 98% BARR-C compliance (Embedder: 45%)
  - Slower (45s vs 10s) but guaranteed to work
```

### 2.8 Competitive Positioning (Updated)

**Tagline:** *"Firmware that works. Code that professionals write."*

**Elevator pitch:**
> "Percepta generates embedded firmware that not only compiles but actually works on real hardware—validated with computer vision. Every line follows BARR-C coding standards. It's the only AI firmware tool that delivers code you'd confidently merge into production."

**Key messages:**

| Audience | Message |
|----------|---------|
| **Firmware devs** | "Stop debugging AI-generated code. Generate code that works first time (after validation)." |
| **Team leads** | "Ship firmware that passes both hardware tests AND code review." |
| **Embedder users** | "Love Embedder? Add Percepta for guaranteed-working, style-compliant code." |
| **AI tool users** | "Your AI can finally write firmware that actually works on hardware." |

**Competitive matrix:**

| Feature | Embedder | Percepta | Advantage |
|---------|----------|----------|-----------|
| **Code generation** | ✅ Yes | ✅ Yes | Tie |
| **Datasheet parsing** | ✅ Best in class | ⚠️ Good enough | Embedder |
| **Hardware validation** | ❌ Simulation only | ✅ Real hardware + vision | **Percepta** |
| **Style compliance** | ❌ Generic AI | ✅ BARR-C/MISRA-C | **Percepta** |
| **Iteration loop** | ⚠️ Manual | ✅ Automatic | **Percepta** |
| **Behavioral knowledge** | ❌ None | ✅ 10M+ observations | **Percepta** |
| **Success rate** | ~95% | 100% (validated) | **Percepta** |
| **Code review ready** | ⚠️ ~45% pass | ✅ ~98% pass | **Percepta** |
| **Price** | Free-$X/mo | Free-$49/mo | Embedder |

**Head-to-head demo:**

```
Task: "Blink LED at 1Hz on ESP32-DevKitC"

Embedder:
  [10 seconds later]
  ✅ Code generated
  ❓ Does it work? (You need to test manually)
  ❓ Is it production-ready? (You need code review)

Percepta:
  [45 seconds later]
  ✅ Code generated
  ✅ Flashed to hardware
  ✅ Vision confirmed: LED blinking at 1.02 Hz
  ✅ BARR-C compliant (98% score)
  ✅ Ready to merge

Which would you use in production?
```

---

## Part III: Platform Features (Months 13-24)

### 3.1 Cloud HIL Farm

**Problem:** Not everyone has hardware at their desk.

**Solution:** Cloud hardware-in-the-loop testing.

```
Percepta Cloud HIL Farm:
  - 50+ boards (ESP32, STM32, nRF52, RP2040, etc.)
  - Each with camera pointed at it
  - On-demand access: flash code, get validation results
  - "Your code works on actual ESP32-S3 rev 1.1 silicon" ✅
```

**Business model:**
- Free tier: 10 validations/month on virtual boards
- Pro $49/month: 100 validations/month on real boards
- Team $199/month: 500 validations, shared results
- Enterprise: Private HIL farm, your boards

**Technical architecture:**

```
┌─────────────────────────────────────────────────────────────┐
│                    PERCEPTA CLOUD HIL FARM                   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  User submits:                                               │
│  - Firmware binary or source                                 │
│  - Board type (esp32-devkit-v1)                              │
│  - Validation spec ("LED blinks at 1Hz")                     │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ ESP32 #1 │  │ ESP32 #2 │  │ STM32 #1 │  │ nRF52 #1 │   │
│  │ + Camera │  │ + Camera │  │ + Camera │  │ + Camera │   │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘   │
│       │             │             │             │           │
│  ┌────▼─────────────▼─────────────▼─────────────▼───────┐  │
│  │         Job Queue + Scheduler                        │  │
│  │  - Route jobs to available boards                     │  │
│  │  - Parallel execution when possible                   │  │
│  │  - Queue during peak times                            │  │
│  └───────────────────────┬──────────────────────────────┘  │
│                          │                                  │
│  Returns:                │                                  │
│  - Validation result     │                                  │
│  - Recorded video        │                                  │
│  - Observation data      │                                  │
│  - "✅ Works on rev 1.1" │                                  │
│                          │                                  │
└──────────────────────────┴──────────────────────────────────┘
```

**Usage:**

```bash
# Submit job to cloud HIL farm
$ percepta cloud validate \
    --code ./src/main.c \
    --board esp32-s3 \
    --spec "LED blinks at 1Hz"

Output:
🚀 Job submitted: job_a1b2c3
📊 Queue position: 2
⏱️  Estimated wait: 30 seconds

[30 seconds later]
✅ Job complete!

Results:
  Board: ESP32-S3-DevKitC-1 (rev 1.1)
  Status: ✅ PASSED
  LED frequency: 1.02 Hz (within tolerance)
  Boot time: 2.1s
  
📹 Video: https://percepta.cloud/jobs/job_a1b2c3/video.mp4
📄 Full report: https://percepta.cloud/jobs/job_a1b2c3/report
```

### 3.2 Community Code Hub

**Think: Hugging Face for embedded firmware**

```
Percepta Hub (hub.percepta.dev):
  - Hardware-verified code patterns
  - Community contributions
  - Search: "LED blink ESP32" → 147 validated implementations
  - Download knowing it ACTUALLY WORKS on hardware
```

**Features:**

1. **Upload & Verify:**
   ```bash
   $ percepta hub publish \
       --code ./my_driver.c \
       --board esp32-devkit-v1 \
       --description "Non-blocking UART driver"
   
   Output:
   📤 Uploading code...
   ✅ Code published
   🔬 Validation job queued...
   
   [2 minutes later]
   ✅ Hardware validation passed!
   🏷️  Tagged as: uart, non-blocking, esp32
   🔗 https://hub.percepta.dev/drivers/uart-nonblocking-esp32
   ```

2. **Search & Download:**
   ```bash
   $ percepta hub search "I2C driver STM32"
   
   Results:
   1. I2C Master (STM32F4) - ✅ Verified on 15 boards
      By: @embedded_pro | Downloads: 1,247
      
   2. I2C DMA Driver (STM32H7) - ✅ Verified on 8 boards
      By: @stm_expert | Downloads: 892
      
   3. I2C Slave (STM32L4) - ✅ Verified on 5 boards
      By: @lowpower_dev | Downloads: 534
   
   $ percepta hub download 1 --board stm32f4-discovery
   
   ✅ Downloaded and adapted for stm32f4-discovery
   📄 Code saved to: ./drivers/i2c_master.c
   📊 This pattern has been validated 15 times on real hardware
   ```

3. **Leaderboard:**
   ```
   Top Contributors (This Month):
   1. @embedded_pro - 23 patterns, 15K downloads
   2. @stm_expert - 18 patterns, 12K downloads
   3. @esp_wizard - 15 patterns, 9K downloads
   ```

**Moat:** Network effects. More users → more patterns → better code → more users.

### 3.3 Team Collaboration Features

**For teams of 5-50 engineers:**

```yaml
# Team dashboard features:
- Shared behavioral knowledge graph
- Team device library (all boards in one place)
- Shared validation history
- Code generation templates
- Style guide enforcement (company-specific)
- SSO integration
- Audit logs
```

**Usage:**

```bash
# Team workspace
$ percepta team init "AcmeCorp Engineering"

$ percepta team add-device \
    --name "lab-esp32-1" \
    --board esp32-devkit-v1 \
    --location "Hardware Lab, Desk 3"

$ percepta team set-style-guide \
    --file ./acmecorp-style-guide.yaml

# Now entire team uses same standards
$ percepta generate "Blink LED" \
    --board lab-esp32-1 \
    --validate

Output:
✅ Code generated using AcmeCorp style guide
✅ Validated on lab-esp32-1 (team device)
📊 Result shared with team workspace
```

### 3.4 CI/CD Integration (Production-Ready)

**GitHub Actions example:**

```yaml
# .github/workflows/firmware-validation.yml
name: Hardware Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Percepta
        uses: percepta/setup-action@v1
        with:
          version: 'latest'
          api-key: ${{ secrets.PERCEPTA_API_KEY }}
      
      - name: Build firmware
        run: make build
      
      - name: Validate on Cloud HIL
        run: |
          percepta cloud validate \
            --code ./build/firmware.bin \
            --board esp32-s3 \
            --spec-file ./tests/hardware-specs.yaml \
            --fail-on-regression
      
      - name: Generate report
        if: always()
        run: |
          percepta report \
            --format markdown \
            --output $GITHUB_STEP_SUMMARY
      
      - name: Upload artifacts
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: validation-failure
          path: ./percepta-artifacts/
```

**Result:** Every PR gets hardware validation before merge.

### 3.5 VSCode Extension

**Features:**
- Inline code generation: Right-click → "Generate with Percepta"
- Real-time validation: See hardware state in sidebar
- Style warnings: BARR-C violations highlighted
- One-click validation: "Test on Hardware" button

```typescript
// VSCode extension example
// src/extension.ts

export function activate(context: vscode.ExtensionContext) {
    // Register code generation command
    const generateCommand = vscode.commands.registerCommand(
        'percepta.generate',
        async () => {
            const editor = vscode.window.activeTextEditor;
            if (!editor) return;
            
            // Get user specification
            const spec = await vscode.window.showInputBox({
                prompt: "Describe what you want to generate",
                placeHolder: "e.g., Blink LED at 1Hz"
            });
            
            if (!spec) return;
            
            // Show progress
            await vscode.window.withProgress({
                location: vscode.ProgressLocation.Notification,
                title: "Generating firmware with Percepta",
                cancellable: false
            }, async (progress) => {
                progress.report({ message: "Generating code..." });
                
                // Call Percepta API
                const result = await perceptaClient.generate({
                    spec: spec,
                    board: getBoardFromConfig(),
                    validate: true
                });
                
                progress.report({ message: "Validating on hardware..." });
                
                // Insert generated code
                editor.edit(editBuilder => {
                    editBuilder.insert(
                        editor.selection.active,
                        result.code
                    );
                });
                
                // Show validation result
                vscode.window.showInformationMessage(
                    `✅ Code generated and validated! ${result.iterations} iterations`
                );
            });
        }
    );
    
    context.subscriptions.push(generateCommand);
    
    // Register real-time validation
    const validationProvider = new PerceptaValidationProvider();
    context.subscriptions.push(
        vscode.workspace.onDidSaveTextDocument(
            doc => validationProvider.validate(doc)
        )
    );
}
```

---

## Part IV: Business Model & Economics

### 4.1 Revenue Streams

| Tier | Price | Target | Features |
|------|-------|--------|----------|
| **Free** | $0 | Hobbyists, students | 10 cloud validations/month, local unlimited |
| **Pro** | $49/month | Individual devs | 100 cloud validations, priority queue, advanced assertions |
| **Team** | $199/month | Small teams (5-10) | 500 validations, shared workspace, SSO |
| **Enterprise** | $999+/month | Large orgs | Unlimited, private HIL farm, SLA, dedicated support |

### 4.2 Unit Economics

**Pro User ($49/month):**
- Revenue: $588/year
- COGS:
  - API costs (Claude): ~$10/month = $120/year
  - HIL farm compute: ~$5/month = $60/year
  - Infrastructure: ~$2/month = $24/year
  - Total COGS: $204/year
- **Gross Margin: 65%**
- CAC: $50 (content marketing)
- Payback: 1 month
- LTV: $1,764 (3-year retention)
- **LTV/CAC: 35x** ← Excellent

**Team Account ($199/month):**
- Revenue: $2,388/year
- COGS: ~$600/year
- **Gross Margin: 75%**
- CAC: $200 (demo + sales)
- Payback: 1 month
- LTV: $7,164 (3-year retention)
- **LTV/CAC: 36x**

### 4.3 Growth Projections

**Year 1:**
- Month 6: Perception only, 1,000 free users
- Month 12: Code gen launched, 5,000 free, 200 Pro = $10K MRR
- ARR: $120K

**Year 2:**
- Month 18: 15,000 free, 800 Pro, 50 Team = $50K MRR
- Month 24: 30,000 free, 1,500 Pro, 100 Team = $90K MRR
- ARR: $1.08M

**Year 3:**
- 100,000 free users
- 5,000 Pro = $245K MRR
- 300 Team = $60K MRR
- 20 Enterprise = $100K MRR
- **ARR: $4.86M**

**Year 4-5:**
- Scale to $10M+ ARR
- Exit: Acquisition $50-100M or continue growing

### 4.4 Funding Strategy

**Bootstrap Phase (Months 0-12):**
- Self-funded or revenue from consulting
- Build perception MVP → prove traction
- Costs: $2-5K/month (APIs, hosting)

**Seed Round (Month 12): $1-2M**
- After: 5,000 users, $10K+ MRR
- Valuation: $5-8M post-money
- Use: Team (3 people), HIL infrastructure, marketing
- Dilution: 20-30%

**Series A (Month 24): $5-10M**
- After: 30,000 users, $100K+ MRR
- Valuation: $25-40M post-money
- Use: Scale team to 15, enterprise sales, international
- Dilution: 20-25%

**Exit Options (Month 36-48):**
- Strategic acquisition: $50-100M
  - Anthropic (perception for Claude)
  - Semiconductor company (e.g., Espressif, STMicro)
  - Embedded tools vendor
- IPO/growth rounds if hitting $20M+ ARR

---

## Part V: Go-to-Market Strategy

### 5.1 Phased Launch Plan

**Phase 1: Stealth Perception (Months 0-6)**
- Build perception MVP in private
- Alpha with 10-20 hand-picked users
- "We're just the eyes" positioning
- **Goal:** 500 users, prove PMF

**Phase 2: The Reveal (Months 7-9)**
- Public announcement: "We generate better code than Embedder"
- Demo video: Side-by-side comparison
- Marketing: "Code that actually works on hardware"
- **Goal:** 200 paying users

**Phase 3: Platform Launch (Months 10-24)**
- HIL farm goes live
- Community hub opens
- Enterprise features
- **Goal:** Industry standard

### 5.2 Marketing Channels

**Content Marketing:**
- Blog: "Why 95% Success Isn't Good Enough for Embedded"
- Blog: "We Tested 1,000 AI-Generated Firmware Files. Here's What Happened."
- YouTube: Weekly embedded dev tips
- Podcast tour: Embedded.fm, etc.

**Community:**
- r/embedded, r/rust, r/esp32
- Rust Embedded Working Group
- Embedder Discord/community
- Embedded conferences

**Partnerships:**
- Embedder: "They generate, we validate"
- Anthropic: "Claude Code + Percepta"
- Semiconductor vendors: "Validated on our boards"

**SEO:**
- Target: "embedded firmware generation"
- Target: "AI embedded code"
- Target: "hardware validation embedded"

### 5.3 Sales Strategy

**Self-serve (Pro tier):**
- Sign up → 10 free validations
- Auto-upgrade prompt after 10th validation
- Conversion rate target: 15%

**Sales-assisted (Team/Enterprise):**
- Demo-first approach
- 14-day free trial with dedicated support
- ROI calculator: "How much time does your team spend debugging firmware?"

**Partner channels:**
- Semiconductor FAEs
- Embedded consulting firms
- University engineering programs

---

## Part VI: Technical Roadmap (Detailed)

### 6.1 Phase 1: Perception MVP (Weeks 1-8) ✅

*See Percepta_PRD_Final.md*

**Deliverables:**
- `percepta observe/assert/diff` working
- SQLite storage
- Basic CLI
- 10 alpha users validating hardware

### 6.2 Phase 2A: Code Gen Infrastructure (Weeks 9-12)

**Week 9-10: Style Checker**
- Implement BARR-C rule engine
- Auto-fix violations
- CLI: `percepta style-check ./src/`
- Test on 1000+ embedded code samples

**Week 11-12: Knowledge Graphs**
- Set up Neo4j for behavioral knowledge
- Set up Qdrant for semantic search
- Import 100K observations from Phase 1
- Build initial quirks database (10 boards)

**Exit criteria:** Style checker 95%+ accurate, knowledge graphs operational

### 6.3 Phase 2B: Code Generation Engine (Weeks 13-16)

**Week 13-14: Base Model Integration**
- Fine-tune Claude on style-compliant code corpus
- Build prompt engineering system
- Test generation quality (compile rate, style compliance)

**Week 14-15: Validation Loop**
- Integrate vision + generation
- Build automatic retry logic
- Test on 5 boards, 10 specifications
- Target: 100% success rate within 5 iterations

**Week 16: Private Beta**
- Invite 20 users to test code generation
- Collect feedback
- Measure: Compile rate, validation success, time per generation

**Exit criteria:** 90%+ compile rate, 80%+ validation success, <2 min per generation

### 6.4 Phase 2C: Public Code Gen Launch (Weeks 17-20)

**Week 17-18: Polish**
- CLI refinement
- Error message improvements
- Documentation (20+ example specs)
- Demo videos (5 different use cases)

**Week 19: Launch Marketing**
- Blog post: "Introducing Percepta Code Generation"
- Hacker News: "We Built AI Firmware Generation That Actually Works"
- Demo at Embedded World (if timing aligns)

**Week 20: Monitoring & Iteration**
- Track: Conversion rate (perception → code gen)
- Track: Generation success rate
- Track: User satisfaction (NPS survey)

**Exit criteria:** 200 paying users, $10K MRR

### 6.5 Phase 3: Platform Features (Months 5-12)

**Month 5-6: Cloud HIL Infrastructure**
- Rent co-location space
- Buy 20 boards (5x ESP32, 5x STM32, 5x nRF52, 5x RP2040)
- Build remote flashing/observation system
- Launch closed beta with 50 users

**Month 7-8: Community Hub**
- Build hub.percepta.dev
- Enable code uploads + verification
- Gamification (leaderboards, badges)
- Seed with 100 initial patterns

**Month 9-10: Team Features**
- Shared workspaces
- SSO integration
- Team analytics dashboard
- Private style guides

**Month 11-12: Enterprise Features**
- On-prem deployment option
- Air-gapped mode
- SLA guarantees
- Dedicated support

**Exit criteria:** $100K MRR, 1000 paying users, recognized as Embedder competitor

### 6.6 Phase 4: Market Domination (Months 13-24)

**Priorities:**
- Scale HIL farm to 50+ boards
- International expansion (EU, Asia)
- Semiconductor partnerships (Espressif, STMicro official integration)
- Embedder acquisition discussions
- Series A fundraise

---

## Part VII: Risk Analysis & Mitigation

### 7.1 Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Generation quality insufficient** | Medium | High | Fine-tune on validated corpus, not just GitHub; validation loop catches issues |
| **Validation loop too slow** | Medium | Medium | Parallel processing, faster flashing, optimize vision inference |
| **HIL farm reliability issues** | Medium | High | Redundant boards, automated health checks, 99% uptime SLA |
| **Style checker false positives** | Low | Medium | Extensive testing, user override options |
| **Knowledge graph doesn't scale** | Low | Medium | Use proven graph DB (Neo4j), shard by board type |

### 7.2 Business Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Embedder builds perception** | Medium | High | 12-month head start, partnership approach first |
| **Market too small** | Low | High | TAM analysis shows 10K-50K potential customers |
| **Can't monetize** | Low | High | Proven willingness to pay (Embedder has paying customers) |
| **Team burnout** | Medium | High | Hire early (Month 7), sustainable pace |
| **Competition from open source** | Low | Medium | Network effects (knowledge graph), proprietary models |

### 7.3 Execution Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Trying to build too much** | High | High | Strict Phase 1 → Phase 2 → Phase 3 discipline |
| **Losing focus** | Medium | High | Weekly review: "Does this serve perception or generation?" |
| **Poor partnerships** | Medium | Medium | Start with small integrations, prove value before going deep |
| **Hiring too fast** | Low | Medium | Bootstrap to $10K MRR before hiring |

---

## Part VIII: Success Metrics & KPIs

### 8.1 Product Metrics

| Metric | Month 6 | Month 12 | Month 24 |
|--------|---------|----------|----------|
| **WAU (Weekly Active Users)** | 500 | 1,500 | 5,000 |
| **DAU/WAU ratio** | 30% | 35% | 40% |
| **Observations/user/week** | 10 | 15 | 20 |
| **Generation success rate** | N/A | 95% | 98% |
| **Avg iterations per generation** | N/A | 2.5 | 2.0 |
| **Knowledge graph size** | 100K obs | 1M obs | 10M obs |

### 8.2 Business Metrics

| Metric | Month 6 | Month 12 | Month 24 |
|--------|---------|----------|----------|
| **MRR** | $0 | $10K | $100K |
| **Paying customers** | 0 | 200 | 1,000 |
| **Free→Paid conversion** | N/A | 10% | 15% |
| **Churn rate** | N/A | <5% | <3% |
| **NPS** | 50+ | 60+ | 70+ |
| **CAC** | N/A | $50 | $40 |
| **LTV** | N/A | $1,500 | $2,000 |

### 8.3 Technical Metrics

| Metric | Month 6 | Month 12 | Month 24 |
|--------|---------|----------|----------|
| **Vision accuracy** | 95% | 96% | 97% |
| **Compile rate** | N/A | 98% | 99% |
| **Style compliance** | N/A | 90% | 95% |
| **API uptime** | 99.5% | 99.9% | 99.95% |
| **Avg response time** | <2s | <1.5s | <1s |
| **HIL farm utilization** | N/A | 40% | 70% |

---

## Part IX: Team & Organization

### 9.1 Team Structure (Month 12)

**Founding Team (3 people):**

**1. Technical Lead / CEO (You)**
- Vision & strategy
- Core platform architecture
- Fundraising & partnerships

**2. ML/Embedded Engineer (Hire Month 7)**
- Code generation engine
- Model fine-tuning
- Knowledge graph

**3. Full-stack Engineer (Hire Month 10)**
- Web dashboard
- HIL farm infrastructure
- VSCode extension

### 9.2 Team Structure (Month 24)

**10-15 people:**

**Engineering (7):**
- Staff Engineer (platform lead)
- ML Engineer (generation)
- Embedded Engineer (validation)
- Full-stack Engineer #1 (dashboard)
- Full-stack Engineer #2 (HIL farm)
- DevOps Engineer (infrastructure)
- QA Engineer (testing)

**Product & Design (2):**
- Product Manager
- Product Designer

**Go-to-Market (3-5):**
- Head of Marketing
- Developer Advocate
- Sales Engineer (enterprise)
- Customer Success Manager

**Operations (1):**
- Operations Manager

### 9.3 Advisory Board

**Ideal advisors:**
- Embedded systems veteran (e.g., former ARM/Qualcomm)
- AI/ML expert (e.g., Anthropic, OpenAI alumni)
- Open source community leader (e.g., Rust Embedded WG)
- GTM expert (e.g., former DevTools exec)

---

## Part X: Conclusion & Next Steps

### The Complete Vision

Percepta is the **only platform** that combines:
1. ✅ Vision-based hardware observation (Phase 1)
2. ✅ AI firmware generation with validation (Phase 2)
3. ✅ Professional code quality (BARR-C/MISRA-C)
4. ✅ Behavioral knowledge graph (what actually works)
5. ✅ Cloud HIL farm (on-demand validation)
6. ✅ Community code hub (network effects)

**Result:** The firmware platform professionals trust.

### Why This Wins

**Three-layer moat:**
1. **Technical:** Vision expertise (12+ months to replicate)
2. **Data:** Behavioral knowledge graph (10M+ observations)
3. **Network:** Community patterns (more users = better code)

**Market timing:**
- Embedder has proven the market (2000+ users)
- MCP ecosystem is exploding (200+ servers)
- AI coding tools are mainstream (Claude Code, Cursor)
- Gap is obvious (no one validates physical behavior)

**Competitive advantages:**
- First-mover in vision-based embedded observation
- Only tool with automatic validation loop
- Only tool that enforces professional code style
- Strong partnership opportunity with Embedder

### Path to $100M Exit

**Month 6:** Perception PMF → 500 users  
**Month 12:** Code gen launch → $10K MRR → Seed round  
**Month 24:** Platform features → $100K MRR → Series A  
**Month 36:** Market leader → $10M ARR → Exit discussions

**Valuation at exit:** $50-100M (5-10x ARR multiple)

### What Could Go Wrong?

1. **Trying to build everything at once** → Solution: Strict phasing
2. **Embedder builds perception first** → Solution: Move fast, partnership
3. **Market smaller than expected** → Solution: $10M ARR is still great
4. **Team burns out** → Solution: Sustainable pace, hire early

### The Discipline Required

**Phase 1 (Months 0-6): ONLY perception**
- Don't build code generation yet
- Don't build HIL farm yet
- Don't hire yet
- Just: observe(), assert(), diff()

**Phase 2 (Months 7-12): Add generation**
- Only after Phase 1 PMF proven
- Start with 1 engineer hire
- Validate before scaling

**Phase 3 (Months 13-24): Scale platform**
- Only after paying customers
- HIL farm, community, enterprise
- Now scale team

**The trap:** Trying to build Phases 1+2+3 simultaneously. That's how you fail.

---

## Immediate Next Steps

### This Week:
1. ✅ Review and approve this PRD
2. ✅ Decide: Perception-first (recommended) or full platform immediately?
3. ✅ Set up repository structure
4. ✅ Recruit first 5 alpha testers

### Next Month:
1. ✅ Ship perception MVP to alpha users
2. ✅ Validate PMF (retention, usage frequency)
3. ✅ Start collecting behavioral knowledge graph data
4. ✅ Begin planning code generation layer (if perception works)

### Month 6 Decision Point:
- ✅ **If perception has PMF:** Proceed to Phase 2 (code generation)
- ❌ **If perception doesn't have PMF:** Pivot or shut down

**Don't skip Phase 1. It's your foundation.**

---

## Contact & Resources

**Project:** Percepta - Vision-Based Embedded AI Platform  
**Stage:** Pre-launch (Planning)  
**Target Launch:** 14 feb 2026  
**Team:** Claude and Utkarsh  
**Contact:** utkarsh@kernex.sbs

**This PRD represents the complete vision. Execute it in phases.**

**Ship Phase 1 first. Prove it works. Then build the platform.**

---

**Built to make embedded AI actually work on real hardware.**  
**Generate professional code. Validate everything. Win the market.**

🚀 **Let's build this.**
