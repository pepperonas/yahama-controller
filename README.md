# Yamaha RX-V577 Advanced Control

Eine professionelle Web-Anwendung zur vollständigen Steuerung des Yamaha RX-V577 AV-Receivers mit erweiterten Funktionen und modernem Dark Theme Interface.

![Yamaha Control Interface](yahama-mockup-1.png)

![Yamaha Control Interface - Extended Features](yamaha-mockup-2.png)

## Features

### 📱 **Navigation & Interface**
- **3 Main Tabs**: Grundsteuerung (Basic), Erweiterte Funktionen (Extended), System Info
- **Multi-Zone Control**: Hauptzone und Zone 2 mit separaten Einstellungen
- **Modern Dark Theme**: Professionelles dunkles UI mit CSS-Variablen
- **Responsive Design**: Optimiert für Desktop, Tablet und Mobile
- **PWA-Ready**: Progressive Web App mit Offline-Funktionalität
- **German Localization**: Vollständig deutsche Benutzeroberfläche

### 🔊 **Audio-Steuerung**
- **Power Control**: Ein/Aus mit visuellem Status-Indikator
- **Volume Control**: -80 dB bis -20 dB (erweitert bis +16 dB)
- **Volume Buttons**: Feineinstellung mit +/- Tasten
- **Mute Toggle**: Audio stumm/laut mit visueller Rückmeldung
- **Extended Volume Mode**: Freischaltung des vollen Lautstärkebereichs

### 📺 **Quellenauswahl**
- **HDMI Inputs**: HDMI 1, 2, 3, 4
- **Analog Inputs**: AV 1, AV 2
- **Digital Sources**: AirPlay, Server
- **Active Source Display**: Visuelle Anzeige der aktiven Quelle

### 🎵 **DSP & Surround Sound**
- **15 DSP Programme**:
  - Straight (Ohne DSP), Surround Decoder
  - Movie, Music, Game
  - Concert Hall, Jazz Club, Rock Concert
  - Stadium, Church, Chamber
  - Drama, Action Game, RPG, Sports
- **Dialogue Level**: Sprachverständlichkeit (-6 bis +6 dB)
- **Live DSP Display**: Anzeige des aktiven Modus

### 🎛️ **7-Band Equalizer**
- **Frequenzbänder**: 63Hz, 160Hz, 400Hz, 1kHz, 2.5kHz, 6.3kHz, 16kHz
- **Einstellbereich**: -6 bis +6 dB pro Band
- **EQ Reset**: Ein-Klick Zurücksetzen auf neutralen Klang
- **Live Preview**: Sofortige Anzeige der Einstellungen

### 🔧 **Erweiterte Audio-Features**
- **Extra Bass**: Verstärkte Basswidergabe
- **Bass/Treble Control**: -10 bis +10 dB Einstellung
- **Compressed Music Enhancer**: Verbesserung komprimierter Audiodateien
- **Pure Direct Mode**: Umgehung der Tonschaltungen
- **Straight Mode**: Direkte Signalverarbeitung
- **Virtual Presence Speaker (VPS)**: Virtueller Surround-Effekt

### 🏠 **Lautsprecher-Konfiguration** (YPAO-gesteuert, Anzeige-only)
- **Subwoofer Level**: Aktuelle YPAO-Einstellung
- **Center Speaker Level**: Center-Lautsprecher Pegel
- **Surround Level**: Surround-Lautsprecher Pegel
- **Dynamic Range Control (DRC)**: Status-Anzeige
- **Lip Sync Delay**: Audio-Verzögerung (YPAO-optimiert)

### 🎬 **Szenen-Steuerung**
- **4 Szenen-Presets**: Schnellzugriff auf gespeicherte Konfigurationen
- **Scene Activation**: Ein-Klick Szenen-Wechsel
- **Custom Configurations**: Persönliche Einstellungen speichern

### 📱 **HDMI-Einstellungen**
- **Audio Format**: PCM, DTS, Dolby Digital, Bitstream
- **HDMI Control**: Ein/Aus für HDMI-CEC Steuerung

