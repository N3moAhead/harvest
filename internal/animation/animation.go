package animation

import (
	"errors"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Package animation provides helpers for managing frame-based animations
// from sprite sheets
//
// It focuses on tracking the current frame of an animation based on time
// and providing the corresponding sub-image from a source spritesheet.
// It does *not* handle drawing the animation to the screen; that responsibility
// remains with the calling code (e.g., using ebiten.DrawImageOptions).
//

var (
	// ErrInvalidParameter is returned when essential animation parameters are invalid (e.g., zero size).
	ErrInvalidParameter = errors.New("animation: invalid parameter provided")
	// ErrNilSourceImage is returned when a nil source image is provided.
	ErrNilSourceImage = errors.New("animation: source image cannot be nil")
)

// Animation represents a single animation sequence from a spritesheet.
type Animation struct {
	// AnimationSpeed defines how many game ticks (Update calls) each frame is displayed.
	// A lower value means a faster animation. Must be greater than 0.
	AnimationSpeed int
	// FrameAmount is the total number of frames in this animation sequence. Must be greater than 0.
	FrameAmount int
	// CurrentFrame is the index (0-based) of the currently active frame.
	CurrentFrame int
	// FrameTimer counts game ticks towards the next frame change.
	FrameTimer int
	// FrameWidth is the width of a single frame in pixels. Must be greater than 0.
	FrameWidth int
	// FrameHeight is the height of a single frame in pixels. Must be greater than 0.
	FrameHeight int
	// SourceStartX is the X coordinate of the top-left corner of the *first* frame
	// of this animation sequence within the SourceImage.
	SourceStartX int
	// SourceStartY is the Y coordinate of the top-left corner of the *first* frame
	// of this animation sequence within the SourceImage.
	// Assumes all frames of this sequence are on the same row.
	SourceStartY int
	// SourceImage is the spritesheet containing all the frames for this animation.
	SourceImage *ebiten.Image
	// Loops determines if the animation should restart from the beginning after reaching the end.
	Loops bool
	// hasFinished is an internal flag primarily for non-looping animations.
	hasFinished bool
}

// NewAnimation creates and validates a new Animation instance.
func NewAnimation(
	sourceImage *ebiten.Image,
	frameWidth, frameHeight int,
	sourceStartX, sourceStartY int,
	frameAmount int,
	animationSpeed int,
	loops bool,
) (*Animation, error) {
	// --- Validation ---
	if sourceImage == nil {
		return nil, ErrNilSourceImage
	}
	if frameWidth <= 0 {
		return nil, fmt.Errorf("%w: frameWidth must be positive, got %d", ErrInvalidParameter, frameWidth)
	}
	if frameHeight <= 0 {
		return nil, fmt.Errorf("%w: frameHeight must be positive, got %d", ErrInvalidParameter, frameHeight)
	}
	if frameAmount <= 0 {
		return nil, fmt.Errorf("%w: frameAmount must be positive, got %d", ErrInvalidParameter, frameAmount)
	}
	if animationSpeed <= 0 {
		return nil, fmt.Errorf("%w: animationSpeed must be positive, got %d", ErrInvalidParameter, animationSpeed)
	}

	return &Animation{
		SourceImage:    sourceImage,
		FrameWidth:     frameWidth,
		FrameHeight:    frameHeight,
		SourceStartX:   sourceStartX,
		SourceStartY:   sourceStartY,
		FrameAmount:    frameAmount,
		AnimationSpeed: animationSpeed,
		Loops:          loops,
		CurrentFrame:   0, // Start at the first frame
		FrameTimer:     0,
		hasFinished:    false,
	}, nil
}

// Update progresses the animation timer and frame.
// It should be called once per game tick (e.g., in your Ebitengine game's Update method).
// Returns true if the animation completed a full cycle *during this update*.
// For non-looping animations, it will return true once when the last frame is reached
// and continue returning true on subsequent calls without updating the frame further.
// For looping animations, it returns true only on the tick when it loops back to the first frame.
func (a *Animation) Update() bool {
	// If the animation doesn't loop and has already finished, do nothing.
	if !a.Loops && a.hasFinished {
		return true // Indicate it's (still) finished.
	}

	animationFinishedThisTick := false
	a.FrameTimer++

	// Check if it's time to advance to the next frame
	if a.FrameTimer >= a.AnimationSpeed {
		a.FrameTimer = 0 // Reset timer

		// Check if the current frame is the last one
		isLastFrame := a.CurrentFrame == a.FrameAmount-1

		if isLastFrame {
			if a.Loops {
				// Loop back to the first frame
				a.CurrentFrame = 0
			} else {
				// Stay on the last frame and mark as finished
				a.hasFinished = true
			}
			animationFinishedThisTick = true // Completed a cycle
		} else {
			// Advance to the next frame
			a.CurrentFrame++
		}
	}

	return animationFinishedThisTick
}

// GetImage returns the sub-image corresponding to the CurrentFrame of the animation.
// It calculates the correct rectangle within the SourceImage based on the animation's
// properties and the CurrentFrame.
// Returns nil if the calculated frame rectangle is outside the bounds of the SourceImage.
func (a *Animation) GetImage() *ebiten.Image {
	// Calculate the X position of the current frame on the spritesheet.
	// Assumes frames are laid out horizontally.
	sx := a.SourceStartX + a.CurrentFrame*a.FrameWidth
	// Y position is constant for this animation sequence.
	sy := a.SourceStartY

	// Define the rectangle for the current frame within the source image.
	frameRect := image.Rect(sx, sy, sx+a.FrameWidth, sy+a.FrameHeight)

	// --- Bounds Check ---
	// Ensure the calculated rectangle is fully within the source image bounds.
	bounds := a.SourceImage.Bounds()
	if !frameRect.In(bounds) {
		// This indicates an issue, potentially with animation setup or source image size.
		fmt.Printf("Error: Frame rectangle %v is out of source image bounds %v\n", frameRect, bounds)
		// Returning nil prevents a panic from SubImage but signals an error state.
		// Consider more robust error handling depending on application needs (e.g., return error).
		return nil
	}

	// Extract the sub-image for the current frame.
	// The result of SubImage is image.Image, so we need to assert it back to *ebiten.Image.
	frameImage := a.SourceImage.SubImage(frameRect).(*ebiten.Image)
	return frameImage
}

// Reset sets the animation back to its first frame and resets the timer.
// Also resets the 'finished' state for non-looping animations.
func (a *Animation) Reset() {
	a.CurrentFrame = 0
	a.FrameTimer = 0
	a.hasFinished = false
}

// IsFinished returns true if a non-looping animation has reached its end.
// For looping animations, this will always return false unless manually set.
func (a *Animation) IsFinished() bool {
	return a.hasFinished
}
