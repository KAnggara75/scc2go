# GitHub Actions SSH Signing - Changes Summary

## What Was Fixed

### 1. SSH Agent Configuration
- ✅ **Fixed passphrase handling**: Using `SSH_ASKPASS` for non-interactive passphrase input
- ✅ **Exported SSH agent variables**: `SSH_AUTH_SOCK` and `SSH_AGENT_PID` to `GITHUB_ENV`
- ✅ **Proper key permissions**: Set correct permissions on SSH keys and known_hosts

### 2. Git Configuration
- ✅ **Correct signing key**: Using `.pub` file for `user.signingkey` (not private key)
- ✅ **SSH remote**: Configured git to use SSH instead of HTTPS
- ✅ **Consistent credentials**: Using `pakaiwa-release-bot` throughout (not github-actions[bot])

### 3. Security Improvements
- ✅ **Environment variables**: Moved secrets to `env:` block to avoid shell injection
- ✅ **Public key generation**: Automatically generates `.pub` from private key
- ✅ **SSH key scanning**: Adds github.com to known_hosts

### 4. Workflow Optimizations
- ✅ **Path filtering**: Only runs on Go file changes
- ✅ **Permissions**: Explicitly set required permissions
- ✅ **Codecov v5**: Updated to latest version with token support

## Key Changes in `.github/workflows/go.yml`

### Before (Issues)
```yaml
# SSH key with base64 encoding (unnecessary complexity)
echo "${{ secrets.RELEASE_BOT_SSH_KEY_B64 }}" | base64 -d > ~/.ssh/id_ed25519

# Passphrase with heredoc (doesn't work in non-interactive)
ssh-add ~/.ssh/id_ed25519 <<< "${{ secrets.RELEASE_BOT_SSH_PASSPHRASE }}"

# Wrong signing key (private key instead of public)
git config --global user.signingkey ~/.ssh/id_ed25519

# Inconsistent user credentials
git config user.name "github-actions[bot]"  # Different from SSH setup
```

### After (Fixed)
```yaml
# Plain SSH key (simpler and more direct)
env:
  SSH_PRIVATE_KEY: ${{ secrets.RELEASE_BOT_SSH_KEY }}
  SSH_PASSPHRASE: ${{ secrets.RELEASE_BOT_SSH_PASSPHRASE }}

# Write key directly
echo "$SSH_PRIVATE_KEY" > ~/.ssh/id_ed25519

# Non-interactive passphrase using SSH_ASKPASS
SSH_ASKPASS=/tmp/ssh-askpass.sh DISPLAY=:0 ssh-add ~/.ssh/id_ed25519 < /dev/null

# Correct signing key (public key)
git config --global user.signingkey ~/.ssh/id_ed25519.pub

# Generate public key
ssh-keygen -y -f ~/.ssh/id_ed25519 > ~/.ssh/id_ed25519.pub

# Consistent user credentials
git config user.name "pakaiwa-release-bot"  # Same as SSH setup

# Explicit SSH remote
git remote set-url origin git@github.com:KAnggara75/scc2go.git
```

## Testing the Changes

### Local Testing (Optional)

You can test the SSH signing locally:

```bash
# Setup (one time)
export SSH_PRIVATE_KEY="$(cat ~/.ssh/pakaiwa_release_bot)"
export SSH_PASSPHRASE="<your-passphrase>"

# Write key to temp file
echo "$SSH_PRIVATE_KEY" > /tmp/test_key
chmod 600 /tmp/test_key

# Generate public key
ssh-keygen -y -f /tmp/test_key > /tmp/test_key.pub

# Configure git
git config user.signingkey /tmp/test_key.pub
git config gpg.format ssh
git config tag.gpgsign true

# Create a test signed tag
git tag -s test-tag -m "Test signed tag"

# Verify
git tag -v test-tag

# Cleanup
git tag -d test-tag
rm /tmp/test_key /tmp/test_key.pub
```

### GitHub Actions Testing

1. **Push a commit** to trigger the workflow:
   ```bash
   git add .
   git commit -m "test: verify SSH signing workflow"
   git push origin main
   ```

2. **Check the workflow**:
   - Go to Actions tab
   - Watch the "Setup SSH signing key" step
   - Verify no errors in SSH agent setup

3. **Verify the tag**:
   - After workflow completes, check the new tag
   - It should show "Verified" badge on GitHub
   - Click the tag to see signature details

## Expected Workflow Output

### Successful SSH Setup
```
✓ Setup SSH directory and key
✓ Start SSH agent and add key with passphrase
✓ Configure git for SSH signing
✓ Create public key for signing
✓ Configure git to use SSH instead of HTTPS
✓ Export SSH_AUTH_SOCK for subsequent steps
```

### Successful Tag Creation
```
✓ Created signed tag v0.2.0
✓ Pushed signed tag: v0.2.0
```

### On GitHub
- Tag shows "Verified" badge
- Signature shows: "pakaiwa-release-bot"
- Release created automatically

## Troubleshooting

If the workflow fails, check:

1. **Secrets are set correctly**:
   - `RELEASE_BOT_SSH_KEY` (plain text with header/footer)
   - `RELEASE_BOT_SSH_PASSPHRASE`
   - `CODECOV_TOKEN`

2. **SSH public key is added to GitHub**:
   - Settings → SSH and GPG keys
   - Key type must be "Signing Key"

3. **Workflow logs**:
   - Check "Setup SSH signing key" step
   - Look for SSH agent errors
   - Verify git config output

## Next Steps

1. ✅ Verify secrets are configured
2. ✅ Add SSH public key to GitHub as Signing Key
3. ✅ Push a commit to trigger the workflow
4. ✅ Verify the tag is created and signed
5. ✅ Check the GitHub release is created

## Documentation

- [SSH_SIGNING_SETUP.md](SSH_SIGNING_SETUP.md) - Detailed setup guide
- [VERSIONING.md](VERSIONING.md) - Semantic versioning guide
- [README.md](README.md) - Project documentation
