import { db } from "./storage";

async function main() {
    let record = await db.table("kv").get("appState");
    if (!record) {
        await db.table("kv").put({ key: "appState" });
        record = await db.table("kv").get("appState");
    }

    const doc = record?.data;
    const store = doc?.getMap("store");

    if (!store.has("status")) {
        store.set("status", "onlineOrOffline")
        store.set("launchedAt", new Date().toLocaleTimeString());
    }
}

main().catch(console.error);