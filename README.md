# 🎛️ Yamaha RX-V577 Controller

> **⚡ Update 2026-06 — Stack & UI**
>
> - **Backend:** schlankes **Go-Binary** (~7 MB, migriert von Node/Express). Serviert die PWA, `/api/health`, `/api/set-receiver-ip` und reverse-proxyt `/api/receiver/*` an den RX-V577. **systemd** `yamaha-controller`. Source: `main.go` (+ `/Users/martin/claude/yamaha-controller-go/`, Cross-Build `GOOS=linux GOARCH=arm64`).
> - **UI:** **Material Design 3 Expressive** + Spring-Animationen.
> - **Deploy:** Binary bauen → `scp` → `sudo systemctl restart yamaha-controller`

<div align="center">

![Yamaha Controller](public/assets/social-preview.png)

![GitHub last commit](https://img.shields.io/github/last-commit/pepperonas/yahama-controller?style=for-the-badge&color=green)
![GitHub repo size](https://img.shields.io/github/repo-size/pepperonas/yahama-controller?style=for-the-badge)
![GitHub code size](https://img.shields.io/github/languages/code-size/pepperonas/yahama-controller?style=for-the-badge)
![GitHub issues](https://img.shields.io/github/issues/pepperonas/yahama-controller?style=for-the-badge)
![GitHub pull requests](https://img.shields.io/github/issues-pr/pepperonas/yahama-controller?style=for-the-badge)
![GitHub stars](https://img.shields.io/github/stars/pepperonas/yahama-controller?style=for-the-badge)
![GitHub forks](https://img.shields.io/github/forks/pepperonas/yahama-controller?style=for-the-badge)
![GitHub watchers](https://img.shields.io/github/watchers/pepperonas/yahama-controller?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-blue.svg?style=for-the-badge)
![Node.js](https://img.shields.io/badge/Node.js-18.x-green?style=for-the-badge&logo=node.js)
![Express.js](https://img.shields.io/badge/Express.js-4.x-000000?style=for-the-badge&logo=express)
![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)
![Status](https://img.shields.io/badge/status-active-success.svg?style=for-the-badge)
![Maintained](https://img.shields.io/badge/Maintained%3F-yes-green.svg?style=for-the-badge)
![Platform](https://img.shields.io/badge/platform-Raspberry%20Pi-C51A4A?style=for-the-badge&logo=raspberry-pi)

<h3>Professional Web Application for Complete Control of Yamaha RX-V577 AV Receivers</h3>

<p>
  <strong>A modern, feature-rich web interface with advanced controls and dark theme support</strong>
</p>

![Yamaha Control Interface](public/assets/yahama-mockup-1.png)

![Yamaha Control Interface - Extended Features](public/assets/yamaha-mockup-2.png)

</div>

---

## 💖 Support This Project

<div align="center">

If you find this project useful, consider supporting its development:

[![PayPal](https://img.shields.io/badge/PayPal-00457C?style=for-the-badge&logo=paypal&logoColor=white)](https://www.paypal.com/donate/?hosted_button_id=YOUR_BUTTON_ID)

<sub>Your support helps maintain and improve this project. Every contribution, no matter how small, is greatly appreciated! ❤️</sub>

</div>

---

## 📋 Table of Contents

- [Features](#-features)
- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
- [Technical Details](#-technical-details)
- [API Documentation](#-api-documentation)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

### 📱 Navigation & Interface
- **Multi-Tab Interface**: Basic controls, extended features, and system information
- **Multi-Zone Control**: Independent control for Main Zone and Zone 2
- **Dual Theme Support**: Dark theme (default) and light theme with one-click toggle
- **Progressive Web App**: Installable on mobile devices with offline functionality
- **Responsive Design**: Optimized for desktop, tablet, and mobile devices
- **Real-time Updates**: Status polling every 5 seconds for live updates

### 🔊 Audio Control
- **Volume Management**: -80 dB to +16 dB range with fine adjustment buttons
- **Power Control**: On/Off with visual status indicators
- **Mute Toggle**: Quick audio muting with visual feedback
- **Extended Volume Mode**: Unlock full volume range capabilities

### 📺 Input Selection
- **HDMI Inputs**: HDMI 1-4 support
- **Analog Inputs**: AV 1-2 channels
- **Digital Sources**: AirPlay, Server, USB, Tuner
- **Active Source Display**: Visual highlighting of current input

### 🎵 DSP & Surround Sound
- **15 DSP Programs**: Including Straight, Surround Decoder, Movie, Music, Game modes
- **Environment Simulations**: Concert Hall, Jazz Club, Rock Concert, Stadium, Church
- **Gaming Modes**: Action Game, RPG, Sports optimizations
- **Dialogue Enhancement**: -6 to +6 dB adjustment for voice clarity

### 🎛️ Advanced Audio Features
- **7-Band Equalizer**: 63Hz, 160Hz, 400Hz, 1kHz, 2.5kHz, 6.3kHz, 16kHz
- **Bass/Treble Control**: -6 to +6 dB adjustment
- **Extra Bass**: Enhanced bass reproduction
- **Compressed Music Enhancer**: Improve compressed audio quality
- **Pure Direct Mode**: Bypass tone circuits for purest sound
- **Virtual Presence Speaker**: Virtual surround effect

### 🏠 Speaker Configuration
- **YPAO Integration**: Display YPAO-calibrated settings
- **Subwoofer Level**: Current YPAO settings display
- **Center/Surround Levels**: Individual speaker level monitoring
- **Dynamic Range Control**: DRC status display
- **Lip Sync Delay**: YPAO-optimized audio delay

### 🎬 Scene Control
- **4 Scene Presets**: Quick access to saved configurations
- **One-Click Activation**: Instant scene switching
- **Custom Configurations**: Save personal settings

### 📊 System Information
- **Firmware Version**: Current system firmware display
- **Temperature Monitoring**: Internal temperature tracking
- **Signal Format**: Active signal information
- **Channel Configuration**: Active channel display
- **Network Details**: IP, MAC, Gateway, signal strength

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/pepperonas/yahama-controller.git
cd yahama-controller

# Install dependencies
npm install

# Start the server
npm start

# Access the application
# Browser: http://localhost:5001
# Network: http://[YOUR-IP]:5001
```

## 📦 Installation

### Prerequisites
- Node.js 18.x or higher
- npm or yarn package manager
- Network access to your Yamaha receiver

### Production Deployment with PM2

```bash
# Install PM2 globally
npm install pm2 -g

# Start the application
pm2 start server.js --name yamaha-controller

# Enable auto-start on system boot
pm2 startup
pm2 save

# View logs
pm2 logs yamaha-controller
```

### Docker Deployment

```dockerfile
# Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 5001
CMD ["node", "server.js"]
```

```bash
# Build and run
docker build -t yamaha-controller .
docker run -d -p 5001:5001 --name yamaha-ctrl yamaha-controller
```

## 🎮 Usage

### Initial Setup

1. **Find Your Receiver's IP Address**:
   - Check your router's admin panel for connected devices
   - Look for "Yamaha" or "RX-V577" in the device list
   - Or use the receiver's network menu to display IP

2. **Connect to Receiver**:
   - Enter the receiver's IP address in the connection panel
   - Click "Connect" to establish connection
   - IP address is saved for future sessions

### Basic Operations

#### Power Control
- Click the round power button to toggle on/off
- Status indicator shows current state (On/Standby)

#### Volume Control
- Use the volume slider for adjustment (-80 dB to +16 dB)
- Click +/- buttons for precise changes
- Toggle mute for quick silence

#### Input Selection
- Click any input button to switch sources
- Active source is highlighted in blue
- Supports HDMI, AV, Audio, AirPlay, Server, USB, and Tuner inputs

#### Zone Control
- Switch between Main Zone and Zone 2 using tabs
- Each zone has independent controls

#### Scene Selection
- Click Scene 1-4 buttons to activate preconfigured settings
- Scenes combine input selection and DSP settings

## 🔧 Technical Details

### Architecture

```
yamaha-controller/
├── public/                 # Static assets
│   ├── assets/            # Images and mockups
│   ├── icons/             # App icons and favicons
│   ├── manifest.json      # PWA manifest
│   └── service-worker.js  # Offline functionality
├── src/                   # Source code (if restructured)
├── index.html             # Main interface
├── server.js              # Express.js server
├── package.json           # Dependencies
├── receiver-config.json   # Saved configuration
└── README.md             # Documentation
```

### Technology Stack

- **Backend**: Node.js with Express.js
- **Frontend**: Vanilla JavaScript with modern ES6+
- **Communication**: XML-based Yamaha protocol over HTTP
- **Process Manager**: PM2 for production
- **Styling**: CSS3 with CSS Variables for theming

### Network Requirements

- Receiver and controller must be on the same network
- HTTP requests to receiver IP on port 80
- No authentication required for local network access
- CORS handled by Express proxy server

## 📚 API Documentation

### XML Command Structure

```xml
<!-- Power Control -->
<YAMAHA_AV cmd="PUT">
  <Main_Zone>
    <Power_Control>
      <Power>On</Power>
    </Power_Control>
  </Main_Zone>
</YAMAHA_AV>

<!-- Volume Adjustment -->
<YAMAHA_AV cmd="PUT">
  <Main_Zone>
    <Volume>
      <Lvl>
        <Val>-200</Val>
        <Exp>1</Exp>
        <Unit>dB</Unit>
      </Lvl>
    </Volume>
  </Main_Zone>
</YAMAHA_AV>

<!-- Input Selection -->
<YAMAHA_AV cmd="PUT">
  <Main_Zone>
    <Input>
      <Input_Sel>HDMI1</Input_Sel>
    </Input>
  </Main_Zone>
</YAMAHA_AV>

<!-- Status Query -->
<YAMAHA_AV cmd="GET">
  <Main_Zone>
    <Basic_Status>GetParam</Basic_Status>
  </Main_Zone>
</YAMAHA_AV>
```

### Server Endpoints

- `GET /` - Serve main application
- `GET /api/health` - Health check endpoint
- `POST /api/receiver/*` - Proxy to receiver
- `GET /api/receiver-ip` - Get saved IP
- `POST /api/receiver-ip` - Save receiver IP

## 🐛 Troubleshooting

### Connection Issues

**Problem**: "Connection failed" error
- Verify receiver IP address is correct
- Ensure receiver is powered on and network-connected
- Check firewall settings aren't blocking connections

**Problem**: Cannot power on via network
- RX-V577 may not respond to network commands in standby over Wi-Fi
- Use Ethernet connection for reliable network wake
- Physical power button or IR remote may be needed for initial power-on

### Status Not Updating

- Check network connectivity
- Verify receiver is powered on
- Status polling occurs every 5 seconds automatically
- Check browser console for error messages

### CORS Errors

If accessing directly via file:// protocol:

```bash
# Option 1: Use the Node.js server (recommended)
npm start

# Option 2: Python HTTP server
python -m http.server 8000

# Option 3: Browser with disabled security (testing only)
chrome --disable-web-security --user-data-dir=/tmp/chrome_dev
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development Setup

```bash
# Install development dependencies
npm install --save-dev

# Run with auto-reload
npm run dev

# Run tests (if available)
npm test

# Lint code
npm run lint
```

## 📄 License

MIT License - Copyright (c) 2025 Martin Pfeffer

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Yamaha and RX-V577 are trademarks of Yamaha Corporation. This is an unofficial, open-source implementation for personal use.

---

<div align="center">
  <sub>Built with ❤️ in Berlin</sub>
</div>