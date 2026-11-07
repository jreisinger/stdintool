Stdintool makes it easy to build handy CLI tools in Go that concurrently process
lines coming from stdin. See the `tool` folder for examples.

To run a tool:

```
$ go run ./tool/lookuphost < /tmp/hosts
```

Stdintool is inspired by (ok, copied from :-) https://youtu.be/woCg2zaIVzQ.
