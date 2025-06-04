package animation_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/N3moAhead/harvest/internal/animation"
)

// newTestAnimation is a helper to create a valid *animation.Animation for testing.
// It uses a source image large enough for the specified frameAmount.
func newTestAnimation(frameAmount int, speed int, loops bool, frameWidth, frameHeight int) (*animation.Animation, error) {
	if frameAmount <= 0 {
		frameAmount = 1 // Default to 1 frame if invalid
	}
	if frameWidth <= 0 {
		frameWidth = 10 // Default width
	}
	if frameHeight <= 0 {
		frameHeight = 10 // Default height
	}
	// Ensure source image is large enough
	sourceImg := newTestSourceImage(frameAmount*frameWidth, frameHeight)
	return animation.NewAnimation(sourceImg, frameWidth, frameHeight, 0, 0, frameAmount, speed, loops)
}

func TestNewAnimationStore(t *testing.T) {
	store := animation.NewAnimationStore()
	if store == nil {
		t.Fatal("NewAnimationStore returned nil")
	}
	if store.Animations == nil {
		t.Error("Expected Animations map to be initialized, got nil")
	}
	if len(store.Animations) != 0 {
		t.Errorf("Expected new store to have 0 animations, got %d", len(store.Animations))
	}
	if store.GetCurrentAnimationName() != "" {
		t.Errorf("Expected currentAnimationName to be empty, got '%s'", store.GetCurrentAnimationName())
	}
}

func TestAnimationStore_AddAnimation(t *testing.T) {
	store := animation.NewAnimationStore()
	anim1, _ := newTestAnimation(1, 1, false, 10, 10)

	t.Run("AddValidAnimation", func(t *testing.T) {
		err := store.AddAnimation("idle", anim1)
		if err != nil {
			t.Fatalf("Failed to add valid animation: %v", err)
		}
		if _, ok := store.Animations["idle"]; !ok {
			t.Error("Animation 'idle' not found in store after adding")
		}
	})

	t.Run("OverwriteAnimation", func(t *testing.T) {
		anim2, _ := newTestAnimation(2, 1, false, 10, 10)
		err := store.AddAnimation("idle", anim2) // Overwrite
		if err != nil {
			t.Fatalf("Failed to overwrite animation: %v", err)
		}
		if store.Animations["idle"] != anim2 {
			t.Error("Animation 'idle' was not overwritten with the new instance")
		}
	})

	t.Run("AddWithEmptyName", func(t *testing.T) {
		err := store.AddAnimation("", anim1)
		if !errors.Is(err, animation.ErrInvalidParameter) {
			t.Errorf("Expected ErrInvalidParameter for empty name, got %v", err)
		}
		expectedMsg := fmt.Sprintf("%s: animation name cannot be empty", animation.ErrInvalidParameter.Error())
		if err == nil || err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%v'", expectedMsg, err)
		}
	})

	t.Run("AddNilAnimation", func(t *testing.T) {
		err := store.AddAnimation("fail", nil)
		if !errors.Is(err, animation.ErrInvalidParameter) {
			t.Errorf("Expected ErrInvalidParameter for nil animation, got %v", err)
		}
		expectedMsg := fmt.Sprintf("%s: cannot add a nil animation with name '%s'", animation.ErrInvalidParameter.Error(), "fail")
		if err == nil || err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%v'", expectedMsg, err)
		}
	})
}

