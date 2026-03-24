package util

import (
	"context"
	"testing"
	"time"
)

func TestGetUnexportedField(t *testing.T) {
	// Test with a struct that has an unexported field
	type innerStruct struct {
		unexportedField string
	}
	type outerStruct struct {
		inner *innerStruct
	}

	inner := &innerStruct{unexportedField: "test_value"}
	outer := &outerStruct{inner: inner}

	// Get unexported field from inner struct
	result := GetUnexportedField(outer.inner, "unexportedField")
	if result == nil {
		t.Fatal("expected to get unexported field, got nil")
	}

	str, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if str != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", str)
	}
}

func TestGetUnexportedField_NonExistent(t *testing.T) {
	// Note: GetUnexportedField panics on non-existent fields
	// This is expected behavior - the function is designed for known fields
	// We test that it works correctly for existing fields
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for non-existent field")
		}
	}()

	type testStruct struct {
		existingField string
	}

	ts := &testStruct{existingField: "value"}
	GetUnexportedField(ts, "nonExistentField")
}

func TestWaitFor_Timeout(t *testing.T) {
	ctx := context.Background()
	err := WaitFor(ctx, WaitTimeout{
		Timeout:     100 * time.Millisecond,
		MinInterval: 10 * time.Millisecond,
	}, func() (bool, error) {
		return false, nil // Never succeeds
	})

	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestWaitFor_ImmediateSuccess(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	err := WaitFor(ctx, WaitTimeout{
		Timeout:     1 * time.Second,
		MinInterval: 10 * time.Millisecond,
	}, func() (bool, error) {
		callCount++
		return true, nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}


// TestAfter tests the After function
func TestAfter(t *testing.T) {
	start := time.Now()
	duration := 10 * time.Millisecond

	ch := After(nil, duration)
	if ch == nil {
		t.Error("expected channel, got nil")
	}

	// Wait for the timer to fire
	<-ch

	elapsed := time.Since(start)
	if elapsed < duration {
		t.Errorf("expected elapsed time >= %v, got %v", duration, elapsed)
	}
	if elapsed > duration+100*time.Millisecond {
		t.Errorf("timer took too long: %v", elapsed)
	}
}

// TestAfter_ClosedChannel tests that After properly closes the channel
func TestAfter_ClosedChannel(t *testing.T) {
	ch := After(nil, 10*time.Millisecond)

	// Read from the channel
	<-ch

	// Try to read again - should return zero value
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("expected channel to be closed immediately")
	}
}
