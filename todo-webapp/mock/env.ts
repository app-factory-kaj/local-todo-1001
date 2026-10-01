// What window._env_ holds in mock mode. This app has no auth dependency and
// no `configurations.env` entries — its one dependency (todo-api) is reached
// same-origin at /api, never through a browser key — so there is nothing to
// carry here; src/env.ts only needs window._env_ to exist.
export const mockEnv = {};
