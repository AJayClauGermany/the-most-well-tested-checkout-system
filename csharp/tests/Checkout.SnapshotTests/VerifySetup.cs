using System.Runtime.CompilerServices;

namespace CheckoutSystem.SnapshotTests;

/// <summary>
/// Verify's global settings. A [ModuleInitializer] runs once, before any test in this
/// assembly, which is where Verify wants to be configured.
/// </summary>
public static class VerifySetup
{
    [ModuleInitializer]
    public static void Initialize()
    {
        // Keep the snapshots in one folder next to the tests, and never pop up a diff tool (CI has none).
        Verifier.UseProjectRelativeDirectory("Snapshots");
        DiffEngine.DiffRunner.Disabled = true;
    }
}
