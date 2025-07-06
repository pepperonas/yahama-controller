# Yamaha Receiver Control

A web application for controlling Yamaha RX-V577 and other compatible network-enabled Yamaha receivers through their XML control protocol.

## Features

### Basic Controls
- **Power Control**: Turn receiver on/off or switch to standby
- **Volume Control**: Adjust volume with slider or buttons, mute/unmute
- **Input Selection**: Switch between HDMI, AV, Audio, and network sources
- **Zone Control**: Support for Main Zone and Zone 2
- **Scene Selection**: Activate preconfigured scenes (1-4)

### Advanced Audio Controls
- **DSP Programs**: 15+ sound programs (Movie, Music, Game, Concert Hall, etc.)
- **Dialogue Level**: Enhance dialogue clarity (-6 to +6 dB)
- **7-Band Equalizer**: Precise frequency control (63Hz to 16kHz)
- **Speaker Configuration**: Individual level control for subwoofer, center, and surround speakers
- **Dynamic Range Control**: Standard/Maximum DRC settings
- **Lip Sync Delay**: Audio delay compensation (0-250ms)

### HDMI & Digital Audio
- **HDMI Audio Format**: PCM, DTS, Dolby Digital, Bitstream
- **HDMI Control**: CEC control on/off

### User Interface
- **Modern Dark Theme**: Material Design inspired with CSS variables
- **Real-time Statistics**: Live display of volume, input, DSP mode, and zone
- **Responsive Design**: Optimized for desktop, tablet, and mobile
- **German Localization**: Complete German interface
- **Status Polling**: Auto-refresh every 5 seconds

## Setup

### Method 1: Using Node.js Server (Recommended)

1. **Install Dependencies**:
   ```bash
   cd yahama-amp
   npm install
   ```

2. **Start the Server**:
   ```bash
   npm start
   ```
   Or for development with auto-reload:
   ```bash
   npm run dev
   ```

3. **Open the Application**:
   - Open your browser and go to `http://localhost:3000`
   - The server handles CORS issues and proxies requests to your receiver

4. **Find Your Receiver's IP Address**:
   - Check your router's admin panel for connected devices
   - Look for "Yamaha" or "RX-V577" in the device list
   - Or use the receiver's network menu to display the IP address

5. **Connect to Receiver**:
   - Enter the receiver's IP address in the connection panel
   - Click "Connect" to establish connection
   - The IP address will be saved for future use

### Method 2: Direct Browser Access (May Have CORS Issues)

1. **Basic Version**:
   - Open `index.html` in your web browser (basic controls)
   - No installation required, but may encounter CORS restrictions

2. **Advanced Version**:
   - Open `index-advanced.html` for full feature set
   - Requires Node.js server for CORS proxy

3. **If CORS Issues Occur**:
   - Use Method 1 (Node.js server) instead
   - Or run a local web server: `python -m http.server 8000`

## File Structure

- `index.html` - Basic receiver control interface
- `index-advanced.html` - Advanced control with full feature set
- `server.js` - Node.js CORS proxy server
- `package.json` - Node.js dependencies
- `README.md` - This documentation

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

## License

This project is provided as-is for educational and personal use. Yamaha and RX-V577 are trademarks of Yamaha Corporation.