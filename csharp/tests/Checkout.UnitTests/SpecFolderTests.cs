using CheckoutSystem.Cli;

namespace CheckoutSystem.UnitTests;

/// <summary>Unit tests for finding the shared spec/ folder by walking up the folder tree.</summary>
public class SpecFolderTests
{
    [Fact]
    public void Find_FromTheTestOutputFolder_FindsTheFolderThatHoldsCatalogJson()
    {
        var folder = SpecFolder.Find(AppContext.BaseDirectory);

        Assert.Equal("spec", Path.GetFileName(folder));
        Assert.True(File.Exists(Path.Combine(folder, "catalog.json")));
        Assert.True(File.Exists(Path.Combine(folder, "checkout-rules.md")));
    }

    [Fact]
    public void Find_StartingDeeperInsideTheRepository_FindsTheSameFolder()
    {
        var folder = SpecFolder.Find(AppContext.BaseDirectory);

        Assert.Equal(folder, SpecFolder.Find(Path.Combine(folder, "test-cases")));
        Assert.Equal(folder, SpecFolder.Find(Path.Combine(folder, "schema", "nested", "deeper")));
    }

    [Fact]
    public void FindForProgram_FromTheTestOutputFolder_FindsTheSameFolder()
    {
        Assert.Equal(SpecFolder.Find(AppContext.BaseDirectory), SpecFolder.FindForProgram());
    }

    [Fact]
    public void Find_NoSpecFolderAbove_GivesUpWithAClearError()
    {
        var exception = Assert.Throws<DirectoryNotFoundException>(() => SpecFolder.Find(Path.GetTempPath()));

        Assert.StartsWith("No spec/catalog.json found in any folder above", exception.Message, StringComparison.Ordinal);
        Assert.Throws<DirectoryNotFoundException>(() => SpecFolder.Find(Path.GetPathRoot(Path.GetTempPath())!));
    }

    [Fact]
    public void Find_SpecFolderWithoutACatalog_IsNotFooled()
    {
        var folder = Directory.CreateTempSubdirectory("not-the-spec-");
        try
        {
            Directory.CreateDirectory(Path.Combine(folder.FullName, "spec"));

            Assert.Throws<DirectoryNotFoundException>(() => SpecFolder.Find(Path.Combine(folder.FullName, "deeper")));
        }
        finally
        {
            folder.Delete(recursive: true);
        }
    }
}
