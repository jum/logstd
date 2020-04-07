# logstd

A simple command to redirect stdout and stderr of a command into the os_log
standard MacOS X logging system.

```sh
logstd [flags] -- shell command
```

To watch the log of the logstd command in a terminal:

```sh
log stream --level debug --predicate 'subsystem == "com.github.jum.logstd"'
```

The following flags are supported:

 * -subsystem string
 
    os_log subsystem (default "com.github.jum.logstd")

* -stderr string

    os_log category for stderr (default "stderr")

* -stdout string

	os_log category for stdout (default "stdout")
    