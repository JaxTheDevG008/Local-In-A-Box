import Dexie from "dexie";
import yDexie from "y-dexie";

// Minimal storage implementation used by the template for TypeScript
// resolution. This mirrors the templated `storage.ts.ejs` behavior
// but keeps things static so the TS server can find `./storage`.
export const db = new Dexie("LIABDB", {
  addons: [yDexie]
});

db.version(1).stores({
  kv: "key, data"
});

console.log("Database shim initialized for templates/default/src/storage.ts");

export default db;
