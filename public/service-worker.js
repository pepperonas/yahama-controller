const CACHE_NAME = 'app-v3';
const urlsToCache = ['/', '/index.html'];

self.addEventListener('install', event => {
  self.skipWaiting();
  event.waitUntil(caches.open(CACHE_NAME).then(cache => cache.addAll(urlsToCache)));
});

self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(keys => Promise.all(
      keys.filter(k => k !== CACHE_NAME).map(k => caches.delete(k))
    )).then(() => self.clients.claim())
  );
});

self.addEventListener('fetch', event => {
  const req = event.request;
  if (req.url.includes('/api/')) return event.respondWith(fetch(req));

  // Network-first for the HTML shell so frontend updates show up immediately;
  // fall back to cache offline.
  const isShell = req.mode === 'navigate' ||
                  req.url.endsWith('/') || req.url.endsWith('/index.html');
  if (isShell) {
    event.respondWith(
      fetch(req).then(res => {
        const copy = res.clone();
        caches.open(CACHE_NAME).then(c => c.put(req, copy));
        return res;
      }).catch(() => caches.match(req).then(r => r || caches.match('/index.html')))
    );
    return;
  }

  // Cache-first for static assets (icons, etc.).
  event.respondWith(caches.match(req).then(r => r || fetch(req)));
});
