// Starts BenchmarkDotNet. Everything after `--` on the command line is for BenchmarkDotNet:
//   --list flat          only list the benchmarks (cheap; what CI does)
//   --job short          a quick run with fewer iterations
//   --filter '*'         run every benchmark without asking
using BenchmarkDotNet.Running;

using CheckoutSystem.Benchmarks;

BenchmarkSwitcher.FromAssembly(typeof(CheckoutBenchmarks).Assembly).Run(args);
