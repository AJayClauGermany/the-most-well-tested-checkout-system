namespace CheckoutSystem.Cli;

/// <summary>
/// Finds the shared <c>spec/</c> folder at the root of the repository.
/// </summary>
public static class SpecFolder
{
    /// <summary>
    /// Walks up from <paramref name="start"/> until it reaches a folder that contains
    /// <c>spec/catalog.json</c>, and returns that <c>spec</c> folder. Looking upwards
    /// (instead of hard-coding <c>../../spec</c>) means the code finds the spec from the
    /// compiled program (<c>bin/Debug/net10.0/</c>), from a test's output folder and from
    /// any working directory inside the repository.
    /// </summary>
    /// <exception cref="DirectoryNotFoundException">No folder above <paramref name="start"/> contains <c>spec/catalog.json</c>.</exception>
    public static string Find(string start) =>
        TryFind(start) ?? throw new DirectoryNotFoundException($"No spec/catalog.json found in any folder above {start}");

    /// <summary>
    /// The spec folder for the running program: first looks up from where the program is
    /// installed (<see cref="AppContext.BaseDirectory"/>), then from the current directory.
    /// </summary>
    public static string FindForProgram() =>
        TryFind(AppContext.BaseDirectory) ?? Find(Directory.GetCurrentDirectory());

    private static string? TryFind(string start)
    {
        for (var folder = new DirectoryInfo(start); folder is not null; folder = folder.Parent)
        {
            var spec = Path.Combine(folder.FullName, "spec");
            if (File.Exists(Path.Combine(spec, "catalog.json")))
            {
                return spec;
            }
        }

        return null;
    }
}