func TestAnimationStore_SetCurrentAnimation(t *testing.T) {
	store := animation.NewAnimationStore()
	animIdle, _ := newTestAnimation(3, 2, true, 10, 10)
	animRun, _ := newTestAnimation(4, 1, true, 10, 10)

	store.AddAnimation("idle", animIdle)
	store.AddAnimation("run", animRun)

	t.Run("SetValidAnimation", func(t *testing.T) {
		ok := store.SetCurrentAnimation("idle")
		if !ok {
			t.Error("SetCurrentAnimation returned false for existing animation 'idle'")
		}
		if store.GetCurrentAnimationName() != "idle" {
			t.Errorf("Expected current animation name 'idle', got '%s'", store.GetCurrentAnimationName())
		}
		if animIdle.CurrentFrame != 0 || animIdle.FrameTimer != 0 {
			t.Error("Newly set animation 'idle' was not reset")
		}
	})

	t.Run("SetNonExistentAnimation", func(t *testing.T) {
		store.SetCurrentAnimation("idle") // Ensure a current anim is set
		initialCurrentName := store.GetCurrentAnimationName()

		ok := store.SetCurrentAnimation("jump") // "jump" doesn't exist
		if ok {
			t.Error("SetCurrentAnimation returned true for non-existent animation 'jump'")
		}
		if store.GetCurrentAnimationName() != initialCurrentName {
			t.Errorf("Expected current animation name to remain '%s', got '%s'", initialCurrentName, store.GetCurrentAnimationName())
		}
	})

	t.Run("SwitchAnimationResetsNewAnimation", func(t *testing.T) {
		// Set and advance "idle"
		store.SetCurrentAnimation("idle")
		animIdle.Update()         // CF=0, FT=1
		animIdle.Update()         // CF=0, FT=0 (advances to frame 1 because speed is 2) -> Oh, speed 2, so FT=0, CF=1
		animIdle.CurrentFrame = 1 // Manually set for clarity, simulating progress
		animIdle.FrameTimer = 1

		// Switch to "run"
		ok := store.SetCurrentAnimation("run")
		if !ok || store.GetCurrentAnimationName() != "run" {
			t.Fatal("Failed to switch to 'run' animation")
		}
		if animRun.CurrentFrame != 0 || animRun.FrameTimer != 0 {
			t.Error("Animation 'run' was not reset upon being set as current")
		}

		// Check "idle" state (should be unchanged as it wasn't active)
		if animIdle.CurrentFrame != 1 || animIdle.FrameTimer != 1 {
			t.Error("Previously active animation 'idle' state changed unexpectedly")
		}

		// Switch back to "idle"
		animRun.Update() // Advance run anim
		animRun.CurrentFrame = 1
		animRun.FrameTimer = 0

		ok = store.SetCurrentAnimation("idle")
		if !ok || store.GetCurrentAnimationName() != "idle" {
			t.Fatal("Failed to switch back to 'idle' animation")
		}
		if animIdle.CurrentFrame != 0 || animIdle.FrameTimer != 0 { // "idle" should be reset now
			t.Errorf("Animation 'idle' was not reset upon being set back as current. CF=%d, FT=%d", animIdle.CurrentFrame, animIdle.FrameTimer)
		}
	})

	t.Run("SetSameAnimationDoesNotReset", func(t *testing.T) {
		store.SetCurrentAnimation("idle") // Resets it
		animIdle.Update()                 // CF=0, FT=1
		animIdle.Update()                 // CF=1, FT=0 (due to speed 2)

		currentFrameBefore := animIdle.CurrentFrame
		frameTimerBefore := animIdle.FrameTimer

		ok := store.SetCurrentAnimation("idle") // Set same again
		if !ok {
			t.Error("SetCurrentAnimation failed when setting the same animation")
		}
		if animIdle.CurrentFrame != currentFrameBefore || animIdle.FrameTimer != frameTimerBefore {
			t.Errorf("Setting the same animation 'idle' incorrectly reset it. CF: %d->%d, FT: %d->%d",
				currentFrameBefore, animIdle.CurrentFrame, frameTimerBefore, animIdle.FrameTimer)
		}
	})
}

func TestAnimationStore_GetCurrentAnimation(t *testing.T) {
	store := animation.NewAnimationStore()
	animIdle, _ := newTestAnimation(1, 1, false, 10, 10)
	store.AddAnimation("idle", animIdle)

	t.Run("NoCurrentAnimationInitially", func(t *testing.T) {
		if current := store.GetCurrentAnimation(); current != nil {
			t.Errorf("Expected GetCurrentAnimation to return nil initially, got %v", current)
		}
	})

	t.Run("GetValidCurrentAnimation", func(t *testing.T) {
		store.SetCurrentAnimation("idle")
		current := store.GetCurrentAnimation()
		if current != animIdle {
			t.Errorf("Expected GetCurrentAnimation to return 'idle' animation, got %v", current)
		}
	})

	t.Run("CurrentAnimationRemovedFromMap", func(t *testing.T) {
		store.SetCurrentAnimation("idle")
		// Simulate removal after setting (e.g., by another part of code or error)
		delete(store.Animations, "idle")

		current := store.GetCurrentAnimation() // This should detect the inconsistency
		if current != nil {
			t.Errorf("Expected GetCurrentAnimation to return nil after removal, got %v", current)
		}
		if store.GetCurrentAnimationName() != "" {
			t.Errorf("Expected currentAnimationName to be cleared after inconsistency, got '%s'", store.GetCurrentAnimationName())
		}
		// Re-add for subsequent tests if needed, or use a fresh store.
		store.AddAnimation("idle", animIdle)
	})
}

func TestAnimationStore_GetAnimation(t *testing.T) {
	store := animation.NewAnimationStore()
	anim1, _ := newTestAnimation(1, 1, false, 10, 10)
	store.AddAnimation("testAnim", anim1)

	t.Run("GetExistingAnimation", func(t *testing.T) {
		anim, ok := store.GetAnimation("testAnim")
		if !ok {
			t.Error("GetAnimation returned false for existing animation 'testAnim'")
		}
		if anim != anim1 {
			t.Error("GetAnimation returned incorrect animation instance for 'testAnim'")
		}
	})

	t.Run("GetNonExistentAnimation", func(t *testing.T) {
		anim, ok := store.GetAnimation("notFound")
		if ok {
			t.Error("GetAnimation returned true for non-existent animation 'notFound'")
		}
		if anim != nil {
			t.Error("GetAnimation returned non-nil animation instance for 'notFound'")
		}
	})
}

