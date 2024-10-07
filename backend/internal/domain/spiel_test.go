package domain

import (
	"testing"
)

func TestNewSpielHat25Karten(t *testing.T) {
	spiel := NewSpiel(TeamRot)

	if len(spiel.Karten) != 25 {
		t.Fatalf("expected 25 cards, got %d", len(spiel.Karten))
	}
}

func TestNewSpielKeineDoppeltenBegriffe(t *testing.T) {
	spiel := NewSpiel(TeamBlau)

	begriffMap := make(map[string]bool)
	for _, karte := range spiel.Karten {
		if begriffMap[karte.Begriff] {
			t.Fatalf("duplicate card found: %s", karte.Begriff)
		}
		begriffMap[karte.Begriff] = true
	}
}

func TestNewSpielFarbenVerteilung(t *testing.T) {
	tests := []struct {
		startendesTeam Teamfarbe
		expectedRote   int
		expectedBlaue  int
	}{
		{TeamRot, 9, 8},
		{TeamBlau, 8, 9},
	}

	for _, tt := range tests {
		spiel := NewSpiel(tt.startendesTeam)

		countColors := map[Kartenfarbe]int{
			KarteRot:     0,
			KarteBlau:    0,
			KarteBeige:   0,
			KarteSchwarz: 0,
		}

		for _, karte := range spiel.Karten {
			countColors[karte.Farbe]++
		}

		if countColors[KarteRot] != tt.expectedRote {
			t.Errorf("expected %d red cards for starting team %v, got %d", tt.expectedRote, tt.startendesTeam, countColors[KarteRot])
		}

		if countColors[KarteBlau] != tt.expectedBlaue {
			t.Errorf("expected %d blue cards for starting team %v, got %d", tt.expectedBlaue, tt.startendesTeam, countColors[KarteBlau])
		}

		if countColors[KarteBeige] != 7 {
			t.Errorf("expected 7 beige cards for starting team %v, got %d", tt.startendesTeam, countColors[KarteBeige])
		}

		if countColors[KarteSchwarz] != 1 {
			t.Errorf("expected 1 black card for starting team %v, got %d", tt.startendesTeam, countColors[KarteSchwarz])
		}

	}
}

func TestNewSpielAktuellerZugFuerStartendesTeam(t *testing.T) {
	var startendesTeam Teamfarbe = TeamBlau
	spiel := NewSpiel(startendesTeam)
	if spiel.AlsNaechstesAmZug != startendesTeam {
		t.Errorf("expected current turn to be with team %v but got team %v", startendesTeam, spiel.AlsNaechstesAmZug)
	}
}

func TestNewSpielStartendesTeamRichtigGesetzt(t *testing.T) {
	var startendesTeam Teamfarbe = TeamRot
	spiel := NewSpiel(startendesTeam)

	if spiel.AlsNaechstesAmZug != startendesTeam {
		t.Errorf("expected starting team to be %v but got %v", startendesTeam, spiel.AlsNaechstesAmZug)
	}
}

func TestNewSpielInitialerSpielstand(t *testing.T) {
	spiel := NewSpiel(TeamRot)

	if spiel.Spielstand[TeamRot] != 0 || spiel.Spielstand[TeamBlau] != 0 {
		t.Errorf("expected initial scores to be zero but got Rot:%d Blau:%d", spiel.Spielstand[TeamRot], spiel.Spielstand[TeamBlau])
	}
}

func TestSetTeamsInitialisiert(t *testing.T) {
	ermittlerRot, _ := NewSpieler("Ermittler Rot", Ermittler)
	ermittlerBlau, _ := NewSpieler("Ermittler Blau", Ermittler)
	chefRot, _ := NewSpieler("Chef Rot", Chef)
	chefBlau, _ := NewSpieler("Chef Blau", Chef)
	teamRot, _ := NewTeam(TeamRot, chefRot, []Spieler{ermittlerRot})
	teamBlau, _ := NewTeam(TeamBlau, chefBlau, []Spieler{ermittlerBlau})
	spiel := NewSpiel(TeamRot)

	err := spiel.SetTeams(teamRot, teamBlau)
	if err != nil {
		t.Errorf("unexpected error on setting the teams : %v", err)
	}
}

