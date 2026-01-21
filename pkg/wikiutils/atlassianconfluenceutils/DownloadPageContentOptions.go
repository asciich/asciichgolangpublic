package atlassianconfluenceutils

type DownloadPageContentOptions struct {
	// If set to true the child pages are downloaded as well.
	Recursive bool

	// If set to true the downloaded pages are converted and stored as MarkDown files instead of HTML.
	ConvertToMdFiles bool
}
