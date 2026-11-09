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
        throw new NotImplementedException();

    /// <summary>
    /// The spec folder for the running program: first looks up from where the program is
    /// installed (<see cref="AppContext.BaseDirectory"/>), then from the current directory.
    /// </summary>
    public static string FindForProgram() =>
        throw new NotImplementedException();
}
