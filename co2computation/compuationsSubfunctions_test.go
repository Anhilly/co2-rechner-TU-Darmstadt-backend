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

	// Normalfall
	t.Run("getEnergieversorgung: ID = 1, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(144)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: ID = 2, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 2
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                   // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(285)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: ID = 3, Jahr = 2020", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 3
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.NoErr(err)                  // Normalfall sollte keinen Error verursachen
		is.Equal(co2Faktor, int32(72)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	// Errortests
	t.Run("getEnergieversorgung: ID = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 0
		var jahr int32 = 2020

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, io.EOF)         // Datenbank wirft EOF fuer unbekannte IDs
		is.Equal(co2Faktor, int32(0)) // tritt ein Fehler auf wird 0 zurückgegeben
	})

	t.Run("getEnergieversorgung: Jahr = 0", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var idEnergieversorgung int32 = 1
		var jahr int32 = 0

		co2Faktor, err := getEnergieCO2Faktor(idEnergieversorgung, jahr)

		is.Equal(err, ErrJahrNichtVorhanden) // Funktion wirft ErrJahrNichtVorhanden fuer unbekanntes Jahr
		is.Equal(co2Faktor, int32(0))        // tritt ein Fehler auf wird 0 zurückgegeben
	})
}
