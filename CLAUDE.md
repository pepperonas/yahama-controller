# 🎛️ Yamaha RX-V577 Controller - Claude Project Documentation

## Projektstatus
- **Status**: ✅ Vollständig funktionsfähig und strukturiert
- **Letztes Update**: Juli 2025
- **Version**: 1.0.0
- **Technologie**: Node.js Express Server mit statischer Frontend-App

## Aktuelle Projektstruktur
```
yamaha-rx-v577-controller/
├── public/                 # Statische Assets
│   ├── assets/            # Bilder und andere Assets
│   │   ├── yahama-mockup-1.png    # Mockup-Bild 1
│   │   └── yamaha-mockup-2.png    # Mockup-Bild 2
│   ├── favicon.ico        # Website Icon
│   ├── favicon-16x16.png  # 16x16 Favicon
│   ├── favicon-32x32.png  # 32x32 Favicon
│   ├── apple-touch-icon.png       # iOS App Icon
│   ├── android-chrome-192x192.png # Android App Icon (192x192)
│   └── android-chrome-512x512.png # Android App Icon (512x512)
├── index.html              # Haupt-Interface mit Advanced Controls
├── server.js               # Express.js Server mit CORS Proxy
├── package.json            # Node.js Dependencies
├── receiver-config.json    # Gespeicherte Receiver-IP (gitignore)
├── deploy-raspi.sh         # Deployment-Script für Raspberry Pi
├── .gitignore             # Git Ignore-Regeln
├── README.md              # Projektdokumentation
└── CLAUDE.md              # Diese Datei
```

## Server-Konfiguration

### PM2 Process Manager
- **Service Name**: `yamaha-controller`
- **Port**: 5001
- **Auto-Start**: ✅ Konfiguriert via SystemD (`pm2-martin.service`)
- **Status**: Online und läuft stabil
- **Zugriff**: http://localhost:5001

### Express.js Server Details
- **Statische Dateien**: Serviert aus `public/` Ordner
- **CORS**: Aktiviert für alle Routen
- **Proxy**: Yamaha Receiver API über `/api/receiver/*`
- **Cache-Headers**: Optimiert für Favicons und Assets (1 Jahr)
- **HTML No-Cache**: Verhindert Browser-Caching für Interface-Updates

## Funktionalitäten

### ✅ Implementierte Features
1. **Vollständige Receiver-Steuerung**
   - Power On/Off mit visueller Rückmeldung
   - Lautstärke-Kontrolle (-80 dB bis +16 dB)
   - Quellenauswahl (HDMI 1-4, AV 1-2, AirPlay, Server)
   - Stummschaltung mit Toggle-Switch

2. **DSP & Audio Processing**
   - 15 DSP Programme (Movie, Music, Game, Concert Hall, etc.)
   - Dialog Level Adjustment (-6 bis +6 dB)
   - 7-Band Equalizer (63Hz - 16kHz)
   - Bass/Treble Controls (-6 bis +6 dB)

3. **Erweiterte Audio-Features**
   - Extra Bass Toggle
   - Compressed Music Enhancer
   - Pure Direct Mode
   - Straight Mode
   - Virtual Presence Speaker (VPS)

4. **System & Netzwerk**
   - IP-Konfiguration mit persistenter Speicherung
   - Real-time Status Polling (5 Sekunden)
   - Multi-Zone Support (Main Zone, Zone 2)
   - Szenen-Presets (1-4)
   - Sleep Timer (30-120 Minuten)

5. **UI/UX**
   - Modern Dark Theme Interface
   - 3 Haupt-Tabs: Basic, Extended, System Info
   - Responsive Design (Desktop/Mobile)
   - Status-Indikatoren und Live-Updates
   - Collapsible Connection Panel

### 📱 Interface-Tabs
1. **Grundsteuerung**: Power, Volume, Input, DSP, Equalizer
2. **Erweiterte Funktionen**: Audio Enhancement, Sleep Timer, Presets
3. **System Info**: Firmware, Temperatur, Netzwerk-Details

