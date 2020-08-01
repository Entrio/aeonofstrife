package game

type Player struct {
	name        string
	id          int64
	currentRoom int64
	hp          int
	maxhp       int
}

func (p *Player) ispassable() bool {
	return true
}

func (p *Player) getName() string {
	return p.name
}

func (p *Player) getID() int64 {
	return p.id
}

func (p *Player) isPlayer() bool {
	return true
}

func (p *Player) isInteractable() bool {
	return true
}

func (p *Player) getRoomID() int64 {
	return p.currentRoom
}

func (p *Player) getCurrentHP() int {
	return p.hp
}

func (p *Player) getMaxHP() int {
	return p.maxhp
}
