// Starts BenchmarkDotNet. Everything after `--` on the command line is for BenchmarkDotNet:
//   --job dry            run each benchmark once, as a quick smoke test (what CI does)
//   --job short          a quick run with fewer iterations
//   --filter '*'         run every benchmark without asking
//   --list flat          only list the benchmarks, run nothing
using BenchmarkDotNet.Running;

using CheckoutSystem.Benchmarks;

var summaries = BenchmarkSwitcher.FromAssembly(typeof(CheckoutBenchmarks).Assembly).Run(args).ToList();

// BenchmarkDotNet prints an error when a benchmark throws, but still ends the program with
// exit code 0. Check the results ourselves, so that a broken benchmark fails the CI step.
var failed = summaries.Any(summary =>
    summary.HasCriticalValidationErrors || summary.Reports.Any(report => !report.Success));
return failed ? 1 : 0;
