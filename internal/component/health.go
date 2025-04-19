package component

type Health struct {
	HP    float64
	MaxHP float64
}

func (h *Health) Damage(amount float64) (alive bool) {
	if amount > h.HP {
		h.HP = 0
		return false
	}
	h.HP -= amount
	return true
}

func (h *Health) Heal(amount float64) {
	if h.HP+amount > h.MaxHP {
		h.HP = h.MaxHP
	} else {
		h.HP += amount
	}
}

func NewHealth(maxHP float64) Health {
	return Health{
		HP:    maxHP,
		MaxHP: maxHP,
	}
}
