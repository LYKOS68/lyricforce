// Rhyme canvas: a pannable, zoomable surface (drag background, mouse wheel)
// with movable, resizable windows — one table per vowel class.
// Layout (view, pane height, windows) is persisted via the UI state store.
import { esc } from "../components.js";

const VOWELS = ["a", "e", "i", "o", "u"];
const STATE_KEY = "rhymes.canvas";
const ZOOM_MIN = 0.3;
const ZOOM_MAX = 2;
const GRID = 24; // background dot spacing at zoom 1
const WIN_DEFAULT = { w: 260, h: 320 };
const WIN_MIN = { w: 180, h: 140 };
const PANE_MIN = 160;

const clamp = (v, lo, hi) => Math.min(hi, Math.max(lo, v));
const num = (v, fallback) => (Number.isFinite(v) ? v : fallback);

function defaultState() {
  return { view: { x: 0, y: 0, zoom: 1 }, height: 340, windows: {} };
}

// Accepts a parsed state document and drops anything malformed.
function sanitize(raw) {
  const s = defaultState();
  if (!raw || typeof raw !== "object") return s;
  s.view.x = num(raw.view?.x, 0);
  s.view.y = num(raw.view?.y, 0);
  s.view.zoom = clamp(num(raw.view?.zoom, 1), ZOOM_MIN, ZOOM_MAX);
  s.height = Math.max(PANE_MIN, num(raw.height, s.height));
  for (const v of VOWELS) {
    const w = raw.windows?.[v];
    if (!w) continue;
    s.windows[v] = {
      x: num(w.x, 0), y: num(w.y, 0),
      w: Math.max(WIN_MIN.w, num(w.w, WIN_DEFAULT.w)),
      h: Math.max(WIN_MIN.h, num(w.h, WIN_DEFAULT.h)),
      z: num(w.z, 1),
    };
  }
  return s;
}

/**
 * Mounts the canvas into `el`; `splitter` resizes the pane height.
 * Returns { refresh } to reload window contents after archive changes.
 */
