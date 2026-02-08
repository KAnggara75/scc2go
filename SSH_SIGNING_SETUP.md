# SSH-Signed Tag Setup Guide

This document explains how the GitHub Actions workflow creates verified (SSH-signed) tags automatically.

## Overview

The workflow uses SSH key signing to create verified tags on GitHub. This provides cryptographic proof that tags were created by an authorized bot account.

## Required Secrets

You need to configure the following secrets in your GitHub repository settings:

### 1. `RELEASE_BOT_SSH_KEY_B64`

This is your SSH private key, base64-encoded.

**To create:**

```bash
# Generate a new SSH key with passphrase
ssh-keygen -t ed25519 -C "release@pakaiwa.dev" -f ~/.ssh/pakaiwa_release_bot

# Enter a strong passphrase when prompted

# Base64 encode the private key
cat ~/.ssh/pakaiwa_release_bot | base64 | tr -d '\n'
```

Copy the output and add it as a repository secret named `RELEASE_BOT_SSH_KEY_B64`.

### 2. `RELEASE_BOT_SSH_PASSPHRASE`

This is the passphrase you used when creating the SSH key.

Add it as a repository secret named `RELEASE_BOT_SSH_PASSPHRASE`.

### 3. `CODECOV_TOKEN`

Your Codecov upload token for coverage reporting.

Get this from [codecov.io](https://codecov.io) after setting up your repository.

## GitHub Configuration

### Add SSH Public Key to GitHub

1. Extract the public key:
   ```bash
   cat ~/.ssh/pakaiwa_release_bot.pub
   ```

2. Go to GitHub Settings → SSH and GPG keys
3. Click "New SSH key"
4. Title: "pakaiwa-release-bot signing key"
5. Key type: **Signing Key** (important!)
6. Paste the public key content
7. Click "Add SSH key"

### Configure Bot Account (Optional but Recommended)

For better security, create a dedicated bot account:

1. Create a new GitHub account (e.g., `pakaiwa-release-bot`)
2. Add the SSH signing key to that account
3. Give the bot account write access to your repository
4. Update the workflow to use a Personal Access Token (PAT) from the bot account

## How It Works

### Workflow Steps

1. **Checkout**: Fetches the repository with full history
2. **Setup SSH Signing Key**:
   - Decodes the base64-encoded SSH private key
   - Sets up SSH agent with the passphrase
   - Configures git to use SSH signing
   - Creates the public key for verification
   - Configures git to use SSH instead of HTTPS
3. **Version Calculation**: Determines the next version based on commits
4. **Tag Creation**: Creates a signed tag using `git tag -s`
5. **Push Tag**: Pushes the signed tag via SSH
6. **GitHub Release**: Creates a release with changelog

### SSH Agent Configuration

The workflow uses `SSH_ASKPASS` to provide the passphrase non-interactively:

```bash
# Create a script that echoes the passphrase
cat > /tmp/ssh-askpass.sh << 'EOF'
#!/bin/bash
echo "$SSH_PASSPHRASE"
EOF

# Use it to add the key
SSH_ASKPASS=/tmp/ssh-askpass.sh DISPLAY=:0 ssh-add ~/.ssh/id_ed25519
```

### Git Configuration

```bash
# Configure user identity
git config --global user.name "pakaiwa-release-bot"
git config --global user.email "release@pakaiwa.dev"

# Configure SSH signing
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub
git config --global tag.gpgsign true

# Use SSH instead of HTTPS
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

## Verification

### Check if Tag is Verified

After the workflow runs:

1. Go to your repository's tags page
2. Click on the tag
3. You should see a "Verified" badge next to the tag

### Verify Locally

```bash
# Fetch the tag
git fetch --tags

# Show tag signature
git tag -v v1.0.0

# Or show tag details
git show v1.0.0 --show-signature
```

## Troubleshooting

### "Permission denied (publickey)"

**Problem**: SSH authentication failed

**Solution**:
- Verify the SSH public key is added to GitHub as a **Signing Key**
- Check that `RELEASE_BOT_SSH_KEY_B64` secret is correctly set
- Ensure the private key matches the public key on GitHub

### "gpg: signing failed: Inappropriate ioctl for device"

**Problem**: SSH agent not properly configured

**Solution**:
- The workflow now uses `SSH_ASKPASS` to avoid this issue
- Ensure `SSH_AUTH_SOCK` is exported to `GITHUB_ENV`

### "failed to push some refs"

**Problem**: Git remote is using HTTPS instead of SSH

**Solution**:
- The workflow now explicitly sets the remote to SSH: `git@github.com:KAnggara75/scc2go.git`
- Verify the remote URL is correct in the workflow

### Tag Created but Not Verified

**Problem**: Tag was created but doesn't show as verified

**Solution**:
- Ensure the SSH key is added as a **Signing Key** (not Authentication Key)
- Verify `git config user.signingkey` points to the `.pub` file
- Check that `tag.gpgsign` is set to `true`

## Security Best Practices

1. **Use a Strong Passphrase**: Protect your SSH private key with a strong passphrase
2. **Rotate Keys Regularly**: Update your SSH keys periodically
3. **Limit Key Scope**: Use a dedicated key only for signing, not for authentication
4. **Monitor Usage**: Review GitHub's audit log for signing key usage
5. **Use a Bot Account**: Create a dedicated bot account for automated releases

## Alternative: GPG Signing

If you prefer GPG signing instead of SSH signing:

1. Generate a GPG key
2. Export the private key and add it as a secret
3. Update the workflow to use GPG instead of SSH
4. Configure git with `gpg.format gpg`

See [GitHub's GPG documentation](https://docs.github.com/en/authentication/managing-commit-signature-verification/generating-a-new-gpg-key) for details.

## References

- [GitHub SSH Signing Documentation](https://docs.github.com/en/authentication/managing-commit-signature-verification/about-commit-signature-verification)
- [Git Tag Signing](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work)
- [SSH Agent Forwarding](https://docs.github.com/en/developers/overview/using-ssh-agent-forwarding)
