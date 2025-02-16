package kubeconfig

import (
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"gopkg.in/yaml.v3"
)

type KubeConfigClusterCluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

type KubeConfigCluster struct {
	Name    string                   `yaml:"name"`
	Cluster KubeConfigClusterCluster `yaml:"cluster"`
}

type KubeConfigContext struct {
	Name    string `yaml:"name"`
	Context struct {
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		User      string `yaml:"user"`
	} `yaml:"context"`
}

type KubeConfigUser struct {
	Name string `yaml:"name"`
	User struct {
		ClientCertificateData string `yaml:"client-certificate-data"`
		ClientKeyData         string `yaml:"client-key-data"`
		Username              string `yaml:"username"`
		Password              string `yaml:"password"`
	} `yaml:"user"`
}

type KubeConfig struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Clusters   []KubeConfigCluster `yaml:"clusters"`
	Contexts   []KubeConfigContext `yaml:"contexts"`
	Users      []KubeConfigUser    `yaml:"users"`
}

func MustLoadFromFilePath(path string, verbose bool) (config *KubeConfig) {
	config, err := LoadFromFilePath(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return config
}

func LoadFromFilePath(path string, verbose bool) (config *KubeConfig, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	file, err := files.GetLocalFileByPath(path)
	if err != nil {
		return nil, err
	}

	return LoadFromFile(file, verbose)
}

func LoadFromFile(file files.File, verbose bool) (config *KubeConfig, err error) {
	if file == nil {
		return nil, tracederrors.TracedErrorNil("file")
	}

	path, err := file.GetPath()
	if err != nil {
		return nil, err
	}

	content, err := file.ReadAsBytes()
	if err != nil {
		return nil, err
	}

	config = new(KubeConfig)

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to load kubeConfig '%s' as yaml: %w", path, err)
	}

	if verbose {
		logging.LogInfof("Loaded kubeConfig '%s'.", path)
	}

	return config, nil
}

func (k *KubeConfig) MustGetServerNames() (serverNames []string) {
	serverNames, err := k.GetServerNames()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return serverNames
}

func (k *KubeConfig) MustGetClusterNames() (clusterNames []string) {
	clusterNames, err := k.GetClusterNames()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return clusterNames
}

func (k *KubeConfig) GetClusterNames() (clusterNames []string, err error) {
	for _, entry := range k.Clusters {
		toAdd := entry.Name
		if toAdd == "" {
			return nil, tracederrors.TracedErrorf("Got empty cluster name toAdd")
		}

		clusterNames = append(clusterNames, toAdd)
	}

	sort.Strings(clusterNames)

	if len(clusterNames) <= 0 {
		return nil, tracederrors.TracedError("No cluster names in config found.")
	}

	return clusterNames, nil
}

func (k *KubeConfig) GetServerNames() (serverNames []string, err error) {
	for _, entry := range k.Clusters {
		toAdd := entry.Cluster.Server
		if toAdd == "" {
			return nil, tracederrors.TracedErrorf("Got empty server name toAdd")
		}

		serverNames = append(serverNames, toAdd)
	}

	sort.Strings(serverNames)

	if len(serverNames) <= 0 {
		return nil, tracederrors.TracedError("No server names in config found.")
	}

	return serverNames, nil
}

func MustMergeConfig(configs ...*KubeConfig) (merged *KubeConfig) {
	merged, err := MergeConfig(configs...)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return merged
}

func MergeConfig(configs ...*KubeConfig) (merged *KubeConfig, err error) {
	if len(configs) <= 0 {
		return nil, tracederrors.TracedError("No KubeConfig elements to merge.")
	}

	merged = configs[0].GetDeepCopy()

	for _, toAdd := range configs {
		err = merged.AddConfig(toAdd)
		if err != nil {
			return nil, err
		}
	}

	return merged, nil
}

