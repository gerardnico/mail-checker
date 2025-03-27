# Mail Checker


Check the spf config of a mail server.


## Example

With the [configuration file](examples/conf.yml)
```bash
mail-checker --conf examples/conf.yml
```

## Conf

[configuration file](examples/conf.yml)

## Reports

### KuberHealthy

If a [KuberHealthy](https://kuberhealthy.github.io/kuberhealthy/) execution is detected, 
we report to it.

### PushGateway

If a [Pushgateway Url is configured](examples/conf.yml), we report to it.

The check metric is named `mail_checker_check` and report a value of:
* `1` if failed
* `0` if successful

The parameters are set in the prometheus label.

## Exit Status

If there is errors, the process will exit with a status of
* `0` if a report has been made, ie the error has already been reported.
* `1` if no reporter has been executed or `mail-checker` encounters an unexpected exception.


