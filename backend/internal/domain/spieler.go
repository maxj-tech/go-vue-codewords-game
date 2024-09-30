package domain

type Rolle int

const (
	Chef Rolle = iota
	Ermittler
)

type Spieler struct {
	Name  string
	Rolle Rolle
}
