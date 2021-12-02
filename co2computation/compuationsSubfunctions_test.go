package co2computation

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/matryer/is"
	"io"
	"testing"
)

func TestComputaionsSubfunctions(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestGetEnergieCO2Faktor", TestGetEnergieCO2Faktor)
}

func TestGetEnergieCO2Faktor(t *testing.T) {
	is := is.NewRelaxed(t)

	// Errortests
	t.Run("getEnergieversorgung: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 0
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, io.EOF)         // Datenbank wirft EOF fuer unbekannte IDs
		is.Equal(co2Faktor, int32(0)) // tritt ein Fehler auf wird 0 zurückgegeben
	})
}
