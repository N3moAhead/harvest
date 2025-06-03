package animation_test

import (
	"errors"
	"fmt"
	"image"
	"testing"

	"github.com/N3moAhead/harvest/internal/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

// newTestSourceImage is a helper to create a dummy *ebiten.Image for testing.
func newTestSourceImage(width, height int) *ebiten.Image {
	return ebiten.NewImage(width, height)
}

func TestNewAnimation(t *testing.T) {
	sourceImg := newTestSourceImage(100, 50) // A valid source image for most tests

	t.Run("ValidAnimation", func(t *testing.T) {
		anim, err := animation.NewAnimation(sourceImg, 10, 10, 0, 0, 5, 2, true)
		if err != nil {
			t.Fatalf("Expected no error for valid animation, got %v", err)
		}
		if anim == nil {
			t.Fatal("Expected animation to be non-nil for valid parameters")
		}
		// Basic check of initialized fields
		if anim.FrameWidth != 10 || anim.FrameHeight != 10 || anim.FrameAmount != 5 ||
			anim.AnimationSpeed != 2 || anim.Loops != true || anim.CurrentFrame != 0 ||
			anim.SourceImage != sourceImg {
			t.Errorf("Animation fields not initialized as expected: %+v", anim)
		}
	})

	// Test cases for invalid parameters
	testCases := []struct {
		name           string
		sourceImage    *ebiten.Image
		frameWidth     int
		frameHeight    int
		sourceStartX   int // Defaulted to 0 for these error tests
		sourceStartY   int // Defaulted to 0 for these error tests
		frameAmount    int
		animationSpeed int
		loops          bool // Defaulted to false for these error tests
		expectedError  error
	}{
		{"NilSourceImage", nil, 10, 10, 0, 0, 1, 1, false, animation.ErrNilSourceImage},
		{"ZeroFrameWidth", sourceImg, 0, 10, 0, 0, 1, 1, false, animation.ErrInvalidParameter},
		{"NegativeFrameWidth", sourceImg, -1, 10, 0, 0, 1, 1, false, animation.ErrInvalidParameter},
		{"ZeroFrameHeight", sourceImg, 10, 0, 0, 0, 1, 1, false, animation.ErrInvalidParameter},
		{"NegativeFrameHeight", sourceImg, 10, -5, 0, 0, 1, 1, false, animation.ErrInvalidParameter},
		{"ZeroFrameAmount", sourceImg, 10, 10, 0, 0, 0, 1, false, animation.ErrInvalidParameter},
		{"NegativeFrameAmount", sourceImg, 10, 10, 0, 0, -2, 1, false, animation.ErrInvalidParameter},
		{"ZeroAnimationSpeed", sourceImg, 10, 10, 0, 0, 1, 0, false, animation.ErrInvalidParameter},
		{"NegativeAnimationSpeed", sourceImg, 10, 10, 0, 0, 1, -3, false, animation.ErrInvalidParameter},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			anim, err := animation.NewAnimation(
				tc.sourceImage, tc.frameWidth, tc.frameHeight,
				tc.sourceStartX, tc.sourceStartY, tc.frameAmount,
				tc.animationSpeed, tc.loops,
			)
			if anim != nil {
				t.Error("Expected animation to be nil on error")
			}
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error to wrap %v, got %v (full error: %s)", tc.expectedError, errors.Unwrap(err), err)
			}

			// Check the specific error message format for ErrInvalidParameter
			if tc.expectedError == animation.ErrInvalidParameter {
				var expectedMsg string
				var invalidValue int
				var paramName string

				if tc.frameWidth <= 0 {
					invalidValue = tc.frameWidth
					paramName = "frameWidth must be positive"
				} else if tc.frameHeight <= 0 {
					invalidValue = tc.frameHeight
					paramName = "frameHeight must be positive"
				} else if tc.frameAmount <= 0 {
					invalidValue = tc.frameAmount
					paramName = "frameAmount must be positive"
				} else if tc.animationSpeed <= 0 {
					invalidValue = tc.animationSpeed
					paramName = "animationSpeed must be positive"
				}
				expectedMsg = fmt.Sprintf("%s: %s, got %d", animation.ErrInvalidParameter.Error(), paramName, invalidValue)
				if err == nil || err.Error() != expectedMsg {
					t.Errorf("Expected error message '%s', got '%v'", expectedMsg, err)
				}
			} else if tc.expectedError == animation.ErrNilSourceImage {
				if err == nil || err.Error() != animation.ErrNilSourceImage.Error() {
					t.Errorf("Expected error message '%s', got '%v'", animation.ErrNilSourceImage.Error(), err)
				}
			}
		})
	}
}

