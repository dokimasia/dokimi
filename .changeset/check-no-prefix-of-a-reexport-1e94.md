---
"go.dokimi.dev/dokimi": patch
---

Fix `sentinel` to leave the check of the prefix of a re-exported error value of another package to the package that declares it.
