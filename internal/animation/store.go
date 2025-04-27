package animation

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// AnimationStore manages a collection of named Animation instances.
// This is useful for entities that can have multiple animations (e.g., "idle", "run", "jump").
type AnimationStore struct {
	// Animations maps a unique string name to its corresponding Animation object.
	Animations map[string]*Animation
	// currentAnimationName holds the key of the currently active animation within the store.
	currentAnimationName string
}

// NewAnimationStore creates a new, empty AnimationStore.
func NewAnimationStore() *AnimationStore {
	return &AnimationStore{
		Animations: make(map[string]*Animation),
		// currentAnimationName is initially empty, indicating no animation is active.
	}
}

// AddAnimation adds an Animation to the store with a given name.
// If an animation with the same name already exists, it will be overwritten.
// It's recommended to check for nil `anim` before calling.
func (as *AnimationStore) AddAnimation(name string, anim *Animation) error {
	if name == "" {
		return fmt.Errorf("%w: animation name cannot be empty", ErrInvalidParameter)
	}
	if anim == nil {
		return fmt.Errorf("%w: cannot add a nil animation with name '%s'", ErrInvalidParameter, name)
	}
	as.Animations[name] = anim
	return nil
}

// SetCurrentAnimation sets the currently active animation by name.
// If the name does not exist in the store, it returns false.
// If the new animation is different from the previous one, it resets the new animation.
// Returns true if the animation was successfully set (i.e., the name exists).
func (as *AnimationStore) SetCurrentAnimation(name string) bool {
	newAnim, ok := as.Animations[name]
	if !ok {
		// Animation name not found in the store.
		// Keep the current animation active (if any).
		return false
	}

	if as.currentAnimationName != name {
		// If switching to a different animation, reset it to start from the beginning.
		newAnim.Reset()
		as.currentAnimationName = name
	}
	return true
}

// GetCurrentAnimation returns the currently active Animation instance.
// Returns nil if no animation has been set using SetCurrentAnimation or
// if the stored currentAnimationName is no longer valid (e.g., animation removed).
func (as *AnimationStore) GetCurrentAnimation() *Animation {
	if as.currentAnimationName == "" {
		return nil // No animation set as current
	}
	anim, ok := as.Animations[as.currentAnimationName]
	if !ok {
		// This shouldn't happen if SetCurrentAnimation is used correctly,
		// but handles the case where the animation might be removed after being set.
		fmt.Printf("Warning: Current animation '%s' not found in store.\n", as.currentAnimationName)
		as.currentAnimationName = "" // Clear invalid name
		return nil
	}
	return anim
}

// GetAnimation retrieves a specific animation by name, regardless of whether it's the current one.
// Returns the animation and true if found, otherwise returns nil and false.
func (as *AnimationStore) GetAnimation(name string) (*Animation, bool) {
	anim, ok := as.Animations[name]
	return anim, ok
}

// Update updates the currently active animation in the store.
// It's a convenience method; you could also get the current animation
// using GetCurrentAnimation() and call Update() on it directly.
// Does nothing if no current animation is set.
// Returns the result of the underlying Animation's Update call (true if it finished a cycle).
func (as *AnimationStore) Update() bool {
	currentAnim := as.GetCurrentAnimation()
	if currentAnim != nil {
		return currentAnim.Update()
	}
	return false // No active animation to update
}

// GetImage returns the current frame image of the currently active animation.
// It's a convenience method. Returns nil if no animation is active or if
// the active animation's GetImage() returns nil (e.g., due to bounds issues).
func (as *AnimationStore) GetImage() *ebiten.Image {
	currentAnim := as.GetCurrentAnimation()
	if currentAnim != nil {
		return currentAnim.GetImage()
	}
	return nil // No active animation
}

func (as *AnimationStore) GetCurrentAnimationName() string {
	return as.currentAnimationName
}
