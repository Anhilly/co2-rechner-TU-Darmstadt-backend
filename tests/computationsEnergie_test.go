package tests

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/co2computation"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/database"
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"testing"
)

func TestComputationsEnergie(t *testing.T) {
	database.ConnectDatabase()
	defer database.DisconnectDatabase()

	t.Run("TestBerechneEnergieverbrauch", TestBerechneEnergieverbrauch)
}

func TestBerechneEnergieverbrauch(t *testing.T) {
	is := is.NewRelaxed(t)

	// normale Berechnungen
	t.Run("BerechneEnergieverbrauch: Slice = nil", func(t *testing.T) {
		is := is.NewRelaxed(t)

		var geaeudeFlaecheDaten []structs.GebaeudeFlaecheAPI = nil
		var jahr int32 = 2020
		var idEnergieversorgung int32 = 1

		emissionen, err := co2computation.BerechneEnergieverbrauch(geaeudeFlaecheDaten, jahr, idEnergieversorgung)

		is.NoErr(err)             // bei normalen Berechnungen sollte kein Error geworfen werden
		is.Equal(emissionen, 0.0) // ohne Eingaben sind Emissionen = 0.0
	})
}
