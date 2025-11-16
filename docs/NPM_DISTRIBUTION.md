# NPM Distribution Guide

This guide explains how to distribute cntm via npm, making it installable with `npx`.

## How It Works

The npm package is a **wrapper** that:
1. Downloads the appropriate Go binary for the user's platform during `npm install`
2. Provides a Node.js wrapper script to execute the binary
3. Exposes a JavaScript API for programmatic usage

This is the same pattern used by tools like `esbuild`, `swc`, and `prisma`.

## Directory Structure

```
npm/
├── package.json          # NPM package metadata
├── install.js            # Post-install script to download binary
├── index.js              # JavaScript API
├── bin/
│   └── cntm.js          # Wrapper script for CLI
├── README.md            # NPM package README
├── test.js              # Simple tests
└── .npmignore           # Files to exclude from npm
```

## Installation Methods for Users

### 1. Using npx (No Installation)

```bash
npx cntm search "code review"
npx cntm install code-reviewer
```

### 2. Using npx with GitHub

```bash
npx github:nghiadoan-work/claude-nia-tool-management-cli
```

### 3. Global Installation

```bash
npm install -g cntm
cntm --help
```

### 4. Local Project Installation

```bash
npm install cntm
npx cntm init
```

## Publishing to NPM

### Prerequisites

1. NPM account: https://www.npmjs.com/signup
2. Login to npm:
   ```bash
   npm login
   ```

### First Time Publishing

1. **Build the Go binaries first:**
   ```bash
   cd /Volumes/ex-macmini-a/claude_projects/agent_skill_cli_go
   ./scripts/build.sh
   ```

2. **Create GitHub release** with the binaries (this is required because install.js downloads from GitHub releases):
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0

   gh release create v1.0.0 \
     --title "cntm v1.0.0" \
     --notes-file docs/RELEASE.md \
     dist/*.tar.gz \
     dist/*.zip \
     dist/checksums.txt
   ```

3. **Publish to npm:**
   ```bash
   cd npm/
   npm publish
   ```

### Updating the Package

When releasing a new version:

1. Update version in both:
   - `package.json` (in root)
   - `npm/package.json`
   - `pkg/version/version.go`

2. Build binaries:
   ```bash
   ./scripts/build.sh
   ```

3. Create GitHub release:
   ```bash
   git tag -a v1.0.1 -m "Release v1.0.1"
   git push origin v1.0.1
   gh release create v1.0.1 ...
   ```

4. Publish to npm:
   ```bash
   cd npm/
   npm publish
   ```

## Automated Publishing with GitHub Actions

Create `.github/workflows/publish.yml`:

```yaml
name: Publish

on:
  release:
    types: [published]

jobs:
  publish-npm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          registry-url: 'https://registry.npmjs.org'

      - name: Publish to NPM
        run: |
          cd npm
          npm publish
        env:
          NODE_AUTH_TOKEN: ${{secrets.NPM_TOKEN}}
```

## Testing the NPM Package Locally

### Test Installation

```bash
cd npm/

# Test the install script
node install.js

# Test the binary wrapper
node bin/cntm.js --help

# Run tests
npm test
```

### Test Programmatic API

```javascript
// test-api.js
const cntm = require('./index.js');

(async () => {
  // Get version
  const version = await cntm.version();
  console.log('Version:', version);

  // Search
  const results = await cntm.search('code');
  console.log('Search results:', results);

  // List installed tools
  const tools = await cntm.list();
  console.log('Installed tools:', tools);
})();
```

### Test with npm link

```bash
cd npm/
npm link

# Now you can use it globally
cntm --help

# Unlink when done
npm unlink -g cntm
```

## Package Scopes (Optional)

If you want to publish under a scope (e.g., `@yourname/cntm`):

1. Update `npm/package.json`:
   ```json
   {
     "name": "@yourname/cntm",
     ...
   }
   ```

2. Publish with public access:
   ```bash
   npm publish --access public
   ```

3. Users install with:
   ```bash
   npm install -g @yourname/cntm
   ```

## Troubleshooting

### Installation Fails

If users report installation failures:

1. **Check GitHub release exists:**
   - Verify the release tag matches the npm version
   - Verify all platform binaries are uploaded

2. **Check download URL:**
   - The install script constructs URLs like:
     `https://github.com/USER/REPO/releases/download/v1.0.0/cntm-1.0.0-darwin-arm64.tar.gz`

3. **Platform not supported:**
   - Add error message suggesting manual installation

### Binary Not Found

If users get "binary not found" errors:

1. Check `bin/` directory was created
2. Verify postinstall script ran
3. Check file permissions (Unix systems)

### Version Mismatch

Keep these versions in sync:
- `npm/package.json` - npm package version
- `pkg/version/version.go` - Go binary version
- Git tags (e.g., `v1.0.0`)
- GitHub release version

## Best Practices

1. **Always test locally before publishing:**
   ```bash
   cd npm/
   npm pack
   # This creates a .tgz file you can test with:
   npm install -g ./cntm-1.0.0.tgz
   ```

2. **Use semantic versioning:**
   - MAJOR.MINOR.PATCH (e.g., 1.0.0)
   - Sync with Git tags

3. **Keep README updated:**
   - Update examples
   - Add new features
   - Update links

4. **Tag releases properly:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

5. **Verify before publishing:**
   - Test on all platforms (or use CI)
   - Check all download URLs work
   - Verify checksums

## NPM Package Statistics

After publishing, you can track:

- Downloads: https://npm-stat.com/charts.html?package=cntm
- Package page: https://www.npmjs.com/package/cntm
- Bundle size: https://bundlephobia.com/package/cntm

## Alternative: GitHub Packages

You can also publish to GitHub Packages instead of (or in addition to) npm:

1. Create `.npmrc`:
   ```
   @yourname:registry=https://npm.pkg.github.com
   ```

2. Update package.json:
   ```json
   {
     "name": "@yourname/cntm",
     "repository": {
       "type": "git",
       "url": "https://github.com/yourname/claude-nia-tool-management-cli.git"
     }
   }
   ```

3. Publish:
   ```bash
   npm publish --registry=https://npm.pkg.github.com
   ```

## Support Matrix

| Platform | Architecture | Supported |
|----------|--------------|-----------|
| macOS    | x64 (Intel)  | ✅ |
| macOS    | ARM64 (M1/M2)| ✅ |
| Linux    | x64          | ✅ |
| Linux    | ARM64        | ✅ |
| Windows  | x64          | ✅ |
| Windows  | ARM64        | ❌ (future) |

## Resources

- [NPM Documentation](https://docs.npmjs.com/)
- [Publishing Packages](https://docs.npmjs.com/packages-and-modules/contributing-packages-to-the-registry)
- [Package.json Guide](https://docs.npmjs.com/cli/v9/configuring-npm/package-json)
- [Semantic Versioning](https://semver.org/)
