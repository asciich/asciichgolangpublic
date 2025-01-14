package asciichgolangpublic

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type PreCommitService struct{}

func NewPreCommitService() (p *PreCommitService) {
	return new(PreCommitService)
}

func PreCommit() (p *PreCommitService) {
	return NewPreCommitService()
}

func (p *PreCommitService) GetAsPreCommitConfigFileOrNilIfContentIsInvalid(file File, verbose bool) (preCommitConfigFile *PreCommitConfigFile, err error) {
	preCommitConfigFile, err = GetPreCommitConfigByFile(file)
	if err != nil {
		return nil, err
	}

	isContentValid, err := preCommitConfigFile.IsValidPreCommitConfigFile(verbose)
	if err != nil {
		return nil, err
	}

	if isContentValid {
		return preCommitConfigFile, nil
	} else {
		return nil, nil
	}
}

func (p *PreCommitService) GetDefaultConfigFileName() (preCommitDefaultName string) {
	return ".pre-commit-config.yaml"
}

func (p *PreCommitService) MustGetAsPreCommitConfigFileOrNilIfContentIsInvalid(file File, verbose bool) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := p.GetAsPreCommitConfigFileOrNilIfContentIsInvalid(file, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func (p *PreCommitService) MustRunInDirectory(directoy Directory, options *PreCommitRunOptions) {
	err := p.RunInDirectory(directoy, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitService) MustRunInGitRepository(gitRepo GitRepository, options *PreCommitRunOptions) {
	err := p.RunInGitRepository(gitRepo, options)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (p *PreCommitService) RunInDirectory(directoy Directory, options *PreCommitRunOptions) (err error) {
	if directoy == nil {
		return errors.TracedErrorNil("directoy")
	}

	if options == nil {
		return errors.TracedErrorNil("options")
	}

	path, err := directoy.GetLocalPath()
	if err != nil {
		return err
	}

	preCommitCommand := []string{
		"bash",
		"-c",
		fmt.Sprintf(
			"cd '%s' && pre-commit run -a",
			path,
		),
	}

	_, err = Bash().RunCommand(
		&RunCommandOptions{
			Command:            preCommitCommand,
			Verbose:            options.Verbose,
			LiveOutputOnStdout: options.Verbose,
		},
	)
	if err != nil {
		return err
	}

	if options.Verbose {
		logging.LogInfof("Pre commit successfully run in '%s'.", path)
	}

	return nil
}

func (p *PreCommitService) RunInGitRepository(gitRepo GitRepository, options *PreCommitRunOptions) (err error) {
	if gitRepo == nil {
		return errors.TracedErrorNil("gitRepo")
	}

	if options == nil {
		return errors.TracedErrorNil("options")
	}

	gitRepoDir, err := gitRepo.GetAsLocalGitRepository()
	if err != nil {
		return err
	}

	err = p.RunInDirectory(gitRepoDir, options)
	if err != nil {
		return err
	}

	path, err := gitRepoDir.GetLocalPath()
	if err != nil {
		return err
	}

	if options.Verbose {
		gitStatusOutput, err := gitRepo.GetGitStatusOutput(options.Verbose)
		if err != nil {
			return err
		}

		logging.LogInfof(
			"Git status of repository '%s' after running pre-commit:\n%s",
			path,
			gitStatusOutput,
		)
	}

	return nil
}
