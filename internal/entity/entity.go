package entity

import (
	"fmt"

	"github.com/N3moAhead/harvest/internal/component"
	"github.com/google/uuid"
)

var id int

// The root class of all entities in the game.
// It can help us that every entity in the game
// has a position with the same keys and each
// Entity has a unique id
type Entity struct {
	ID  uuid.UUID
	Pos component.Vector2D
}

func NewEntity(posX, posY float64) *Entity {
	newID, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("Error while generating a new uuid: %v", err))
	}

	return &Entity{
		ID:  newID,
		Pos: component.NewVector2D(posX, posY),
	}
}

func (e *Entity) GetId() uuid.UUID {
	return e.ID
}

func (e *Entity) GetIdString() string {
	return e.ID.String()
}

func (e *Entity) GetPosition() component.Vector2D {
	return e.Pos
}
