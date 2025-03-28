

## Idea

Go Sdk Nix with is at:
``` bash
which go
# then share/go
# /nix/store/x09080wj4ja1zi7kycszya4cm7f1hq46-go-1.23.6/share/go
```

## Nix

We use Nix to install and pin the build tools to a version.

## Dev tools

For the devtools because they are used globally.
You may need to install them with:
* a personal nix shell
* or an OS package manager 
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

## Test Pushgateway

```bash
kubee kubectl port-forward -n monitoring svc/pushgateway 9091
# should work
echo "test_metric 3.14" | curl --data-binary @- http://127.0.0.1:9091/metrics/job/test
```