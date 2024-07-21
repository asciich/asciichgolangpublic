package asciichgolangpublic

import (
	"math"
	"sort"
	"strings"
)

type SlicesService struct {
}

func NewSlicesService() (s *SlicesService) {
	return new(SlicesService)
}

func Slices() (slices *SlicesService) {
	return new(SlicesService)
}

func (o *SlicesService) RemoveDuplicatedStrings(sliceOfStrings []string) (cleaned []string) {
	if sliceOfStrings == nil {
		return []string{}
	}

	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	for _, entry := range sliceOfStrings {
		if o.ContainsString(cleaned, entry) {
			continue
		}

		cleaned = append(cleaned, entry)
	}

	return cleaned
}

func (o *SlicesService) RemoveLastElementIfEmptyString(sliceOfStrings []string) (cleanedUp []string) {
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

func (o *SlicesService) TrimAllPrefix(sliceOfStrings []string, prefixToRemove string) (sliceOfStringsWithPrefixRemoved []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	if len(prefixToRemove) <= 0 {
		return sliceOfStrings
	}

	sliceOfStringsWithPrefixRemoved = []string{}
	for _, sliceToCheck := range sliceOfStrings {
		sliceToCheck = Strings().TrimAllPrefix(sliceToCheck, prefixToRemove)

		sliceOfStringsWithPrefixRemoved = append(sliceOfStringsWithPrefixRemoved, sliceToCheck)
	}

	return sliceOfStringsWithPrefixRemoved
}

func (s *SlicesService) AddPrefixToEachString(stringSlices []string, prefix string) (output []string) {
	if len(stringSlices) <= 0 {
		return []string{}
	}

	output = []string{}
	for _, part := range stringSlices {
		output = append(output, prefix+part)
	}

	return output
}

func (s *SlicesService) AddSuffixToEachString(stringSlices []string, suffix string) (output []string) {
	if len(stringSlices) <= 0 {
		return []string{}
	}

	output = []string{}
	for _, part := range stringSlices {
		output = append(output, part+suffix)
	}

	return output
}

func (s *SlicesService) ContainsInt(intSlice []int, intToSearch int) (containsInt bool) {
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

func (s *SlicesService) ContainsString(sliceOfStrings []string, toCheck string) (contains bool) {
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

func (s *SlicesService) GetDeepCopyOfStringsSlice(sliceOfStrings []string) (deepCopy []string) {
	if len(sliceOfStrings) <= 0 {
		return []string{}
	}

	deepCopy = []string{}
	deepCopy = append(deepCopy, sliceOfStrings...)

	return deepCopy
}

func (s *SlicesService) GetIntSliceInitialized(nValues int, initValue int) (initializedSlice []int) {
	initializedSlice = []int{}
	if nValues <= 0 {
		return initializedSlice
	}

	for i := 0; i < nValues; i++ {
		initializedSlice = append(initializedSlice, initValue)
	}

	return initializedSlice
}

func (s *SlicesService) GetIntSliceInitializedWithZeros(nValues int) (initializedSlice []int) {
	return s.GetIntSliceInitialized(nValues, 0)
}

func (s *SlicesService) GetStringElementsNotInOtherSlice(toCheck []string, other []string) (elementsNotInOther []string) {
	if len(toCheck) <= 0 {
		return []string{}
	}

	elementsNotInOther = []string{}
	for _, elementToCheck := range toCheck {
		if !s.ContainsString(other, elementToCheck) {
			elementsNotInOther = append(elementsNotInOther, elementToCheck)
		}
	}

	return elementsNotInOther
}

func (s *SlicesService) MaxIntValuePerIndex(intSlice1 []int, intSlice2 []int) (maxValues []int) {
	maxLen := Math().MaxInt(len(intSlice1), len(intSlice2))

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

		valueToAdd := Math().MaxInt(slice1Value, slice2Value)
		maxValues = append(maxValues, valueToAdd)
	}

	return maxValues
}

func (s *SlicesService) MustRemoveStringsWhichContains(sliceToRemoveStringsWhichContains []string, searchString string) (cleanedUpSlice []string) {
	cleanedUpSlice, err := s.RemoveStringsWhichContains(sliceToRemoveStringsWhichContains, searchString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return cleanedUpSlice
}

func (s *SlicesService) RemoveEmptyStrings(sliceOfStrings []string) (sliceOfStringsWithoutEmptyStrings []string) {
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

func (s *SlicesService) RemoveMatchingStrings(sliceToRemoveMatching []string, matchingStringToRemove string) (cleanedUpSlice []string) {
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

func (s *SlicesService) RemoveStringEntryAtIndex(elements []string, indexToRemove int) (elementsWithIndexRemoved []string) {
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

func (s *SlicesService) RemoveStringsWhichContains(sliceToRemoveStringsWhichContains []string, searchString string) (cleanedUpSlice []string, err error) {
	if len(searchString) <= 0 {
		return nil, TracedError("searchString is empty string")
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

func (s *SlicesService) SortStringSlice(sliceOfStrings []string) (sorted []string) {
	sorted = s.GetDeepCopyOfStringsSlice(sliceOfStrings)

	sort.Strings(sorted)

	return sorted
}

func (s *SlicesService) SortStringSliceAndRemoveEmpty(input []string) (sortedAndWithoutEmptyStrings []string) {
	if len(input) <= 0 {
		return []string{}
	}

	sortedAndWithoutEmptyStrings = s.RemoveEmptyStrings(input)
	sortedAndWithoutEmptyStrings = s.SortStringSlice(sortedAndWithoutEmptyStrings)

	return sortedAndWithoutEmptyStrings
}

func (s *SlicesService) SortVersionStringSlice(input []string) (sorted []string) {
	return s.SortStringSlice(input)
}

func (s *SlicesService) SplitStrings(input []string, splitAt string) (splitted []string) {
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

func (s *SlicesService) SplitStringsAndRemoveEmpty(input []string, splitAt string) (splitted []string) {
	splitted = s.SplitStrings(input, splitAt)
	splitted = s.RemoveEmptyStrings(splitted)
	return splitted
}

func (s *SlicesService) ToLower(input []string) (lower []string) {
	lower = []string{}

	for _, i := range input {
		lower = append(lower, strings.ToLower(i))
	}

	return lower
}

func (s *SlicesService) TrimPrefix(sliceOfStrings []string, prefixToRemove string) (sliceOfStringsWithPrefixRemoved []string) {
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

func (s *SlicesService) TrimSpace(toTrim []string) (trimmed []string) {
	if len(toTrim) <= 0 {
		return []string{}
	}

	trimmed = []string{}
	for _, t := range toTrim {
		trimmed = append(trimmed, strings.TrimSpace(t))
	}

	return trimmed
}
