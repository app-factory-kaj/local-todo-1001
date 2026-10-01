import createClient from "openapi-fetch";
import type { paths } from "./generated/todo-api";

// Same-origin: nginx's /api location proxies to the todo-api sibling (see
// nginx/default.conf and nginx/15-aep-api-proxy.sh). Never the sibling's
// injected address directly — that is pod env for nginx, not a browser key.
export const todoApi = createClient<paths>({ baseUrl: "/api" });
