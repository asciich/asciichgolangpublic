package vectordatabasegeneric

import "strings"

func SplitText(text string, chunkSize, overlap int) []string {
	separators := []string{"\n\n", "\n", ". ", " ", ""}
	return RecursiveSplit(text, separators, chunkSize, overlap)
}

func RecursiveSplit(text string, separators []string, chunkSize, overlap int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	sep := separators[0]
	remainingSeps := separators[1:]

	parts := strings.Split(text, sep)
	var chunks []string
	var current strings.Builder

	for _, part := range parts {
		// If adding this part would exceed chunkSize and we have content, flush current
		if current.Len()+len(part)+len(sep) > chunkSize && current.Len() > 0 {
			chunk := strings.TrimSpace(current.String())
			if len(chunk) > chunkSize && len(remainingSeps) > 0 {
				chunks = append(chunks, RecursiveSplit(chunk, remainingSeps, chunkSize, overlap)...)
			} else {
				chunks = append(chunks, chunk)
			}

			// Handle overlap
			overlapText := current.String()
			current.Reset()
			if overlap > 0 && len(overlapText) > overlap {
				current.WriteString(overlapText[len(overlapText)-overlap:])
			}
		}

		// If the part alone exceeds chunkSize, recursively split it
		if len(part) > chunkSize {
			// Flush current first if non-empty
			if current.Len() > 0 {
				chunk := strings.TrimSpace(current.String())
				if chunk != "" {
					if len(chunk) > chunkSize && len(remainingSeps) > 0 {
						chunks = append(chunks, RecursiveSplit(chunk, remainingSeps, chunkSize, overlap)...)
					} else {
						chunks = append(chunks, chunk)
					}
				}
				current.Reset()
			}

			// Recursively split the oversized part with remaining separators
			if len(remainingSeps) > 0 {
				chunks = append(chunks, RecursiveSplit(part, remainingSeps, chunkSize, overlap)...)
			} else {
				// Last resort: hard split by chunkSize
				for i := 0; i < len(part); i += chunkSize - overlap {
					end := i + chunkSize
					if end > len(part) {
						end = len(part)
					}
					chunk := strings.TrimSpace(part[i:end])
					if chunk != "" {
						chunks = append(chunks, chunk)
					}
					if end == len(part) {
						break
					}
				}
			}
			continue
		}

		if current.Len() > 0 {
			current.WriteString(sep)
		}
		current.WriteString(part)
	}

	// Flush remaining content
	if current.Len() > 0 {
		chunk := strings.TrimSpace(current.String())
		if chunk != "" {
			if len(chunk) > chunkSize && len(remainingSeps) > 0 {
				chunks = append(chunks, RecursiveSplit(chunk, remainingSeps, chunkSize, overlap)...)
			} else if len(chunk) > chunkSize {
				// Hard split
				for i := 0; i < len(chunk); i += chunkSize - overlap {
					end := i + chunkSize
					if end > len(chunk) {
						end = len(chunk)
					}
					c := strings.TrimSpace(chunk[i:end])
					if c != "" {
						chunks = append(chunks, c)
					}
					if end == len(chunk) {
						break
					}
				}
			} else {
				chunks = append(chunks, chunk)
			}
		}
	}

	return chunks
}
