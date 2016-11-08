# neupeumeu

A client for npm written in golang.

## Npm registry

All examples are using [beulogue](https://www.npmjs.com/package/beulogue)

- Get a package info: https://registry.npmjs.org/beulogue
- Get info for a specific version: https://registry.npmjs.org/beulogue/4.0.2
- Get info for the latest version: https://registry.npmjs.org/beulogue/latest

## Todo

- [x] List things in cwd package.json
- [x] Download deps in neupeumeu's cache based on that semver thing
- [ ] Don't re-download if matching version is already there
- [x] Install deps in package.json
- [ ] Install deps of deps of deps... in package.json
- [ ] Install a named package and update package.json (deps / devDeps)
- [ ] Cache
- [x] Check deps' shasum
- [x] Semver
- [ ] Run things
- [ ] Native modules (node-gyp...)

![](purple.png)
