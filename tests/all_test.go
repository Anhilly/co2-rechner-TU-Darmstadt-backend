package tests

import (
	"testing"
)

func TestAll(t *testing.T) {
	// Tests für die find Funktionen der Datenbank
	t.Run("TestFind", TestFind)

	// Tests für die Berechnungvorschrift von Dienstreisen, Pendelwege und IT-Geraeten
	t.Run("TestComputations", TestComputations)
	t.Run("TestComputationsEnergie", TestComputationsEnergie)
}
