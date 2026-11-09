using System.Text.Json;
using System.Text.Json.Nodes;

using Json.Schema;

namespace CheckoutSystem.TestSupport;

/// <summary>
/// JSON Schema validation with JsonSchema.Net (https://docs.json-everything.net/schema/basics/).
///
/// The schemas in spec/schema/ are the CONTRACT every language's output must honour.
/// A "contract test" checks a real receipt (or saved order, or CLI error) against them.
/// The schemas use JSON Schema draft 2020-12, which JsonSchema.Net supports.
/// </summary>
public static class SchemaValidator
{
    /// <summary>
    /// The schemas are registered under made-up URIs that share one "folder", so that
    /// order.schema.json's <c>"$ref": "receipt.schema.json"</c> resolves to its neighbour,
    /// exactly as it would on a web server.
    /// </summary>
    private static readonly Uri SchemaBase = new("https://checkout.example/spec/schema/");

    private static readonly Dictionary<string, JsonSchema> Schemas = LoadSchemas();

    private static readonly EvaluationOptions Options = new()
    {
        // "List" output: every failing keyword, not just "invalid", so tests can show what is wrong.
        OutputFormat = OutputFormat.List,
    };

    /// <summary>
    /// Validates <paramref name="json"/> against one of the shared schemas.
    /// </summary>
    /// <param name="schemaName">A file name in spec/schema/, e.g. "receipt.schema.json".</param>
    /// <param name="json">The JSON text to check.</param>
    /// <returns>
    /// An empty list when the JSON is valid, otherwise one "where: what" message per problem
    /// (an empty list is easy to assert with <c>Assert.Empty</c> and prints nicely when it fails).
    /// </returns>
    public static IReadOnlyList<string> Errors(string schemaName, string json)
    {
        using var document = JsonDocument.Parse(json);
        return Errors(schemaName, document.RootElement);
    }

    /// <summary>Validates an already parsed <see cref="JsonNode"/>.</summary>
    public static IReadOnlyList<string> Errors(string schemaName, JsonNode? json) =>
        Errors(schemaName, JsonSerializer.SerializeToElement(json));

    /// <summary>Validates an object by serializing it the way the program does (snake_case).</summary>
    public static IReadOnlyList<string> ErrorsFor<T>(string schemaName, T value) =>
        Errors(schemaName, JsonSerializer.SerializeToElement(value, CheckoutJson.Options));

    private static List<string> Errors(string schemaName, JsonElement json)
    {
        if (!Schemas.TryGetValue(schemaName, out var schema))
        {
            throw new ArgumentException($"{schemaName} is not one of the shared schemas", nameof(schemaName));
        }

        var results = schema.Evaluate(json, Options);
        return CollectErrors(results).ToList();
    }

    private static Dictionary<string, JsonSchema> LoadSchemas()
    {
        // One registry knows every schema, so `$ref`s between them can be resolved.
        // receipt.schema.json comes first because order.schema.json refers to it.
        var buildOptions = new BuildOptions { SchemaRegistry = new SchemaRegistry() };
        var schemas = new Dictionary<string, JsonSchema>(StringComparer.Ordinal);
        foreach (var name in new[] { "receipt.schema.json", "order.schema.json", "cli-error.schema.json", "basket.schema.json" })
        {
            var path = Path.Combine(SpecLoader.SpecDirectory, "schema", name);
            var schema = JsonSchema.FromFile(path, buildOptions, new Uri(SchemaBase, name));
            buildOptions.SchemaRegistry.Register(schema);
            schemas[name] = schema;
        }

        return schemas;
    }

    private static IEnumerable<string> CollectErrors(EvaluationResults results)
    {
        if (results.Errors is not null)
        {
            foreach (var (keyword, message) in results.Errors)
            {
                var location = results.InstanceLocation.ToString();
                yield return $"{(location.Length == 0 ? "/" : location)} {keyword}: {message}";
            }
        }

        foreach (var detail in results.Details ?? [])
        {
            foreach (var error in CollectErrors(detail))
            {
                yield return error;
            }
        }
    }
}
