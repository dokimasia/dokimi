package bad

import (
	"testing"

	"go.dokimi.dev/assert/prop"

	"example.com/shop/store"
)

// AroundOfOne checks that Get returns the item inside an environment of one
// parameter.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfOne
func AroundOfOne(h any) {}

// AroundOfText checks that Get returns the item inside an environment of a
// text.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfText
func AroundOfText(h any) {}

// AroundOfACount checks that Get returns the item inside an environment of a
// count.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfACount
func AroundOfACount(h any) {}

// AroundOfACall checks that Get returns the item inside an environment whose
// run takes a count.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfACall
func AroundOfACall(h any) {}

// AroundOfAFailingRun checks that Get returns the item inside an environment
// whose run returns an error.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfAFailingRun
func AroundOfAFailingRun(h any) {}

// aroundOfACall takes a run of a count.
func aroundOfACall(*testing.T, func(int)) {}

// aroundOfAFailingRun takes a run that returns an error.
func aroundOfAFailingRun(*testing.T, func(*testing.T) error) {}

// AroundOfAnError checks that Get returns the item inside an environment that
// returns an error.
//
//dokimi:check store.Store.Get defect=broken around=aroundOfAnError
func AroundOfAnError(h any) {}

// aroundOfOne takes the test alone.
func aroundOfOne(*testing.T) {}

// aroundOfText takes a text in place of the test.
func aroundOfText(string, func(*testing.T)) {}

// aroundOfACount takes a count in place of the run.
func aroundOfACount(*testing.T, int) {}

// aroundOfAnError returns an error.
func aroundOfAnError(*testing.T, func(*testing.T)) error { return nil }

// Empty is a declaration under the directive suite without callables.
//
//dokimi:suite
type Empty interface{}

// Normalize is a function under the directive suite, which a check names
// without a method.
//
//dokimi:suite
func Normalize(text string) string { return text }

// TrimsTheText checks that Trim of Normalize trims the text.
//
//dokimi:check Normalize.Trim defect=broken
func TrimsTheText(h any) {}

// NamesNoCallable checks that Empty has a callable.
//
//dokimi:check Empty defect=broken
func NamesNoCallable(h any) {}

// GetsNothing checks that Get of a declaration without the directive suite
// finds nothing.
//
//dokimi:check Lister.List defect=broken
func GetsNothing(h any) {}

// FetchesTheItem checks that Fetch returns the item.
//
//dokimi:check store.Store.Fetch defect=broken
func FetchesTheItem(h any) {}

// KeepsAnItem checks that Get returns the item of one of two stores.
//
//dokimi:check Store.Get defect=broken
func KeepsAnItem(h any) {}

// KeepsTheName checks that Get of the twin store returns the name.
//
//dokimi:check twin.Store.Get defect=broken
func KeepsTheName(h any) {}

// PutsTheName checks that Put of the twin store keeps the name.
//
//dokimi:check twin.Store.Put defect=broken
func PutsTheName(h any) {}

// FindsNothing checks that Find of a declaration without callables finds
// nothing.
//
//dokimi:check Empty.Find defect=broken
func FindsNothing(h any) {}

// Undocumented returns the item of the store.
//
//dokimi:check store.Store.Get defect=broken
func Undocumented(h any) {}

// TakesTwo checks that Get takes two handles.
//
//dokimi:check store.Store.Get defect=broken
func TakesTwo(h, other any) {}

// TakesThree checks that Get takes a handle and two cases.
//
//dokimi:check store.Store.Get defect=broken
func TakesThree(h, c, d any) {}

// TakesCases checks that Get takes a list of cases.
//
//dokimi:check store.Store.Get defect=broken
func TakesCases(h any, c []prop.Case) {}

// TakesAnotherCase checks that Get takes a case of another package.
//
//dokimi:check store.Store.Get defect=broken
func TakesAnotherCase(h any, c *Case) {}

// TakesAnOption checks that Get takes an option of the property engine.
//
//dokimi:check store.Store.Get defect=broken
func TakesAnOption(h any, c *prop.Option) {}

// Case is a type with the name of the case of the property engine.
type Case struct{}

// ReturnsOne checks that Get returns a result.
//
//dokimi:check store.Store.Get defect=broken
func ReturnsOne(h any) error { return nil }

// BuildsNothing checks that Get returns the item.
//
//dokimi:check store.Store.Get defect=shapeless
func BuildsNothing(h any) {}

