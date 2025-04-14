package component

type Health struct {
	HP    int
	MaxHP int
}

func (h *Health) Damage(amount int) (alive bool) {
	if amount > h.HP {
		h.HP = 0
		return false
	}
	h.HP -= amount
	return true
}

func (h *Health) Heal(amount int) {
	if h.HP+amount > h.MaxHP {
		h.HP = h.MaxHP
	} else {
		h.HP += amount
	}
}

func NewHealth(maxHP int) Health {
	return Health{
		HP:    maxHP,
		MaxHP: maxHP,
	}
}
