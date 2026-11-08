// Vitest settings. Everything the test runner needs to know is in this one file.
import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    // Where the tests are. The kinds of test each have their own folder - see README.md.
    include: ["tests/**/*.test.ts"],

    // Type-level tests (tests/types/*.test-d.ts) run with `npm run test:types`.
    typecheck: {
      include: ["tests/types/**/*.test-d.ts"],
      tsconfig: "./tsconfig.json",
    },

    // Benchmarks (tests/bench/*.bench.ts) run with `npm run bench`.
    benchmark: {
      include: ["tests/bench/**/*.bench.ts"],
      // checkout() calls into other modules on every run, and Vitest's module loader wraps
      // every import in a getter. Vitest warns that this adds overhead; for our purposes
      // (showing the tool, not measuring nanoseconds) that is fine.
      suppressExportGetterWarnings: true,
    },

    // node:sqlite still prints "ExperimentalWarning" on Node 22; keep the test output clean.
    execArgv: ["--disable-warning=ExperimentalWarning"],

    // Two reports: the usual terminal output, and a JUnit XML file for CI dashboards.
    reporters: ["default", ["junit", { outputFile: "reports/junit.xml" }]],

    // Coverage: `npm run test:coverage`. The thresholds make the run fail when coverage drops.
    coverage: {
      provider: "v8",
      include: ["src/**/*.ts"],
      // cli.ts is the few lines that connect command.ts to stdin/stdout. It runs in a separate
      // process (tests/e2e/), which in-process coverage cannot see.
      exclude: ["src/cli.ts"],
      reporter: ["text", "html", "lcov"],
      reportsDirectory: "reports/coverage",
      thresholds: {
        statements: 100,
        branches: 100,
        functions: 100,
        lines: 100,
      },
    },
  },
});
