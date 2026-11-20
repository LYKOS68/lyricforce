// Settings: version, database and log locations.
import { esc, pageHeader } from "../components.js";

function formatBytes(n) {
  if (!n) return "0 B";
  const u = ["B", "KB", "MB", "GB"];
  const i = Math.min(Math.floor(Math.log(n) / Math.log(1024)), u.length - 1);
  return `${(n / 1024 ** i).toFixed(i ? 1 : 0)} ${u[i]}`;
}

export function render(ctx) {
  return `
    ${pageHeader("Einstellungen")}
    <dl class="kv" id="info">
      <dt>Version</dt><dd>${esc(ctx.info.version)}</dd>
      <dt>Datenbank</dt><dd id="db">…</dd>
      <dt>Logdatei</dt><dd>${esc(ctx.info.logPath || "–")}</dd>
    </dl>`;
}

export async function mount(root, { api }) {
  const el = root.querySelector("#db");
  try {
    const s = await api.databaseStatus();
    el.innerHTML = s.connected
      ? `${esc(s.path)} <span class="faint">· Schema v${s.schemaVersion} · ${formatBytes(s.sizeBytes)}</span>`
      : `<span class="err">nicht verbunden${s.error ? ` – ${esc(s.error)}` : ""}</span>`;
  } catch (err) {
    el.innerHTML = `<span class="err">${esc(err.message)}</span>`;
  }
}
