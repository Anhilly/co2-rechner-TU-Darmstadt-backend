package database

import (
	"github.com/Anhilly/co2-rechner-TU-Darmstadt-backend/structs"
	"github.com/matryer/is"
	"testing"
)

func TestDatabaseSubfunctions(t *testing.T) {
	//is := is.NewRelaxed(t)

	t.Run("TestIntInSlice", TestIntInSlice)
	t.Run("TestJahrInVersogerSlice", TestJahrInVersogerSlice)
}

func TestIntInSlice(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("intInSlice: True", func(t *testing.T) {
		is := is.NewRelaxed(t)

		a := int32(10)
		b := []int32{10, 20, 1}
		is.True(intInSlice(a, b))
	})

	t.Run("intInSlice: False", func(t *testing.T) {
		is := is.NewRelaxed(t)

		a := int32(100)
		b := []int32{10, 20, 1}
		is.True(!intInSlice(a, b))
	})
}

func TestJahrInVersogerSlice(t *testing.T) {
	is := is.NewRelaxed(t)

	// Normalfall
	t.Run("jahrInVersorgerSlice: True", func(t *testing.T) {
		is := is.NewRelaxed(t)

		jahr := int32(2000)
		versorger := []structs.Versoger{
			{Jahr: 2000, IDVertrag: 0},
			{Jahr: 2003, IDVertrag: 0},
			{Jahr: 2005, IDVertrag: 0},
		}
		is.True(jahrInVersorgerSlice(jahr, versorger))
	})

	t.Run("jahrInVersorgerSlice: True", func(t *testing.T) {
		is := is.NewRelaxed(t)

		jahr := int32(2005)
		versorger := []structs.Versoger{
			{Jahr: 2000, IDVertrag: 0},
			{Jahr: 2003, IDVertrag: 0},
			{Jahr: 2005, IDVertrag: 0},
		}
		is.True(jahrInVersorgerSlice(jahr, versorger))
	})

	t.Run("jahrInVersorgerSlice: False", func(t *testing.T) {
		is := is.NewRelaxed(t)

		jahr := int32(2004)
		versorger := []structs.Versoger{
			{Jahr: 2000, IDVertrag: 0},
			{Jahr: 2003, IDVertrag: 0},
			{Jahr: 2005, IDVertrag: 0},
		}
		is.True(!jahrInVersorgerSlice(jahr, versorger))
	})
}
