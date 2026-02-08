# Quick Setup Guide - SSH Signed Tags

## TL;DR - 3 Steps to Get Started

### 1. Generate SSH Key

```bash
ssh-keygen -t ed25519 -C "release@pakaiwa.dev" -f ~/.ssh/pakaiwa_release_bot
# Enter a strong passphrase when prompted
```

### 2. Add Secrets to GitHub

Go to your repository → Settings → Secrets and variables → Actions → New repository secret

**Secret 1: `RELEASE_BOT_SSH_KEY`**
```bash
cat ~/.ssh/pakaiwa_release_bot
```
Copy the entire output (including header/footer lines) and paste as the secret value.

**Secret 2: `RELEASE_BOT_SSH_PASSPHRASE`**
Enter the passphrase you used when creating the key.

**Secret 3: `CODECOV_TOKEN`** (optional, for coverage reporting)
Get from [codecov.io](https://codecov.io) after setting up your repository.

### 3. Add SSH Public Key to GitHub

```bash
cat ~/.ssh/pakaiwa_release_bot.pub
```

1. Go to GitHub Settings → SSH and GPG keys
2. Click "New SSH key"
3. Title: `pakaiwa-release-bot signing key`
4. Key type: **Signing Key** ⚠️ (Important!)
5. Paste the public key
6. Click "Add SSH key"

## That's It! 🎉

Now when you push commits to main/master:

```bash
git add .
git commit -m "feat: add new feature"
git push origin main
```

The workflow will:
- ✅ Run tests with coverage
- ✅ Calculate next version based on commit message
- ✅ Create a **verified** signed tag
- ✅ Create a GitHub release with changelog

## Commit Message Format

Use [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: description` → Minor version bump (v0.1.0 → v0.2.0)
- `fix: description` → Patch version bump (v0.1.0 → v0.1.1)
- `feat!: description` → Major version bump (v0.1.0 → v1.0.0)

## Verify It Works

After pushing, check:
1. **Actions tab** - Workflow should complete successfully
2. **Tags** - New tag should show "Verified" badge
3. **Releases** - New release should be created with changelog

## Troubleshooting

### Tag not verified?
- Make sure you added the key as **Signing Key** (not Authentication Key)
- Verify the secret includes the header/footer lines

### Workflow fails at SSH setup?
- Check that `RELEASE_BOT_SSH_KEY` is the full private key (plain text)
- Verify `RELEASE_BOT_SSH_PASSPHRASE` matches the key passphrase

### Permission denied?
- Ensure the public key on GitHub matches your private key
- Verify you have write access to the repository

## Full Documentation

For detailed information, see:
- [SSH_SIGNING_SETUP.md](SSH_SIGNING_SETUP.md) - Complete setup guide
- [VERSIONING.md](VERSIONING.md) - Semantic versioning details
- [SSH_SIGNING_CHANGES.md](SSH_SIGNING_CHANGES.md) - What was changed
