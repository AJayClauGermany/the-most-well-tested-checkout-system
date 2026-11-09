# Contributing

Thanks for helping! This project exists to **show and compare ways of testing
software**, using a small checkout system written in six languages. Contributions
of all sizes are welcome – a single new test case is a great first pull request.

## The golden rule: readable first

Our readers are people who are new to programming (late high school, early college,
self-taught). When you write code or tests here, optimise for **understanding**, not
for cleverness or performance:

* Name tests after the behaviour they check: `three_coffees_get_one_free`, not `test_promo_3`.
* Use the **Arrange → Act → Assert** shape and keep each test about one idea.
* Avoid logic (loops, `if`s) inside a test unless the test is *about* data-driven testing.
* Add a short comment saying **what kind of test it is** and **why** it is useful.
* Prefer the language's most common tools. Exotic tools are fine when the point of the
  contribution is to showcase that tool – say so in the README.

## Ways to contribute

| You want to…                                      | Go to                                     |
|---------------------------------------------------|-------------------------------------------|
| Add an example to the shared spec                 | [`spec/test-cases/`](spec/test-cases/) or [`spec/features/`](spec/features/) |
| Show a new kind of test in a language             | that language's folder and its `README.md` |
| Implement part of the checkout (test-first!)      | that language's folder                    |
| Improve explanations                              | [`docs/`](docs/) (published with GitHub Pages) |

The **[testing matrix](docs/testing-matrix.md)** is the list of what we want to show, with
the tool we recommend for each kind of test in each language. Pick an empty cell, or open
an issue with your own idea. Each language README mirrors its column with practical hints.

When you add a kind of test, also tick it in the matrix – in the same pull request.

## How the pieces fit together

```
spec/        rules + shared test data (the same for every language)
e2e/         one end-to-end runner for every language's CLI
python/      ┐
typescript/  │
go/          │ one folder per language, each with its own README,
java/        │ build tool, tests and .gitignore
csharp/      │
rust/        ┘
docs/        the documentation website (Jekyll / GitHub Pages)
.github/     CI workflows (one per language, plus one for the spec)
```

## Working on a change

1. Create a branch.
2. Read [`spec/checkout-rules.md`](spec/checkout-rules.md) if you haven't yet.
3. Make your change. If you change *behaviour*, change the rules first, then the shared
   test cases, then the code. (That is test-driven development in a nutshell.)
   When you touch `spec/checkout-rules.md`, check the two questions in
   [Keep the spec declarative](#keep-the-spec-declarative).
4. Run the tests and checks for every language you touched – each language README
   lists the exact commands.
5. Open a pull request and fill in the template.

## Keep the spec declarative

The shared spec says **what** the checkout does, never **how** a program should be
built. That is what lets six languages each solve it in their own idiomatic way and still
be compared. Before you change `spec/checkout-rules.md`, `spec/features/` or
`spec/properties.md`, ask:

1. **Does it describe *what*, not *how*?** Inputs, outputs, error codes, amounts. Not
   steps, loops, classes, method names or the order in which things are computed.
2. **Could two languages implement it differently and both pass every shared test?**
   If only one program structure can satisfy the sentence, it belongs in a 💡 *Hint*, in a
   language README, or nowhere.

Names like `save` and `find` are vocabulary so scenarios and properties can refer to a
capability; they are not required identifiers. The Gherkin step sentences are the one
place where fixed wording is required – and they describe behaviour, not code.

## Continuous integration (CI)

CI runs on **pull requests only** (and can be started by hand from the *Actions* tab).
Each workflow only runs when its own folder or the shared `spec/` changes, so a Rust-only
change will not wait for the Java build.

| Workflow                       | Runs when these change                     |
|--------------------------------|--------------------------------------------|
| `spec.yml`                     | `spec/**`, `e2e/**`                        |
| `<language>.yml`               | `<language>/**`, `spec/**` or `e2e/**`     |

## Tests first – and the git history shows it

Every language has a test that runs **all** shared spec cases against `checkout(...)`.
While a checkout did not exist yet, that test was skipped (`skip`, `@Disabled`,
`#[ignore]`, …); the first commit of each implementation removed the skip and added the
whole test suite, red, and the next commit made it green. Run
`git log --oneline -- <language folder>` to see it.

Keep working that way: when you change behaviour, change the rules and the shared cases
first, watch the six implementations go red, then fix them. Coverage and mutation testing
run in every language's CI, so a change that is not tested shows up there.
