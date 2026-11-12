import { newMemory } from "./memory";
import { runStoreSuite, storeChecks } from "./store.suite";

runStoreSuite(newMemory, {
  subject: "memory.Memory",
  excuse: {
    [storeChecks.size.smoke()]: "size reads no input",
  },
});
