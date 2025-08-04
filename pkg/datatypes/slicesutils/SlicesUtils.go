package slicesutils

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

func RemoveDuplicatedStrings(sliceOfStrings []string) (cleaned []string) {
	if sliceOfStrings == nil {
		return []string{}
	}

	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	for _, entry := range sliceOfStrings {
		if ContainsString(cleaned, entry) {
			continue
		}

		cleaned = append(cleaned, entry)
	}

	return cleaned
}

func RemoveLastElementIfEmptyString(sliceOfStrings []string) (cleanedUp []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	if len(sliceOfStrings) == 1 {
		if sliceOfStrings[0] == "" {
			return []string{}
		}

		return sliceOfStrings
	}

	if sliceOfStrings[len(sliceOfStrings)-1] == "" {
		return sliceOfStrings[:len(sliceOfStrings)-1]
	}

	return sliceOfStrings
}

func StringSlicesEqual(input1 []string, input2 []string) (slicesEqual bool) {
	if input1 == nil {
		return false
	}

	if input2 == nil {
		return false
	}

	if len(input1) != len(input2) {
		return false
	}

	for i, toCeck := range input1 {
		if toCeck != input2[i] {
			return false
		}
	}

	return true
}

func TrimAllPrefix(sliceOfStrings []string, prefixToRemove string) (sliceOfStringsWithPrefixRemoved []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	if len(prefixToRemove) <= 0 {
		return sliceOfStrings
	}

	sliceOfStringsWithPrefixRemoved = []string{}
	for _, sliceToCheck := range sliceOfStrings {
		for {
			if strings.HasPrefix(sliceToCheck, prefixToRemove) {
				sliceToCheck = strings.TrimPrefix(sliceToCheck, prefixToRemove)
				continue
			}

			break
		}

		sliceOfStringsWithPrefixRemoved = append(sliceOfStringsWithPrefixRemoved, sliceToCheck)
	}

	return sliceOfStringsWithPrefixRemoved
}

func AddPrefixToEachString(stringSlices []string, prefix string) (output []string) {
	if len(stringSlices) <= 0 {
		return []string{}
	}

	output = []string{}
	for _, part := range stringSlices {
		output = append(output, prefix+part)
	}

	return output
}

func AddSuffixToEachString(stringSlices []string, suffix string) (output []string) {
	if len(stringSlices) <= 0 {
		return []string{}
	}

	output = []string{}
	for _, part := range stringSlices {
		output = append(output, part+suffix)
	}

	return output
}

func AtLeastOneElementStartsWith(elements []string, toCheck string) (atLeastOneElementStartsWith bool) {
	for _, e := range elements {
		if strings.HasPrefix(e, toCheck) {
			return true
		}
	}

	return false
}

func ByteSlicesEqual(input1 []byte, input2 []byte) (slicesEqual bool) {
	if input1 == nil {
		return false
	}

	if input2 == nil {
		return false
	}

	if len(input1) != len(input2) {
		return false
	}

	for i, toCeck := range input1 {
		if toCeck != input2[i] {
			return false
		}
	}

	return true
}

func ContainsAllStrings(input []string, toCheck []string) (containsAllStrings bool) {
	if len(input) <= 0 {
		return false
	}

	if len(toCheck) <= 0 {
		return false
	}

	for _, c := range toCheck {
		if !ContainsString(input, c) {
			return false
		}
	}

	return true
}

func ContainsEmptyString(input []string) (containsEmptyString bool) {
	for _, i := range input {
		if i == "" {
			return true
		}
	}

	return false
}

func ContainsInt(intSlice []int, intToSearch int) (containsInt bool) {
	if len(intSlice) <= 0 {
		return false
	}

	for _, i := range intSlice {
		if i == intToSearch {
			return true
		}
	}

	return false
}

func ContainsNoEmptyStrings(input []string) (containsNoEmptyString bool) {
	return !ContainsEmptyString(input)
}

func ContainsOnlyUniqeStrings(input []string) (containsOnlyUniqeStrings bool) {
	for _, i := range input {
		if CountStrings(input, i) > 1 {
			return false
		}
	}

	return true
}

/* TODO move to SSH
func ContainsSshPublicKeyWithSameKeyMaterial(sshKeys []*SSHPublicKey, keyToSearch *SSHPublicKey) (contains bool) {
	if len(sshKeys) <= 0 {
		return false
	}

	if keyToSearch == nil {
		return false
	}

	keyMaterialToSearch, err := keyToSearch.GetKeyMaterialAsString()
	if err != nil {
		return false
	}

	for _, toCheck := range sshKeys {
		keyMaterialToCheck, err := toCheck.GetKeyMaterialAsString()
		if err != nil {
			continue
		}

		if keyMaterialToCheck == keyMaterialToSearch {
			return true
		}
	}

	return false
}
*/