func TestAnimationStore_Update(t *testing.T) {
	store := animation.NewAnimationStore()
	// Looping animation: 2 frames, speed 1. Will loop every 2 updates.
	animLoop, _ := newTestAnimation(2, 1, true, 10, 10)
	store.AddAnimation("loop", animLoop)

	t.Run("UpdateWithNoCurrentAnimation", func(t *testing.T) {
		finished := store.Update()
		if finished {
			t.Error("Update returned true when no current animation was set")
		}
	})

	t.Run("UpdateCurrentAnimation", func(t *testing.T) {
		store.SetCurrentAnimation("loop")

		// Tick 1: animLoop CF=0, FT=0 -> anim.Update() -> CF=1, FT=0, returns false
		finished1 := store.Update()
		if finished1 {
			t.Error("Update 1: Store.Update returned true prematurely")
		}
		if animLoop.CurrentFrame != 1 || animLoop.FrameTimer != 0 {
			t.Errorf("Update 1: animLoop state incorrect. CF=%d (exp 1), FT=%d (exp 0)", animLoop.CurrentFrame, animLoop.FrameTimer)
		}

		// Tick 2: animLoop CF=1, FT=0 -> anim.Update() -> CF=0, FT=0, returns true (looped)
		finished2 := store.Update()
		if !finished2 {
			t.Error("Update 2: Store.Update returned false at loop point")
		}
		if animLoop.CurrentFrame != 0 || animLoop.FrameTimer != 0 {
			t.Errorf("Update 2: animLoop state incorrect after loop. CF=%d (exp 0), FT=%d (exp 0)", animLoop.CurrentFrame, animLoop.FrameTimer)
		}
	})
}

func TestAnimationStore_GetImage(t *testing.T) {
	store := animation.NewAnimationStore()
	// Animation with 2 frames, 10x10 each. Source image is 20x10.
	animValid, _ := newTestAnimation(2, 1, true, 10, 10)
	store.AddAnimation("valid", animValid)

	// Animation that will cause GetImage to return nil (frame out of bounds)
	sourceImgSmall := newTestSourceImage(5, 5) // Too small for 10x10 frame
	animOutOfBounds, _ := animation.NewAnimation(sourceImgSmall, 10, 10, 0, 0, 1, 1, false)
	store.AddAnimation("outOfBounds", animOutOfBounds)

	t.Run("GetImageWithNoCurrentAnimation", func(t *testing.T) {
		img := store.GetImage()
		if img != nil {
			t.Error("GetImage returned non-nil when no current animation was set")
		}
	})

	t.Run("GetImageFromValidCurrentAnimation", func(t *testing.T) {
		store.SetCurrentAnimation("valid")
		img := store.GetImage()
		if img == nil {
			t.Fatal("GetImage returned nil for a valid current animation")
		}
		if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
			t.Errorf("GetImage returned image with incorrect dimensions: %v", img.Bounds())
		}

		// Advance frame and check again
		store.Update() // animValid CF=1
		imgFrame2 := store.GetImage()
		if imgFrame2 == nil {
			t.Fatal("GetImage returned nil for frame 2 of valid animation")
		}
		// Note: SubImage bounds are relative to the original image, but Dx/Dy should be frame size
		if imgFrame2.Bounds().Dx() != 10 || imgFrame2.Bounds().Dy() != 10 {
			t.Errorf("GetImage (frame 2) returned image with incorrect dimensions: %v", imgFrame2.Bounds())
		}
	})

	t.Run("GetImageWhenAnimationReturnsNil", func(t *testing.T) {
		store.SetCurrentAnimation("outOfBounds")
		img := store.GetImage()
		if img != nil {
			t.Error("GetImage returned non-nil when underlying animation's GetImage should return nil")
		}
	})
}

func TestAnimationStore_GetCurrentAnimationName(t *testing.T) {
	store := animation.NewAnimationStore()
	anim1, _ := newTestAnimation(1, 1, false, 10, 10)
	store.AddAnimation("anim1", anim1)

	t.Run("InitialNameIsEmpty", func(t *testing.T) {
		if name := store.GetCurrentAnimationName(); name != "" {
			t.Errorf("Expected initial current animation name to be empty, got '%s'", name)
		}
	})

	t.Run("NameAfterSet", func(t *testing.T) {
		store.SetCurrentAnimation("anim1")
		if name := store.GetCurrentAnimationName(); name != "anim1" {
			t.Errorf("Expected current animation name 'anim1', got '%s'", name)
		}
	})

	t.Run("NameAfterFailedSet", func(t *testing.T) {
		store.SetCurrentAnimation("anim1") // Ensure it's "anim1"
		store.SetCurrentAnimation("nonExistent")
		if name := store.GetCurrentAnimationName(); name != "anim1" {
			t.Errorf("Expected current animation name to remain 'anim1' after failed set, got '%s'", name)
		}
	})
}
