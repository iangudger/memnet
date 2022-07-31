# Portable in-memory Go "net" package implementations [![GoDoc](https://godoc.org/github.com/iangudger/memnet?status.png)](https://godoc.org/github.com/iangudger/memnet)

Allows portable and hermetic testing of code which uses the go `net` package. Useful for all cases where a Unix socket would (or could) be used, but you don't actually need to cross process boundaries. Fully portable and should work even on operating systems without native Unix socket support.

## Unbuffered
Simple and light-weight implementation built on `net.Pipe`. Prefer this implementation if it works for you. Supports `Listen`, `Dial` and pairs. Without buffering, there is no difference between stream and message oriented sockets, so there is only one type of unbuffered connection. Unlike `net.Pipe`, these implement both `net.Conn` and `net.PacketConn`.

## Buffered
Buffered types emulate Linux's version of Unix sockets and are based on [gVisor](https://gvisor.dev)'s [Unix socket implementation](https://cs.opensource.google/gvisor/gvisor/+/master:pkg/sentry/socket/unix/). These are heavier weight than the unbuffered implementations, but should prevent deadlocks in applications which depend on socket buffering. These should work for all applications which do not depend on IP addresses or ports.

The modifications to gVisor's Unix socket implementation have been intentionally kept minimal, will all the complexity of adapting it to `net` interfaces in the top-level package.

## Need more?
[`memipnet`](https://github.com/iangudger/memipnet) provides full TCP and UDP loopback emulation.
