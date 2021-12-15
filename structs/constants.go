package structs

import "time"

const ( // Konstanten fuer Collection Namen
	KaeltezaehlerCol      = "kaeltezaehler"
	StromzaehlerCol       = "stromzaehler"
	WaermezaehlerCol      = "waermezaehler"
	DienstreisenCol       = "dienstreisen"
	EnergieversorgungCol  = "energieversorgung"
	ITGeraeteCol          = "itGeraete"
	PendelwegCol          = "pendelweg"
	GebaeudeCol           = "gebaeude"
	UmfrageCol            = "umfrage"
	MitarbeiterUmfrageCol = "mitarbeiterUmfrage"
	NutzerdatenCol        = "nutzerdaten"
)

const TimeoutDuration time.Duration = 5 * time.Second // Timeout Zeit fuer Datenbank-Kontext

const DumpPath = "/autoDump/" // Pfad fuer die automatischen Dumps

const ( // nach IDs in der Datenbank
	IDEnergieversorgungWaerme int32 = 1
	IDEnergieversorgungStrom  int32 = 2
	IDEnergieversorgungKaelte int32 = 3
)

const ( // nach IDs in der Datenbank
	IDDienstreiseBahn     int32 = 1
	IDDienstreiseAuto     int32 = 2
	IDDienstreiseFlugzeug int32 = 3
)

const ( // fuer Einheiten
	EinheitkWh     = "kWh"
	EinheitMWh     = "MWh"
	Einheitqm      = "m^2"
	EinheitgkWh    = "g/kWh"
	EinheitgPkm    = "g/Pkm"
	EinheitgStueck = "g/Stueck"
)

const ( // fuer Zaehertypen
	ZaehlertypWaerme = "Waerme"
	ZaehlertypKaelte = "Kaelte"
	ZaehlertypStrom  = "Strom"
)

const (
	ResponseSuccess = "success"
	ResponseError   = "error"
)
