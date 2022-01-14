package tests

import (
	"testing"
)

func TestAll(t *testing.T) {
	// Tests für die find Funktionen der Datenbank
	t.Run("TestFind", TestFind)
	t.Run("TestAdd", TestAdd)
	t.Run("TestInsert", TestInsert)
	t.Run("TestUpdate", TestUpdate)
	t.Run("TestDelete", TestDelete)

	// Tests für die Berechnungsvorschrift von Dienstreisen, Pendelwege und IT-Geraeten
	t.Run("TestComputations", TestComputations)
	t.Run("TestComputationsEnergie", TestComputationsEnergie)
}
