package loose

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// printerSubjectEmbedded is Printer under the name of the field that
// printerSubject embeds.
type printerSubjectEmbedded = Printer

// printerSubject is a value of Printer that only TestPrinterSuite builds.
type printerSubject struct{ printerSubjectEmbedded }

// TestPrinterSuite checks the suite of Printer and the two forms of its
// application.
func TestPrinterSuite(t *testing.T) {
	t.Parallel()

	t.Run("PrinterSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PrinterSuite())
	})

	t.Run("RunPrinterSuite", func(t *testing.T) {
		t.Parallel()
		var built *printerSubject
		RunPrinterSuite(
			t,
			func() *printerSubject {
				built = &printerSubject{}
				return built
			},
			conform.With(PrinterCheck{
				ID:        "Printer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PrinterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Printer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PrinterHandle) Printer { return &printerSubject{} },
			}),
		)
	})

	t.Run("RunPrinterSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *printerSubject
		RunPrinterSuiteFromHandle(
			t,
			func(*PrinterHandle) *printerSubject {
				built = &printerSubject{}
				return built
			},
			conform.With(PrinterCheck{
				ID:        "Printer/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PrinterHandle) {
					got := h.Subject()
					assert.Equal(h, got, Printer(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PrinterHandle) Printer { return &printerSubject{} },
			}),
		)
	})
}

// printerIdle returns the planted defect of a ceiling of Printer with a
// spend of nothing, a subject that answers each call at once.
func printerIdle(h *PrinterHandle) Printer {
	return PrinterSuite().Calls[0].Over(h, func() {})
}

// TestPrinterSuiteBudgets applies the allocation ceilings of the suite of
// Printer to printerIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestPrinterSuiteBudgets(t *testing.T) {
	RunPrinterBudgets(
		t,
		printerIdle,
		conform.Ceilings(PrinterChecks.Print.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzPrinterSuitePrint registers the fuzz target of Printer.Print on
// printerIdle, with every check of the suite excused.
func FuzzPrinterSuitePrint(f *testing.F) {
	FuzzPrinterSuite(
		f,
		PrinterChecks.Print.Callable,
		printerIdle,
	)
}

// BenchmarkPrinterSuite measures each callable of the suite of Printer on
// printerIdle.
func BenchmarkPrinterSuite(b *testing.B) {
	b.Run("Print", func(b *testing.B) {
		BenchPrinterSuite(b, PrinterChecks.Print.Callable, printerIdle)
	})
}

// recorderSubjectEmbedded is Recorder under the name of the field that
// recorderSubject embeds.
type recorderSubjectEmbedded = Recorder

// recorderSubject is a value of Recorder that only TestRecorderSuite
// builds.
type recorderSubject struct{ recorderSubjectEmbedded }

// TestRecorderSuite checks the suite of Recorder and the two forms of its
// application.
func TestRecorderSuite(t *testing.T) {
	t.Parallel()

	t.Run("RecorderSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, RecorderSuite())
	})

	t.Run("RunRecorderSuite", func(t *testing.T) {
		t.Parallel()
		var built *recorderSubject
		RunRecorderSuite(
			t,
			func() *recorderSubject {
				built = &recorderSubject{}
				return built
			},
			conform.With(RecorderCheck{
				ID:        "Recorder/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RecorderHandle) {
					got := h.Subject()
					assert.Equal(h, got, Recorder(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RecorderHandle) Recorder { return &recorderSubject{} },
			}),
		)
	})

	t.Run("RunRecorderSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *recorderSubject
		RunRecorderSuiteFromHandle(
			t,
			func(*RecorderHandle) *recorderSubject {
				built = &recorderSubject{}
				return built
			},
			conform.With(RecorderCheck{
				ID:        "Recorder/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *RecorderHandle) {
					got := h.Subject()
					assert.Equal(h, got, Recorder(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*RecorderHandle) Recorder { return &recorderSubject{} },
			}),
		)
	})
}

// recorderIdle returns the planted defect of a ceiling of Recorder with a
// spend of nothing, a subject that answers each call at once.
func recorderIdle(h *RecorderHandle) Recorder {
	return RecorderSuite().Calls[0].Over(h, func() {})
}

// TestRecorderSuiteBudgets applies the allocation ceilings of the suite of
// Recorder to recorderIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestRecorderSuiteBudgets(t *testing.T) {
	RunRecorderBudgets(
		t,
		recorderIdle,
		conform.Ceilings(RecorderChecks.Field.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzRecorderSuiteField registers the fuzz target of Recorder.Field on
// recorderIdle, with every check of the suite excused.
func FuzzRecorderSuiteField(f *testing.F) {
	FuzzRecorderSuite(
		f,
		RecorderChecks.Field.Callable,
		recorderIdle,
	)
}

// BenchmarkRecorderSuite measures each callable of the suite of Recorder
// on recorderIdle.
func BenchmarkRecorderSuite(b *testing.B) {
	b.Run("Field", func(b *testing.B) {
		BenchRecorderSuite(b, RecorderChecks.Field.Callable, recorderIdle)
	})
}

// journalSubjectEmbedded is Journal under the name of the field that
// journalSubject embeds.
type journalSubjectEmbedded = Journal

// journalSubject is a value of Journal that only TestJournalSuite builds.
type journalSubject struct{ journalSubjectEmbedded }

// TestJournalSuite checks the suite of Journal and the two forms of its
// application.
func TestJournalSuite(t *testing.T) {
	t.Parallel()

	t.Run("JournalSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, JournalSuite())
	})

	t.Run("RunJournalSuite", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		RunJournalSuite(
			t,
			func() *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*JournalHandle) Journal { return &journalSubject{} },
			}),
		)
	})

	t.Run("RunJournalSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *journalSubject
		RunJournalSuiteFromHandle(
			t,
			func(*JournalHandle) *journalSubject {
				built = &journalSubject{}
				return built
			},
			conform.With(JournalCheck{
				ID:        "Journal/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *JournalHandle) {
					got := h.Subject()
					assert.Equal(h, got, Journal(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*JournalHandle) Journal { return &journalSubject{} },
			}),
		)
	})
}

