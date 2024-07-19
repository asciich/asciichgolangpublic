package asciichgolangpublic

type GitlabService struct{}

func Gitlab() (g *GitlabService) {
	return new(GitlabService)
}

func NewGitlabService() (g *GitlabService) {
	return new(GitlabService)
}

func (g *GitlabService) GetDefaultGitlabCiYamlFileName() (fileName string) {
	return ".gitlab-ci.yml"
}
