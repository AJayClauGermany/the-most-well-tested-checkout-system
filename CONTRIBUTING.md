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

Each language README has a **"Testing ideas we'd love contributions for"** list.
Pick one, or open an issue with your own idea.

## How the pieces fit together

```
spec/        rules + shared test data (the same for every language)
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
4. Run the tests and checks for every language you touched – each language README
   lists the exact commands.
5. Open a pull request and fill in the template.

## Continuous integration (CI)

CI runs on **pull requests only** (and can be started by hand from the *Actions* tab).
Each workflow only runs when its own folder or the shared `spec/` changes, so a Rust-only
change will not wait for the Java build.

| Workflow                       | Runs when these change                     |
|--------------------------------|--------------------------------------------|
| `spec.yml`                     | `spec/**`                                  |
| `<language>.yml`               | `<language>/**` or `spec/**`               |

## Skipped tests are on purpose

Every language has a test that runs **all** shared spec cases against `checkout(...)`.
Until the checkout is implemented, that test is skipped (`skip`, `@Disabled`,
`#[ignore]`, …). When you start implementing, remove the skip and watch the
tests go from red to green one by one.
