package component

type Knockback struct {
	Dir      Vector2D
	Distance float64
	Flying   bool
}

const (
	dampingFactor      = 0.3
	stopFlyingDistance = 0.5
)

func (k *Knockback) Init(from *Vector2D, to *Vector2D, dist float64) {
	if !k.Flying {
		dir := to.Sub(*from).Normalize()
		k.Dir = dir
		k.Distance = dist
		k.Flying = true
	}
}

// Modifies the given position by the remaining knockback
func (k *Knockback) Update(pos *Vector2D) {
	if k.Flying {
		flyDistance := k.Distance * dampingFactor
		k.Distance -= flyDistance
		if flyDistance < stopFlyingDistance {
			k.Flying = false
		}
		*pos = pos.Add(k.Dir.Mul(flyDistance))
	}
}
