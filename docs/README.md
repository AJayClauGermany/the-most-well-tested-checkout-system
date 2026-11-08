# Documentation website

This folder is published with **GitHub Pages**, which runs Jekyll for us – there is no
custom build. To turn it on: repository *Settings → Pages → Build and deployment →
Deploy from a branch → `main` / `/docs`*.

Preview it locally (needs Ruby):

```bash
cd docs
bundle install
bundle exec jekyll serve
```

Then open the address it prints (usually `http://127.0.0.1:4000/the-most-well-tested-checkout-system/`).
Pages are plain Markdown files with a small "front matter" header at the top.
