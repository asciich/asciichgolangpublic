package kubeconfigutils

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
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

type KubeConfig struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Clusters   []KubeConfigCluster `yaml:"clusters"`
	Contexts   []KubeConfigContext `yaml:"contexts"`
	Users      []KubeConfigUser    `yaml:"users"`
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

func (k *KubeConfig) GetClusterServerUrlAsString(clusterName string) (string, error) {
	cluster, err := k.GetClusterEntryByName(clusterName)
	if err != nil {
		return "", err
	}

	return cluster.GetServerUrlAsString()
}

func (k *KubeConfig) GetUserNameByContextName(contextName string) (userName string, err error) {
	if contextName == "" {
		return "", tracederrors.TracedErrorEmptyString("contextName")
	}

	contextEntry, err := k.GetContextEntryByName(contextName)
	if err != nil {
		return "", err
	}

	userName, err = contextEntry.GetUserName()
	if err != nil {
		return "", err
	}

	return userName, nil
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

func (k *KubeConfig) GetContextEntryByName(name string) (kubeConfigContext *KubeConfigContext, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	for _, c := range k.Contexts {
		if c.Name == name {
			kubeConfigContext = new(KubeConfigContext)
			*kubeConfigContext = c
			return kubeConfigContext, nil
		}
	}

	return nil, tracederrors.TracedErrorf("Context by name '%s' not found.", name)
}

func (k *KubeConfig) GetUserEntryByName(name string) (kubeConfigContext *KubeConfigUser, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	for _, c := range k.Users {
		if c.Name == name {
			kubeConfigContext = new(KubeConfigUser)
			*kubeConfigContext = c
			return kubeConfigContext, nil
		}
	}

	return nil, tracederrors.TracedErrorf("User by name '%s' not found.", name)
}

func (k *KubeConfig) GetClusterAndContextAndUserEntryByName(name string) (cluster *KubeConfigCluster, kubeConfigContext *KubeConfigContext, user *KubeConfigUser, err error) {
	if name == "" {
		return nil, nil, nil, tracederrors.TracedErrorEmptyString("name")
	}

	cluster, err = k.GetClusterEntryByName(name)
	if err != nil {
		return nil, nil, nil, err
	}

	kubeConfigContext, err = k.GetContextEntryByName(name)
	if err != nil {
		return nil, nil, nil, err
	}

	userName, err := kubeConfigContext.GetUserName()
	if err != nil {
		return nil, nil, nil, err
	}

	user, err = k.GetUserEntryByName(userName)
	if err != nil {
		return nil, nil, nil, err
	}

	return cluster, kubeConfigContext, user, nil
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

	for i, c := range k.Clusters {
		if c.Name == nameToAdd {
			k.Clusters[i] = *cluster
			return nil
		}
	}

	k.Clusters = append(k.Clusters, *cluster)

	return nil
}

func (k *KubeConfig) AddContextEntry(kubeConfigContext *KubeConfigContext) (err error) {
	if kubeConfigContext == nil {
		return tracederrors.TracedErrorNil("context")
	}

	nameToAdd := kubeConfigContext.Name

	for i, c := range k.Contexts {
		if c.Name == nameToAdd {
			k.Contexts[i] = *kubeConfigContext
			return nil
		}
	}

	k.Contexts = append(k.Contexts, *kubeConfigContext)

	return nil
}

func (k *KubeConfig) AddUserEntry(user *KubeConfigUser) (err error) {
	if user == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	nameToAdd := user.Name

	for i, c := range k.Users {
		if c.Name == nameToAdd {
			k.Users[i] = *user
			return nil
		}
	}

	k.Users = append(k.Users, *user)

	return nil
}

func (k *KubeConfig) AddClusterAndContextAndUserEntry(cluster *KubeConfigCluster, kubeConfigContext *KubeConfigContext, user *KubeConfigUser) (err error) {
	if cluster == nil {
		return tracederrors.TracedErrorNil("cluster")
	}

	if kubeConfigContext == nil {
		return tracederrors.TracedErrorNil("context")
	}

	if user == nil {
		return tracederrors.TracedErrorNil("user")
	}

	err = k.AddClusterEntry(cluster)
	if err != nil {
		return err
	}

	err = k.AddContextEntry(kubeConfigContext)
	if err != nil {
		return err
	}

	err = k.AddUserEntry(user)
	if err != nil {
		return err
	}

	return nil
}

// This function does an exec to "kubectl" using the given config file "path".
// Useful to validate if a written config "path" is understood by "kubectl".
func IsFilePathLoadableByKubectl(path string, verbose bool) (isLoadable bool, err error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString(path)
	}

	stdout, err := commandexecutor.Bash().RunCommandAndGetStdoutAsString(
		contextutils.GetVerbosityContextByBool(verbose),
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

// Use exec to invoke a "kubectl config get-context" with the given config "path".
// Useful to validate if the config is understood correctly by kubectl.
func ListContextNamesUsingKubectl(path string, verbose bool) (contextNames []string, err error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString(path)
	}

	contextNames, err = commandexecutor.Bash().RunCommandAndGetStdoutAsLines(
		contextutils.GetVerbosityContextByBool(verbose),
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

func (k *KubeConfig) GetClientKeyDataForUser(name string) (string, error) {
	user, err := k.GetUserEntryByName(name)
	if err != nil {
		return "", err
	}

	return user.GetClientKeyData()
}

func (k *KubeConfigCluster) GetServerUrlAsString() (string, error) {
	if k.Cluster.Server == "" {
		return "", tracederrors.TracedError("Kluster.Server not set")
	}

	return k.Cluster.Server, nil
}

func GetDefaultKubeConfigPath(ctx context.Context) (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get users home: %s", err)
	}

	kubeConfigPath := filepath.Join(dirname, ".kube", "config")

	logging.LogInfoByCtxf(ctx, "Default kube config path is: '%s'.", kubeConfigPath)

	return kubeConfigPath, nil
}
