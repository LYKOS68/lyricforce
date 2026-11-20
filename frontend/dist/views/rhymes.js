// Rhyme archive: add, search, filter by vowel class, list, delete,
// plus the rhyme canvas below the list.
import { esc, errorText, pageHeader } from "../components.js";
import { createCanvas } from "./canvas.js";

const VOWELS = ["a", "e", "i", "o", "u"];

// The page fills the window height so the list and canvas can share it.
export const fill = true;

export function render() {
  return `
    ${pageHeader("Reimarchiv", `<span class="faint mono" id="total"></span>`)}

    <form class="add" id="add" autocomplete="off">
      <input class="input" name="text" placeholder="Neues Wort" aria-label="Wort">
      <div class="seg" role="radiogroup" aria-label="Vokal">
        ${VOWELS.map((v) => `
          <label title="Vokal ${v.toUpperCase()}"><input type="radio" name="vowel" value="${v}"><span>${v.toUpperCase()}</span></label>`).join("")}
      </div>
      <input class="input" name="rhyme" placeholder="Reimklang ab betontem Vokal, z. B. acht" aria-label="Reimklang">
      <button class="btn primary" type="submit">Hinzufügen</button>
    </form>
    <p class="msg" id="msg" role="status"></p>

    <div class="toolbar">
      <input class="input search" id="q" type="search" placeholder="Suchen – Wort oder Reimklang" aria-label="Suchen">
      <div class="seg filter" id="filter" aria-label="Vokal-Filter"></div>
    </div>

    <div class="list-pane">
      <table class="table">
        <thead><tr><th>Wort</th><th>Reimklang</th><th>Vokal</th><th></th></tr></thead>
        <tbody id="rows"></tbody>
      </table>
    </div>
    <div class="splitter" id="splitter" role="separator" aria-orientation="horizontal" title="Höhe des Canvas ändern"></div>
    <div class="canvas" id="canvas"></div>`;
}

export function mount(root, { api }) {
  const $ = (sel) => root.querySelector(sel);
  const form = $("#add");
  const filter = { query: "", vowel: "" };
  let searchTimer;
  const canvas = createCanvas($("#canvas"), $("#splitter"), api);

  function setMsg(text, isError = false) {
    const el = $("#msg");
    el.textContent = text;
    el.classList.toggle("err", isError);
  }

  function renderFilter(counts) {
    const total = VOWELS.reduce((n, v) => n + (counts[v] ?? 0), 0);
    $("#total").textContent = `${total} ${total === 1 ? "Wort" : "Wörter"}`;
    const btn = (value, label, n) => `
      <button type="button" data-vowel="${value}" class="${filter.vowel === value ? "on" : ""}">
        ${label} <span class="n">${n}</span></button>`;
    $("#filter").innerHTML = btn("", "Alle", total) + VOWELS.map((v) => btn(v, v.toUpperCase(), counts[v] ?? 0)).join("");
    return total;
  }

  function renderRows(words, total) {
    const rows = $("#rows");
    if (!words.length) {
      const text = total === 0 ? "Archiv ist leer – erstes Wort oben hinzufügen." : "Keine Treffer.";
      rows.innerHTML = `<tr><td colspan="4" class="empty">${text}</td></tr>`;
      return;
    }
    let prevRhyme = null;
    rows.innerHTML = words.map((w) => {
      // A rule between rhyme groups makes groups visible without extra UI.
      const groupStart = prevRhyme !== null && w.rhyme !== prevRhyme;
      prevRhyme = w.rhyme;
      return `
        <tr class="${groupStart ? "group" : ""}">
          <td class="word">${esc(w.text)}</td>
          <td class="mono">${w.rhyme ? esc(w.rhyme) : `<span class="faint">–</span>`}</td>
          <td class="mono">${w.vowel.toUpperCase()}</td>
          <td class="actions"><button class="del" data-id="${w.id}" title="Löschen">×</button></td>
        </tr>`;
    }).join("");
  }

  async function refresh() {
    try {
      const [words, counts] = await Promise.all([api.listWords(filter), api.countWordsByVowel()]);
      renderRows(words, renderFilter(counts));
    } catch (err) {
      setMsg(errorText(err), true);
    }
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const text = form.elements.text.value;
    const vowel = form.elements.vowel.value;
    try {
      const w = await api.addWord(text, vowel, form.elements.rhyme.value);
      setMsg(`„${w.text}“ hinzugefügt.`);
      form.elements.text.value = "";
      form.elements.rhyme.value = "";
      form.elements.text.focus();
      refresh();
      canvas.refresh();
    } catch (err) {
      setMsg(errorText(err), true);
    }
  });

  $("#q").addEventListener("input", (e) => {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => { filter.query = e.target.value; refresh(); }, 120);
  });

  $("#filter").addEventListener("click", (e) => {
    const b = e.target.closest("button[data-vowel]");
    if (!b) return;
    filter.vowel = b.dataset.vowel;
    // Preselect the filtered vowel for the next entry.
    if (filter.vowel) form.querySelector(`input[name=vowel][value=${filter.vowel}]`).checked = true;
    refresh();
  });

  // Delete needs two clicks: the first arms the button, the second deletes.
  $("#rows").addEventListener("click", async (e) => {
    const b = e.target.closest("button.del");
    if (!b) return;
    if (!b.classList.contains("armed")) {
      root.querySelectorAll("button.del.armed").forEach((x) => { x.classList.remove("armed"); x.textContent = "×"; });
      b.classList.add("armed");
      b.textContent = "Löschen?";
      return;
    }
    try {
      await api.deleteWord(Number(b.dataset.id));
      setMsg("");
      refresh();
      canvas.refresh();
    } catch (err) {
      setMsg(errorText(err), true);
    }
  });

  form.elements.text.focus();
  refresh();
}
