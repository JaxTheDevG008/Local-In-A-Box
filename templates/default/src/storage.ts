import Dexie from "dexie";
import yDexie from "y-dexie";
import type * as Y from "yjs";

export const db = new Dexie("{{PROJECT_NAME}}DB", {
    addons: [yDexie]
});

db.version(1).stores({
    kv: "key, data: Y.Doc",
});

console.log("Database initialized successfully.");