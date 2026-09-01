# slicesutils

Package `slicesutils` provides utility functions for working with slices.

## Functions

- `AddPrefixToEachString` - Adds a prefix to each string in a slice
- `AddSuffixToEachString` - Adds a suffix to each string in a slice
- `AtLeastOneElementStartsWith` - Checks if at least one element starts with a given prefix
- `ByteSlicesEqual` - Compares two byte slices for equality
- `ContainsAllStrings` - Checks if a slice contains all strings from another slice
- `ContainsEmptyString` - Checks if a slice contains an empty string
- `ContainsInt` - Checks if an int slice contains a specific int
- `ContainsNoEmptyStrings` - Checks if a slice contains no empty strings
- `ContainsOnlyUniqeStrings` - Checks if a slice contains only unique strings
- `ContainsSshPublicKeyWithSameKeyMaterial` - Checks for SSH public keys with same key material
- `ContainsString` - Checks if a string slice contains a specific string
- `ContainsStringIgnoreCase` - Checks if a string slice contains a string (case-insensitive)
- `CountStrings` - Counts occurrences of a string in a slice
- `DiffStringSlices` - Returns elements unique to each slice
- `GetDeepCopyOfByteSlice` - Creates a deep copy of a byte slice
- `GetDeepCopyOfStringsSlice` - Creates a deep copy of a string slice
- `GetInitializedIntSlice` - Creates an initialized int slice with a given value
- `GetInitializedIntSliceWithZeros` - Creates an int slice initialized with zeros
- `GetSortedDeepCopyOfStringsSlice` - Creates a sorted deep copy of a string slice
- `GetStringElementsNotInOtherSlice` - Returns elements not present in another slice
- `MaxIntValuePerIndex` - Returns max values per index from two slices
- `RemoveDuplicatedStrings` - Removes duplicate strings from a slice
- `RemoveEmptyStrings` - Removes empty strings from a slice
- `RemoveEmptyStringsAtEnd` - Removes empty strings at the end of a slice
- `RemoveLastElementIfEmptyString` - Removes the last element if it's an empty string
- `RemoveMatchingStrings` - Removes strings matching a regex pattern
- `RemoveString` - Removes a specific string from a slice
- `RemoveStringEntryAtIndex` - Removes the element at a specific index
- `RemoveStringsWhichContains` - Removes strings containing a substring
- `SortStringSliceAndRemoveDuplicates` - Sorts and removes duplicates
- `SortStringSliceAndRemoveEmpty` - Sorts and removes empty strings
- `SplitStrings` - Splits all strings by a delimiter
- `SplitStringsAndRemoveEmpty` - Splits strings and removes empty results
- `StringSlicesEqual` - Compares two string slices for equality
- `ToLower` - Converts all strings to lowercase
- `TrimAllPrefix` - Removes a prefix from all strings
- `TrimPrefix` - Removes a prefix from all strings
- `TrimSpace` - Trims space from all strings

## Specifications

For specifications see [slicesutils.spec.md](slicesutils.spec.md)
