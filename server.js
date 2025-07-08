const express = require('express');
const cors = require('cors');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');
const fs = require('fs');

const app = express();
const PORT = process.env.PORT || 5001;

// Enable CORS for all routes
app.use(cors());

// Serve static files from current directory with proper headers for PWA
app.use(express.static(path.join(__dirname), {
    setHeaders: (res, path) => {
        // Set appropriate cache headers for PWA files
        if (path.endsWith('manifest.json')) {
            res.setHeader('Content-Type', 'application/manifest+json');
            res.setHeader('Cache-Control', 'no-cache');
        }
        if (path.endsWith('sw.js')) {
            res.setHeader('Content-Type', 'application/javascript');
            res.setHeader('Cache-Control', 'no-cache');
        }
        // Cache static assets for PWA
        if (path.endsWith('.png') || path.endsWith('.ico')) {
            res.setHeader('Cache-Control', 'public, max-age=31536000'); // 1 year
        }
    }
}));

// Config file path
const CONFIG_FILE = path.join(__dirname, 'receiver-config.json');

// Load receiver IP from file
function loadReceiverIP() {
    try {
        if (fs.existsSync(CONFIG_FILE)) {
            const config = JSON.parse(fs.readFileSync(CONFIG_FILE, 'utf8'));
            return config.receiverIP || null;
        }
    } catch (error) {
        console.error('Error loading receiver IP:', error);
    }
    return null;
}

// Save receiver IP to file
function saveReceiverIP(ip) {
    try {
        const config = { receiverIP: ip };
        fs.writeFileSync(CONFIG_FILE, JSON.stringify(config, null, 2));
        console.log(`Receiver IP saved to file: ${ip}`);
    } catch (error) {
        console.error('Error saving receiver IP:', error);
    }
}

// Store receiver IP (loaded from file)
let receiverIP = loadReceiverIP();

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
    saveReceiverIP(ip);
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
    res.sendFile(path.join(__dirname, 'index.html'));
});

// Health check endpoint
app.get('/api/health', (req, res) => {
    res.json({ 
        status: 'OK', 
        receiverIP: receiverIP,
        timestamp: new Date().toISOString()
    });
});

app.listen(PORT, '0.0.0.0', () => {
    console.log(`Yamaha Receiver Control Server running on port ${PORT}`);
    console.log(`Open http://localhost:${PORT} in your browser`);
    if (receiverIP) {
        console.log(`Loaded receiver IP from config: ${receiverIP}`);
    }
});