func TestAnimation_Update(t *testing.T) {
	sourceImg := newTestSourceImage(100, 20) // Suitable for 5 frames of 20x20

	t.Run("FrameProgression", func(t *testing.T) {
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 5, 3, true) // speed 3

		// Tick 1
		finished := anim.Update()
		if finished || anim.CurrentFrame != 0 || anim.FrameTimer != 1 {
			t.Errorf("Tick 1: finished=%v, CurrentFrame=%d (exp 0), FrameTimer=%d (exp 1)", finished, anim.CurrentFrame, anim.FrameTimer)
		}

		// Tick 2
		finished = anim.Update()
		if finished || anim.CurrentFrame != 0 || anim.FrameTimer != 2 {
			t.Errorf("Tick 2: finished=%v, CurrentFrame=%d (exp 0), FrameTimer=%d (exp 2)", finished, anim.CurrentFrame, anim.FrameTimer)
		}

		// Tick 3 (frame changes here)
		finished = anim.Update()
		if finished || anim.CurrentFrame != 1 || anim.FrameTimer != 0 { // `finished` should be false as it's a looping anim not at end of cycle
			t.Errorf("Tick 3: finished=%v, CurrentFrame=%d (exp 1), FrameTimer=%d (exp 0)", finished, anim.CurrentFrame, anim.FrameTimer)
		}
	})

	t.Run("LoopingAnimation", func(t *testing.T) {
		frameAmount := 3
		animationSpeed := 2
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, frameAmount, animationSpeed, true)

		totalTicksForLoop := frameAmount * animationSpeed
		loopedCorrectly := false
		for i := 0; i < totalTicksForLoop; i++ {
			finishedThisTick := anim.Update()
			isLastTickOfLoop := (i == totalTicksForLoop-1)

			if isLastTickOfLoop {
				if !finishedThisTick {
					t.Errorf("Tick %d (loop point): Expected Update to return true", i)
				}
				if anim.CurrentFrame != 0 {
					t.Errorf("Tick %d (loop point): Expected CurrentFrame to be 0, got %d", i, anim.CurrentFrame)
				}
				loopedCorrectly = true
			} else {
				if finishedThisTick {
					t.Errorf("Tick %d: Update returned true prematurely", i)
				}
			}
		}
		if !loopedCorrectly {
			t.Error("Looping animation did not report loop completion correctly or at the right time.")
		}
		if anim.CurrentFrame != 0 { // Verify final state after loop
			t.Errorf("After loop, expected CurrentFrame 0, got %d", anim.CurrentFrame)
		}
	})

	t.Run("NonLoopingAnimation_Behavior", func(t *testing.T) {
		frameAmount := 3
		animationSpeed := 2
		// Enough image width for 3 frames of 20px: 3*20 = 60px.
		nonLoopSourceImg := newTestSourceImage(60, 20)
		anim, _ := animation.NewAnimation(nonLoopSourceImg, 20, 20, 0, 0, frameAmount, animationSpeed, false)

		// Ticks to reach the start of the last frame's display
		// (frameAmount - 1) * animationSpeed ticks pass, then CurrentFrame becomes frameAmount-1.
		// Example: 3 frames, speed 2. (3-1)*2 = 4 ticks.
		// Tick 0: CF=0, FT=1
		// Tick 1: CF=0, FT=2 -> CF=1, FT=0
		// Tick 2: CF=1, FT=1
		// Tick 3: CF=1, FT=2 -> CF=2, FT=0. CurrentFrame is now 2 (frameAmount-1).
		ticksToStartLastFrameDisplay := (frameAmount - 1) * animationSpeed
		for i := 0; i < ticksToStartLastFrameDisplay; i++ {
			finished := anim.Update()
			if finished || anim.IsFinished() {
				t.Errorf("Tick %d: Animation reported finished prematurely (finished=%v, IsFinished()=%v)", i, finished, anim.IsFinished())
			}
		}

		if anim.CurrentFrame != frameAmount-1 {
			t.Errorf("After %d ticks, expected CurrentFrame %d (last frame), got %d", ticksToStartLastFrameDisplay, frameAmount-1, anim.CurrentFrame)
		}
		if anim.IsFinished() { // Should not be finished yet, just started displaying the last frame
			t.Errorf("After %d ticks, IsFinished was true prematurely (should only be true after last frame duration)", ticksToStartLastFrameDisplay)
		}

		// Now, run updates for the duration of the last frame.
		// It will be marked finished on the last tick of its display.
		finishedOnCorrectTick := false
		for i := 0; i < animationSpeed; i++ {
			currentOverallTick := ticksToStartLastFrameDisplay + i
			finished := anim.Update()

			isEndOfLastFrameDuration := (i == animationSpeed-1)
			if isEndOfLastFrameDuration {
				if !finished {
					t.Errorf("Tick %d (finish point): Expected Update to return true", currentOverallTick)
				}
				if !anim.IsFinished() {
					t.Errorf("Tick %d (finish point): Expected IsFinished to be true", currentOverallTick)
				}
				finishedOnCorrectTick = true
			} else { // During display of last frame, but not its end
				if finished || anim.IsFinished() {
					t.Errorf("Tick %d (during last frame): Animation reported finished prematurely (finished=%v, IsFinished()=%v)", currentOverallTick, finished, anim.IsFinished())
				}
			}
			if anim.CurrentFrame != frameAmount-1 { // Should stay on the last frame
				t.Errorf("Tick %d: CurrentFrame changed from last frame; expected %d, got %d", currentOverallTick, frameAmount-1, anim.CurrentFrame)
			}
		}
		if !finishedOnCorrectTick {
			t.Error("Animation did not report 'finished' on the correct tick for the last frame.")
		}

		// Behavior after finishing: should remain finished and on the last frame
		finishedAfterCompletion := anim.Update() // One more tick
		if !finishedAfterCompletion || !anim.IsFinished() || anim.CurrentFrame != frameAmount-1 {
			t.Errorf("Post-finish: state incorrect. Update returned %v (exp true), IsFinished()=%v (exp true), CF=%d (exp %d)",
				finishedAfterCompletion, anim.IsFinished(), anim.CurrentFrame, frameAmount-1)
		}
	})
}

