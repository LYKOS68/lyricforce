// Shared UI helpers.

export function esc(v) {
  return String(v ?? "").replace(/[&<>"']/g, (c) => (
    { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" }[c]
  ));
}

export function pageHeader(title, aside = "") {
  return `<header class="page-head"><h1>${esc(title)}</h1>${aside}</header>`;
}

// Maps stable backend error messages to user-facing text.
const ERROR_TEXT = {
  "word is empty": "Bitte ein Wort eingeben.",
  "invalid vowel": "Bitte einen Vokal wählen.",
  "word already exists": "Dieses Wort ist bereits im Archiv.",
  "word not found": "Wort nicht gefunden.",
  "database not available": "Datenbank nicht verfügbar.",
};

export function errorText(err) {
  const msg = String(err?.message ?? err ?? "");
  return ERROR_TEXT[msg] ?? msg;
}
