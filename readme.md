# Snapper

Package snapper prints snapshots of values as they would appear in code. Inspired by snapshot testing from the frontend ecosystem.

Snapper supports the following types:

* All primitives
* Structs
* Pointers to structs
* Arrays/slices

It does NOT support:

* Pointers to primitives. Go doesn't have a straightforward way to represent these without relying on an external variable declaration or function call
* Private fields in structs (they will be skipped)
* Channels
* Probably more things that I'm forgetting right now

Note that pointers to the same struct will be represented as two separately instantiated
structs.

# Example

```go
snapper.Snap(thing) // Print to stdout
snapper.Fsnap(writer, thing) // Write to writer
str :=  snapper.Ssnap(thing) // Return string
```

# Why?

When writing a test case for a function that returns a big result (say, a slice of structs with
many fields) it's tedious to type out the entire test case. With snapper you can run your function,
print out a snapshot of the result, ensure it's correct, then copy and paste the snapshot directly
into your test file.
