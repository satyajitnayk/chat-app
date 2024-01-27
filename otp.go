package main

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OTP struct {
	Key     string
	Created time.Time
}

type RetentionMap map[string]OTP

// NewRetentionMap will create a new retentionmap and start the retention given the set period
func NewRetentionMap(ctx context.Context, retentionPeriod time.Duration) RetentionMap {
	rm := make(RetentionMap)

	go rm.InvalidateOTPs(ctx, retentionPeriod)

	return rm
}

// NewOTP creates and adds a new otp to the map
func (rm RetentionMap) NewOTP() OTP {
	o := OTP{
		Key:     uuid.NewString(),
		Created: time.Now(),
	}

	rm[o.Key] = o
	return o
}

// VerifyOTP will make sure a OTP exists
// and return true if so
// It will also delete the key so it cant be reused
func (rm RetentionMap) VerifyOTP(otp string) bool {
	if _, ok := rm[otp]; !ok {
		return false
	}
	// delete otp if it is used
	delete(rm, otp)
	return true
}

// Retention will make sure old OTPs are removed
// Is Blocking, so run as a Goroutine
func (rm RetentionMap) InvalidateOTPs(ctx context.Context, retentionPeriod time.Duration) {
	// ticker to recheck otp
	ticker := time.NewTicker(400 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			for _, otp := range rm {
				// check expired otps
				if otp.Created.Add(retentionPeriod).After(time.Now()) {
					delete(rm, otp.Key)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
