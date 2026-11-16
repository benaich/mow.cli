# CLI extensions

This branch extends [jawher/mow.cli](https://github.com/jawher/mow.cli) v1.2.0.
The upstream module path and license are preserved.

- Standalone `-` values for string options, including `--file -` and `-f -`.
- `Cmd.Output`, an inherited per-app writer for help and parser diagnostics.
- `Cmd.HelpPrinter`, an inherited optional custom help renderer.
- `Cmd.Describe`, offline command declaration metadata without argument parsing,
  environment reads or command actions.

These extensions support offline command discovery, completion and structured
error output. The existing default help and parsing behavior is retained.
Regression tests accompany the extensions.

## Upstream contribution

- [ ] Discuss these features and suitable APIs with Jawher.
- [ ] Propose the extensions upstream.
- [ ] Switch consumers back to an upstream release when it supports their needs.
