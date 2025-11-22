# Copilot Instructions — Go

Purpose

These instructions tell GitHub Copilot how to produce idiomatic, safe, well-tested, and maintainable Go code for this repository.

Guidelines

- Be idiomatic: prefer Go idioms over patterns from other languages. Use short, clear names, and prefer small functions that do one thing.
- Formatting: always format generated code using `gofmt`/`gofumpt` conventions. Keep import groups (`std`, third-party, local`) separated.
- Imports: prefer `goimports` behavior — remove unused imports and add missing ones.
- Package boundaries: keep packages focused; avoid creating circular dependencies. Use meaningful package names (not `util`).
- Error handling: return errors instead of panicking. Wrap errors with context using `fmt.Errorf("%w")` or `errors.Join`/`errors.As` as appropriate. Prefer sentinel errors sparingly.
- Context: accept `context.Context` on functions that perform IO, long-running work, or blocking operations; propagate contexts and respect cancellation.
- Concurrency: use goroutines, channels, `sync` primitives, and worker pools only when necessary. Prefer simpler patterns first. Avoid sharing mutable state without synchronization.
- Pointers vs values: prefer values for small structs; use pointers when mutability or large memory copying matters.
- Interfaces: define interfaces at the consumer side (dependency inversion). Keep interfaces small — one or a few methods.
- Testing: generate unit tests alongside code. Use table-driven tests and `t.Run`. Add examples and property tests when helpful.
- Benchmarks: add `Benchmark` tests for hotspots when relevant.
-- Documentation: add clear, concise `//` comments for exported functions/types and examples in `Example` tests. Comments should explain the reasoning, intent, invariants, or tradeoffs behind the code (the "why"), not merely restate what the code does. Keep comments short and focused — prefer one-line rationale or a short paragraph for non-obvious decisions.
- Security: validate input, avoid command injection, sanitize interactions with the network and file system.
- Performance: prefer clarity; optimize after measuring. Avoid premature micro-optimizations.
- Logging: use structured logging and context-based request IDs where applicable.
- Dependencies: add modules with `go.mod`. Prefer standard library; pin third-party modules to minimal required versions.

Prompting Rules (how to interpret prompts in this repo)

- If the user asks for an implementation of a package-level function, generate the function with tests and an example if feasible.
- If the user asks to refactor, provide a minimal, safe refactor with tests that preserve behavior.
- When adding public API, include documentation comments and a basic usage example.
- When asked for code fixes, provide a brief explanation of the root cause and a compact patch.

Do This (examples)

- "Write an idiomatic function to parse a poker hand string into a struct, return typed errors, and include table-driven tests."
- "Add unit tests for `dealer.Deal` covering edge cases: empty deck, insufficient cards, and concurrency safety." 
- "Refactor `pkg/db` to use context-aware methods and add an interface for easier mocking in tests."
 - "When documenting non-obvious behavior, explain the rationale: e.g., '// We use a fixed-size pool to limit memory growth and to bound latency under load.' Keep this concise."

Don't Do This

- Do not generate code that imports unused packages.
- Do not use panics for normal error handling.
- Do not leak goroutines; ensure goroutines exit on context cancellation or done channels.
- Do not generate large monolithic functions — split into testable units.
 - Do not write comments that simply restate the code (e.g., `// increment i` above `i++`). Avoid long prose that duplicates implementation details.

Editor & Tooling

- Ensure generated code is `gofmt`-compatible.
- Prefer `gofumpt` and `go vet` rules where they improve clarity.
- Respect repository `golangci-lint` settings if present.

Review & Commit Tips

- Keep commits small and focused. Include tests in the same commit as implementation changes.
- Use descriptive commit messages: "pkg/player: add ParseHand with tests".

Quick Checklist (for Copilot and reviewers)

- [ ] Code is gofmt-formatted
- [ ] No unused imports or variables
- [ ] Errors are returned, wrapped, and handled
- [ ] Public APIs have docs and examples
- [ ] Tests cover normal and edge cases
- [ ] No goroutine leaks; contexts respected
 - [ ] Comments explain the reasoning and intent (the "why"), not just the "what"

Examples of prompts to get best results

- Good: "Implement `ParseHand(s string) (Hand, error)` in package `poker` with table-driven tests and examples." 
- Good: "Refactor `game.Start` to be cancelable with `context.Context`. Add tests for cancellation." 
- Bad: "Quickly hack a parser in any style." (too vague — request clarification)


---

If you'd like, I can adjust tone (more strict/lenient), add repository-specific linter rules, or generate example code and tests for a particular function in this repo.
