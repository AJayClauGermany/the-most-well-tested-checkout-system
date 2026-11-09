// The command line program (rules §8): basket JSON in on stdin, receipt JSON out on stdout.
//
//   echo '{"items":[{"sku":"APPLE","quantity":1}]}' | dotnet run --project src/Checkout.Cli
//
// All the logic is in CheckoutCommand.cs; this file only connects it to the terminal, which is
// why it is left out of the in-process coverage report and tested end-to-end instead.

using System.Diagnostics.CodeAnalysis;
using System.Text.Json;

using CheckoutSystem.Cli;

[ExcludeFromCodeCoverage(Justification = "Runs in a separate process; covered by Checkout.EndToEndTests and e2e/run.sh")]
internal static class Program
{
    private static int Main()
    {
        var input = Console.In.ReadToEnd();
        try
        {
            var result = CheckoutCommand.Run(input, CheckoutCommand.LoadStandardCatalog());
            // "\n" rather than WriteLine, so the output is the same on Windows (no "\r\n").
            Console.Out.Write(result.Output + "\n");
            return result.ExitCode;
        }
        catch (JsonException error)
        {
            // Not a basket at all: nothing on stdout, an explanation on stderr, exit code 2.
            Console.Error.Write($"checkout: {error.Message}\n");
            return 2;
        }
    }
}