func ContainsString(sliceOfStrings []string, toCheck string) (contains bool) {
	if len(sliceOfStrings) <= 0 {
		return false
	}

	for _, stringToCheck := range sliceOfStrings {
		if stringToCheck == toCheck {
			return true
		}
	}

	return false
}

func ContainsStringIgnoreCase(sliceOfStrings []string, toCheck string) (contains bool) {
	if len(sliceOfStrings) <= 0 {
		return false
	}

	toCheckLower := strings.ToLower(toCheck)
	for _, stringToCheck := range sliceOfStrings {
		if strings.ToLower(stringToCheck) == toCheckLower {
			return true
		}
	}

	return false
}

func CountStrings(input []string, toSearch string) (count int) {
	count = 0
	for _, i := range input {
		if i == toSearch {
			count += 1
		}
	}

	return count
}

func DiffStringSlices(a []string, b []string) (aNotInB []string, bNotInA []string) {
	aNotInB = []string{}
	bNotInA = []string{}

	for _, toCheck := range a {
		if ContainsString(b, toCheck) {
			continue
		}

		aNotInB = append(aNotInB, toCheck)
	}

	for _, toCheck := range b {
		if ContainsString(a, toCheck) {
			continue
		}

		bNotInA = append(bNotInA, toCheck)
	}

	sort.Strings(aNotInB)
	sort.Strings(bNotInA)

	return aNotInB, bNotInA
}

func GetDeepCopyOfByteSlice(input []byte) (deepCopy []byte) {
	if input == nil {
		return nil
	}

	deepCopy = make([]byte, len(input))
	copy(deepCopy, input)

	return deepCopy
}

func GetDeepCopyOfStringsSlice(sliceOfStrings []string) (deepCopy []string) {
	if sliceOfStrings == nil {
		return nil
	}

	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	deepCopy = make([]string, len(sliceOfStrings))
	copy(deepCopy, sliceOfStrings)

	return deepCopy
}

func GetInitializedIntSlice(nValues int, initValue int) []int {
	if nValues < 0 {
		return make([]int, 0)
	}

	initializedSlice := make([]int, nValues)
	for i := range initializedSlice {
		initializedSlice[i] = initValue
	}
	return initializedSlice
}

func GetInitializedIntSliceWithZeros(nValues int) (initializedSlice []int) {
	if nValues < 0 {
		return make([]int, 0)
	}
	return make([]int, nValues)
}

func GetStringElementsNotInOtherSlice(toCheck []string, other []string) (elementsNotInOther []string) {
	if len(toCheck) <= 0 {
		return []string{}
	}

	elementsNotInOther = []string{}
	for _, elementToCheck := range toCheck {
		if !ContainsString(other, elementToCheck) {
			elementsNotInOther = append(elementsNotInOther, elementToCheck)
		}
	}

	return elementsNotInOther
}

func maxInt(x, y int) (res int) {
	if x > y {
		return x
	}

	return y
}

func MaxIntValuePerIndex(intSlice1 []int, intSlice2 []int) (maxValues []int) {
	maxLen := maxInt(len(intSlice1), len(intSlice2))

	maxValues = []int{}
	for i := 0; i < maxLen; i++ {
		slice1Value := math.MinInt
		slice2Value := math.MinInt

		if i < len(intSlice1) {
			slice1Value = intSlice1[i]
		}

		if i < len(intSlice2) {
			slice2Value = intSlice2[i]
		}

		valueToAdd := maxInt(slice1Value, slice2Value)
		maxValues = append(maxValues, valueToAdd)
	}

	return maxValues
}

func RemoveEmptyStrings(sliceOfStrings []string) (sliceOfStringsWithoutEmptyStrings []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	sliceOfStringsWithoutEmptyStrings = []string{}
	for _, stringToCheck := range sliceOfStrings {
		if len(stringToCheck) <= 0 {
			continue
		}

		sliceOfStringsWithoutEmptyStrings = append(sliceOfStringsWithoutEmptyStrings, stringToCheck)
	}

	return sliceOfStringsWithoutEmptyStrings
}

