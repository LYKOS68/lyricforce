// Application shell: navigation and routing (#/<module-id>).
import { api } from "./api.js";
import { esc, pageHeader } from "./components.js";

import * as rhymes from "./views/rhymes.js";
import * as settings from "./views/settings.js";

// Views export render(ctx) → HTML and optionally mount(root, ctx) and fill (page uses full height).
// Modules without a view are planned and get a one-line page.
const VIEWS = { rhymes, settings };

const state = { modules: [], info: null };
const ctx = {
  api,
  get info() { return state.info; },
};

const plannedView = {
  render: (mod) => `${pageHeader(mod.title)}<p class="faint">Geplant – noch nicht verfügbar.</p>`,
};

function renderNav(activeId) {
  document.getElementById("nav").innerHTML = state.modules.map((m) => `
    <a class="nav-item${m.id === activeId ? " active" : ""}${m.status === "planned" ? " planned" : ""}" href="#/${m.id}">
      <span>${esc(m.title)}</span>
      ${m.status === "planned" ? `<span class="nav-tag">geplant</span>` : ""}
    </a>`).join("");
}

function route() {
  const id = location.hash.replace(/^#\//, "");
  const mod = state.modules.find((m) => m.id === id) ?? state.modules[0];
  renderNav(mod.id);
  const main = document.getElementById("main");
  const view = VIEWS[mod.id];
  main.innerHTML = `<div class="page${view?.fill ? " fill" : ""}">${view ? view.render(ctx) : plannedView.render(mod)}</div>`;
  main.scrollTop = 0;
  view?.mount?.(main, ctx);
}

async function boot() {
  try {
    [state.modules, state.info] = await Promise.all([api.modules(), api.appInfo()]);
  } catch (err) {
    document.getElementById("main").innerHTML =
      `<div class="page"><p class="err">Start fehlgeschlagen: ${esc(err.message)}</p></div>`;
    return;
  }
  window.addEventListener("hashchange", route);
  route();
}

boot();
