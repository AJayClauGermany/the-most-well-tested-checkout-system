using BenchmarkDotNet.Attributes;

using CheckoutSystem.Cli;

namespace CheckoutSystem.Benchmarks;

/// <summary>
/// A MICRO-BENCHMARK of <see cref="Checkout.Calculate"/>. BenchmarkDotNet runs each
/// [Benchmark] method many times, in a separate warmed-up process, and reports the mean time
/// and the memory allocated per call. The checkout is tiny, so one benchmark class is enough
/// to show the tool; this is not a performance project (see docs/testing-matrix.md).
/// </summary>
[MemoryDiagnoser]
public class CheckoutBenchmarks
{
    private static readonly Catalog Catalog = CheckoutCommand.LoadStandardCatalog();

    private static readonly Basket OneApple = new([new BasketItem("APPLE", 1)]);

    private static readonly Basket Everything = new(
        [
            new BasketItem("HEADPHONES", 1),
            new BasketItem("COFFEE", 6),
            new BasketItem("NOTEBOOK", 12),
            new BasketItem("PEN", 5),
            new BasketItem("BREAD", 1),
            new BasketItem("APPLE", 3),
        ],
        Coupon: "WELCOME5",
        Shipping: "express");

    [Benchmark(Baseline = true, Description = "one apple")]
    public Receipt OneAppleBasket() => Checkout.Calculate(OneApple, Catalog);

    [Benchmark(Description = "every product, both promotions, a coupon and express shipping")]
    public Receipt EverythingBasket() => Checkout.Calculate(Everything, Catalog);
}
