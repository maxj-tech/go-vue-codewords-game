package domain

import (
	"testing"
)

func TestNewTeam(t *testing.T) {
	const teamFarbe = TeamBlau
	chef, _ := NewSpieler("Chef", Chef)
	ermittler, _ := NewSpieler("Ermittler", Ermittler)
	team, err := NewTeam(teamFarbe, chef, []Spieler{ermittler})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if team.Farbe != teamFarbe {
		t.Fatalf("expected %v, got %v", teamFarbe, team.Farbe)
	}

	if team.Chef.ID() != chef.ID() {
		t.Fatalf("expected chef.id %v, got %v", chef.ID(), team.Chef.ID())
	}

	if storedErmittler, exists := team.Ermittler[ermittler.ID()]; !exists || storedErmittler.ID() != ermittler.ID() {
		t.Fatalf("expected Ermittler with ID %v, got %v", ermittler.ID(), storedErmittler.ID())
	}

}

func TestNewTeamErrorErmittlerAlsChef(t *testing.T) {
	const teamFarbe = TeamRot
	ermittler, _ := NewSpieler("Ermittler", Ermittler)
	_, err := NewTeam(teamFarbe, ermittler, []Spieler{ermittler}) // Ermittler als Chef !
	if err == nil {
		t.Fatalf("expected an error, but got none")
	}
}

func TestNewTeamErrorKeinErmittler(t *testing.T) {
	const teamFarbe = TeamBlau
	chef, _ := NewSpieler("Chef", Chef)
	_, err := NewTeam(teamFarbe, chef, []Spieler{})
	if err == nil {
		t.Fatalf("expected an error, but got none")
	}
}

func TestNewTeamErrorDoppelterErmittler(t *testing.T) {
	const teamFarbe = TeamBlau
	chef, _ := NewSpieler("Chef", Chef)
	ermittler, _ := NewSpieler("Ermittler", Ermittler)
	_, err := NewTeam(teamFarbe, chef, []Spieler{ermittler, ermittler})
	if err == nil {
		t.Fatalf("expected an error, but got none")
	}
}

func TestNewTeamErrorUngueltigeFarbe(t *testing.T) {
	chef, _ := NewSpieler("Chef", Chef)
	ermittler, _ := NewSpieler("Ermittler", Ermittler)
	farben := []int{-1, 2, 42, 1337}
	for _, farbe := range farben {
		_, err := NewTeam(Teamfarbe(farbe), chef, []Spieler{ermittler})
		if err == nil {
			t.Fatalf("expected error for invalid Teamfarbe %v, but didn't get one", farbe)
		}
	}
}
