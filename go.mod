// Reviewed: 2025-10-29
// Review notes (short):
// - Good separation: cards, deck, hand, cmd packages.
// - Add unit tests for hand evaluation and deck operations.
// - Make Shuffle seed injectable for deterministic tests (accept rand.Source).
// - Add package docs and examples for flexibility to support other poker variants.
// - Run `gofmt`, `go vet` and `go test` across packages.

module github.com/dangogh/GoPoker

go 1.25