// journalIdle returns the planted defect of a ceiling of Journal with a
// spend of nothing, a subject that answers each call at once.
func journalIdle(h *JournalHandle) Journal {
	return JournalSuite().Calls[0].Over(h, func() {})
}

// TestJournalSuiteBudgets applies the allocation ceilings of the suite of
// Journal to journalIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestJournalSuiteBudgets(t *testing.T) {
	RunJournalBudgets(
		t,
		journalIdle,
		conform.Ceilings(JournalChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzJournalSuitePut registers the fuzz target of Journal.Put on
// journalIdle, with every check of the suite excused.
func FuzzJournalSuitePut(f *testing.F) {
	FuzzJournalSuite(
		f,
		JournalChecks.Put.Callable,
		journalIdle,
	)
}

// BenchmarkJournalSuite measures each callable of the suite of Journal on
// journalIdle.
func BenchmarkJournalSuite(b *testing.B) {
	b.Run("Put", func(b *testing.B) {
		BenchJournalSuite(b, JournalChecks.Put.Callable, journalIdle)
	})
}

// plainSubjectEmbedded is Plain under the name of the field that
// plainSubject embeds.
type plainSubjectEmbedded = Plain

// plainSubject is a value of Plain that only TestPlainSuite builds.
type plainSubject struct{ plainSubjectEmbedded }

// TestPlainSuite checks the suite of Plain and the two forms of its
// application.
func TestPlainSuite(t *testing.T) {
	t.Parallel()

	t.Run("PlainSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, PlainSuite())
	})

	t.Run("RunPlainSuite", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuite(
			t,
			func() *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) Plain { return &plainSubject{} },
			}),
		)
	})

	t.Run("RunPlainSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *plainSubject
		RunPlainSuiteFromHandle(
			t,
			func(*PlainHandle) *plainSubject {
				built = &plainSubject{}
				return built
			},
			conform.With(PlainCheck{
				ID:        "Plain/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *PlainHandle) {
					got := h.Subject()
					assert.Equal(h, got, Plain(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*PlainHandle) Plain { return &plainSubject{} },
			}),
		)
	})
}

// plainIdle returns the planted defect of a ceiling of Plain with a spend
// of nothing, a subject that answers each call at once.
func plainIdle(h *PlainHandle) Plain {
	return PlainSuite().Calls[0].Over(h, func() {})
}

