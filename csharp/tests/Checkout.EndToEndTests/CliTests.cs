using System.Diagnostics;
using System.Reflection;
using System.Text.Json;
using System.Text.Json.Nodes;

namespace CheckoutSystem.EndToEndTests;

/// <summary>
/// END-TO-END TESTS: run the real command line program (rules §8) as a separate process,
/// exactly as a user (or the shared runner in e2e/) would, and look only at what comes out.
///
/// The program is started as <c>dotnet Checkout.Cli.dll</c>. Because the test only sees stdout,
/// stderr and the exit code, it would pass unchanged for a program written in any language –
/// that is the idea behind e2e/run.sh at the root of the repository.
/// </summary>
public class CliTests
{
    /// <summary>Where the build put Checkout.Cli.dll (see the AssemblyMetadata in the .csproj).</summary>
    private static readonly string CliPath = Path.GetFullPath(
        typeof(CliTests).Assembly.GetCustomAttributes<AssemblyMetadataAttribute>().Single(attribute => attribute.Key == "CheckoutCliPath").Value!);

    public static TheoryData<string> CaseIds() =>
        new(SpecLoader.LoadCases().Select(specCase => specCase.Id));

    [Fact]
    public void TheProgram_HasBeenBuilt()
    {
        Assert.True(File.Exists(CliPath), $"Checkout.Cli.dll not found at {CliPath}; build src/Checkout.Cli first");
    }

    // Every shared case, through the real program.
    [Theory]
    [MemberData(nameof(CaseIds))]
    public void Cli_SharedSpecCase_PrintsTheExpectedReceiptOrError(string caseId)
    {
        var specCase = SpecLoader.LoadCase(caseId);

        var result = RunCli(JsonSerializer.Serialize(specCase.Basket, CheckoutJson.Options) + "\n");

        if (specCase.ExpectedError is { } expectedError)
        {
            Assert.Equal(1, result.ExitCode);
            Assert.Equal($$$"""{"error":"{{{expectedError.ToSpecCode()}}}"}""", result.Stdout.Trim());
        }
        else
        {
            Assert.Equal(0, result.ExitCode);
            var receipt = JsonNode.Parse(result.Stdout)!;
            Assert.Equal(specCase.Expected, new Totals(
                receipt["subtotal"]!.GetValue<int>(),
                receipt["promotion_discount"]!.GetValue<int>(),
                receipt["coupon_discount"]!.GetValue<int>(),
                receipt["shipping"]!.GetValue<int>(),
                receipt["tax"]!.GetValue<int>(),
                receipt["total"]!.GetValue<int>()));
            Assert.Empty(SchemaValidator.Errors("receipt.schema.json", result.Stdout));
        }
    }

    [Fact]
    public void Cli_ValidBasket_PrintsExactlyOneLineOfJson()
    {
        var result = RunCli("""{"items":[{"sku":"APPLE","quantity":1}]}""");

        Assert.EndsWith("\n", result.Stdout, StringComparison.Ordinal);
        Assert.Single(result.Stdout.TrimEnd('\n').Split('\n'));
        Assert.Equal("", result.Stderr);
    }

    [Fact]
    public void Cli_InputThatIsNotJson_ExitsWith2AndExplainsOnStderr()
    {
        var result = RunCli("this is not json");

        Assert.Equal(2, result.ExitCode);
        Assert.Equal("", result.Stdout);
        Assert.StartsWith("checkout: ", result.Stderr, StringComparison.Ordinal);
    }

    [Fact]
    public void Cli_JsonThatIsNotABasket_ExitsWith2()
    {
        var result = RunCli("""{"hello": "world"}""");

        Assert.Equal(2, result.ExitCode);
        Assert.Equal("", result.Stdout);
        Assert.StartsWith("checkout: ", result.Stderr, StringComparison.Ordinal);
    }

    [Fact]
    public void Cli_EmptyInput_ExitsWith2()
    {
        var result = RunCli("");

        Assert.Equal(2, result.ExitCode);
        Assert.Equal("", result.Stdout);
    }

    private sealed record CliResult(string Stdout, string Stderr, int ExitCode);

    /// <summary>Runs the program once, with <paramref name="input"/> on standard input.</summary>
    private static CliResult RunCli(string input)
    {
        var startInfo = new ProcessStartInfo("dotnet", [CliPath])
        {
            RedirectStandardInput = true,
            RedirectStandardOutput = true,
            RedirectStandardError = true,
            UseShellExecute = false,
        };
        using var process = Process.Start(startInfo)!;

        process.StandardInput.Write(input);
        process.StandardInput.Close();
        var stdout = process.StandardOutput.ReadToEndAsync();
        var stderr = process.StandardError.ReadToEndAsync();
        process.WaitForExit();

        return new CliResult(stdout.Result, stderr.Result, process.ExitCode);
    }
}
