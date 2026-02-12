//go:build darwin

package camera

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AVFoundation -framework CoreMedia -framework CoreVideo -framework Foundation

#import <AVFoundation/AVFoundation.h>
#import <CoreMedia/CoreMedia.h>
#import <CoreVideo/CoreVideo.h>
#include <stdlib.h>

// Capture a single frame from the default camera
// Returns JPEG data and size, caller must free the returned buffer
unsigned char* captureFrame(int* outSize, char** outError) {
    @autoreleasepool {
        // Get default video device
        AVCaptureDevice *device = [AVCaptureDevice defaultDeviceWithMediaType:AVMediaTypeVideo];
        if (!device) {
            *outError = strdup("No camera device found");
            return NULL;
        }

        NSError *error = nil;
        AVCaptureDeviceInput *input = [AVCaptureDeviceInput deviceInputWithDevice:device error:&error];
        if (!input) {
            *outError = strdup([[error localizedDescription] UTF8String]);
            return NULL;
        }

        // Create capture session
        AVCaptureSession *session = [[AVCaptureSession alloc] init];
        session.sessionPreset = AVCaptureSessionPresetPhoto;

        if (![session canAddInput:input]) {
            *outError = strdup("Cannot add camera input to session");
            return NULL;
        }
        [session addInput:input];

        // Create output
        AVCapturePhotoOutput *output = [[AVCapturePhotoOutput alloc] init];
        if (![session canAddOutput:output]) {
            *outError = strdup("Cannot add photo output to session");
            return NULL;
        }
        [session addOutput:output];

        // Start session
        [session startRunning];

        // Give camera time to initialize and adjust settings
        [NSThread sleepForTimeInterval:0.5];

        // Capture photo synchronously using a semaphore
        dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
        __block NSData *jpegData = nil;
        __block NSString *captureError = nil;

        AVCapturePhotoSettings *settings = [AVCapturePhotoSettings photoSettings];

        // Create a delegate class inline
        Class delegateClass = NSClassFromString(@"PhotoCaptureDelegate");
        if (!delegateClass) {
            delegateClass = objc_allocateClassPair([NSObject class], "PhotoCaptureDelegate", 0);
            class_addProtocol(delegateClass, @protocol(AVCapturePhotoCaptureDelegate));

            IMP captureImp = imp_implementationWithBlock(^(id self, AVCapturePhoto *photo, NSError *err) {
                if (err) {
                    captureError = [err localizedDescription];
                } else {
                    jpegData = [photo fileDataRepresentation];
                }
                dispatch_semaphore_signal(semaphore);
            });

            class_addMethod(delegateClass,
                          @selector(captureOutput:didFinishProcessingPhoto:error:),
                          captureImp,
                          "v@:@@@");

            objc_registerClassPair(delegateClass);
        }

        id delegate = [[delegateClass alloc] init];
        [output capturePhotoWithSettings:settings delegate:delegate];

        // Wait for capture to complete (5 second timeout)
        dispatch_time_t timeout = dispatch_time(DISPATCH_TIME_NOW, 5 * NSEC_PER_SEC);
        long result = dispatch_semaphore_wait(semaphore, timeout);

        [session stopRunning];

        if (result != 0) {
            *outError = strdup("Camera capture timeout");
            return NULL;
        }

        if (captureError) {
            *outError = strdup([captureError UTF8String]);
            return NULL;
        }

        if (!jpegData) {
            *outError = strdup("No image data captured");
            return NULL;
        }

        // Copy JPEG data to malloc'd buffer that Go can free
        *outSize = (int)[jpegData length];
        unsigned char *buffer = (unsigned char*)malloc(*outSize);
        memcpy(buffer, [jpegData bytes], *outSize);

        return buffer;
    }
}
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/perceptumx/percepta/internal/core"
)

// AVFoundationCamera implements core.CameraDriver for macOS using AVFoundation
type AVFoundationCamera struct {
	devicePath string // Not used on macOS, kept for interface compatibility
}

// NewAVFoundationCamera creates a new macOS AVFoundation camera driver
func NewAVFoundationCamera(devicePath string) core.CameraDriver {
	return &AVFoundationCamera{devicePath: devicePath}
}

func (c *AVFoundationCamera) Open() error {
	// AVFoundation doesn't require explicit open - sessions are created per-capture
	// Just verify we can access the camera
	return nil
}

func (c *AVFoundationCamera) CaptureFrame() ([]byte, error) {
	var size C.int
	var errMsg *C.char
	defer func() {
		if errMsg != nil {
			C.free(unsafe.Pointer(errMsg))
		}
	}()

	// Capture frame using AVFoundation
	buffer := C.captureFrame(&size, &errMsg)
	if buffer == nil {
		if errMsg != nil {
			return nil, fmt.Errorf("camera capture failed: %s", C.GoString(errMsg))
		}
		return nil, fmt.Errorf("camera capture failed: unknown error")
	}
	defer C.free(unsafe.Pointer(buffer))

	// Copy C buffer to Go slice
	frame := C.GoBytes(unsafe.Pointer(buffer), size)
	return frame, nil
}

func (c *AVFoundationCamera) Close() error {
	// No cleanup needed - AVFoundation sessions are cleaned up automatically
	return nil
}
