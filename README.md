# Kubee Mail Checker


`Mail Checker` is a synthetic monitoring test utility, part of the [kubee platform](https://github.com/EraldyHq/kubee)

It checks the `spf` DNS configuration.

It does not support `yet` the following checks:
* Dkim
* Reverse Dns
* Dmarc

## Example

With the [configuration file](examples/conf.yml)
```bash
mail-checker --conf examples/conf.yml
```

## Conf

[configuration file](examples/conf.yml)

## Reports

`Mail-Checker` can reports the results of the check to the following platforms.

### PushGateway

If a [Pushgateway Url is configured](examples/conf.yml), we report to it with the following metrics.

The metric created is called `kubee_check` metrics and has the following labels:
* `type`: type of check. `spf`
* `domain`: the domain name checked
* `job`: `mail-checker`
* `mailer`: the mailer ip
* `value`:
    * `0` if successful
    * `1` if failed

### KuberHealthy

If a [KuberHealthy](https://kuberhealthy.github.io/kuberhealthy/) execution is detected,
we report to it.

## Exit Status

If there is errors, the process will exit with a status of
* `0` if a report has been made, ie the error has already been reported.
* `1` if no reporter has been executed or `mail-checker` encounters an unexpected exception.


