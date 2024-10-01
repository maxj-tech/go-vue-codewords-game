package domain

import (
	"github.com/google/uuid"
	"testing"
)

func TestNeuerSpieler(t *testing.T) {
	// Teste, ob ein neuer Spieler korrekt erstellt wird
	spieler, _ := NewSpieler("Spieler1", Ermittler)

	if spieler.Name != "Spieler1" {
		t.Fatalf("expected name Spieler1, got %s", spieler.Name)
	}

	if spieler.Rolle != Ermittler {
		t.Fatalf("expected role Ermittler, got %v", spieler.Rolle)
	}

	// Teste, ob die ID korrekt erstellt wird
	if spieler.id == uuid.Nil {
		t.Fatalf("expected non-nil UUID, got %v", spieler.id)
	}

	// prüfe ob die ID eine gültige UUID ist
	_, err := uuid.Parse(spieler.id.String())
	if err != nil {
		t.Fatalf("expected valid UUID, got %v", spieler.id)
	}
}

func TestNeuerSpielerMitLeeremNamen(t *testing.T) {
	tests := []string{"", " ", "\t", "\n"}

	for _, name := range tests {
		spieler, err := NewSpieler(name, Chef)
		if err == nil || spieler.Name != "" {
			t.Fatalf("expected error for empty or whitespace name, got %v and name %s", err, spieler.Name)
		}
	}
}
