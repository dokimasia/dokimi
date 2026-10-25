package ledger

// KeepsTheSum checks that Balance returns the sum of the credits of the
// account.
//
//dokimi:check Ledger.Balance defect=dropsCredits
func KeepsTheSum(h *LedgerHandle) {}

// CreditsInPlace checks that Credit allocates nothing once the account has a
// balance.
//
//dokimi:check Ledger.Credit defect=copiesBalances counted=true
func CreditsInPlace(h *LedgerHandle) {}

// dropsCredits returns a Ledger whose Credit drops the amount.
func dropsCredits(*LedgerHandle) Ledger { return nil }

// copiesBalances returns a Ledger whose Credit copies every balance.
func copiesBalances(*LedgerHandle) Ledger { return nil }
