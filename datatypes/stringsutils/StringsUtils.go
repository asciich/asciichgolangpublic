package stringsutils

import (
	"encoding/hex"
	"fmt"
	"log"
	"regexp"

	"strconv"
	"strings"
	"unicode"
)

type StringsService struct{}

func NewStringsService() (s *StringsService) {
	return new(StringsService)
}

func Strings() (stringsService *StringsService) {
	return new(StringsService)
}

func ContainsAtLeastOneSubstring(input string, substrings []string) (atLeastOneSubstringFound bool) {
	for _, substring := range substrings {
		if strings.Contains(input, substring) {
			return true
		}
	}

	return false
}

func ContainsAtLeastOneSubstringIgnoreCase(input string, substring []string) (atLeastOneSubstringFound bool) {
	lowerInput := strings.ToLower(input)

	for _, s := range substring {
		if strings.Contains(lowerInput, strings.ToLower(s)) {
			return true
		}
	}

	return false
}

func ContainsCommentOnly(input string) (containsCommentOnly bool) {
	if strings.TrimSpace(input) == "" {
		return false
	}

	withoutComment := RemoveCommentsAndTrimSpace(input)
	return withoutComment == ""
}

func ContainsIgnoreCase(input string, substring string) (contains bool) {
	return strings.Contains(
		strings.ToLower(input),
		strings.ToLower(substring),
	)
}

func ContainsLine(input string, line string) (containsLine bool) {
	if input == "" {
		return false
	}

	for _, l := range SplitLines(input, false) {
		if l == line {
			return true
		}
	}

	return false
}

func CountLines(input string) (nLines int) {
	if len(input) <= 0 {
		return 0
	}

	nLines = strings.Count(input, "\n") + 1
	return nLines
}

func EnsureEndsWithExactlyOneLineBreak(input string) (ensuredLineBreak string) {
	if len(input) <= 0 {
		return "\n"
	}

	ensuredLineBreak = TrimSuffixUntilAbsent(input, "\n")
	ensuredLineBreak = EnsureEndsWithLineBreak(ensuredLineBreak)

	return ensuredLineBreak
}

func EnsureEndsWithLineBreak(input string) (ensuredLineBreak string) {
	ensuredLineBreak = EnsureSuffix(input, "\n")
	return ensuredLineBreak
}

func EnsureFirstCharLowercase(input string) (firstCharUppercase string) {
	if input == "" {
		return ""
	}

	if len(input) == 1 {
		return strings.ToLower(input)
	}

	firstChar := string(input[0])
	firstCharUppercase = strings.ToLower(firstChar) + input[1:]

	return firstCharUppercase
}

func EnsureFirstCharUppercase(input string) (firstCharUppercase string) {
	if input == "" {
		return ""
	}

	if len(input) == 1 {
		return strings.ToUpper(input)
	}

	firstChar := string(input[0])
	firstCharUppercase = strings.ToUpper(firstChar) + input[1:]

	return firstCharUppercase
}

func EnsurePrefix(input string, prefix string) (ensuredPrefix string) {
	if strings.HasPrefix(input, prefix) {
		return input
	}
	return prefix + input
}

func EnsureSuffix(input string, suffix string) (ensuredSuffix string) {
	if strings.HasSuffix(input, suffix) {
		return input
	}
	return input + suffix
}

func FirstCharToUpper(input string) (output string) {
	if len(input) <= 0 {
		return ""
	}

	if len(input) == 1 {
		return strings.ToUpper(input)
	}

	firstChar := string(input[0])
	suffix := input[1:]

	output = strings.ToUpper(firstChar) + suffix

	return output
}

func GetAsKeyValues(input string) (output map[string]string, err error) {
	if len(strings.TrimSpace(input)) <= 0 {
		return map[string]string{}, nil
	}

	output = map[string]string{}

	delimiter := ""

	for _, line := range SplitLines(input, true) {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if delimiter == "" {
			if strings.Contains(line, "=") {
				delimiter = "="
			} else if strings.Contains(line, ":") {
				delimiter = ":"
			} else {
				return nil, fmt.Errorf(
					"unable to find delimiter for getting key values in line: '%s'",
					line,
				)
			}
		}

		splitted := strings.SplitN(line, delimiter, 2)

		if len(splitted) != 2 {
			return nil, fmt.Errorf(
				"unable to split line '%s' into key values",
				line,
			)
		}

		key := strings.TrimSpace(splitted[0])
		if key == "" {
			return nil, fmt.Errorf(
				"key is empty string after evaluation of line '%s' for key values",
				line,
			)
		}

		value := strings.TrimSpace(splitted[1])

		output[key] = value
	}

	return output, nil
}

func GetFirstLine(input string) (firstLine string) {
	if len(input) <= 0 {
		return ""
	}

	lines := SplitLines(input, false)
	if len(lines) <= 0 {
		return ""
	}

	return lines[0]
}

