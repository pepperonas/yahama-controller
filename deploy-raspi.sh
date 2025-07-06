#!/bin/bash

# Yamaha Receiver App - Raspberry Pi Deployment Script
# This script deploys the Yamaha Receiver Control App to a Raspberry Pi

set -e  # Exit on any error

echo "🎵 Yamaha Receiver App - Raspberry Pi Deployment"
echo "================================================="

# Configuration
APP_NAME="yamaha-receiver"
APP_DIR="$HOME/yamaha-amp"
PORT=5001
PM2_APP_NAME="yamaha-amp"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_step() {
    echo -e "${BLUE}📋 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running on Raspberry Pi
print_step "Checking system..."
if ! grep -q "Raspberry Pi\|BCM" /proc/cpuinfo 2>/dev/null; then
    print_warning "This doesn't appear to be a Raspberry Pi, but continuing anyway..."
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed!"
    echo "Please install Node.js first:"
    echo "curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -"
    echo "sudo apt-get install -y nodejs"
    exit 1
fi

NODE_VERSION=$(node --version)
print_success "Node.js found: $NODE_VERSION"

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed!"
    exit 1
fi

NPM_VERSION=$(npm --version)
print_success "npm found: $NPM_VERSION"

# Check if PM2 is installed globally
if ! command -v pm2 &> /dev/null; then
    print_step "Installing PM2 globally..."
    npm install -g pm2
    print_success "PM2 installed"
else
    print_success "PM2 already installed"
fi

# Create app directory
print_step "Creating application directory: $APP_DIR"
mkdir -p "$APP_DIR"
cd "$APP_DIR"

# Copy application files
print_step "Copying application files..."

# Create package.json
cat > package.json << EOF
{
  "name": "yamaha-receiver-control",
  "version": "1.0.0",
  "description": "Web application for controlling Yamaha RX-V577 and compatible receivers",
  "main": "server.js",
  "scripts": {
    "start": "node server.js",
    "dev": "nodemon server.js",
    "pm2:start": "pm2 start ecosystem.config.js",
    "pm2:stop": "pm2 stop $PM2_APP_NAME",
    "pm2:restart": "pm2 restart $PM2_APP_NAME",
    "pm2:logs": "pm2 logs $PM2_APP_NAME"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "http-proxy-middleware": "^2.0.6"
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  },
  "keywords": [
    "yamaha",
    "receiver",
    "control",
    "web",
    "app",
    "raspberry-pi"
  ],
  "author": "",
  "license": "MIT"
}
EOF

# Create server.js with custom port
cat > server.js << 'EOF'
const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 5001;

// Enable CORS for all routes
app.use(cors());

// Serve static files from current directory
app.use(express.static(path.join(__dirname)));

// Store receiver IP (in production, this could be in a database)
let receiverIP = null;

// Middleware to parse JSON bodies
app.use(express.json());

// Endpoint to set receiver IP
app.post('/api/set-receiver-ip', (req, res) => {
    const { ip } = req.body;
    if (!ip) {
        return res.status(400).json({ error: 'IP address is required' });
    }
    
    // Basic IP validation
    const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    if (!ipRegex.test(ip)) {
        return res.status(400).json({ error: 'Invalid IP address format' });
    }
    
    receiverIP = ip;
    console.log(`Receiver IP set to: ${receiverIP}`);
    res.json({ success: true, ip: receiverIP });
});

// Proxy middleware for Yamaha receiver
const receiverProxy = createProxyMiddleware({
    target: 'http://placeholder', // Will be dynamically set
    changeOrigin: true,
    pathRewrite: {
        '^/api/receiver': '' // Remove /api/receiver prefix
    },
    router: (req) => {
        if (!receiverIP) {
            return null;
        }
        return `http://${receiverIP}`;
    },
    onError: (err, req, res) => {
        console.error('Proxy error:', err.message);
        res.status(500).json({ 
            error: 'Failed to connect to receiver', 
            message: err.message 
        });
    },
    onProxyReq: (proxyReq, req, res) => {
        console.log(`Proxying request to: http://${receiverIP}${req.url}`);
    }
});

// Use the proxy for receiver requests
app.use('/api/receiver', receiverProxy);

// Serve the main HTML file
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'index-advanced.html'));
});

// Serve advanced version explicitly
app.get('/advanced', (req, res) => {
    res.sendFile(path.join(__dirname, 'index-advanced.html'));
});

// Serve basic version
app.get('/basic', (req, res) => {
    res.sendFile(path.join(__dirname, 'index.html'));
});

// Health check endpoint
app.get('/api/health', (req, res) => {
    res.json({ 
        status: 'OK', 
        receiverIP: receiverIP,
        timestamp: new Date().toISOString(),
        version: '1.0.0',
        platform: 'Raspberry Pi'
    });
});

// Get system info
app.get('/api/system', (req, res) => {
    const os = require('os');
    res.json({
        hostname: os.hostname(),
        platform: os.platform(),
        arch: os.arch(),
        uptime: os.uptime(),
        loadavg: os.loadavg(),
        freemem: os.freemem(),
        totalmem: os.totalmem(),
        nodeVersion: process.version
    });
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`🎵 Yamaha Receiver Control Server running on port ${PORT}`);
    console.log(`🌐 Access URLs:`);
    console.log(`   Local:    http://localhost:${PORT}`);
    console.log(`   Network:  http://${getLocalIP()}:${PORT}`);
    console.log(`   Advanced: http://${getLocalIP()}:${PORT}/advanced`);
    console.log(`   Basic:    http://${getLocalIP()}:${PORT}/basic`);
});

