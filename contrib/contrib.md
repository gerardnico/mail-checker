

## Idea

Go Sdk Nix with is at:
```
/nix/store/x09080wj4ja1zi7kycszya4cm7f1hq46-go-1.23.6/share/go
```

## Nix

We use Nix to install and pin the build tools to a version.

## Dev tools

For the devtools because they are used globally.
You may need to install them with an OS package manager 
if you want to commit from another location for instance.

Example:
```bash
brew install node
npm install -g @commitlint/cli @commitlint/config-conventional
```

## Github

Repo set 
* in [](../.envrc) for `git cliff`
* in [jreleaser](../jreleaser.yml)
