# Troubleshooting Guide

Common issues and solutions for cntm.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Network Errors](#network-errors)
- [Authentication Errors](#authentication-errors)
- [Tool Installation Errors](#tool-installation-errors)
- [Lock File Issues](#lock-file-issues)
- [Cache Problems](#cache-problems)
- [Permission Errors](#permission-errors)

---

## Installation Issues

### "Tool not found in registry"

**Symptom:**
```
Error: tool 'xyz' not found in registry
```

**Solutions:**

1. **Search for similar tools:**
   ```bash
   cntm search xyz
   ```

2. **Check tool name spelling**

3. **Verify registry URL:**
   ```bash
   cat .claude-config.yaml | grep url
   ```

4. **Refresh cache:**
   ```bash
   rm -rf .claude/.cache
   cntm search xyz
   ```

### "Already installed"

**Symptom:**
```
Warning: Tool xyz is already installed (version 1.0.0)
Hint: Use --force to reinstall
```

**Solutions:**

1. **Force reinstall:**
   ```bash
   cntm install --force xyz
   ```

2. **Update instead:**
   ```bash
   cntm update xyz
   ```

3. **Check installed version:**
   ```bash
   cntm list
   ```

### "Integrity check failed"

**Symptom:**
```
Error: Integrity check failed for xyz.zip
Hint: The downloaded file may be corrupted. Try again...
```

**Solutions:**

1. **Retry the installation:**
   ```bash
   cntm install xyz
   ```

2. **Clear cache and retry:**
   ```bash
   rm -rf .claude/.cache
   cntm install xyz
   ```

3. **Check network connection**

4. **Report to tool author if persists**

---

## Network Errors

### "Connection timeout"

**Symptom:**
```
Error: Network error during download
Hint: Check your internet connection and try again
```

**Solutions:**

1. **Check internet connection:**
   ```bash
   ping github.com
   ```

2. **Use authentication to increase limits:**
   ```bash
   export CNTM_GITHUB_TOKEN=ghp_xxx
   cntm install xyz
   ```

3. **Try again later** (GitHub may be down)

### "Rate limit exceeded"

**Symptom:**
```
Error: GitHub API rate limit exceeded
Remaining: 0/60, resets at: 2024-01-15 10:30:00
```

**Solutions:**

1. **Authenticate to increase limit (60 → 5000):**
   ```bash
   export CNTM_GITHUB_TOKEN=ghp_xxx
   cntm install xyz
   ```

2. **Wait for rate limit reset**

3. **Use cached data:**
   ```bash
   cntm list --remote  # Uses cache if available
   ```

---

## Authentication Errors

### "Authentication failed"

**Symptom:**
```
Error: Authentication failed
Hint: Check your GitHub token in the config file or CNTM_GITHUB_TOKEN environment variable
```

**Solutions:**

1. **Set GitHub token:**
   ```bash
   export CNTM_GITHUB_TOKEN=ghp_xxxxxxxxxxxxx
   ```

2. **Generate new token:**
   - Go to: https://github.com/settings/tokens
   - Click "Generate new token (classic)"
   - Select scopes: `repo` or `public_repo`
   - Copy and export token

3. **Verify token is set:**
   ```bash
   echo $CNTM_GITHUB_TOKEN
   ```

4. **Check token hasn't expired**

### "Permission denied (private registry)"

**Symptom:**
```
Error: 404 Not Found
```

**Solutions:**

1. **Ensure token has `repo` scope** (not just `public_repo`)

2. **Verify you have access to the repository**

3. **Check registry URL is correct:**
   ```bash
   cat .claude-config.yaml | grep url
   ```

---

## Tool Installation Errors

### "Failed to create directory"

**Symptom:**
```
Error: Permission denied: cannot write to .claude/agents/
```

**Solutions:**

1. **Check directory permissions:**
   ```bash
   ls -la .claude
   ```

2. **Fix permissions:**
   ```bash
   chmod -R u+w .claude
   ```

3. **Use custom path:**
   ```bash
   cntm install --path ~/tools xyz
   ```

### "ZIP extraction failed"

**Symptom:**
```
Error: Failed to extract ZIP file
```

**Solutions:**

1. **Check disk space:**
   ```bash
   df -h
   ```

2. **Retry installation:**
   ```bash
   cntm install --force xyz
   ```

3. **Check ZIP isn't corrupted** (integrity error would show)

### "Invalid tool structure"

**Symptom:**
```
Error: Tool structure is invalid
```

**Solutions:**

1. **Contact tool author** (tool ZIP is malformed)

2. **Try different version:**
   ```bash
   cntm install xyz@1.0.0
   ```

---

## Lock File Issues

### "Lock file corrupted"

**Symptom:**
```
Error: Failed to parse lock file
```

**Solutions:**

1. **Backup lock file:**
   ```bash
   cp .claude-lock.json .claude-lock.json.backup
   ```

2. **Reinitialize (WARNING: loses installation history):**
   ```bash
   rm .claude-lock.json
   cntm init
   # Reinstall tools manually
   ```

3. **Manually fix JSON:**
   ```bash
   cat .claude-lock.json | jq .  # Check for syntax errors
   ```

### "Lock file not found"

**Symptom:**
```
Error: .claude-lock.json not found
Hint: Run 'cntm init' to initialize the project
```

**Solutions:**

1. **Initialize project:**
   ```bash
   cntm init
   ```

2. **Check you're in correct directory:**
   ```bash
   pwd
   ls -la | grep .claude
   ```

---

## Cache Problems

### "Stale cache data"

**Symptom:**
Tool appears in `cntm search` but fails to install (was removed from registry).

**Solutions:**

1. **Clear cache:**
   ```bash
   rm -rf .claude/.cache
   ```

2. **Refresh cache:**
   ```bash
   cntm list --remote
   ```

### "Cache permission errors"

**Symptom:**
```
Error: Failed to write to cache
```

**Solutions:**

1. **Fix cache permissions:**
   ```bash
   chmod -R u+w .claude/.cache
   ```

2. **Disable cache temporarily:**
   ```bash
   export CNTM_CACHE_ENABLED=false
   cntm search xyz
   ```

---

## Permission Errors

### "Cannot write to .claude directory"

**Symptom:**
```
Error: Permission denied: cannot write to .claude/agents/xyz/
```

**Solutions:**

1. **Check ownership:**
   ```bash
   ls -la .claude
   ```

2. **Fix ownership:**
   ```bash
   sudo chown -R $USER:$USER .claude
   ```

3. **Fix permissions:**
   ```bash
   chmod -R u+w .claude
   ```

### "Cannot create .claude directory"

**Symptom:**
```
Error: Failed to create .claude directory
```

**Solutions:**

1. **Check parent directory permissions:**
   ```bash
   ls -la .
   ```

2. **Run init with sudo (not recommended):**
   ```bash
   cntm init
   sudo chown -R $USER:$USER .claude
   ```

---

## Common Error Messages

### Error Message Reference

| Error | Likely Cause | Solution |
|-------|--------------|----------|
| "Tool not found" | Wrong name or not in registry | `cntm search <name>` |
| "Already installed" | Tool exists locally | Use `--force` flag |
| "Authentication failed" | No/invalid GitHub token | Set `CNTM_GITHUB_TOKEN` |
| "Rate limit exceeded" | Too many unauthenticated requests | Authenticate with token |
| "Network error" | Connection issues | Check internet, retry |
| "Integrity check failed" | Corrupted download | Retry installation |
| "Permission denied" | File permission issues | Fix with `chmod`/`chown` |
| "Lock file not found" | Not initialized | Run `cntm init` |

---

## Debugging Tips

### Enable Verbose Logging

```bash
# Set environment variable for debug output
export CNTM_DEBUG=true
cntm install xyz
```

### Check System Status

```bash
# Verify cntm installation
cntm --version

# Check config
cat .claude-config.yaml

# List installed tools
cntm list

# Check outdated tools
cntm outdated

# Verify cache
ls -la .claude/.cache
```

### Network Diagnostics

```bash
# Test GitHub connectivity
ping github.com

# Test HTTPS
curl -I https://github.com

# Check GitHub status
curl https://www.githubstatus.com/api/v2/status.json
```

### Clean Reinstall

```bash
# Backup
tar -czf claude-backup.tar.gz .claude

# Clean
rm -rf .claude

# Reinitialize
cntm init

# Reinstall tools
cntm install tool1 tool2 tool3
```

---

## Getting Help

### Before Asking for Help

1. **Search existing issues:**
   - GitHub Issues
   - Documentation

2. **Gather information:**
   ```bash
   cntm --version
   uname -a  # OS info
   cat .claude-config.yaml
   cat .claude-lock.json | jq .
   ```

3. **Try debug mode:**
   ```bash
   export CNTM_DEBUG=true
   cntm <command> 2>&1 | tee debug.log
   ```

### How to Report Issues

Include:
- cntm version
- Operating system
- Full error message
- Steps to reproduce
- Config file (without tokens!)
- Debug log

Example:

```
**Environment:**
- cntm version: 1.0.0
- OS: macOS 14.0
- Go version: 1.21

**Issue:**
Installation fails with integrity error

**Steps to reproduce:**
1. cntm init
2. cntm install code-reviewer

**Error:**
Error: Integrity check failed for code-reviewer.zip

**Config:**
registry:
  url: "https://github.com/claude-tools/registry"

**Debug log:**
[attached debug.log]
```

---

## FAQ

**Q: Can I use cntm offline?**
A: Partially. You can use cached data for search/browse, but cannot install new tools without internet.

**Q: How do I update cntm itself?**
A: Download the latest release from GitHub and replace your binary.

**Q: Can I have multiple registries?**
A: Not simultaneously in v1.0. Use `--config` flag to switch between registry configs.

**Q: Where are tools installed?**
A: Default: `.claude/<type>/<name>/`. Customize with `--path` flag.

**Q: How do I uninstall cntm?**
A: Remove the binary and optionally `~/.claude-config.yaml`.

**Q: Is it safe to commit .claude-lock.json?**
A: Yes! It tracks installed tools for reproducibility.

**Q: Why is search slow?**
A: First search fetches registry. Subsequent searches use cache. Use `cntm list --remote` to warm up cache.

---

## Still Having Issues?

- Check [GitHub Issues](https://github.com/your-org/cntm/issues)
- Read [Configuration Guide](CONFIGURATION.md)
- Review [Command Reference](COMMANDS.md)
- Ask for help in [Discussions](https://github.com/your-org/cntm/discussions)
