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

// Serve the advanced HTML file as main
app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'index-advanced.html'));
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
});