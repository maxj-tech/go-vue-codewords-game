package domain

import (
	"github.com/google/uuid"
	"testing"
)

func TestNeuerSpieler(t *testing.T) {
	// Teste, ob ein neuer Spieler korrekt erstellt wird
	spieler := NewSpieler("Spieler1", Ermittler)

	if spieler.Name != "Spieler1" {
		t.Fatalf("expected name Spieler1, got %s", spieler.Name)
	}

	if spieler.Rolle != Ermittler {
		t.Fatalf("expected role Ermittler, got %v", spieler.Rolle)
	}

	// Teste, ob die ID korrekt erstellt wird
	if spieler.Id == uuid.Nil {
		t.Fatalf("expected non-nil UUID, got %v", spieler.Id)
	}

	// prüfe ob die ID eine gültige UUID ist
	_, err := uuid.Parse(spieler.Id.String())
	if err != nil {
		t.Fatalf("expected valid UUID, got %v", spieler.Id)
	}
}

// weitere: Name darf nicht leer oder whitespace sein
