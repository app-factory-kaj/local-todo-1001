// Typed read of window._env_, mounted by the platform at /env-config.js.
//
// This app has no auth dependency and no `configurations.env` entries, and its
// one dependency (todo-api) is a sibling `component`, reached same-origin at
// `/api` — not through a browser-visible URL. So there is nothing to declare
// here; the type is empty and the only job left is to prove the file loaded.
type Env = Record<string, never>;

declare global {
  interface Window {
    _env_: Env;
  }
}

if (!window._env_) {
  throw new Error(
    "window._env_ not set — /env-config.js failed to load. " +
      "The platform mounts this file; if you see this locally, host " +
      "/env-config.js from your dev server.",
  );
}

export const env: Env = window._env_;
