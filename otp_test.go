package main

import (
	"context"
	"testing"
	"time"
)

func TestRetentionMap_VerifyOTP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rm := NewRetentionMap(ctx, 1*time.Second)
	otp := rm.NewOTP()

	if ok := rm.VerifyOTP(otp.Key); !ok {
		t.Error("failed to verify otp key that exists")
	}

	if ok := rm.VerifyOTP(otp.Key); ok {
		t.Error("Reusing a OTP should not succeed")
	}
}

func TestOTP_InvalidateOTP(t *testing.T) {
	// create a context with cancel to stop goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create a retention map & add few otps with few seconds in b/w
	rm := NewRetentionMap(ctx, 1*time.Second)

	rm.NewOTP()
	rm.NewOTP()

	time.Sleep(2 * time.Second)

	otp := rm.NewOTP()

	// makesure only 1 otp left
	if len(rm) != 1 {
		t.Error("Failed to clean up")
	}

	if rm[otp.Key] != otp {
		t.Error("the key should be present")
	}
}
