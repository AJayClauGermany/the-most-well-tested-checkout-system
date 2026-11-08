---
layout: page
title: Contributing
---

The full guide is in
[`CONTRIBUTING.md`]({{ site.repository_files }}/CONTRIBUTING.md). The short version:

1. **Readable first.** Write tests a beginner can understand: clear names,
   Arrange → Act → Assert, one idea per test, a comment saying what kind of test it is.
2. **Rules before code.** To change behaviour, update
   [`spec/checkout-rules.md`]({{ site.repository_files }}/spec/checkout-rules.md) first,
   then the shared test cases, then the implementations.
3. **Run it locally.** Each language README lists the exact commands CI runs.
4. **Open a pull request.** CI runs only for the folders you changed.

## Good first contributions

* Add a boundary case to one of the [shared test case files]({{ site.repository_files }}/spec/test-cases/).
* Pick an item from a language's *"Testing ideas we'd love contributions for"* list.
* Improve an explanation on this website – these pages live in the
  [`docs/`]({{ site.repository_files }}/docs/) folder.

## Editing this website

Each page is a Markdown file in `docs/` with a small header:

```markdown
---
layout: page
title: My new page
---

Your content here.
```

To show it in the top menu, add it to `header_pages` in `docs/_config.yml`.
GitHub Pages rebuilds the site automatically when the change is merged.
