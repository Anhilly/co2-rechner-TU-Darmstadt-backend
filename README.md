# Backend des CO2 Rechners der TU Darmstadt

Das Backend des CO2 Rechner, welcher im Rahmen des Bachelorpraktikums für das Büro für Nachhaltigkeit und das Institut für Fluidsystemtechnik entwickelt wurde.

Dieses Projekt ermöglicht eine Erfassung von CO2 Emissionen von TU Einheiten und anschließende Auswertung.

## Implementierte Funktionen

- Anmeldungssteuerung
- Emissionsberechnung
- Datenbankinteraktionen

## Verwandte Projekte

Das Backend kann nicht unabhängig betrieben werden.

Das Frontend, welches eine webbasierte Interaktionsseite bietet, ist hier einsehbar: [Github](https://github.com/Lithium-1Hauptgruppe/CO2-Rechner-TU-Darmstadt-Frontend)
Die Interaktionen zwischen Frontend und Backend sind in einer REST artigen API definiert, welche hier eingesehen werden kann: [Github](https://github.com/Anhilly/CO2-Rechner-api)

## Abhängigkeiten

Das Projekt ist in der Sprache Go geschrieben.
Der CO2 Rechner verwendet folgende Abhängigkeiten um die Funktionalität bereitzustellen:

- [Go Lang Version 1.17](https://go.dev/) - Go Entwicklungssprache
- [go chi Version 5.0.7](github.com/go-chi/chi) - Go Router für HTTP Dienste 
- [go cors Version 1.2.0](github.com/go-chi/cors) - Go Middleware CORS
- [UUID Version 1.3.0](github.com/google/uuid) - Eindeutige ID Generierung
- [is Version 1.4.0](github.com/matryer/is) - Test Framework
- [errors Version 0.9.1](github.com/pkg/errors) - Vereinfachte Fehlerbehandlung
- [mongo-driver Version 1.8.0](go.mongodb.org/mongo-driver) - Mongodb Treiber für Go
- [crypto Version 0.0.0-20201216223049-8b5274cf687f](golang.org/x/crypto) - Verschlüsselungsalgorithmen

## Entwicklungssetup

Nach Download des Repositorys muss eine neue Datei mit Datenbankinformationen angelegt werden.
Als Vorlage dient die Datei database/db_config_example.go, aus der die Datei database/db_config.go erstellt werden muss.