## Technische Implementierung

### Yamaha XML API Integration
- **Protokoll**: HTTP POST requests mit XML-Payloads
- **Endpoint**: `/YamahaRemoteControl/ctrl` auf Receiver IP
- **Status Polling**: Automatisch alle 5 Sekunden
- **Error Handling**: Comprehensive mit User-Feedback

### Frontend-Technologie
- **Vanilla JavaScript**: Keine Frameworks, optimierte Performance
- **CSS Custom Properties**: Dark Theme mit CSS-Variablen
- **Responsive Grid**: Auto-fit layouts für verschiedene Screengrößen
- **PWA-Ready**: Favicon-Sets und Cache-Headers vorbereitet

### Server-Features
- **CORS Proxy**: Umgeht Browser-Sicherheitsbeschränkungen
- **IP Validation**: Regex-basierte Validierung für Receiver-IPs
- **Config Persistence**: JSON-basierte Speicherung der Receiver-IP
- **Health Check**: `/api/health` Endpoint für Monitoring

## Git Repository

### Repository-Details
- **URL**: https://github.com/pepperonas/yahama-controller.git
- **Branch**: main
- **Letzter Commit**: "Update project to Yamaha RX-V577 Controller"

### .gitignore Konfiguration
- Node.js Dependencies (`node_modules/`)
- Runtime Config (`receiver-config.json`)
- Core Dumps (`core`)
- OS-spezifische Dateien (`.DS_Store`, `Thumbs.db`)
- IDE-Dateien (`.vscode/`, `.idea/`)
- Logs und Cache-Dateien

## Deployment & Betrieb

### Raspberry Pi Deployment
- **Auto-Start**: SystemD Service konfiguriert
- **Process Manager**: PM2 für Stability und Monitoring
- **Network Access**: Läuft auf 0.0.0.0:5001 für LAN-Zugriff
- **Logging**: PM2 integrierte Log-Verwaltung

### Übliche Befehle
```bash
# PM2 Management
pm2 status                    # Status aller Services
pm2 logs yamaha-controller    # Log-Ausgabe anzeigen
pm2 restart yamaha-controller # Service neu starten
pm2 save                      # Aktuelle Konfiguration speichern

# Manual Testing
npm start                     # Direkter Start für Development
curl http://localhost:5001/api/health  # Health Check
```

## Bekannte Kompatibilität

### Getestete Receiver
- **Yamaha RX-V577**: ✅ Vollständig kompatibel
- **Andere RX-V Serie**: Sehr wahrscheinlich kompatibel
- **Yamaha HTR Serie**: Mit Netzwerk-Features kompatibel

### Browser-Support
- **Chrome/Chromium**: ✅ Beste Kompatibilität
- **Firefox**: ✅ Vollständig kompatibel
- **Safari**: ⚠️ CORS-Einschränkungen möglich
- **Edge**: ✅ Vollständig kompatibel

## Entwicklungshinweise

### Code-Qualität
- **Error Handling**: Comprehensive try-catch blocks
- **User Feedback**: Toast-Messages für alle Aktionen
- **Performance**: Optimierte Status-Updates und Caching
- **Security**: IP-Validation und No-Script-Injection

### Erweiterbarkeit
- **Modular JavaScript**: Gut strukturierte YamahaAdvancedReceiver Klasse
- **CSS-Variablen**: Einfache Theme-Anpassungen
- **Config-System**: JSON-basiert, leicht erweiterbar
- **API-Ready**: RESTful Design für zukünftige Features

### Debugging
- **PM2 Logs**: `pm2 logs yamaha-controller`
- **Browser Console**: JavaScript Error-Logging
- **Network Tab**: XML-Requests zur Receiver-API überwachen
- **Health Endpoint**: `/api/health` für Server-Status

---

**Projekt erfolgreich implementiert und deployment-ready!** 🎉