// TestPlainSuiteBudgets applies the allocation ceilings of the suite of
// Plain to plainIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestPlainSuiteBudgets(t *testing.T) {
	RunPlainBudgets(
		t,
		plainIdle,
		conform.Ceilings(PlainChecks.Print.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzPlainSuitePrint registers the fuzz target of Plain.Print on
// plainIdle, with every check of the suite excused.
func FuzzPlainSuitePrint(f *testing.F) {
	FuzzPlainSuite(
		f,
		PlainChecks.Print.Callable,
		plainIdle,
	)
}

// BenchmarkPlainSuite measures each callable of the suite of Plain on
// plainIdle.
func BenchmarkPlainSuite(b *testing.B) {
	b.Run("Print", func(b *testing.B) {
		BenchPlainSuite(b, PlainChecks.Print.Callable, plainIdle)
	})
}

// zonedSubjectEmbedded is Zoned under the name of the field that
// zonedSubject embeds.
type zonedSubjectEmbedded = Zoned

// zonedSubject is a value of Zoned that only TestZonedSuite builds.
type zonedSubject struct{ zonedSubjectEmbedded }

// TestZonedSuite checks the suite of Zoned and the two forms of its
// application.
func TestZonedSuite(t *testing.T) {
	t.Parallel()

	t.Run("ZonedSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, ZonedSuite())
	})

	t.Run("RunZonedSuite", func(t *testing.T) {
		t.Parallel()
		var built *zonedSubject
		RunZonedSuite(
			t,
			func() *zonedSubject {
				built = &zonedSubject{}
				return built
			},
			conform.With(ZonedCheck{
				ID:        "Zoned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ZonedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Zoned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ZonedHandle) Zoned { return &zonedSubject{} },
			}),
		)
	})

	t.Run("RunZonedSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *zonedSubject
		RunZonedSuiteFromHandle(
			t,
			func(*ZonedHandle) *zonedSubject {
				built = &zonedSubject{}
				return built
			},
			conform.With(ZonedCheck{
				ID:        "Zoned/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *ZonedHandle) {
					got := h.Subject()
					assert.Equal(h, got, Zoned(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*ZonedHandle) Zoned { return &zonedSubject{} },
			}),
		)
	})
}

// zonedIdle returns the planted defect of a ceiling of Zoned with a spend
// of nothing, a subject that answers each call at once.
func zonedIdle(h *ZonedHandle) Zoned {
	return ZonedSuite().Calls[0].Over(h, func() {})
}

// TestZonedSuiteBudgets applies the allocation ceilings of the suite of
// Zoned to zonedIdle with every check of the suite excused, and a ceiling
// of 0 to each callable without a budget whose calls build no input of
// their own.
func TestZonedSuiteBudgets(t *testing.T) {
	RunZonedBudgets(
		t,
		zonedIdle,
		conform.Ceilings(ZonedChecks.In.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ZonedChecks.Scale.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.Ceilings(ZonedChecks.Ratio.Callable, conform.Budget{Allocs: new(uint64(0))}),
	)
}

// FuzzZonedSuiteIn registers the fuzz target of Zoned.In on zonedIdle,
// with every check of the suite excused.
func FuzzZonedSuiteIn(f *testing.F) {
	FuzzZonedSuite(
		f,
		ZonedChecks.In.Callable,
		zonedIdle,
	)
}

// FuzzZonedSuiteScale registers the fuzz target of Zoned.Scale on
// zonedIdle, with every check of the suite excused.
func FuzzZonedSuiteScale(f *testing.F) {
	FuzzZonedSuite(
		f,
		ZonedChecks.Scale.Callable,
		zonedIdle,
	)
}

// FuzzZonedSuiteRatio registers the fuzz target of Zoned.Ratio on
// zonedIdle, with every check of the suite excused.
func FuzzZonedSuiteRatio(f *testing.F) {
	FuzzZonedSuite(
		f,
		ZonedChecks.Ratio.Callable,
		zonedIdle,
	)
}

// BenchmarkZonedSuite measures each callable of the suite of Zoned on
// zonedIdle.
func BenchmarkZonedSuite(b *testing.B) {
	b.Run("In", func(b *testing.B) {
		BenchZonedSuite(b, ZonedChecks.In.Callable, zonedIdle)
	})
	b.Run("Scale", func(b *testing.B) {
		BenchZonedSuite(b, ZonedChecks.Scale.Callable, zonedIdle)
	})
	b.Run("Ratio", func(b *testing.B) {
		BenchZonedSuite(b, ZonedChecks.Ratio.Callable, zonedIdle)
	})
}
