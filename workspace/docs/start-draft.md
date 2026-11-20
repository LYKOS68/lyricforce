# Was ist LyricForce?

LyricForce ist eine eigenständige Windows-Anwendung für Deutschrap, Songwriting und kreatives Arbeiten mit Sprache. Im Mittelpunkt steht der Aufbau eines persönlichen, strukturierten Reimarchivs, in dem Reimwörter, Reimketten, Phrasen und relatable Sätze gesammelt, organisiert und gezielt wiedergefunden werden können. Statt verstreute Wortlisten und spontane Einfälle zu verwalten, entsteht ein zentraler Ort, an dem sprachliches Material dauerhaft erhalten bleibt und jederzeit nutzbar ist.

Das zweite Herzstück von LyricForce ist ein kreativer Generator, der eigenständig neue Reime und relatable Sätze entwickeln kann. Dabei soll nicht einfach irgendein Satz entstehen, der sich zufällig reimt, sondern ein grammatikalisch korrekter, verständlicher und inhaltlich nachvollziehbarer Gedanke. Endwörter, Reimstrukturen, Silbenbegrenzungen, Kontraste und Gegensätze sollen gezielt berücksichtigt werden können. Die generierten Ergebnisse sollen anschließend bearbeitet, verglichen und direkt in das eigene Archiv übernommen werden können.

LyricForce richtet sich an Rapper, Songwriter und kreative Sprachmenschen, die ihre sprachlichen Möglichkeiten systematisch erweitern und ihr eigenes Reim- und Ideenmaterial aufbauen möchten. Die Kernidee besteht darin, ein lebendiges Spracharchiv mit einem intelligenten Kreativwerkzeug zu verbinden, sodass gespeicherte Inhalte neue Ideen ermöglichen und neue Ideen wiederum das Archiv bereichern. LyricForce soll damit langfristig mehr sein als ein Reimlexikon: eine persönliche digitale Werkbank für Sprache, Gedanken, Reimkunst und Songwriting.

# Technologie-Entscheidung (vorläufig)

## Bevorzugte Kombination: Go + Wails + SQLite

Die aktuell bevorzugte technische Grundlage ist **Go + Wails + SQLite**:

- **Go** übernimmt die gesamte Anwendungslogik, Datenverarbeitung und Generator-Logik. Go ist bereits installiert und eingerichtet, sodass keine zweite Programmiersprache oder zusätzliche Entwicklungsumgebung eingeführt werden muss.
- **Wails** verbindet die Go-Backend-Logik mit einer modernen, frei gestaltbaren Benutzeroberfläche auf Basis von HTML/CSS/JavaScript. Es ist kein Electron-Ersatz im negativen Sinne, sondern eine schlanke Brücke zwischen Go und Web-UI-Technologien.
- **SQLite** speichert das gesamte Reimarchiv dauerhaft und lokal auf dem Gerät des Nutzers. Keine Cloud-Abhängigkeit, keine externe Datenbankinfrastruktur.

## Bewertung der Alternativen

- **C#** wäre ebenfalls eine sehr gute Wahl für eine klassische Windows-Desktop-Anwendung (z.B. WPF oder WinUI), würde aber eine komplett andere technische Grundlage und zusätzliche Einrichtung bedeuten. Derzeit nicht bevorzugt.
- **Fyne** (natives Go-GUI-Framework) scheidet aus: Bei einer anspruchsvollen, individuell gestalteten Oberfläche mit Seitenleiste, Karten, Suchbereichen und mehreren Arbeitsbereichen stößt Fyne schnell an seine Grenzen.

## UI-Anforderung bestätigt

Das gewünschte Erscheinungsbild von LyricForce ist das einer **digitalen Werkbank**: Seitenleiste, Suchbereich, Karten, Ergebnislisten und mehrere Arbeitsbereiche. Dieser Werkbank-Charakter erfordert eine bewusst gestaltete, frei layoutbare Oberfläche — keine Standard-Widget-UI. Wails erfüllt diese Anforderung, Fyne nicht.

## Status

Diese Technologie-Entscheidung ist **noch nicht festgeschrieben**. Sie wird anhand der tatsächlichen Anforderungen weiter geprüft, bevor eine endgültige Festlegung erfolgt.

# Funktionsumfang und Entwicklungsstrategie

## Langfristiger Funktionsumfang

- **Reimarchiv:** Einzelwörter, Reimketten, Mehrsilber und Endreime
- **Relatable-Archiv:** Vollständige, grammatikalisch korrekte Sätze mit eigenständiger Aussage
- **Kontrastgenerator:** Gegensätze, Widersprüche, Perspektivwechsel
- **Silbenanalyse:** Automatische Silbenzählung und metrische Filter
- **Suchsystem:** Suche nach Endungen, Reimstruktur, Thema und Silbenanzahl

## MVP-Fokus: Einsilbige Wörter

Der Einstieg erfolgt bewusst klein und kontrolliert:

- Einsilbige Wörter als erste Reimkategorie
- Einfache Vokalklänge als Ausgangspunkt
- Schrittweise Erweiterung um Diphthonge und komplexere Lautstrukturen

## Frontend-Platzhalter-Prinzip

Das Frontend wird **etappenweise** gebaut. Bereiche, die noch nicht implementiert sind, werden mit **expliziten, sichtbaren Platzhaltern** markiert — kein verstecktes "kommt später". Dieser Ansatz legt den Fokus gezielt auf den aktiven Arbeitsbereich und macht den Entwicklungsstand jederzeit transparent.

## Phonetische Reimlogik (fachliche Grundlage)

Die Reimengine arbeitet **phonetisch**, nicht orthografisch. Buchstabenvergleiche werden von Anfang an vermieden, um einen späteren Neuaufbau der Grundlogik zu verhindern.

Konkrete fachliche Festlegungen:

- **„ie" ist kein Diphthong**, sondern die Schreibweise für den langen Vokal /iː/.
- **Diphthonge** im Deutschen sind: /aɪ̯/ (ei, ai), /aʊ̯/ (au) und /ɔʏ̯/ (eu, äu).
- Die Reimengine berücksichtigt die **Lautfolge ab dem betonten Vokal**, nicht nur den Vokal selbst.
- Lautklassen und Reimstrukturen werden schrittweise erweitert, ohne die Grundlogik ersetzen zu müssen.

## Entwicklungsreihenfolge (grob)

1. Einsilbige Wörter, einfache Vokalklänge
2. Diphthonge (/aɪ̯/, /aʊ̯/, /ɔʏ̯/)
3. Mehrsilbige Reime und Endreime
4. Silbenanalyse und metrische Filter
5. Relatable-Archiv und Kontrastgenerator
6. Vollständiges Suchsystem
