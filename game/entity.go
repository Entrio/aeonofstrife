package game

type entity interface {
	getName() string
	getID() int64
	isPlayer() bool
	isInteractable() bool
	ispassable() bool
	getRoomID() int64
	getCurrentHP() int
	getMaxHP() int
}