func (k *KubeConfig) GetDeepCopy() (copy *KubeConfig) {
	copy = new(KubeConfig)

	*copy = *k

	return copy
}

func (k *KubeConfig) GetClusterEntryByName(name string) (cluster *KubeConfigCluster, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	for _, c := range k.Clusters {
		if c.Name == name {
			cluster = new(KubeConfigCluster)
			*cluster = c
			return cluster, nil
		}
	}

	return nil, tracederrors.TracedErrorf("Cluster by name '%s' not found.", name)
}

func (k *KubeConfig) GetContextEntryByName(name string) (context *KubeConfigContext, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	for _, c := range k.Contexts {
		if c.Name == name {
			context = new(KubeConfigContext)
			*context = c
			return context, nil
		}
	}

	return nil, tracederrors.TracedErrorf("Context by name '%s' not found.", name)
}

func (k *KubeConfig) GetUserEntryByName(name string) (context *KubeConfigUser, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	for _, c := range k.Users {
		if c.Name == name {
			context = new(KubeConfigUser)
			*context = c
			return context, nil
		}
	}

	return nil, tracederrors.TracedErrorf("User by name '%s' not found.", name)
}

func (k *KubeConfig) GetClusterAndContextAndUserEntryByName(name string) (cluster *KubeConfigCluster, context *KubeConfigContext, user *KubeConfigUser, err error) {
	if name == "" {
		return nil, nil, nil, tracederrors.TracedErrorEmptyString("name")
	}

	cluster, err = k.GetClusterEntryByName(name)
	if err != nil {
		return nil, nil, nil, err
	}

	context, err = k.GetContextEntryByName(name)
	if err != nil {
		return nil, nil, nil, err
	}

	user, err = k.GetUserEntryByName(name)
	if err != nil {
		return nil, nil, nil, err
	}

	return cluster, context, user, nil
}