### ⏰ **Sleep Timer**
- **Timer-Optionen**: Aus, 30, 60, 90, 120 Minuten
- **Automatisches Ausschalten**: Programmierbare Abschaltung

### 🎯 **Quick Setup Presets**
- **Movie Mode**: Optimiert für Kinoerlebnis
- **Music Mode**: Verbessert für Musikwiedergabe
- **Gaming Mode**: Angepasst für Gaming-Audio
- **Preset Save/Load**: Benutzerdefinierte Konfigurationen

### 📊 **System-Informationen**
- **Firmware Version**: Aktuelle System-Firmware
- **Gerätetemperatur**: Interne Temperaturüberwachung
- **Signal Format**: Aktuelle Signalinformationen
- **Kanal-Konfiguration**: Aktive Kanal-Anzeige
- **Netzwerk-Details**: IP, MAC, Gateway, Signalstärke

### 🌐 **Verbindungsmanagement**
- **IP-Konfiguration**: Manuelle Receiver-IP Einstellung
- **Connection Status**: Echtzeit-Verbindungsüberwachung
- **Auto-Connect**: Automatische Wiederverbindung beim Start
- **Connection Panel Toggle**: Einklappbare Verbindungssteuerung
- **Real-Time Polling**: 5-Sekunden Status-Updates

## Installation & Setup

### Automatische Installation mit PM2 (Empfohlen)

```bash
# Repository klonen
git clone https://github.com/your-username/yahama-amp.git
cd yahama-amp

# Dependencies installieren
npm install

# PM2 Process starten
pm2 start server.js --name yahama-amp

# PM2 Auto-Start aktivieren
pm2 startup
pm2 save
```

### Manuelle Installation

1. **Dependencies installieren**:
   ```bash
   npm install
   ```

2. **Server starten**:
   ```bash
   npm start
   ```

3. **Anwendung öffnen**:
   - Browser: `http://localhost:5001`
   - Netzwerk: `http://[RaspberryPI-IP]:5001`

4. **Find Your Receiver's IP Address**:
   - Check your router's admin panel for connected devices
   - Look for "Yamaha" or "RX-V577" in the device list
   - Or use the receiver's network menu to display the IP address

5. **Connect to Receiver**:
   - Enter the receiver's IP address in the connection panel
   - Click "Connect" to establish connection
   - The IP address will be saved in `receiver-config.json` for future use

### Alternative: Direct Browser Access (May Have CORS Issues)

1. **Direct Access**:
   - Open `index.html` in your web browser
   - No installation required, but may encounter CORS restrictions

2. **If CORS Issues Occur**:
   - Use Node.js server method (recommended) instead
   - Or run a local web server: `python -m http.server 8000`

## Projektstruktur

```
yahama-amp/
├── index.html              # Haupt-Interface mit Advanced Controls
├── server.js               # Express.js Server mit CORS Proxy
├── package.json            # Node.js Dependencies
├── receiver-config.json    # Gespeicherte Receiver-IP-Adresse
├── deploy-raspi.sh         # Deployment-Script für Raspberry Pi
├── favicon.*               # Favicon-Dateien (verschiedene Größen)
├── android-chrome-*.png    # Android App Icons
├── apple-touch-icon.png    # iOS App Icon
└── README.md              # Diese Dokumentation
```

## Technische Spezifikationen

- **Server**: Node.js mit Express.js
- **Dependencies**: Express, CORS, HTTP-Proxy-Middleware
- **Port**: 5001 (Standard)
- **Netzwerk**: Läuft auf 0.0.0.0 für Netzwerkzugriff
- **Process Manager**: PM2 für Produktionsumgebung
- **Interface**: Single Page Application mit Dark Theme
- **Config**: Automatische IP-Speicherung in JSON-Datei

## Usage

### Initial Connection
1. Enter your receiver's IP address (e.g., `192.168.1.100`)
2. Click "Connect" - the status indicator will show "Connected" when successful
3. Control panels will appear once connected

### Power Control
- Click the circular power button to toggle power on/off
- Status displays current power state (On/Standby)

