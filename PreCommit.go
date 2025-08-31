package asciichgolangpublic

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type PreCommitService struct{}

func NewPreCommitService() (p *PreCommitService) {
	return new(PreCommitService)
}

func PreCommit() (p *PreCommitService) {
	return NewPreCommitService()
}

func (p *PreCommitService) GetAsPreCommitConfigFileOrNilIfContentIsInvalid(file filesinterfaces.File, verbose bool) (preCommitConfigFile *PreCommitConfigFile, err error) {
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

func (p *PreCommitService) MustGetAsPreCommitConfigFileOrNilIfContentIsInvalid(file filesinterfaces.File, verbose bool) (preCommitConfigFile *PreCommitConfigFile) {
	preCommitConfigFile, err := p.GetAsPreCommitConfigFileOrNilIfContentIsInvalid(file, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return preCommitConfigFile
}

func (p *PreCommitService) RunInDirectory(ctx context.Context, directoy filesinterfaces.Directory, options *PreCommitRunOptions) (err error) {
	if directoy == nil {
		return tracederrors.TracedErrorNil("directoy")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
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

	_, err = commandexecutorbashoo.Bash().RunCommand(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: preCommitCommand,
		},
	)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pre commit successfully run in '%s'.", path)

	return nil
}

func (p *PreCommitService) RunInGitRepository(ctx context.Context, gitRepo GitRepository, options *PreCommitRunOptions) (err error) {
	if gitRepo == nil {
		return tracederrors.TracedErrorNil("gitRepo")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	gitRepoDir, err := gitRepo.GetAsLocalGitRepository()
	if err != nil {
		return err
	}

	err = p.RunInDirectory(ctx, gitRepoDir, options)
	if err != nil {
		return err
	}

	path, err := gitRepoDir.GetLocalPath()
	if err != nil {
		return err
	}

	gitStatusOutput, err := gitRepo.GetGitStatusOutput(ctx)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(
		ctx,
		"Git status of repository '%s' after running pre-commit:\n%s",
		path,
		gitStatusOutput,
	)

	return nil
}
