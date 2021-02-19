# Contributing Guidelines

_These are mostly for myself, so I don't forget the project conventions._

## Commits

Make sure that [pre-commit](https://pre-commit.com) is installed and running for every commit. Try `pre-commit install` to make this happen.

Follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard for commit messages and use the following commit types (based on [Angular](https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelines)):

- build: Changes that affect the build system or external dependencies (example scopes: gulp, broccoli, npm) [Angular]
- ci: Changes to our CI configuration files and scripts (example scopes: Travis, Circle, BrowserStack, SauceLabs) [Angular]
- docs: Documentation only changes [Angular]
- feat: A new feature [Angular]
- fix: A bug fix [Angular]
- perf: A code change that improves performance [Angular]
- refactor: A code change that neither fixes a bug nor adds a feature [Angular]
- style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc) [Angular]
- test: Adding missing tests or correcting existing tests [Angular]

## Languages Used

- `client/` - TypeScript, see `client/package.json` for the language version
- `data/` - Python, see `data/.python-version` for the language version
- `server/` - Go, see `server/go.mod` for the language version

## License

All contributions must be explicitly dedicated to the public domain. The [Unlicense](https://unlicense.org) recommends submitting a signed variant of the following statement:

> I dedicate any and all copyright interest in this software to the
> public domain. I make this dedication for the benefit of the public at
> large and to the detriment of my heirs and successors. I intend this
> dedication to be an overt act of relinquishment in perpetuity of all
> present and future rights to this software under copyright law.
