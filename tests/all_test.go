package tests

import "testing"

func TestAll(t *testing.T) {
	// Tests für die find Funktionen der Datenbank
	t.Run("TestFind", TestFind)

	// Tests für die Berechnungvorschrift der IT-Gearaete
	t.Run("TestComputations", TestComputations)
}
