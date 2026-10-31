package loose_test

import (
	"testing"

	"example.com/shop/loose"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// printerSubjectEmbedded is loose.Printer under the name of the field that
// printerSubject embeds.
type printerSubjectEmbedded = loose.Printer

// printerSubject is a value of Printer that only TestPrinterSuite builds.
type printerSubject struct{ printerSubjectEmbedded }

// TestPrinterSuite checks the suite of Printer and the two forms of its
// application.
func TestPrinterSuite(t *testing.T) {
	t.Parallel()

	t.Run("PrinterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, loose.PrinterSuite())
	})

	t.Run("RunPrinterSuite", func(t *testing.T) {
		t.Parallel()
		var built *printerSubject
		loose.RunPrinterSuite(
			t,
			func() *printerSubject {
				built = &printerSubject{}
				return built
			},
			conform.With(loose.PrinterCheck{
				ID:        "Printer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.PrinterHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Printer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.PrinterHandle) loose.Printer { return &printerSubject{} },
			}),
			conform.Excuse(
				loose.PrinterChecks.Print.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPrinterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *printerSubject
		loose.RunPrinterSuiteFromHandle(
			t,
			func(*loose.PrinterHandle) *printerSubject {
				built = &printerSubject{}
				return built
			},
			conform.With(loose.PrinterCheck{
				ID:        "Printer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.PrinterHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Printer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.PrinterHandle) loose.Printer { return &printerSubject{} },
			}),
			conform.Excuse(
				loose.PrinterChecks.Print.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// printerIdle returns the planted defect of a ceiling of Printer with a
// spend of nothing, a subject that answers each call at once.
func printerIdle(h *loose.PrinterHandle) loose.Printer {
	return loose.PrinterSuite().Calls[0].Over(h, func() {})
}

// TestPrinterSuiteBudgets applies the allocation ceilings of the suite of
// Printer to printerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPrinterSuiteBudgets(t *testing.T) {
	loose.RunPrinterBudgets(
		t,
		printerIdle,
		conform.Ceilings(
			loose.PrinterChecks.Print.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			loose.PrinterChecks.Print.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPrinterSuitePrint registers the fuzz target of Printer.Print on
// printerIdle, with every check of the suite excused.
func FuzzPrinterSuitePrint(f *testing.F) {
	loose.FuzzPrinterSuite(
		f,
		loose.PrinterChecks.Print.Callable,
		printerIdle,
		conform.Excuse(
			loose.PrinterChecks.Print.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPrinterSuite measures each callable of the suite of Printer on
// printerIdle.
func BenchmarkPrinterSuite(b *testing.B) {
	b.Run("Print", func(b *testing.B) {
		loose.BenchPrinterSuite(b, loose.PrinterChecks.Print.Callable, printerIdle)
	})
}

// recorderSubjectEmbedded is loose.Recorder under the name of the field
// that recorderSubject embeds.
type recorderSubjectEmbedded = loose.Recorder

// recorderSubject is a value of Recorder that only TestRecorderSuite
// builds.
type recorderSubject struct{ recorderSubjectEmbedded }

// TestRecorderSuite checks the suite of Recorder and the two forms of its
// application.
func TestRecorderSuite(t *testing.T) {
	t.Parallel()

	t.Run("RecorderSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, loose.RecorderSuite())
	})

	t.Run("RunRecorderSuite", func(t *testing.T) {
		t.Parallel()
		var built *recorderSubject
		loose.RunRecorderSuite(
			t,
			func() *recorderSubject {
				built = &recorderSubject{}
				return built
			},
			conform.With(loose.RecorderCheck{
				ID:        "Recorder/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.RecorderHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Recorder(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.RecorderHandle) loose.Recorder { return &recorderSubject{} },
			}),
			conform.Excuse(
				loose.RecorderChecks.Field.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunRecorderSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *recorderSubject
		loose.RunRecorderSuiteFromHandle(
			t,
			func(*loose.RecorderHandle) *recorderSubject {
				built = &recorderSubject{}
				return built
			},
			conform.With(loose.RecorderCheck{
				ID:        "Recorder/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.RecorderHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Recorder(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.RecorderHandle) loose.Recorder { return &recorderSubject{} },
			}),
			conform.Excuse(
				loose.RecorderChecks.Field.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// recorderIdle returns the planted defect of a ceiling of Recorder with a
// spend of nothing, a subject that answers each call at once.
func recorderIdle(h *loose.RecorderHandle) loose.Recorder {
	return loose.RecorderSuite().Calls[0].Over(h, func() {})
}

// TestRecorderSuiteBudgets applies the allocation ceilings of the suite of
// Recorder to recorderIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRecorderSuiteBudgets(t *testing.T) {
	loose.RunRecorderBudgets(
		t,
		recorderIdle,
		conform.Ceilings(
			loose.RecorderChecks.Field.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			loose.RecorderChecks.Field.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzRecorderSuiteField registers the fuzz target of Recorder.Field on
// recorderIdle, with every check of the suite excused.
func FuzzRecorderSuiteField(f *testing.F) {
	loose.FuzzRecorderSuite(
		f,
		loose.RecorderChecks.Field.Callable,
		recorderIdle,
		conform.Excuse(
			loose.RecorderChecks.Field.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkRecorderSuite measures each callable of the suite of Recorder
// on recorderIdle.
func BenchmarkRecorderSuite(b *testing.B) {
	b.Run("Field", func(b *testing.B) {
		loose.BenchRecorderSuite(b, loose.RecorderChecks.Field.Callable, recorderIdle)
	})
}

// journalSubjectEmbedded is loose.Journal under the name of the field that
// journalSubject embeds.
type journalSubjectEmbedded = loose.Journal

// journalSubject is a value of Journal that only TestJournalSuite builds.
type journalSubject struct{ journalSubjectEmbedded }

// TestJournalSuite checks the suite of Journal and the two forms of its
// application.
func TestJournalSuite(t *testing.T) {
	t.Parallel()

	t.Run("JournalSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, loose.JournalSuite())
	})

	t.Run("RunJournalSuite", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		loose.RunJournalSuite(
			t,
			func() *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(loose.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.JournalHandle) loose.Journal { return &journalSubject{} },
			}),
			conform.Excuse(
				loose.JournalChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunJournalSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		loose.RunJournalSuiteFromHandle(
			t,
			func(*loose.JournalHandle) *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(loose.JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.JournalHandle) loose.Journal { return &journalSubject{} },
			}),
			conform.Excuse(
				loose.JournalChecks.Put.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// journalIdle returns the planted defect of a ceiling of Journal with a
// spend of nothing, a subject that answers each call at once.
func journalIdle(h *loose.JournalHandle) loose.Journal {
	return loose.JournalSuite().Calls[0].Over(h, func() {})
}

// TestJournalSuiteBudgets applies the allocation ceilings of the suite of
// Journal to journalIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestJournalSuiteBudgets(t *testing.T) {
	loose.RunJournalBudgets(
		t,
		journalIdle,
		conform.Ceilings(loose.JournalChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			loose.JournalChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzJournalSuitePut registers the fuzz target of Journal.Put on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuitePut(f *testing.F) {
	loose.FuzzJournalSuite(
		f,
		loose.JournalChecks.Put.Callable,
		journalIdle,
		conform.Excuse(
			loose.JournalChecks.Put.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkJournalSuite measures each callable of the suite of Journal on
// journalIdle.
func BenchmarkJournalSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		loose.BenchJournalSuite(b, loose.JournalChecks.Put.Callable, journalIdle)
	})
}

// plainSubjectEmbedded is loose.Plain under the name of the field that
// plainSubject embeds.
type plainSubjectEmbedded = loose.Plain

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, loose.PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		loose.RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(loose.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.PlainHandle) loose.Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				loose.PlainChecks.Print.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		loose.RunPlainSuiteFromHandle(
			t,
			func(*loose.PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(loose.PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.PlainHandle) loose.Plain { return &plainSubject{} },
			}),
			conform.Excuse(
				loose.PlainChecks.Print.Smoke(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *loose.PlainHandle) loose.Plain {
	return loose.PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	loose.RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(loose.PlainChecks.Print.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Excuse(
			loose.PlainChecks.Print.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzPlainSuitePrint registers the fuzz target of Plain.Print on
// plainIdle, with every check of the suite excused.
func FuzzPlainSuitePrint(f *testing.F) {
	loose.FuzzPlainSuite(
		f,
		loose.PlainChecks.Print.Callable,
		plainIdle,
		conform.Excuse(
			loose.PlainChecks.Print.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Print", func(b *testing.B) {
		loose.BenchPlainSuite(b, loose.PlainChecks.Print.Callable, plainIdle)
	})
}

// zonedSubjectEmbedded is loose.Zoned under the name of the field that
// zonedSubject embeds.
type zonedSubjectEmbedded = loose.Zoned

// zonedSubject is a value of Zoned that only TestZonedSuite builds.
type zonedSubject struct{ zonedSubjectEmbedded }

// TestZonedSuite checks the suite of Zoned and the two forms of its
// application.
func TestZonedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ZonedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, loose.ZonedSuite())
	})

	t.Run("RunZonedSuite", func(t *testing.T) {
		t.Parallel()
		var built *zonedSubject
		loose.RunZonedSuite(
			t,
			func() *zonedSubject {
				built = &zonedSubject{}
				return built
			},
			conform.With(loose.ZonedCheck{
				ID:        "Zoned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.ZonedHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Zoned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.ZonedHandle) loose.Zoned { return &zonedSubject{} },
			}),
		)
	})

	t.Run("RunZonedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *zonedSubject
		loose.RunZonedSuiteFromHandle(
			t,
			func(*loose.ZonedHandle) *zonedSubject {
				built = &zonedSubject{}
				return built
			},
			conform.With(loose.ZonedCheck{
				ID:        "Zoned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *loose.ZonedHandle) {
					got := h.Subject()
					assert.Equal(h, got, loose.Zoned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*loose.ZonedHandle) loose.Zoned { return &zonedSubject{} },
			}),
		)
	})
}

// zonedIdle returns the planted defect of a ceiling of Zoned with a spend
// of nothing, a subject that answers each call at once.
func zonedIdle(h *loose.ZonedHandle) loose.Zoned {
	return loose.ZonedSuite().Calls[0].Over(h, func() {})
}

// TestZonedSuiteBudgets applies the allocation ceilings of the suite of
// Zoned to zonedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestZonedSuiteBudgets(t *testing.T) {
	loose.RunZonedBudgets(
		t,
		zonedIdle,
		conform.Ceilings(loose.ZonedChecks.In.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(loose.ZonedChecks.Scale.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(loose.ZonedChecks.Ratio.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzZonedSuiteIn registers the fuzz target of Zoned.In on zonedIdle,
// with every check of the suite excused.
func FuzzZonedSuiteIn(f *testing.F) {
	loose.FuzzZonedSuite(
		f,
		loose.ZonedChecks.In.Callable,
		zonedIdle,
	)
}

// FuzzZonedSuiteScale registers the fuzz target of Zoned.Scale on
// zonedIdle, with every check of the suite excused.
func FuzzZonedSuiteScale(f *testing.F) {
	loose.FuzzZonedSuite(
		f,
		loose.ZonedChecks.Scale.Callable,
		zonedIdle,
	)
}

// FuzzZonedSuiteRatio registers the fuzz target of Zoned.Ratio on
// zonedIdle, with every check of the suite excused.
func FuzzZonedSuiteRatio(f *testing.F) {
	loose.FuzzZonedSuite(
		f,
		loose.ZonedChecks.Ratio.Callable,
		zonedIdle,
	)
}

// BenchmarkZonedSuite measures each callable of the suite of Zoned on
// zonedIdle.
func BenchmarkZonedSuite(b *testing.B) {
	b.Run("In", func(b *testing.B) {
		loose.BenchZonedSuite(b, loose.ZonedChecks.In.Callable, zonedIdle)
	})
	b.Run("Scale", func(b *testing.B) {
		loose.BenchZonedSuite(b, loose.ZonedChecks.Scale.Callable, zonedIdle)
	})
	b.Run("Ratio", func(b *testing.B) {
		loose.BenchZonedSuite(b, loose.ZonedChecks.Ratio.Callable, zonedIdle)
	})
}