function getLocalIP() {
    const os = require('os');
    const networkInterfaces = os.networkInterfaces();
    
    for (const interfaceName in networkInterfaces) {
        const interfaces = networkInterfaces[interfaceName];
        for (const iface of interfaces) {
            if (iface.family === 'IPv4' && !iface.internal) {
                return iface.address;
            }
        }
    }
    return '0.0.0.0';
}
EOF

# Create PM2 ecosystem configuration
cat > ecosystem.config.js << EOF
module.exports = {
  apps: [{
    name: '$PM2_APP_NAME',
    script: './server.js',
    instances: 1,
    autorestart: true,
    watch: false,
    max_memory_restart: '1G',
    env: {
      NODE_ENV: 'production',
      PORT: $PORT
    },
    error_file: './logs/err.log',
    out_file: './logs/out.log',
    log_file: './logs/combined.log',
    time: true
  }]
};
EOF

# Create logs directory
mkdir -p logs

print_success "Application files created"

# Install dependencies
print_step "Installing Node.js dependencies..."
npm install
print_success "Dependencies installed"

# Copy HTML files (if they exist in the source directory)
SOURCE_DIR="/Users/martin/WebstormProjects/mrx3k1/yahama-amp"

if [ -f "$SOURCE_DIR/index.html" ]; then
    print_step "Copying HTML files from source..."
    cp "$SOURCE_DIR/index.html" .
    if [ -f "$SOURCE_DIR/index-advanced.html" ]; then
        cp "$SOURCE_DIR/index-advanced.html" .
    fi
    print_success "HTML files copied"
else
    print_warning "Source HTML files not found, creating basic HTML file..."
    
    # Create a basic HTML file if source doesn't exist
    cat > index.html << 'EOF'
<!DOCTYPE html>
<html lang="de">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Yamaha Receiver Control</title>
</head>
<body>
    <h1>Yamaha Receiver Control</h1>
    <p>Basic version - Please copy the full HTML files from your development environment.</p>
    <p><a href="/advanced">Advanced Version</a> | <a href="/basic">Basic Version</a></p>
</body>
</html>
EOF

    cp index.html index-advanced.html
fi

# Stop existing PM2 process if running
print_step "Stopping existing PM2 process (if any)..."
pm2 stop $PM2_APP_NAME 2>/dev/null || true
pm2 delete $PM2_APP_NAME 2>/dev/null || true

# Start with PM2
print_step "Starting application with PM2..."
pm2 start ecosystem.config.js
pm2 save
print_success "Application started with PM2"

# Enable PM2 startup script
print_step "Setting up PM2 auto-startup..."
pm2 startup systemd -u $USER --hp $HOME
print_success "PM2 startup configured"

# Create systemctl service for easier management
print_step "Creating systemctl service..."
sudo tee /etc/systemd/system/yamaha-amp.service > /dev/null << EOF
[Unit]
Description=Yamaha Receiver Control App
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$APP_DIR
ExecStart=/usr/bin/pm2 start ecosystem.config.js --no-daemon
ExecStop=/usr/bin/pm2 stop $PM2_APP_NAME
ExecReload=/usr/bin/pm2 restart $PM2_APP_NAME
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable yamaha-amp.service
print_success "Systemctl service created and enabled"

# Create update script
cat > update.sh << 'EOF'
#!/bin/bash
echo "🔄 Updating Yamaha Receiver App..."
cd ~/yamaha-amp
npm update
pm2 restart yamaha-amp
echo "✅ Update complete!"
EOF

chmod +x update.sh

# Create backup script
cat > backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="$HOME/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="yamaha-amp_backup_$TIMESTAMP.tar.gz"

echo "📦 Creating backup..."
mkdir -p "$BACKUP_DIR"
tar -czf "$BACKUP_DIR/$BACKUP_FILE" -C "$HOME" yamaha-amp
echo "✅ Backup created: $BACKUP_DIR/$BACKUP_FILE"
EOF

chmod +x backup.sh

# Get local IP
LOCAL_IP=$(hostname -I | awk '{print $1}')

print_success "Deployment completed successfully!"
echo ""
echo "🎵 Yamaha Receiver Control App is now running!"
echo "================================================="
echo -e "${GREEN}📍 Access URLs:${NC}"
echo "   Local:    http://localhost:$PORT"
echo "   Network:  http://$LOCAL_IP:$PORT"
echo "   Advanced: http://$LOCAL_IP:$PORT/advanced"
echo "   Basic:    http://$LOCAL_IP:$PORT/basic"
echo ""
echo -e "${BLUE}🔧 Management Commands:${NC}"
echo "   Status:   pm2 status $PM2_APP_NAME"
echo "   Logs:     pm2 logs $PM2_APP_NAME"
echo "   Restart:  pm2 restart $PM2_APP_NAME"
echo "   Stop:     pm2 stop $PM2_APP_NAME"
echo "   Update:   ./update.sh"
echo "   Backup:   ./backup.sh"
echo ""
echo -e "${BLUE}🗂️  Files created in: $APP_DIR${NC}"
echo ""
echo -e "${GREEN}✅ Setup complete! Your Yamaha Receiver Control App is ready to use.${NC}"
echo -e "${YELLOW}📝 Don't forget to copy your HTML files if they weren't found automatically.${NC}"