package generic_test

// [Keyed] is generic, and its type parameters have no witnesses that a
// test outside its package can refer to, so the file has no checks of
// KeyedBuilder. A test in the package of a consumer instantiates it.
