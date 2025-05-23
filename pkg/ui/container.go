package ui

type ContainerDirection string

const (
	Undefined ContainerDirection = "undefined"
	Row       ContainerDirection = "row"
	Col       ContainerDirection = "col"
)

type ContainerOptions struct {
	Direction ContainerDirection
	Gap       float64
}

type Container struct {
	BaseElement
	direction ContainerDirection
	gap       float64
	xOffset   float64
	yOffset   float64
}

func NewContainer(x, y float64, op *ContainerOptions) *Container {
	// Setting the default to Row
	if op.Direction == Undefined {
		op.Direction = Row
	}
	return &Container{
		// Width and height will be defined as zero
		BaseElement: *NewBaseElement(x, y, 0, 0),
		direction:   op.Direction,
		gap:         op.Gap,
	}
}

// AddChild places the children based on the config
// beside or below each other, with a defined gap between them
func (c *Container) AddChild(newChild UIElement) {
	offsetX, offsetY := 0.0, 0.0
	numExistingChildren := len(c.Children)

	totalGapOffset := float64(numExistingChildren) * c.gap

	if c.direction == Row {
		currentTotalWidth := 0.0
		for _, existingChild := range c.Children {
			width, _ := existingChild.GetSize()
			currentTotalWidth += width
		}
		offsetX = currentTotalWidth + totalGapOffset
	} else {
		currentTotalHeight := 0.0
		for _, existingChild := range c.Children {
			_, height := existingChild.GetSize()
			currentTotalHeight += height
		}
		offsetY = currentTotalHeight + totalGapOffset
	}

	newChild.SetPosition(c.X+offsetX, c.Y+offsetY)
	c.BaseElement.AddChild(newChild)
}

var _ UIElement = (*Container)(nil)