func GetFirstLineAndTrimSpace(input string) (firstLine string) {
	if len(input) <= 0 {
		return ""
	}

	firstLine = GetFirstLine(input)
	firstLine = strings.TrimSpace(firstLine)

	return firstLine
}

func GetFirstLineWithoutCommentAndTrimSpace(input string) (firstLine string) {
	withoutComment := RemoveComments(input)
	firstLine = GetFirstLineAndTrimSpace(withoutComment)

	return firstLine
}

func GetNumberOfLinesWithPrefix(content string, prefix string, trimLines bool) (numberOfLinesWithPrefix int) {
	if content == "" {
		return 0
	}

	numberOfLinesWithPrefix = 0

	for _, line := range SplitLines(content, false) {
		lineToUse := line

		if trimLines {
			lineToUse = strings.TrimSpace(lineToUse)
		}

		if strings.HasPrefix(lineToUse, prefix) {
			numberOfLinesWithPrefix += 1
		}
	}

	return numberOfLinesWithPrefix
}

func GetValueAsInt(input string, key string) (value int, err error) {
	if key == "" {
		return -1, fmt.Errorf("key is empty string")
	}

	valueString, err := GetValueAsString(input, key)
	if err != nil {
		return -1, err
	}

	value, err = strconv.Atoi(valueString)
	if err != nil {
		return -1, fmt.Errorf(
			"unalbe to parse '%s' as string",
			valueString,
		)
	}

	return value, nil
}

func GetValueAsString(input string, key string) (value string, err error) {
	if key == "" {
		return "", fmt.Errorf("key is empty string")
	}

	keyValues, err := GetAsKeyValues(input)
	if err != nil {
		return "", err
	}

	value, ok := keyValues[key]
	if !ok {
		return "", fmt.Errorf(
			"key not found: %s",
			key,
		)
	}

	return value, nil
}

func HasAtLeastOnePrefix(toCheck string, prefixes []string) (hasPrefix bool) {
	if toCheck == "" {
		return false
	}

	for _, prefix := range prefixes {
		if strings.HasPrefix(toCheck, prefix) {
			return true
		}
	}

	return false
}

func HasPrefixIgnoreCase(input string, prefix string) (hasPrefix bool) {
	if prefix == "" {
		return true
	}

	inputLower := strings.ToLower(input)
	prefixLower := strings.ToLower(prefix)

	return strings.HasPrefix(inputLower, prefixLower)
}

func HexStringToBytes(hexString string) (output []byte, err error) {
	if hexString == "" {
		return []byte{}, nil
	}

	hexStringToParse := strings.TrimPrefix(hexString, "0x")
	hexStringToParse = strings.TrimPrefix(hexStringToParse, "0X")

	if len(hexString) == 1 {
		hexStringToParse = "0" + hexStringToParse
	}

	output, err = hex.DecodeString(hexStringToParse)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to convert hexString to bytes: %w",
			err,
		)
	}

	return output, nil
}

func IsComment(input string) (isComment bool) {
	if input == "" {
		return false
	}

	for _, line := range SplitLines(input, true) {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		if strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		return false
	}

	return true
}

func IsFirstCharLowerCase(input string) (isFirstCharLowerCase bool) {
	if input == "" {
		return false
	}

	firstChar := rune(input[0])

	if !unicode.IsLetter(firstChar) {
		return false
	}

	return unicode.IsLower(firstChar)
}

func IsFirstCharUpperCase(input string) (isFirstCharUpperCase bool) {
	if input == "" {
		return false
	}

	firstChar := rune(input[0])

	if !unicode.IsLetter(firstChar) {
		return false
	}

	return unicode.IsUpper(firstChar)
}

func MatchesRegex(input string, regex string) (matches bool, err error) {
	matches, err = regexp.Match(regex, []byte(input))
	if err != nil {
		return false, fmt.Errorf("match regex failed: '%w'", err)
	}

	return matches, nil
}

func MustGetAsKeyValues(input string) (output map[string]string) {
	output, err := GetAsKeyValues(input)
	if err != nil {
		log.Panic(err)
	}

	return output
}

func MustGetValueAsInt(input string, key string) (value int) {
	value, err := GetValueAsInt(input, key)
	if err != nil {
		log.Panic(err)
	}

	return value
}

func MustGetValueAsString(input string, key string) (value string) {
	value, err := GetValueAsString(input, key)
	if err != nil {
		log.Panic(err)
	}

	return value
}

func MustHexStringToBytes(hexString string) (output []byte) {
	output, err := HexStringToBytes(hexString)
	if err != nil {
		log.Panic(err)
	}

	return output
}

func MustMatchesRegex(input string, regex string) (matches bool) {
	matches, err := MatchesRegex(input, regex)
	if err != nil {
		log.Panic(err)
	}

	return matches
}

