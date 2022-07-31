# Portable in-memory Go "net" package implementations [![GoDoc](https://godoc.org/github.com/iangudger/memnet?status.png)](https://godoc.org/github.com/iangudger/memnet)

Allows portable and hermetic testing and connecting. Useful for all cases where a Unix socket would be used, but you don't actually need to cross process boundaries.

Buffered types emulate Linux's version of Unix sockets and are based on [gVisor](https://gvisor.dev)'s [Unix socket implementation](https://cs.opensource.google/gvisor/gvisor/+/master:pkg/sentry/socket/unix/).

## Status
- [x] Unbuffered packet socketpair
- [x] Unbuffered listener
- [x] Buffered stream socketpair
- [x] Buffered packet socketpair
- [x] Buffered listener

The underlying implementations (`net.Pipe` and gVisor's Unix sockets) are pretty battle tested at this point. The modifications to gVisor's Unix socket implementation have been intentionally kept minimal, will all the complexity of adapting it to `net` interfaces in this package. This compatibility layer is new and may contain bugs. Please report any that you find.
