/* Vertrags-Pins der Yamaha-Oberflaeche. Ausfuehren: node --test tests/
 *
 * Erste Suite dieses Repos — angelegt mit der Phase-3-Adoption der geteilten
 * UI-Schicht (2026-08-30). Pins auf Textebene, kommentarfrei wo Abwesenheit
 * geprueft wird (die Doku dieses Hauses zitiert entfernte Muster woertlich).
 */
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const HERE = dirname(fileURLToPath(import.meta.url));
const HTML = readFileSync(join(HERE, '..', 'index.html'), 'utf8');
const CSS = [...HTML.matchAll(/<style>([\s\S]*?)<\/style>/g)].map(m => m[1]).join('\n');
const CSS_PUR = CSS.replace(/\/\*[\s\S]*?\*\//g, '');
const JS_PUR = HTML.replace(/\/\*[\s\S]*?\*\//g, '').replace(/^\s*\/\/[^\n]*$/gm, '');

// ---- Buttons sprechen den Suite-Dialekt -------------------------------------

test('die 8 Quellen-Buttons sind sh-pill-Auswahl-Chips', () => {
  const inputs = [...HTML.matchAll(/<button[^>]*data-input="[^"]+"[^>]*>/g)];
  assert.equal(inputs.length, 8);
  for (const m of inputs) assert.match(m[0], /class="[^"]*sh-pill/, m[0]);
});

test('der aktive Quellen-Chip wird ueber .active gefuellt (JS-Sync)', () => {
  // Das alte Paar btn-primary/btn-outline bleibt (Altvertraege), aber die
  // sh-pill-Fuellung haengt an .active — ohne diese Zeile bliebe die Auswahl
  // unsichtbar:
  assert.match(JS_PUR, /btn\.classList\.toggle\('active',\s*isActive\)/);
});

test('Aktions-Buttons tragen sh-btn mit den richtigen Varianten', () => {
  assert.match(HTML, /id="vol-down"[^>]*|class="[^"]*sh-btn tonal sm[^"]*"[^>]*id="vol-down"/);
  for (const id of ['vol-down', 'vol-up']) {
    const m = HTML.match(new RegExp('<button[^>]*id="' + id + '"[^>]*>'));
    assert.match(m[0], /sh-btn tonal sm/, id);
  }
  const scenes = [...HTML.matchAll(/<button[^>]*data-scene="\d"[^>]*>/g)];
  assert.equal(scenes.length, 4);
  for (const m of scenes) assert.match(m[0], /sh-btn outline/, m[0]);
  assert.match(HTML.match(/<button[^>]*id="eq-reset"[^>]*>/)[0], /sh-btn tonal/);
  assert.match(HTML.match(/<button[^>]*id="connect-btn"[^>]*>/)[0], /sh-btn/);
});

// ---- Toggles: sh-switch mit klassengetriebener Brueckenlogik ----------------

test('alle 5 Toggles sind sh-switch mit Wanne und Daumen', () => {
  const toggles = [...HTML.matchAll(/<div[^>]*class="[^"]*toggle-switch[^"]*"[^>]*>/g)];
  assert.equal(toggles.length, 5);
  for (const m of toggles) {
    assert.match(m[0], /sh-switch/, m[0]);
    assert.match(m[0], /role="switch"/, m[0]);
    assert.match(m[0], /tabindex="0"/, m[0]);
  }
  assert.equal([...HTML.matchAll(/sh-switch-track/g)].length >= 5, true);
});

test('die Bruecke synct aria-checked aus .active und macht Enter/Space zum Klick', () => {
  assert.match(JS_PUR, /aria-checked/);
  assert.match(JS_PUR, /MutationObserver[\s\S]{0,400}sh-switch|sh-switch[\s\S]{0,400}MutationObserver/);
  assert.match(JS_PUR, /key === 'Enter' \|\| e\.key === ' '|keydown/);
});

test('kein eigenes Toggle-Skinning mehr (::after-Daumen ist Geschichte)', () => {
  assert.ok(!CSS_PUR.includes('.toggle-switch::after'),
    'der ::after-Daumen kollidiert mit dem sh-switch-Span-Daumen');
});

// ---- Slider -----------------------------------------------------------------

test('die 4 horizontalen Slider tragen sh-slider; die 7 vertikalen EQ-Fader bleiben Identitaet', () => {
  // .eq-slider ist writing-mode:vertical-lr — der Dialekt-Slider ist horizontal
  // gebaut; vertikale Fader sind dokumentierte App-Identitaet (Masterprompt 0.3).
  const sliders = [...HTML.matchAll(/<input[^>]*type="range"[^>]*>/g)];
  assert.equal(sliders.length, 11);
  const horiz = sliders.filter(m => m[0].includes('class="slider'));
  assert.equal(horiz.length, 4);
  for (const m of horiz) assert.match(m[0], /sh-slider/, m[0]);
  for (const m of sliders.filter(x => x[0].includes('eq-slider')))
    assert.ok(!/sh-slider/.test(m[0]), 'EQ-Fader sind vertikal — kein sh-slider');
  assert.match(JS_PUR, /setProperty\('--sh-slider-fill'/);
  assert.ok(!/style\.background\s*=\s*'linear-gradient/.test(JS_PUR),
    'der alte Inline-Gradient wuerde die sh-Wanne uebermalen');
});

test('kein eigenes Slider-Skinning mehr', () => {
  assert.ok(!CSS_PUR.includes('.slider::-webkit-slider-thumb'));
  assert.ok(!CSS_PUR.includes('.slider::-moz-range-thumb'));
});

// ---- Tabs -------------------------------------------------------------------

test('Reiter sind sh-pills in sh-tabs-Containern, der tote Indikator bleibt versteckt', () => {
  assert.match(HTML, /id="main-tabs"[^>]*class="[^"]*sh-tabs|class="[^"]*sh-tabs[^"]*"[^>]*id="main-tabs"/);
  for (const m of [...HTML.matchAll(/<button[^>]*class="[^"]*\b(main-tab|tab)\b[^"]*"[^>]*data-(tab|zone)=/g)]) {
    assert.match(m[0], /sh-pill/, m[0]);
  }
  // Das JS schreibt weiter auf .tab-ind — die Regel haelt ihn unsichtbar:
  assert.match(CSS_PUR, /\.tab-ind\{\s*display:\s*none/);
});

// ---- Ripple: nur noch die geteilte Implementierung --------------------------

test('kein eigenes Ripple mehr (ui.js ist die eine Quelle)', () => {
  assert.ok(!CSS_PUR.includes('md-ripple'), 'md-ripple-CSS lebt noch');
  assert.ok(!JS_PUR.includes('md-ripple'), 'md-ripple-JS lebt noch');
});

// ---- Skin-Dedup -------------------------------------------------------------

test('die alten Button-Skins sind geloescht, das Layout-Skelett bleibt', () => {
  assert.ok(!/\.btn-primary\s*\{/.test(CSS_PUR), 'btn-primary-Skin lebt noch');
  assert.ok(!/\.btn-danger\s*\{/.test(CSS_PUR), 'btn-danger-Skin lebt noch');
  assert.ok(!/\.btn-outline\s*\{/.test(CSS_PUR), 'btn-outline-Skin lebt noch');
  // Skelett: Icon+Label brauchen flex+gap, das bringt sh-btn nicht mit.
  assert.match(CSS_PUR, /\.btn\s*\{[^}]*inline-flex/);
});