func RemoveCommentMarkers(input string) (commentContent string) {
	if input == "" {
		return ""
	}

	commentContent = ""
	for i, line := range SplitLines(input, false) {
		if i > 0 {
			commentContent += "\n"
		}

		if IsComment(line) {
			trimmedLine := strings.TrimSpace(line)
			commentContentLine := strings.TrimPrefix(trimmedLine, "#")
			commentContentLine = strings.TrimPrefix(commentContentLine, "//")
			commentContentLine = strings.TrimSpace(commentContentLine)
			commentContent += commentContentLine
			continue
		}

		commentContent += line
	}

	return commentContent
}

func RemoveCommentMarkersAndTrimSpace(input string) (commentContent string) {
	commentContent = RemoveCommentMarkers(input)
	commentContent = TrimSpaceForEveryLine(commentContent)
	return commentContent
}

func RemoveComments(input string) (contentWithoutComments string) {
	if len(input) <= 0 {
		return ""
	}

	contentWithoutComments = ""
	for _, line := range SplitLines(input, false) {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		if len(contentWithoutComments) > 0 {
			contentWithoutComments += "\n"
		}

		contentWithoutComments += line
	}

	return contentWithoutComments
}

func RemoveCommentsAndTrimSpace(input string) (output string) {
	output = RemoveComments(input)
	output = strings.TrimSpace(output)
	return output
}

func RemoveLinesWithPrefix(input string, prefixToRemove string) (output string) {
	lines := SplitLines(input, false)

	outputLines := []string{}
	for _, l := range lines {
		if strings.HasPrefix(l, prefixToRemove) {
			continue
		}

		outputLines = append(outputLines, l)
	}

	output = strings.Join(outputLines, "\n")

	return output
}

func RemoveSurroundingQuotationMarks(input string) (output string) {
	if len(input) <= 0 {
		return ""
	}

	output = input
	if strings.HasPrefix(output, "\"") {
		if strings.HasSuffix(output, "\"") {
			output = strings.TrimPrefix(output, "\"")
			output = strings.TrimSuffix(output, "\"")
			return output
		}
	}

	if strings.HasPrefix(output, "'") {
		if strings.HasSuffix(output, "'") {
			output = strings.TrimPrefix(output, "'")
			output = strings.TrimSuffix(output, "'")
			return output
		}
	}

	return output
}

func RemoveTailingNewline(input string) (cleaned string) {
	if len(input) == 0 {
		return ""
	}

	return strings.TrimSuffix(input, "\n")
}

func RepeatReplaceAll(input string, search string, replaceWith string) (replaced string) {
	if len(input) <= 0 {
		return input
	}

	if len(search) <= 0 {
		return input
	}

	replaced = input
	for strings.Contains(replaced, search) {
		replaced = strings.ReplaceAll(replaced, search, replaceWith)
	}

	return replaced
}

func RightFillWithSpaces(input string, fillLength int) (filled string) {
	if fillLength <= 0 {
		return input
	}

	if len(input) >= fillLength {
		return input
	}

	charsToAdd := fillLength - len(input)
	filled = input + strings.Repeat(" ", charsToAdd)
	return filled
}

func SplitAndGetLastElement(input string, token string) (lastElement string) {
	splitted := strings.Split(input, token)

	if len(splitted) <= 0 {
		return ""
	}

	return splitted[len(splitted)-1]
}

func SplitAtSpacesAndRemoveEmptyStrings(input string) (splitted []string) {
	splitted = []string{}
	for _, s := range strings.Split(input, " ") {
		if s == "" {
			continue
		}

		splitted = append(splitted, s)
	}

	return splitted
}

func SplitFirstLineAndContent(input string) (firstLine string, contentWithoutFirstLine string) {
	if len(input) <= 0 {
		return "", ""
	}

	splitted := strings.SplitN(input, "\n", 2)
	if len(splitted) == 0 {
		return "", ""
	} else if len(splitted) == 1 {
		return splitted[0], ""
	} else {
		return splitted[0], splitted[1]
	}
}

func SplitLines(input string, removeLastLineIfEmpty bool) (splittedLines []string) {
	if len(input) <= 0 {
		return []string{}
	}

	toSplit := strings.ReplaceAll(input, "\r\n", "\n")
	splittedLines = strings.Split(toSplit, "\n")

	if removeLastLineIfEmpty {
		if len(splittedLines) > 1 {
			if splittedLines[len(splittedLines)-1] == "" {
				splittedLines = splittedLines[:len(splittedLines)-1]
			}
		}
	}

	return splittedLines
}