func RemoveEmptyStringsAtEnd(input []string) (withoutEmptyStringsAtEnd []string) {
	if len(input) <= 0 {
		return []string{}
	}

	lastIndex := len(input) - 1
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == "" {
			lastIndex = i
		} else {
			break
		}
	}

	withoutEmptyStringsAtEnd = input[:lastIndex]

	return withoutEmptyStringsAtEnd
}

func RemoveMatchingStrings(sliceToRemoveMatching []string, matchingStringToRemove string) (cleanedUpSlice []string) {
	if len(sliceToRemoveMatching) <= 0 {
		return []string{}
	}

	cleanedUpSlice = []string{}
	for _, s := range sliceToRemoveMatching {
		if s == matchingStringToRemove {
			continue
		}

		cleanedUpSlice = append(cleanedUpSlice, s)
	}

	return cleanedUpSlice
}

func RemoveString(elements []string, toRemove string) (cleanedUpElements []string) {
	cleanedUpElements = []string{}

	for _, e := range elements {
		if e == toRemove {
			continue
		}

		cleanedUpElements = append(cleanedUpElements, e)
	}

	return cleanedUpElements
}

func RemoveStringEntryAtIndex(elements []string, indexToRemove int) (elementsWithIndexRemoved []string) {
	if len(elements) <= 0 {
		return []string{}
	}

	elementsWithIndexRemoved = []string{}
	for i, element := range elements {
		if i == indexToRemove {
			continue
		}

		elementsWithIndexRemoved = append(elementsWithIndexRemoved, element)
	}

	return elementsWithIndexRemoved
}

func RemoveStringsWhichContains(sliceToRemoveStringsWhichContains []string, searchString string) (cleanedUpSlice []string, err error) {
	if len(searchString) <= 0 {
		return nil, fmt.Errorf("searchString is empty string")
	}

	if len(sliceToRemoveStringsWhichContains) <= 0 {
		return []string{}, nil
	}

	cleanedUpSlice = []string{}
	for _, s := range sliceToRemoveStringsWhichContains {
		if strings.Contains(s, searchString) {
			continue
		}

		cleanedUpSlice = append(cleanedUpSlice, s)
	}

	return cleanedUpSlice, nil
}

func GetSortedDeepCopyOfStringsSlice(sliceOfStrings []string) (sorted []string) {
	sorted = GetDeepCopyOfStringsSlice(sliceOfStrings)

	sort.Strings(sorted)

	return sorted
}

func SortStringSliceAndRemoveDuplicates(input []string) (output []string) {
	if len(input) <= 0 {
		return []string{}
	}

	sorted := GetSortedDeepCopyOfStringsSlice(input)
	return RemoveDuplicatedStrings(sorted)
}

func SortStringSliceAndRemoveEmpty(input []string) (sortedAndWithoutEmptyStrings []string) {
	if len(input) <= 0 {
		return []string{}
	}

	sortedAndWithoutEmptyStrings = RemoveEmptyStrings(input)
	sort.Strings(sortedAndWithoutEmptyStrings)

	return sortedAndWithoutEmptyStrings
}

func SplitStrings(input []string, splitAt string) (splitted []string) {
	if len(input) <= 0 {
		return []string{}
	}

	splitted = []string{}
	for _, toSplit := range input {
		toAdd := strings.Split(toSplit, splitAt)
		splitted = append(splitted, toAdd...)
	}

	return splitted
}

func SplitStringsAndRemoveEmpty(input []string, splitAt string) (splitted []string) {
	splitted = SplitStrings(input, splitAt)
	splitted = RemoveEmptyStrings(splitted)
	return splitted
}

func ToLower(input []string) (lower []string) {
	lower = []string{}

	for _, i := range input {
		lower = append(lower, strings.ToLower(i))
	}

	return lower
}

func TrimPrefix(sliceOfStrings []string, prefixToRemove string) (sliceOfStringsWithPrefixRemoved []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	if len(prefixToRemove) <= 0 {
		return sliceOfStrings
	}

	sliceOfStringsWithPrefixRemoved = []string{}
	for _, stringToCheck := range sliceOfStrings {
		stringToCheck = strings.TrimPrefix(stringToCheck, prefixToRemove)

		sliceOfStringsWithPrefixRemoved = append(sliceOfStringsWithPrefixRemoved, stringToCheck)
	}

	return sliceOfStringsWithPrefixRemoved
}

func TrimSpace(toTrim []string) (trimmed []string) {
	if len(toTrim) <= 0 {
		return []string{}
	}

	trimmed = []string{}
	for _, t := range toTrim {
		trimmed = append(trimmed, strings.TrimSpace(t))
	}

	return trimmed
}
