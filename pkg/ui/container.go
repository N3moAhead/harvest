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

func (c *Container) getNextPosition(numExistingChildren int) (float64, float64) {
	offsetX, offsetY := 0.0, 0.0

	totalGapOffset := float64(numExistingChildren) * c.gap

	if c.direction == Row {
		currentTotalWidth := 0.0
		for i, existingChild := range c.Children {
			// If set position is called on an container we
			// can only count some not all elements
			if i >= numExistingChildren {
				break
			}
			width, _ := existingChild.GetSize()
			currentTotalWidth += width
		}
		offsetX = currentTotalWidth + totalGapOffset
	} else {
		currentTotalHeight := 0.0
		for i, existingChild := range c.Children {
			if i >= numExistingChildren {
				break
			}
			_, height := existingChild.GetSize()
			currentTotalHeight += height
		}
		offsetY = currentTotalHeight + totalGapOffset
	}

	posX := c.X + offsetX
	posY := c.Y + offsetY

	return posX, posY
}

// AddChild places the children based on the config
// beside or below each other, with a defined gap between them
func (c *Container) AddChild(newChild UIElement) {
	numExistingChildren := len(c.Children)
	posX, posY := c.getNextPosition(numExistingChildren)
	newChild.SetPosition(posX, posY)
	c.BaseElement.AddChild(newChild)
}

func (c *Container) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
	for i, child := range c.Children {
		posX, posY := c.getNextPosition(i)
		child.SetPosition(posX, posY)
	}
}

var _ UIElement = (*Container)(nil)
