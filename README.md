# GoHAProxyUtils
Utils for the haproxy

## Update HAProxy Allowed IPs

Allows to merge the current allowed and blocked IPs with the new ones. Provides the result that should be used.
There is a warning when the same IP is to be added both to contains and not_present. If only the "new" configs are provided, it can be used to check for this warning in an existing config.

## Find Duplicate IPs

Finds duplicate entries in the list, showing which IPs are duplicate and the unique list to be used.

## FAQ

**How do I provide the list? Which format?**

The IPs should be 1 per line. The spaces and "-" are removed so there is no need to clean this up.

**Is the order guaranteed?**

No. Since the order doesn't matter to the haproxy, this is not a problem.

**How can I access the service?**

Access [http://localhost:8080/](http://localhost:8080/) in a browser
