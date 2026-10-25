package nested

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// courierSubjectEmbedded is Courier under the name of the field that
// courierSubject embeds.
type courierSubjectEmbedded = Courier

// courierSubject is a value of Courier that only TestCourierSuite builds.
type courierSubject struct{ courierSubjectEmbedded }

// TestCourierSuite checks the suite of Courier and the two forms of its
// application.
func TestCourierSuite(t *testing.T) {
	t.Parallel()

	t.Run("CourierSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CourierSuite())
	})

	t.Run("RunCourierSuite", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		RunCourierSuite(
			t,
			func() *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CourierHandle) Courier { return &courierSubject{} },
			}),
		)
	})

	t.Run("RunCourierSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *courierSubject
		RunCourierSuiteFromHandle(
			t,
			func(*CourierHandle) *courierSubject {
				built = &courierSubject{}
				return built
			},
			conform.With(CourierCheck{
				ID:        "Courier/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CourierHandle) {
					got := h.Subject()
					assert.Equal(h, got, Courier(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CourierHandle) Courier { return &courierSubject{} },
			}),
		)
	})
}

// courierIdle returns the planted defect of a ceiling of Courier with a
// spend of nothing, a subject that answers each call at once.
func courierIdle(h *CourierHandle) Courier {
	return CourierSuite().Calls[0].Over(h, func() {})
}

// TestCourierSuiteBudgets applies the allocation ceilings of the suite of
// Courier to courierIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCourierSuiteBudgets(t *testing.T) {
	RunCourierBudgets(
		t,
		courierIdle,
		conform.Ceilings(CourierChecks.Ship.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(CourierChecks.Send.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzCourierSuiteShip registers the fuzz target of Courier.Ship on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteShip(f *testing.F) {
	FuzzCourierSuite(
		f,
		CourierChecks.Ship.Callable,
		courierIdle,
	)
}

// FuzzCourierSuiteSend registers the fuzz target of Courier.Send on
// courierIdle, with every check of the suite excused.
func FuzzCourierSuiteSend(f *testing.F) {
	FuzzCourierSuite(
		f,
		CourierChecks.Send.Callable,
		courierIdle,
	)
}

// BenchmarkCourierSuite measures each callable of the suite of Courier on
// courierIdle.
func BenchmarkCourierSuite(b *testing.B) {
	b.Run("Ship", func(b *testing.B) {
		BenchCourierSuite(b, CourierChecks.Ship.Callable, courierIdle)
	})
	b.Run("Send", func(b *testing.B) {
		BenchCourierSuite(b, CourierChecks.Send.Callable, courierIdle)
	})
}

// TestHireSuite checks the suite of Hire and the two forms of its
// application.
func TestHireSuite(t *testing.T) {
	t.Parallel()

	t.Run("HireSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, HireSuite())
	})

	t.Run("Hire", func(t *testing.T) {
		t.Parallel()
		RunHireSuite(t, Hire)
	})

	t.Run("RunHireSuite", func(t *testing.T) {
		t.Parallel()
		RunHireSuite(
			t,
			Hire,
			conform.With(HireCheck{
				ID:        "Hire/receives-the-function",
				Statement: "the subject of a check is the function of the application",
				Apply: func(h *HireHandle) {
					assert.Equal(h, h.Subject(), HireFunc(Hire), h.Statement())
				},
				Defect: func(*HireHandle) HireFunc { return nil },
			}),
		)
	})

	t.Run("RunHireSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		RunHireSuiteFromHandle(
			t,
			func(*HireHandle) HireFunc {
				return Hire
			},
			conform.With(HireCheck{
				ID:        "Hire/builds-the-subject",
				Statement: "the subject of a check is the function that the factory returns",
				Apply: func(h *HireHandle) {
					assert.Equal(h, h.Subject(), HireFunc(Hire), h.Statement())
				},
				Defect: func(*HireHandle) HireFunc { return nil },
			}),
		)
	})
}

// BenchmarkHire measures Hire.
func BenchmarkHire(b *testing.B) {
	BenchHireSuite(b, Hire)
}

// FuzzHire fuzzes the property checks of Hire.
func FuzzHire(f *testing.F) {
	FuzzHireSuite(f, Hire)
}