func TestAnimation_GetImage(t *testing.T) {
	t.Run("ValidFrame", func(t *testing.T) {
		sourceImg := newTestSourceImage(60, 20) // Enough for 3 frames of 20x20
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, true)

		img0 := anim.GetImage()
		if img0 == nil || img0.Bounds().Dx() != 20 || img0.Bounds().Dy() != 20 {
			t.Fatalf("Frame 0: GetImage returned nil or incorrect bounds. Got: %p, Bounds: %v", img0, img0.Bounds())
		}

		anim.Update() // Advance to Frame 1 (speed 1)
		img1 := anim.GetImage()
		if img1 == nil || img1.Bounds().Dx() != 20 || img1.Bounds().Dy() != 20 {
			t.Fatalf("Frame 1: GetImage returned nil or incorrect bounds. Got: %p, Bounds: %v", img1, img1.Bounds())
		}
	})

	t.Run("FrameOutOfBounds_Width", func(t *testing.T) {
		// Source: 50x20. Frame: 20x20. Configured for 3 frames.
		// Frame 0 (idx 0): sx=0, rect=(0,0)-(20,20) -> OK
		// Frame 1 (idx 1): sx=20, rect=(20,0)-(40,20) -> OK
		// Frame 2 (idx 2): sx=40, rect=(40,0)-(60,20) -> Out of bounds (source width 50)
		sourceImg := newTestSourceImage(50, 20)
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, true)

		_ = anim.GetImage() // Frame 0
		anim.Update()
		_ = anim.GetImage() // Frame 1
		anim.Update()       // To Frame 2

		img := anim.GetImage() // Attempt to get Frame 2
		if img != nil {
			// Calculate expected frame rectangle for error message clarity
			expectedRect := image.Rect(
				anim.SourceStartX+anim.CurrentFrame*anim.FrameWidth,
				anim.SourceStartY,
				anim.SourceStartX+anim.CurrentFrame*anim.FrameWidth+anim.FrameWidth,
				anim.SourceStartY+anim.FrameHeight,
			)
			t.Errorf("Frame 2: Expected GetImage to return nil due to X out of bounds, but got an image. Expected FrameRect: %v, SourceBounds: %v",
				expectedRect, sourceImg.Bounds())
		}
	})

	t.Run("FrameOutOfBounds_Height", func(t *testing.T) {
		sourceImg := newTestSourceImage(20, 10) // Source height 10
		// Frame height 20. Frame 0: sy=0, rect=(0,0)-(20,20) -> Out of bounds (source height 10)
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 1, 1, true)

		img := anim.GetImage()
		if img != nil {
			expectedRect := image.Rect(anim.SourceStartX, anim.SourceStartY, anim.SourceStartX+anim.FrameWidth, anim.SourceStartY+anim.FrameHeight)
			t.Errorf("Frame 0: Expected GetImage to return nil due to Y out of bounds, but got an image. Expected FrameRect: %v, SourceBounds: %v",
				expectedRect, sourceImg.Bounds())
		}
	})

	t.Run("FrameOutOfBounds_SourceStartXY", func(t *testing.T) {
		sourceImg := newTestSourceImage(50, 50)

		// StartX 40, FrameWidth 20. Frame 0 x-coords: 40 to 60. Out of bounds for width 50.
		animX, _ := animation.NewAnimation(sourceImg, 20, 20, 40, 0, 1, 1, true)
		imgX := animX.GetImage()
		if imgX != nil {
			t.Errorf("SourceStart X: Expected GetImage nil due to out of bounds, got image.")
		}

		// StartY 40, FrameHeight 20. Frame 0 y-coords: 40 to 60. Out of bounds for height 50.
		animY, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 40, 1, 1, true)
		imgY := animY.GetImage()
		if imgY != nil {
			t.Errorf("SourceStart Y: Expected GetImage nil due to out of bounds, got image.")
		}
	})
}