export function createCanvas(el, splitter, api) {
  let state = defaultState();
  let topZ = 1;
  let saveTimer;

  el.innerHTML = `
    <div class="cv-world" id="cv-world"></div>
    <div class="cv-toolbar">
      <div class="seg" id="cv-open" aria-label="Vokal-Fenster öffnen">
        ${VOWELS.map((v) => `<button type="button" data-open="${v}" title="Fenster ${v.toUpperCase()} öffnen">${v.toUpperCase()}</button>`).join("")}
      </div>
      <button type="button" class="btn small" id="cv-fit" title="Alle Fenster in den sichtbaren Bereich einpassen">Einpassen</button>
      <span class="cv-zoom mono" id="cv-zoom"></span>
    </div>
    <p class="cv-hint" id="cv-hint">Vokal oben links wählen, um ein Fenster zu öffnen. Hintergrund ziehen verschiebt, Mausrad zoomt.</p>`;
  const world = el.querySelector("#cv-world");

  // ---------- persistence ----------
  function save() {
    clearTimeout(saveTimer);
    saveTimer = setTimeout(() => {
      api.setUIState(STATE_KEY, JSON.stringify(state)).catch(() => {});
    }, 400);
  }

  // ---------- rendering ----------
  function applyView() {
    const { x, y, zoom } = state.view;
    world.style.transform = `translate(${x}px, ${y}px) scale(${zoom})`;
    el.style.backgroundSize = `${GRID * zoom}px ${GRID * zoom}px`;
    el.style.backgroundPosition = `${x}px ${y}px`;
    el.querySelector("#cv-zoom").textContent = `${Math.round(zoom * 100)} %`;
  }

  function applyHeight() {
    const max = Math.max(PANE_MIN, el.parentElement.clientHeight - 220);
    state.height = clamp(state.height, PANE_MIN, max);
    el.style.height = `${state.height}px`;
  }

  function updateChrome() {
    const open = Object.keys(state.windows);
    el.querySelectorAll("[data-open]").forEach((b) => b.classList.toggle("on", open.includes(b.dataset.open)));
    el.querySelector("#cv-hint").hidden = open.length > 0;
  }

  function winEl(v) {
    return world.querySelector(`.win[data-v="${v}"]`);
  }

  function position(v) {
    const w = state.windows[v];
    Object.assign(winEl(v).style, {
      left: `${w.x}px`, top: `${w.y}px`, width: `${w.w}px`, height: `${w.h}px`, zIndex: w.z,
    });
  }

  function renderWindow(v) {
    const div = document.createElement("div");
    div.className = "win";
    div.dataset.v = v;
    div.innerHTML = `
      <div class="win-head">
        <span class="win-letter">${v.toUpperCase()}</span>
        <span class="win-count mono"></span>
        <button type="button" class="win-close" title="Fenster schließen">×</button>
      </div>
      <div class="win-body"></div>
      <div class="win-resize" title="Größe ändern"></div>`;
    world.appendChild(div);
    position(v);
  }

  async function loadWindow(v) {
    const div = winEl(v);
    if (!div) return;
    try {
      const words = await api.listWords({ query: "", vowel: v });
      div.querySelector(".win-count").textContent = `${words.length} ${words.length === 1 ? "Wort" : "Wörter"}`;
      if (!words.length) {
        div.querySelector(".win-body").innerHTML = `<p class="win-empty">Keine Wörter mit ${v.toUpperCase()}.</p>`;
        return;
      }
      let prev = null;
      div.querySelector(".win-body").innerHTML = `<table class="win-table">${words.map((w) => {
        const group = prev !== null && w.rhyme !== prev;
        prev = w.rhyme;
        return `<tr class="${group ? "group" : ""}"><td>${esc(w.text)}</td><td class="mono">${w.rhyme ? esc(w.rhyme) : `<span class="faint">–</span>`}</td></tr>`;
      }).join("")}</table>`;
    } catch (err) {
      div.querySelector(".win-body").innerHTML = `<p class="win-empty err">${esc(err?.message ?? err)}</p>`;
    }
  }

  // ---------- window actions ----------
  function bringToFront(v) {
    const w = state.windows[v];
    if (w.z === topZ) return;
    w.z = ++topZ;
    position(v);
    save();
  }

  function centerOn(v) {
    const w = state.windows[v];
    const { zoom } = state.view;
    state.view.x = el.clientWidth / 2 - (w.x + w.w / 2) * zoom;
    state.view.y = el.clientHeight / 2 - (w.y + w.h / 2) * zoom;
    applyView();
  }

  function openWindow(v) {
    if (state.windows[v]) {
      bringToFront(v);
      centerOn(v);
      save();
      return;
    }
    // Place new windows in the visible area, staggered by vowel.
    const { x, y, zoom } = state.view;
    const i = VOWELS.indexOf(v);
    state.windows[v] = {
      x: Math.round(-x / zoom + 24 + i * (WIN_DEFAULT.w + 20)),
      y: Math.round(-y / zoom + 64),
      ...WIN_DEFAULT,
      z: ++topZ,
    };
    renderWindow(v);
    loadWindow(v);
    updateChrome();
    save();
  }

  function closeWindow(v) {
    winEl(v)?.remove();
    delete state.windows[v];
    updateChrome();
    save();
  }

  function fit() {
    const wins = Object.values(state.windows);
    if (!wins.length) {
      state.view = { x: 0, y: 0, zoom: 1 };
    } else {
      const pad = 32;
      const minX = Math.min(...wins.map((w) => w.x));
      const minY = Math.min(...wins.map((w) => w.y));
      const maxX = Math.max(...wins.map((w) => w.x + w.w));
      const maxY = Math.max(...wins.map((w) => w.y + w.h));
      const zoom = clamp(Math.min(
        (el.clientWidth - 2 * pad) / (maxX - minX),
        (el.clientHeight - 2 * pad - 40) / (maxY - minY),
      ), ZOOM_MIN, 1);
      state.view = {
        zoom,
        x: (el.clientWidth - (maxX - minX) * zoom) / 2 - minX * zoom,
        y: 40 + (el.clientHeight - 40 - (maxY - minY) * zoom) / 2 - minY * zoom,
      };
    }
    applyView();
    save();
  }

  function zoomAt(mx, my, factor) {
    const v = state.view;
    const zoom = clamp(v.zoom * factor, ZOOM_MIN, ZOOM_MAX);
    // Keep the point under the cursor fixed.
    v.x = mx - (mx - v.x) * (zoom / v.zoom);
    v.y = my - (my - v.y) * (zoom / v.zoom);
    v.zoom = zoom;
    applyView();
    save();
  }

  // ---------- pointer interaction ----------
  // Generic drag: calls onMove(dx, dy) with screen-pixel deltas since the last event.
  function drag(e, target, onMove, cls) {
    e.preventDefault();
    target.setPointerCapture(e.pointerId);
    let lx = e.clientX;
    let ly = e.clientY;
    if (cls) el.classList.add(cls);
    const move = (ev) => {
      onMove(ev.clientX - lx, ev.clientY - ly);
      lx = ev.clientX;
      ly = ev.clientY;
    };
    const up = () => {
      target.removeEventListener("pointermove", move);
      target.removeEventListener("pointerup", up);
      target.removeEventListener("pointercancel", up);
      if (cls) el.classList.remove(cls);
      save();
    };
    target.addEventListener("pointermove", move);
    target.addEventListener("pointerup", up);
    target.addEventListener("pointercancel", up);
  }

  el.addEventListener("pointerdown", (e) => {
    if (e.button !== 0 || e.target.closest(".cv-toolbar")) return;
    const div = e.target.closest(".win");
    if (!div) {
      drag(e, el, (dx, dy) => { state.view.x += dx; state.view.y += dy; applyView(); }, "panning");
      return;
    }
    const v = div.dataset.v;
    const w = state.windows[v];
    bringToFront(v);
    const z = () => state.view.zoom;
    if (e.target.closest(".win-resize")) {
      // Track the unclamped size so the corner stays under the cursor after hitting the minimum.
      let rw = w.w;
      let rh = w.h;
      drag(e, e.target, (dx, dy) => {
        rw += dx / z(); rh += dy / z();
        w.w = Math.round(Math.max(WIN_MIN.w, rw));
        w.h = Math.round(Math.max(WIN_MIN.h, rh));
        position(v);
      });
    } else if (e.target.closest(".win-head") && !e.target.closest(".win-close")) {
      drag(e, e.target.closest(".win-head"), (dx, dy) => {
        w.x += dx / z(); w.y += dy / z();
        position(v);
      }, "moving");
    }
  });

  el.addEventListener("wheel", (e) => {
    // Inside a scrollable window body the wheel scrolls the list; Ctrl+wheel always zooms.
    const body = e.target.closest(".win-body");
    if (body && !e.ctrlKey && body.scrollHeight > body.clientHeight) return;
    e.preventDefault();
    const r = el.getBoundingClientRect();
    zoomAt(e.clientX - r.left, e.clientY - r.top, Math.exp(-e.deltaY * 0.0015));
  }, { passive: false });

  el.addEventListener("click", (e) => {
    const open = e.target.closest("[data-open]");
    if (open) return openWindow(open.dataset.open);
    if (e.target.closest("#cv-fit")) return fit();
    const close = e.target.closest(".win-close");
    if (close) closeWindow(close.closest(".win").dataset.v);
  });

  splitter.addEventListener("pointerdown", (e) => {
    if (e.button !== 0) return;
    drag(e, splitter, (_, dy) => { state.height -= dy; applyHeight(); });
  });

  // ---------- init ----------
  async function init() {
    try {
      const raw = await api.getUIState(STATE_KEY);
      if (raw) state = sanitize(JSON.parse(raw));
    } catch {
      state = defaultState();
    }
    topZ = Math.max(1, ...Object.values(state.windows).map((w) => w.z));
    applyHeight();
    applyView();
    for (const v of Object.keys(state.windows)) renderWindow(v);
    updateChrome();
    refresh();
  }

  function refresh() {
    for (const v of Object.keys(state.windows)) loadWindow(v);
  }

  init();
  return { refresh };
}