### Volume Control
- Use the volume slider to adjust level (-80 dB to +16 dB)
- Click + or - buttons for precise adjustments
- Click "Mute" to toggle mute state

### Input Selection
- Click any input button to switch sources
- Active input is highlighted in blue
- Supports HDMI, AV, Audio, AirPlay, Server, USB, and Tuner inputs

### Zone Control
- Switch between Main Zone and Zone 2 using the tabs
- Each zone has independent control

### Scene Selection
- Click Scene 1-4 buttons to activate preconfigured settings
- Scenes combine input selection and DSP settings

## Supported Receivers

This application is designed for the Yamaha RX-V577 but should work with other Yamaha receivers that support the XML control protocol, including:

- RX-V series (RX-V473, RX-V573, RX-V673, RX-V773)
- HTR series with network capability
- Other network-enabled Yamaha AV receivers

## Technical Details

### XML Commands
The application uses Yamaha's XML control protocol via HTTP POST requests to `/YamahaRemoteControl/ctrl`.

Example commands:
```xml
<!-- Power On -->
<YAMAHA_AV cmd="PUT"><Main_Zone><Power_Control><Power>On</Power></Power_Control></Main_Zone></YAMAHA_AV>

<!-- Volume Set -->
<YAMAHA_AV cmd="PUT"><Main_Zone><Volume><Lvl><Val>-200</Val><Exp>1</Exp><Unit>dB</Unit></Lvl></Volume></Main_Zone></YAMAHA_AV>

<!-- Input Select -->
<YAMAHA_AV cmd="PUT"><Main_Zone><Input><Input_Sel>HDMI1</Input_Sel></Input></Main_Zone></YAMAHA_AV>

<!-- Get Status -->
<YAMAHA_AV cmd="GET"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>
```

### Network Requirements
- Receiver and controlling device must be on the same network
- HTTP requests to receiver IP address on port 80
- No authentication required for local network access

## Troubleshooting

### Connection Issues
- **"Connection Failed"**: Verify the IP address is correct
- **"Network error - check CORS settings"**: 
  - Modern browsers may block cross-origin requests
  - Consider using a local web server or proxy
  - Some browsers work better than others for local network requests

### Power On Issues
- **Power on doesn't work**: 
  - RX-V577 may not respond to network commands when in standby mode over Wi-Fi
  - Try using Ethernet connection for more reliable power on
  - Physical power button or IR remote may be needed for initial power on

### Status Not Updating
- Check network connectivity
- Verify receiver is powered on
- Status polling occurs every 5 seconds automatically

## CORS and Security

Due to browser security restrictions, you may encounter CORS (Cross-Origin Resource Sharing) issues when making requests to the receiver. Solutions:

1. **Use a local web server**:
   ```bash
   # Python 3
   python -m http.server 8000
   
   # Node.js
   npx http-server
   ```

2. **Browser flags** (for testing only):
   - Chrome: `--disable-web-security --user-data-dir=/tmp/chrome_dev`
   - Not recommended for regular use

3. **Proxy setup**: Configure a local proxy to forward requests to the receiver

## Browser Compatibility

- Chrome/Chromium: Best compatibility
- Firefox: Good compatibility
- Safari: Limited due to stricter CORS policies
- Edge: Good compatibility

## PM2 Management

```bash
# Status anzeigen
pm2 status

# Logs anzeigen
pm2 logs yahama-amp

# Neustart
pm2 restart yahama-amp

# Stoppen
pm2 stop yahama-amp

# Entfernen
pm2 delete yahama-amp
```

## Raspberry Pi Deployment

Das Projekt läuft optimal auf einem Raspberry Pi als dedizierter Yamaha-Controller:

1. **Raspberry Pi OS** mit Node.js und PM2
2. **Netzwerkzugriff** über lokale IP-Adresse
3. **Auto-Start** beim Systemstart durch PM2
4. **Log-Management** durch PM2 integriert

## Lizenz

Martin Pfeffer 2025 Berlin, Lizenz MIT

Yamaha und RX-V577 sind Markenzeichen der Yamaha Corporation.