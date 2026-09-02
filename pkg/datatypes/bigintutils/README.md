# bigints

This package provides helper functions to handle the [big.Int] datatype. 
While [big.Int] works with pointers the provided functions focus on storing the values as string since very common for input and output.

But there are still convenience functions offered operating with `*big.Int` pointers. The end with the suffix `Ints`.

## Specifications

For specifications see [bigintutils.spec.md](bigintutils.spec.md)
