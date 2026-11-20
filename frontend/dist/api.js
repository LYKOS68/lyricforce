// Thin wrapper around the Go bindings (window.go.main.App).
// All backend access from the UI goes through this file.

function backend() {
  const app = window.go?.main?.App;
  if (!app) throw new Error("Go backend not connected");
  return app;
}

export const api = {
  appInfo: () => backend().GetAppInfo(),
  modules: () => backend().GetModules(),
  databaseStatus: () => backend().GetDatabaseStatus(),
  listWords: (filter) => backend().ListWords(filter),
  countWordsByVowel: () => backend().CountWordsByVowel(),
  addWord: (text, vowel, rhyme) => backend().AddWord(text, vowel, rhyme),
  deleteWord: (id) => backend().DeleteWord(id),
  getUIState: (key) => backend().GetUIState(key),
  setUIState: (key, value) => backend().SetUIState(key, value),
};
