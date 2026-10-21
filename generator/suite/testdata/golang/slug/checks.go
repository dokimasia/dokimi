package slug

// KeepsTheWords checks that Slug keeps each word of the title.
//
//dokimi:check Slug defect=dropsWords
func KeepsTheWords(h *SlugHandle) {}

// KeepsTheOrder checks that Slug keeps the words of the title in their order.
// Its callable has the name of the package.
//
//dokimi:check slug.Slug defect=dropsWords
func KeepsTheOrder(h *SlugHandle) {}

// dropsWords returns a Slug that drops the words of the title.
func dropsWords(*SlugHandle) SlugFunc { return nil }
