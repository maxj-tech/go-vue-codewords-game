package domain

type Teamfarbe int

const (
	TeamRot  Teamfarbe = Teamfarbe(Rot)
	TeamBlau Teamfarbe = Teamfarbe(Blau)
)

type Team struct {
	Farbe   Teamfarbe
	Spieler []Spieler
	Chef    Spieler // todo erstmal nur einen Chef je Team
}