func TestAnimation_Reset(t *testing.T) {
	sourceImg := newTestSourceImage(60, 20)
	// 3 frames, speed 2, no loop. Finishes after 3*2 = 6 updates.
	anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 2, false)

	// Advance animation to the end
	for i := 0; i < 6; i++ { // 6 updates to reach finished state
		anim.Update()
	}

	// Pre-reset checks: after 6 updates, animation is finished, on the last frame.
	// CurrentFrame = 2 (frameAmount-1), FrameTimer = 0 (reset when frame advanced to last and finished condition met), hasFinished = true.
	if anim.CurrentFrame != 2 || !anim.IsFinished() || anim.FrameTimer != 0 {
		t.Fatalf("Pre-reset state incorrect: CF=%d (exp 2), Finished=%v (exp true), FT=%d (exp 0)",
			anim.CurrentFrame, anim.IsFinished(), anim.FrameTimer)
	}

	// Call Update one more time; for a finished non-looping anim, FrameTimer should not change from 0.
	anim.Update()
	if anim.FrameTimer != 0 {
		// This confirms FrameTimer is not incremented if !Loops && hasFinished
		t.Fatalf("Pre-reset (after extra update on finished anim): Expected FrameTimer 0, got %d", anim.FrameTimer)
	}

	anim.Reset()

	// Post-reset checks
	if anim.CurrentFrame != 0 || anim.FrameTimer != 0 || anim.IsFinished() {
		t.Errorf("Post-reset state incorrect: CF=%d (exp 0), FT=%d (exp 0), Finished=%v (exp false)",
			anim.CurrentFrame, anim.FrameTimer, anim.IsFinished())
	}
}

func TestAnimation_IsFinished(t *testing.T) {
	sourceImg := newTestSourceImage(60, 20) // Suitable for 3 frames of 20x20

	t.Run("NonLooping_Initial", func(t *testing.T) {
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, false)
		if anim.IsFinished() {
			t.Error("Non-looping: Expected IsFinished false initially")
		}
	})

	t.Run("NonLooping_BecomesFinished", func(t *testing.T) {
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, false) // 3 frames, speed 1
		// Takes 3 updates to finish (1 update per frame change, 1 for last frame duration)
		// Update 1: CF0 -> CF1 (FT0)
		// Update 2: CF1 -> CF2 (FT0)
		// Update 3: CF2 -> CF2 (FT0), finished=true
		for i := 0; i < 2; i++ {
			anim.Update()
			if anim.IsFinished() {
				t.Fatalf("Tick %d: Non-looping: Expected IsFinished false before completion", i)
			}
		}
		anim.Update() // This is the update that marks it as finished
		if !anim.IsFinished() {
			t.Error("Non-looping: Expected IsFinished true after completion")
		}
		anim.Update() // Another update, should remain finished
		if !anim.IsFinished() {
			t.Error("Non-looping: Expected IsFinished to remain true after subsequent updates")
		}
	})

	t.Run("Looping_AlwaysFalse", func(t *testing.T) {
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, true) // Looping
		if anim.IsFinished() {
			t.Error("Looping: Expected IsFinished false initially")
		}
		for i := 0; i < 5; i++ { // Multiple updates, including loops
			anim.Update()
			if anim.IsFinished() { // For a looping animation, IsFinished (hasFinished) should always be false
				t.Errorf("Tick %d: Looping: Expected IsFinished to remain false", i)
			}
		}
	})

	t.Run("NonLooping_ResetClearsFinished", func(t *testing.T) {
		anim, _ := animation.NewAnimation(sourceImg, 20, 20, 0, 0, 3, 1, false)
		for i := 0; i < 3; i++ { // Run 3 updates to ensure it finishes
			anim.Update()
		}
		if !anim.IsFinished() {
			t.Fatal("Non-looping: Animation should be finished before reset")
		}
		anim.Reset()
		if anim.IsFinished() {
			t.Error("Non-looping: Expected IsFinished false after reset")
		}
	})
}
