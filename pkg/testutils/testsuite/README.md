# testsuite package

Provides a simple way to declare test suites and run them.

## Available test types

| test_type       | Comment                                                    | Example                                                        |
| --------------- | ---------------------------------------------------------- | -------------------------------------------------------------- | 
| `command`       | Run an arbitrary command and pass if the exit code is `0`. | [Test Google reachable](./Example_TestGoolgeReachable_test.go) |
| `tcp_port_open` | Passes when the TCP `port` on `host` is open.              | [Test Google reachable](./Example_TestGoolgeReachable_test.go) |