func TestSetTeamsErrorGleicheTeams(t *testing.T) {
	ermittlerRot, _ := NewSpieler("Ermittler Rot", Ermittler)
	chefRot, _ := NewSpieler("Chef Rot", Chef)
	teamRot, _ := NewTeam(TeamRot, chefRot, []Spieler{ermittlerRot})
	spiel := NewSpiel(TeamRot)

	err := spiel.SetTeams(teamRot, teamRot)
	if err == nil {
		t.Errorf("expected error about equal teams but didn't")
	}
}

func TestSetTeamsErrorVerschiedeneTeamsGleicheFarbe(t *testing.T) {
	ermittler1, _ := NewSpieler("Ermittler1", Ermittler)
	ermittler2, _ := NewSpieler("Ermittle1", Ermittler)
	chef1, _ := NewSpieler("Chef1", Chef)
	chef2, _ := NewSpieler("Chef2", Chef)
	team1, _ := NewTeam(TeamRot, chef1, []Spieler{ermittler1})
	team2, _ := NewTeam(TeamRot, chef2, []Spieler{ermittler2}) // gleiche Farbe!
	spiel := NewSpiel(TeamRot)

	if err := spiel.SetTeams(team1, team2); err == nil {
		t.Errorf("expected error about equal teams but didn't")
	}
}

func TestSetTeamsDisjunkteErmittlerBug(t *testing.T) {
	spiel := NewSpiel(TeamRot)

	chef1, _ := NewSpieler("Chef1", Chef)
	ermittler1a, _ := NewSpieler("Ermittler1a", Ermittler)
	ermittler1b, _ := NewSpieler("Ermittler1b", Ermittler)

	chef2, _ := NewSpieler("Chef2", Chef)
	ermittler3a, _ := NewSpieler("Ermittler3a", Ermittler) // Neuer Ermittler für das zweite Team

	team1, _ := NewTeam(TeamRot, chef1, []Spieler{ermittler1a, ermittler1b})
	team2, _ := NewTeam(TeamBlau, chef2, []Spieler{ermittler3a, ermittler1a}) // Überlappender Ermittlerspieler

	if err := spiel.SetTeams(team1, team2); err == nil {
		t.Fatalf("Fehler über einen Ermittlerspieler in beiden Teams erwartet, aber nicht erhalten")
	}
}

// todo karten nicht verteilt, Teams nicht gesetzt
func TestSpielStarten(t *testing.T) {
	// Karten verteilt, Teams gesetzt?
}

// todo weitere Fälle: nur Zahl schriftlich, beides mündlich
func TestChefGibtHinweisBeidesSchriftlich(t *testing.T) {
	t.Errorf("not yet implemented")
}

// todo irl kommt es manchmal vor, dass der Chef seinen Hinweis vom gegnerischen Chef prüfen lässt.

func TestSpielabbruch(t *testing.T) {
	t.Errorf("not yet implemented")
}

func TestChefKannNichtAufKartenTippen(t *testing.T) {
	// Punkt für Team, Zug geht weiter
	t.Errorf("not yet implemented")
}

func TestErmittlerTipptKarteDerEigenenFarbe(t *testing.T) {
	// Punkt für Team, Zug geht weiter
	t.Errorf("not yet implemented")
}

func TestErmittlerTipptKarteDerGegnerischenFarbe(t *testing.T) {
	// Zug beendet, Punkt für Gegner
	t.Errorf("not yet implemented")
}

func TestErmittlerTipptKarteDerNeutralenFarbe(t *testing.T) {
	// Zug beendet, kein Punkt
	t.Errorf("not yet implemented")
}

func TestErmittlerTipptKarteSchwarzeFarbe(t *testing.T) {
	// Spielende, Gegner gewinnt
	t.Errorf("not yet implemented")
}
