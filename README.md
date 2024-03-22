# Backend des CO2 Rechners der TU Darmstadt

Das Backend des CO2 Rechner, welcher im Rahmen des Bachelorpraktikums für das Büro für Nachhaltigkeit und das Institut für Fluidsystemtechnik entwickelt wurde.

Dieses Projekt ermöglicht eine Erfassung von CO2 Emissionen von TU Einheiten und anschließende Auswertung.

## Implementierte Funktionen

- Anmeldungssteuerung
- Emissionsberechnung
- Datenbankinteraktionen

## Verwandte Projekte

Das Backend kann nicht unabhängig betrieben werden.

Das Frontend, welches eine webbasierte Interaktionsseite bietet, ist hier einsehbar: [Github](https://github.com/felix-marx/CO2-Rechner-TU-Darmstadt-Frontend)  
Die Interaktionen zwischen Frontend und Backend sind in einer REST artigen API definiert, welche hier eingesehen werden kann: [Github](https://github.com/Anhilly/CO2-Rechner-api)

## Abhängigkeiten

Das Projekt ist in der Sprache Go geschrieben.
Der CO2-Rechner verwendet die folgenden direkten Abhängigkeiten, um die Funktionalität bereitzustellen:

- [Go Lang Version 1.17](https://go.dev/) - Go Entwicklungssprache
- [go chi Version 5.0.7](https://github.com/go-chi/chi) - Go Router für HTTP Dienste
- [is Version 1.4.0](https://github.com/matryer/is) - Test Framework
- [mongo-driver Version 1.8.0](https://go.mongodb.org/mongo-driver) - Mongodb Treiber für Go
- [lumberjack Version 2.0.0](https://gopkg.in/natefinch/lumberjack.v2) - Logger
- [gocloak Version 13.8.0](https://github.com/Nerzal/gocloak/) - Go Keycloak Client

## Entwicklungssetup

Für das Entwicklungssetup wird eine lokale Installation von Go benötigt. Die Docker Compose Datei enthält die folgenden Container für die Entwicklung:
- MongoDB als Datenbank
- NGINX als Webserver und Reverse Proxy
- Keycloak zur Authentifizierung und Kommunikation mit externen Diensten
- Postgres als Datenbank für Keycloak

Das Frontend und Backend müssen unabhängig von der Docker Compose lokal gestartet werden. 
Die default Konfiguration erwartet das Frontend unter `localhost:8081` und das Backend unter `localhost:3000`.

Fürs Setup muss eine dump der Datenbank in den Ordner `development/dump` gelegt werden. 
Zusätzlich muss eine `config.go` Datei erstellt und in den Ordner `config` abgelegt werden. 
Ein Beispiel für die `config.go` Datei ist in `config/config.go.example` zu finden. 

Ob das Backend in `prod` oder `dev` Modus startet, wird über die Variable `mode` in `main.go` gesteuert. 
Die Variable kann entweder manuell oder per symbol substitution während link time gesetzt werden. 
Für symbol substitution muss die folgende Flag gesetzt werden:
- `-ldflags "-X main.mode=dev"` für dev mode
- `-ldflags "-X main.mode=prod"` für prod mode