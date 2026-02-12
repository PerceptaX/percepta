// Example C code with BARR-C violations for testing

#include <stdint.h>

// VIOLATION: Function naming - should be LED_Init
void initLED() {
    // VIOLATION: Variable naming - should be led_pin
    int ledPin = 13;

    // VIOLATION: Using unsigned char instead of uint8_t
    unsigned char brightness = 255;
}

// VIOLATION: Pointer without const
void processData(uint8_t* data) {
    // Implementation
}

// GOOD: Proper BARR-C compliant code
void LED_Init(void) {
    uint8_t led_pin = 13;
    const uint8_t max_brightness = 255;
}

void LED_SetBrightness(const uint8_t* brightness_ptr) {
    // Implementation
}
