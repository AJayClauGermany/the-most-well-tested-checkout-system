using ArchUnitNET.Domain;
using ArchUnitNET.Fluent;
using ArchUnitNET.Loader;
using ArchUnitNET.xUnitV3;

using CheckoutSystem.Cli;
using CheckoutSystem.Storage;

using static ArchUnitNET.Fluent.ArchRuleDefinition;

namespace CheckoutSystem.ArchitectureTests;

/// <summary>
/// ARCHITECTURE TESTS: rules about the shape of the code, not its behaviour.
///
/// The checkout (the "domain", namespace CheckoutSystem in src/Checkout) is pure
/// calculation: it must not read files, open databases or know about the command line.
/// Those live in Checkout.Storage and Checkout.Cli, which depend on the domain – never the
/// other way round. A human reviewer might miss a stray <c>File.ReadAllText</c>; ArchUnitNET
/// reads the compiled assemblies and does not.
/// </summary>
public class DependencyRulesTests
{
    // Loading the assemblies once for all tests in this class.
    private static readonly Architecture Architecture = new ArchLoader()
        .LoadAssemblies(typeof(Checkout).Assembly, typeof(JsonFileOrderRepository).Assembly, typeof(CheckoutCommand).Assembly)
        .Build();

    [Fact]
    public void Domain_DoesNotDependOnSystemIO()
    {
        IArchRule rule = Types().That().ResideInAssembly(typeof(Checkout).Assembly)
            .Should().NotDependOnAnyTypesThat().ResideInNamespaceOrBelow("System.IO")
            .Because("the checkout is a calculation; files and folders belong to Checkout.Storage and the CLI");

        rule.Check(Architecture);
    }

    [Fact]
    public void Domain_DoesNotDependOnTheCliOrTheStorage()
    {
        IArchRule rule = Types().That().ResideInAssembly(typeof(Checkout).Assembly)
            .Should().NotDependOnAnyTypesThat().ResideInAssembly(typeof(CheckoutCommand).Assembly)
            .AndShould().NotDependOnAnyTypesThat().ResideInAssembly(typeof(JsonFileOrderRepository).Assembly)
            .Because("dependencies point inwards: the CLI and the stores use the checkout, not the reverse");

        rule.Check(Architecture);
    }

    [Fact]
    public void Domain_DoesNotDependOnTheConsole()
    {
        IArchRule rule = Types().That().ResideInAssembly(typeof(Checkout).Assembly)
            .Should().NotDependOnAnyTypesThat().HaveFullName("System.Console")
            .Because("printing is the CLI's job");

        rule.Check(Architecture);
    }

    [Fact]
    public void Storage_DoesNotDependOnTheCli()
    {
        IArchRule rule = Types().That().ResideInAssembly(typeof(JsonFileOrderRepository).Assembly)
            .Should().NotDependOnAnyTypesThat().ResideInAssembly(typeof(CheckoutCommand).Assembly)
            .Because("a store should be usable from any program, not only this command line");

        rule.Check(Architecture);
    }

    [Fact]
    public void OrderRepositories_ImplementTheInterface()
    {
        IArchRule rule = Classes().That().HaveNameEndingWith("OrderRepository")
            .Should().ImplementInterface(typeof(IOrderRepository))
            .Because("rules §6 names the capabilities save and find; every store must offer both");

        rule.Check(Architecture);
    }
}
