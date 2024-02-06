package asciichgolangpublic

import (
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

func (s *StringsService) ContainsAtLeastOneSubstring(input string, substrings []string) (atLeastOneSubstringFound bool) {
	for _, substring := range substrings {
		if strings.Contains(input, substring) {
			return true
		}
	}

	return false
}

func (s *StringsService) ContainsAtLeastOneSubstringIgnoreCase(input string, substring []string) (atLeastOneSubstringFound bool) {
	return s.ContainsAtLeastOneSubstring(
		strings.ToLower(input),
		Slices().ToLower(substring),
	)
}

func (s *StringsService) ContainsCommentOnly(input string) (containsCommentOnly bool) {
	if strings.TrimSpace(input) == "" {
		return false
	}

	withoutComment := s.RemoveCommentsAndTrimSpace(input)
	return withoutComment == ""
}

func (s *StringsService) CountLines(input string) (nLines int) {
	if len(input) <= 0 {
		return 0
	}

	nLines = strings.Count(input, "\n") + 1
	return nLines
}

func (s *StringsService) EnsureEndsWithExactlyOneLineBreak(input string) (ensuredLineBreak string) {
	if len(input) <= 0 {
		return "\n"
	}

	ensuredLineBreak = Strings().TrimSuffixUntilAbsent(input, "\n")
	ensuredLineBreak = Strings().EnsureEndsWithLineBreak(ensuredLineBreak)

	return ensuredLineBreak
}

func (s *StringsService) EnsureEndsWithLineBreak(input string) (ensuredLineBreak string) {
	ensuredLineBreak = Strings().EnsureSuffix(input, "\n")
	return ensuredLineBreak
}

func (s *StringsService) EnsureFirstCharLowercase(input string) (firstCharUppercase string) {
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

func (s *StringsService) EnsureFirstCharUppercase(input string) (firstCharUppercase string) {
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

func (s *StringsService) EnsureSuffix(input string, suffix string) (ensuredSuffix string) {
	if strings.HasSuffix(input, suffix) {
		return input
	} else {
		return input + suffix
	}
}

func (s *StringsService) FirstCharToUpper(input string) (output string) {
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

func (s *StringsService) GetFirstLine(input string) (firstLine string) {
	if len(input) <= 0 {
		return ""
	}

	lines := s.SplitLines(input)
	if len(lines) <= 0 {
		return ""
	}

	return lines[0]
}

func (s *StringsService) GetFirstLineAndTrimSpace(input string) (firstLine string) {
	if len(input) <= 0 {
		return ""
	}

	firstLine = s.GetFirstLine(input)
	firstLine = strings.TrimSpace(firstLine)

	return firstLine
}

func (s *StringsService) GetFirstLineWithoutCommentAndTrimSpace(input string) (firstLine string) {
	withoutComment := s.RemoveComments(input)
	firstLine = s.GetFirstLineAndTrimSpace(withoutComment)

	return firstLine
}

func (s *StringsService) GetNumberOfLinesWithPrefix(content string, prefix string, trimLines bool) (numberOfLinesWithPrefix int) {
	if content == "" {
		return 0
	}

	numberOfLinesWithPrefix = 0

	for _, line := range s.SplitLines(content) {
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

func (s *StringsService) HasAtLeastOnePrefix(toCheck string, prefixes []string) (hasPrefix bool) {
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

func (s *StringsService) HasPrefixIgnoreCase(input string, prefix string) (hasPrefix bool) {
	if prefix == "" {
		return true
	}

	inputLower := strings.ToLower(input)
	prefixLower := strings.ToLower(prefix)

	return strings.HasPrefix(inputLower, prefixLower)
}

func (s *StringsService) IsComment(input string) (isComment bool) {
	if input == "" {
		return false
	}

	for _, line := range s.SplitLines(s.RemoveTailingNewline(input)) {
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

func (s *StringsService) IsFirstCharLowerCase(input string) (isFirstCharLowerCase bool) {
	if input == "" {
		return false
	}

	firstChar := rune(input[0])

	if !unicode.IsLetter(firstChar) {
		return false
	}

	return unicode.IsLower(firstChar)
}

func (s *StringsService) IsFirstCharUpperCase(input string) (isFirstCharUpperCase bool) {
	if input == "" {
		return false
	}

	firstChar := rune(input[0])

	if !unicode.IsLetter(firstChar) {
		return false
	}

	return unicode.IsUpper(firstChar)
}

func (s *StringsService) RemoveCommentMarkers(input string) (commentContent string) {
	if input == "" {
		return ""
	}

	commentContent = ""
	for i, line := range s.SplitLines(input) {
		if i > 0 {
			commentContent += "\n"
		}

		if s.IsComment(line) {
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

func (s *StringsService) RemoveCommentMarkersAndTrimSpace(input string) (commentContent string) {
	commentContent = s.RemoveCommentMarkers(input)
	commentContent = s.TrimSpaceForEveryLine(commentContent)
	return commentContent
}

func (s *StringsService) RemoveComments(input string) (contentWithoutComments string) {
	if len(input) <= 0 {
		return ""
	}

	contentWithoutComments = ""
	for _, line := range Strings().SplitLines(input) {
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

func (s *StringsService) RemoveCommentsAndTrimSpace(input string) (output string) {
	output = s.RemoveComments(input)
	output = strings.TrimSpace(output)
	return output
}

func (s *StringsService) RemoveSurroundingQuotationMarks(input string) (output string) {
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

	return output
}

func (s *StringsService) RemoveTailingNewline(input string) (cleaned string) {
	if len(input) == 0 {
		return ""
	}

	return strings.TrimSuffix(input, "\n")
}

func (s *StringsService) RepeatReplaceAll(input string, search string, replaceWith string) (replaced string) {
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

func (s *StringsService) RightFillWithSpaces(input string, fillLength int) (filled string) {
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

func (s *StringsService) SplitAtSpacesAndRemoveEmptyStrings(input string) (splitted []string) {
	splitted = strings.Split(input, " ")
	splitted = Slices().RemoveEmptyStrings(splitted)
	return splitted
}

func (s *StringsService) SplitFirstLineAndContent(input string) (firstLine string, contentWithoutFirstLine string) {
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

func (s *StringsService) SplitLines(input string) (splittedLines []string) {
	if len(input) <= 0 {
		return []string{}
	}

	splittedLines = strings.Split(input, "\n")
	return splittedLines
}

func (s *StringsService) SplitWords(input string) (words []string) {
	words = []string{input}
	for _, splitChar := range []string{",", ".", "{", "}", "(", ")", "[", "]", "\t", "\n", " "} {
		words = Slices().SplitStringsAndRemoveEmpty(words, splitChar)
	}

	return words
}

func (s *StringsService) ToPascalCase(input string) (pascalCase string) {
	splitted := strings.Split(input, " ")
	toJoin := []string{}

	for _, part := range splitted {
		toJoin = append(toJoin, s.FirstCharToUpper(part))
	}

	pascalCase = strings.Join(toJoin, "")

	return pascalCase
}

func (s *StringsService) ToSnakeCase(input string) (snakeCase string) {
	splitted := strings.Split(input, " ")
	snakeCase = strings.Join(splitted, "_")
	snakeCase = strings.ToLower(snakeCase)
	return snakeCase
}

func (s *StringsService) TrimAllPrefix(stringToCheck string, prefixToRemove string) (trimmedString string) {
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

func (s *StringsService) TrimAllSuffix(stringToCheck string, suffixToRemove string) (trimmedString string) {
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

func (s *StringsService) TrimPrefixIgnoreCase(input string, prefix string) (trimmed string) {
	hasPrefix := s.HasPrefixIgnoreCase(input, prefix)

	if !hasPrefix {
		return input
	}

	trimmed = input[len(prefix):]
	return trimmed
}

func (s *StringsService) TrimSpaceForEveryLine(input string) (trimmedForEveryLine string) {
	lines := s.SplitLines(input)
	toJoin := Slices().TrimSpace(lines)
	return strings.Join(toJoin, "\n")
}

func (s *StringsService) TrimSpacesLeft(input string) (trimmedLeft string) {
	return strings.TrimLeft(input, "\t \n")
}

func (s *StringsService) TrimSuffixAndSpace(input string, suffix string) (output string) {
	output = strings.TrimSuffix(input, suffix)
	output = strings.TrimSpace(output)

	return output
}

func (s *StringsService) TrimSuffixUntilAbsent(input string, suffixToRemove string) (withoutSuffix string) {
	withoutSuffix = input
	for strings.HasSuffix(withoutSuffix, suffixToRemove) {
		withoutSuffix = strings.TrimSuffix(withoutSuffix, suffixToRemove)
	}
	return withoutSuffix
}
