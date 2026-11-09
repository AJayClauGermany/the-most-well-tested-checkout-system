namespace CheckoutSystem.IntegrationTests;

/// <summary>
/// A fresh, empty folder under the system's temp directory, deleted when disposed.
///
/// xUnit creates a new instance of a test class for every test, and calls <c>Dispose()</c>
/// afterwards when the class implements <see cref="IDisposable"/>. So a test class that
/// creates one of these in its constructor gets its own folder per test and a clean-up for
/// free – the xUnit way of "setup and teardown".
/// </summary>
public sealed class TemporaryFolder : IDisposable
{
    private readonly DirectoryInfo _directory = Directory.CreateTempSubdirectory("checkout-tests-");

    public string Path => _directory.FullName;

    /// <summary>A path inside the folder (the file or sub-folder is not created).</summary>
    public string PathOf(params string[] parts) => System.IO.Path.Combine([Path, .. parts]);

    public void Dispose() => _directory.Delete(recursive: true);
}
