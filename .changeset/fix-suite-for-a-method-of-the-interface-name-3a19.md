---
"go.dokimi.dev/dokimi": patch
---

Fix `suite` to embed the subject type through an alias, so an interface with a method of its own name gets a suite that compiles.