func (k *KubeConfig) AddConfig(toAdd *KubeConfig) (err error) {
	if toAdd == nil {
		return tracederrors.TracedErrorNil("toAdd")
	}

	namesToAdd, err := toAdd.GetClusterNames()
	if err != nil {
		return err
	}

	for _, name := range namesToAdd {
		cluster, context, user, err := toAdd.GetClusterAndContextAndUserEntryByName(name)
		if err != nil {
			return err
		}

		err = k.AddClusterAndContextAndUserEntry(cluster, context, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KubeConfig) AddClusterEntry(cluster *KubeConfigCluster) (err error) {
	if cluster == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	nameToAdd := cluster.Name

	for _, c := range k.Clusters {
		if c.Name == nameToAdd {
			c = *cluster
			return nil
		}
	}

	k.Clusters = append(k.Clusters, *cluster)

	return nil
}

func (k *KubeConfig) AddContextEntry(context *KubeConfigContext) (err error) {
	if context == nil {
		return tracederrors.TracedErrorNil("context")
	}

	nameToAdd := context.Name

	for _, c := range k.Contexts {
		if c.Name == nameToAdd {
			c = *context
			return nil
		}
	}

	k.Contexts = append(k.Contexts, *context)

	return nil
}

func (k *KubeConfig) AddUserEntry(user *KubeConfigUser) (err error) {
	if user == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	nameToAdd := user.Name

	for _, c := range k.Users {
		if c.Name == nameToAdd {
			c = *user
			return nil
		}
	}

	k.Users = append(k.Users, *user)

	return nil
}

func (k *KubeConfig) AddClusterAndContextAndUserEntry(cluster *KubeConfigCluster, context *KubeConfigContext, user *KubeConfigUser) (err error) {
	if cluster == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	if context == nil {
		return tracederrors.TracedErrorNil("context")
	}

	if user == nil {
		return tracederrors.TracedErrorNil("user")
	}

	err = k.AddClusterEntry(cluster)
	if err != nil {
		return err
	}

	err = k.AddContextEntry(context)
	if err != nil {
		return err
	}

	err = k.AddUserEntry(user)
	if err != nil {
		return err
	}

	return nil
}

func MustIsFilePathLoadableByKubectl(path string, verbose bool) (isLoadable bool) {
	isLoadable, err := IsFilePathLoadableByKubectl(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return isLoadable
}

// This function does an exex to "kubectl" using the given config file "path".
// Useful to validate if a written config "path" is understood by "kubectl".
func IsFilePathLoadableByKubectl(path string, verbose bool) (isLoadable bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString(path)
	}

	stdout, err := commandexecutor.Bash().RunCommandAndGetStdoutAsString(
		&parameteroptions.RunCommandOptions{
			Command: []string{"KUBECONFIG=" + path, "bash", "-c", "kubectl config get-contexts &> /dev/null && echo YES || echo NO"},
		},
	)
	if err != nil {
		return false, err
	}

	stdout = strings.TrimSpace(stdout)
	if stdout == "YES" {
		isLoadable = true
	} else if stdout == "NO" {
		isLoadable = false
	} else {
		return false, tracederrors.TracedErrorf("Unexpected output: '%s'", stdout)
	}

	if verbose {
		if isLoadable {
			logging.LogInfof("Kube config '%s' is loadable by kubectl.", path)
		} else {
			logging.LogInfof("Kube config '%s' is not loadable by kubectl.", path)
		}
	}

	return isLoadable, nil
}

func (k *KubeConfig) MustWriteToTemporaryFileAndGetPath(verbose bool) (tempFilePath string) {
	tempFilePath, err := k.WriteToTemporaryFileAndGetPath(verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return tempFilePath
}

func (k *KubeConfig) WriteToTemporaryFileAndGetPath(verbose bool) (tempFilePath string, err error) {
	tempFilePath, err = tempfiles.CreateEmptyTemporaryFileAndGetPath(verbose)
	if err != nil {
		return "", err
	}

	err = k.WriteToFileByPath(tempFilePath, verbose)
	if err != nil {
		return "", err
	}

	return tempFilePath, nil

}

func (k *KubeConfig) GetAsYamlString() (yamlSring string, err error) {
	content, err := yaml.Marshal(k)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to marshal KubeConfig as yaml: %w", err)
	}

	return string(content), nil
}

func (k *KubeConfig) WriteToFileByPath(path string, verbose bool) (err error) {
	if path == "" {
		return tracederrors.TracedErrorEmptyString(path)
	}

	outFile, err := files.GetLocalFileByPath(path)
	if err != nil {
		return err
	}

	return k.WriteToFile(outFile, verbose)
}

func (k *KubeConfig) WriteToFile(outFile files.File, verbose bool) (err error) {
	if outFile == nil {
		return tracederrors.TracedErrorNil("outfile")
	}

	path, err := outFile.GetPath()
	if err != nil {
		return err
	}

	content, err := k.GetAsYamlString()
	if err != nil {
		return err
	}

	err = outFile.WriteString(content, verbose)
	if err != nil {
		return err
	}

	if verbose {
		logging.LogChangedf(
			"Wrote KubeConfig to '%s'", path,
		)
	}

	return nil
}

func MustListContextNamesUsingKubectl(path string, verbose bool) (contextNames []string) {
	contextNames, err := ListContextNamesUsingKubectl(path, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return contextNames
}

// Use exec to invoke a "kubectl config get-context" with the given config "path".
// Useful to validate if the config is understood correctly by kubectl.
func ListContextNamesUsingKubectl(path string, verbose bool) (contextNames []string, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString(path)
	}

	contextNames, err = commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
		&parameteroptions.RunCommandOptions{
			Command: []string{"KUBECONFIG=" + path, "bash", "-c", "kubectl config get-contexts -o name"},
		},
	)
	if err != nil {
		return nil, err
	}

	sort.Strings(contextNames)

	return contextNames, nil
}
