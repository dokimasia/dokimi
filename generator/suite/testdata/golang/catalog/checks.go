package catalog

// FindsTheAddedProduct checks that Find returns the product that Add listed.
//
//dokimi:check Catalog.Find defect=amnesiac
func FindsTheAddedProduct(h *CatalogHandle) {}

// amnesiac returns a Catalog that lists no product.
func amnesiac(*CatalogHandle) CatalogMethods { return nil }