// BuildsTwo checks that Get returns the item.
//
//dokimi:check store.Store.Get defect=paired
func BuildsTwo(h any) {}

// BuildsTheSubject checks that Normalize returns the text.
//
//dokimi:check Normalize defect=broken
func BuildsTheSubject(h any) {}

// ReceivesTheFunction checks that Normalize keeps the text.
//
//dokimi:check Normalize defect=broken
func ReceivesTheFunction(h any) {}

// LatencyBudget checks that Get returns within its budget.
//
//dokimi:check store.Store.Get defect=broken
func LatencyBudget(h any) {}

// RestartsAlone checks that Get returns the item after a restart.
//
//dokimi:check store.Store.Get defect=broken restart=restartsAlone
func RestartsAlone(h any) {}

// CrashesWithoutTheCall checks that Put leaves the whole item after a crash.
//
//dokimi:check store.Store.Put defect=broken crash=crashesWithoutTheCall
func CrashesWithoutTheCall(h any) {}

// CrashesTheText checks that Put leaves the whole item after a crash of a text.
//
//dokimi:check store.Store.Put defect=broken crash=crashesText
func CrashesTheText(h any) {}

// ReferencesTwo checks that Get returns the item of the reference.
//
//dokimi:check store.Store.Get defect=broken reference=referencesTwo
func ReferencesTwo(h any) {}

// OptionsOfAnApply checks that Get returns the item under options.
//
//dokimi:check store.Store.Get defect=broken options=sharesOfAnApply
func OptionsOfAnApply(h any) {}

// OptionsWithoutTheHandle checks that Get returns each item.
//
//dokimi:check store.Store.Get defect=broken options=sharesWithoutTheHandle
func OptionsWithoutTheHandle(h any, c *prop.Case) {}

// OptionsOfTwo checks that Get returns each item under options and an error.
//
//dokimi:check store.Store.Get defect=broken options=sharesAndAnError
func OptionsOfTwo(h any, c *prop.Case) {}

// OptionsOfText checks that Get returns each item under the labels of text.
//
//dokimi:check store.Store.Get defect=broken options=labels
func OptionsOfText(h any, c *prop.Case) {}

// OptionsOfOne checks that Get returns each item under one option.
//
//dokimi:check store.Store.Get defect=broken options=oneShare
func OptionsOfOne(h any, c *prop.Case) {}

// OptionsOfCases checks that Get returns each item under cases of the
// property engine.
//
//dokimi:check store.Store.Get defect=broken options=casesOfTheEngine
func OptionsOfCases(h any, c *prop.Case) {}

// OptionsOfThisPackage checks that Get returns each item under options of this
// package.
//
//dokimi:check store.Store.Get defect=broken options=optionsOfThisPackage
func OptionsOfThisPackage(h any, c *prop.Case) {}

// Option is an option of this package, which has the name of the options of
// the property engine.
type Option struct{}

// casesOfTheEngine returns cases of the property engine in place of options.
func casesOfTheEngine(any) []prop.Case { return nil }

// optionsOfThisPackage returns options of this package.
func optionsOfThisPackage(any) []Option { return nil }

// oneShare returns one option.
func oneShare(any) prop.Option {
	var share prop.Option
	return share
}

// broken returns a store that breaks every check.
func broken(any) store.Store { return nil }

// sharesOfAnApply returns the options of a check without a property run.
func sharesOfAnApply(any) []prop.Option { return nil }

// sharesWithoutTheHandle returns options without the handle of the check.
func sharesWithoutTheHandle() []prop.Option { return nil }

// sharesAndAnError returns options and an error.
func sharesAndAnError(any) ([]prop.Option, error) { return nil, nil }

// labels returns the labels of the outcomes as text.
func labels(any) []string { return nil }

// restartsAlone returns a store without the planted defect that it restarts.
func restartsAlone(any) store.Store { return nil }

// crashesWithoutTheCall returns a store without the call that it crashes.
func crashesWithoutTheCall(any, store.Store) store.Store { return nil }

// crashesText returns a store after the crash of a text in place of a call.
func crashesText(any, store.Store, string) store.Store { return nil }

// referencesTwo returns a reference and an error.
func referencesTwo(any) (store.Store, error) { return nil, nil }

// shapeless returns a store without the handle of a check.
func shapeless() store.Store { return nil }

// paired returns a store and an error.
func paired(any) (store.Store, error) { return nil, nil }