func SplitWords(input string) (words []string) {
	words = []string{input}
	for _, splitChar := range []string{",", ".", "{", "}", "(", ")", "[", "]", "\t", "\n", " "} {
		nextWords := []string{}
		for _, w := range words {
			splitted := strings.Split(w, splitChar)
			for _, s := range splitted {
				if s != "" {
					nextWords = append(nextWords, s)
				}
			}
		}

		words = nextWords
	}

	return words
}

func ToPascalCase(input string) (pascalCase string) {
	splitted := strings.Split(input, " ")
	toJoin := []string{}

	for _, part := range splitted {
		toJoin = append(toJoin, FirstCharToUpper(part))
	}

	pascalCase = strings.Join(toJoin, "")

	return pascalCase
}

func ToSnakeCase(input string) (snakeCase string) {
	splitted := strings.Split(input, " ")
	snakeCase = strings.Join(splitted, "_")
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}

func TrimAllLeadingAndTailingNewLines(input string) (output string) {
	output = TrimAllLeadingNewLines(input)
	output = TrimAllTailingNewLines(output)
	return output
}

func TrimAllLeadingNewLines(input string) (output string) {
	return TrimAllPrefix(input, "\n")
}

func TrimAllPrefix(stringToCheck string, prefixToRemove string) (trimmedString string) {
	if len(stringToCheck) <= 0 {
		return ""
	}

	if len(prefixToRemove) <= 0 {
		return stringToCheck
	}

	trimmedString = stringToCheck
	for strings.HasPrefix(trimmedString, prefixToRemove) {
		trimmedString = strings.TrimPrefix(trimmedString, prefixToRemove)
	}

	return trimmedString
}

func TrimAllSuffix(stringToCheck string, suffixToRemove string) (trimmedString string) {
	if len(stringToCheck) <= 0 {
		return ""
	}

	if len(suffixToRemove) <= 0 {
		return stringToCheck
	}

	trimmedString = stringToCheck
	for strings.HasSuffix(trimmedString, suffixToRemove) {
		trimmedString = strings.TrimSuffix(trimmedString, suffixToRemove)
	}

	return trimmedString
}

func TrimAllTailingNewLines(input string) (output string) {
	return TrimAllSuffix(input, "\n")
}

func TrimPrefixAndSuffix(input string, prefix string, suffix string) (output string) {
	output = strings.TrimPrefix(input, prefix)
	output = strings.TrimSuffix(output, suffix)
	return output
}

func TrimPrefixIgnoreCase(input string, prefix string) (trimmed string) {
	hasPrefix := HasPrefixIgnoreCase(input, prefix)

	if !hasPrefix {
		return input
	}

	trimmed = input[len(prefix):]
	return trimmed
}

func TrimSpaceForEveryLine(input string) (trimmedForEveryLine string) {
	for i, l := range SplitLines(input, false) {
		if i > 0 {
			trimmedForEveryLine += "\n"
		}
		trimmedForEveryLine += l
	}

	return trimmedForEveryLine
}

func TrimSpacesLeft(input string) (trimmedLeft string) {
	return strings.TrimLeft(input, "\t \n")
}

func TrimSpacesRight(input string) (trimmedLeft string) {
	return strings.TrimRight(input, "\t \n")
}

func TrimSuffixAndSpace(input string, suffix string) (output string) {
	output = strings.TrimSuffix(input, suffix)
	output = strings.TrimSpace(output)

	return output
}

func TrimSuffixUntilAbsent(input string, suffixToRemove string) (withoutSuffix string) {
	withoutSuffix = input
	for strings.HasSuffix(withoutSuffix, suffixToRemove) {
		withoutSuffix = strings.TrimSuffix(withoutSuffix, suffixToRemove)
	}
	return withoutSuffix
}

// Convert "input" to a hex string while every char is separated by the "delimiter".
func ToHexString(input string, delimiter string) (hexString string) {
	if delimiter == "" {
		return hex.EncodeToString([]byte(input))
	}

	for i, c := range input {
		toAdd := hex.EncodeToString([]byte(string(c)))
		if i == 0 {
			hexString += toAdd
		} else {
			hexString += delimiter + toAdd
		}
	}
	return hexString
}

// Convert "input" to a hex string slices while every element of the hexStringSlice is the value of a char in hex as string.
func ToHexStringSlice(input string) (hexStringSlice []string) {
	hexStringSlice = []string{}
	for _, c := range input {
		hexStringSlice = append(
			hexStringSlice,
			hex.EncodeToString([]byte(string(c))),
		)
	}
	return hexStringSlice
}

// Returns true if s1 comes before s2 in the alphabeth.
// Alphabethical order is:
// - empty string
// - numbers [0-9]
// - chars [a-z]
func IsBeforeInAlphabeth(s1 string, s2 string) (isBefore bool) {
	if s1 == "" && s2 == "" {
		return false
	}

	if s2 == "" {
		return false
	}

	if strings.Compare(s1, s2) >= 0 {
		return false
	}

	return true
}
