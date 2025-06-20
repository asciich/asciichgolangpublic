package kubernetesimplementationindependend_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesimplementationindependend"
)

func getReplicaSet(name string, namespace string) (yaml string) {
	yaml += "---\n"
	yaml += "apiVersion: apps/v1\n"
	yaml += "kind: ReplicaSet\n"
	yaml += "metadata:\n"
	yaml += "    name: " + name + "\n"
	if namespace != "" {
		yaml += "    namespace: " + namespace + "\n"
	}
	yaml += "    labels:\n"
	yaml += "        app: guestbook\n"
	yaml += "        tier: frontend\n"
	yaml += "spec:\n"
	yaml += "    # modify replicas according to your case\n"
	yaml += "    replicas: 3\n"
	yaml += "    selector:\n"
	yaml += "        matchLabels:\n"
	yaml += "            tier: frontend\n"
	yaml += "    template:\n"
	yaml += "        metadata:\n"
	yaml += "            labels:\n"
	yaml += "                tier: frontend\n"
	yaml += "        spec:\n"
	yaml += "            containers:\n"
	yaml += "                - name: php-redis\n"
	yaml += "                  image: us-docker.pkg.dev/google-samples/containers/gke/gb-frontend:v5\n"

	return yaml
}

func getNginxDeployment(name string, namespace string) (yaml string) {
	yaml += "---\n"
	yaml += "apiVersion: apps/v1\n"
	yaml += "kind: Deployment\n"
	yaml += "metadata:\n"
	yaml += "  name: " + name + "\n"
	if namespace != "" {
		yaml += "  namespace: " + namespace + "\n"
	}
	yaml += "  labels:\n"
	yaml += "    app: nginx\n"
	yaml += "spec:\n"
	yaml += "  replicas: 3\n"
	yaml += "  selector:\n"
	yaml += "    matchLabels:\n"
	yaml += "      app: nginx\n"
	yaml += "  template:\n"
	yaml += "    metadata:\n"
	yaml += "      labels:\n"
	yaml += "        app: nginx\n"
	yaml += "    spec:\n"
	yaml += "      containers:\n"
	yaml += "      - name: nginx\n"
	yaml += "        image: nginx:1.14.2\n"
	yaml += "        ports:\n"
	yaml += "        - containerPort: 80\n"

	return yaml
}

func Test_unmarshalReplicaset(t *testing.T) {
	t.Run("namespace", func(t *testing.T) {
		u, err := kubernetesimplementationindependend.UnmarshalResourceYaml(getReplicaSet("abc", "def"))
		require.NoError(t, err)
		require.Len(t, u, 1)
		require.EqualValues(t, "abc", u[0].Name())
		require.EqualValues(t, "def", u[0].Namespace())
		require.EqualValues(t, "ReplicaSet", u[0].Kind())
		require.EqualValues(t, "apps/v1", u[0].ApiVersion())
	})

	t.Run("no namespace", func(t *testing.T) {
		u, err := kubernetesimplementationindependend.UnmarshalResourceYaml(getReplicaSet("abc", ""))
		require.NoError(t, err)
		require.Len(t, u, 1)
		require.EqualValues(t, "abc", u[0].Name())
		require.EqualValues(t, "", u[0].Namespace())
		require.EqualValues(t, "ReplicaSet", u[0].Kind())
		require.EqualValues(t, "apps/v1", u[0].ApiVersion())
	})
}

func Test_unmarshalNginxDeplyoments(t *testing.T) {
	t.Run("namespace", func(t *testing.T) {
		u, err := kubernetesimplementationindependend.UnmarshalResourceYaml(getNginxDeployment("abc", "def"))
		require.NoError(t, err)
		require.Len(t, u, 1)
		require.EqualValues(t, "abc", u[0].Name())
		require.EqualValues(t, "def", u[0].Namespace())
		require.EqualValues(t, "Deployment", u[0].Kind())
	})

	t.Run("no namespace", func(t *testing.T) {
		u, err := kubernetesimplementationindependend.UnmarshalResourceYaml(getNginxDeployment("abc", ""))
		require.NoError(t, err)
		require.Len(t, u, 1)
		require.EqualValues(t, "abc", u[0].Name())
		require.EqualValues(t, "", u[0].Namespace())
		require.EqualValues(t, "Deployment", u[0].Kind())
	})
}

func TestSortResourcesYaml(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		sorted, err := kubernetesimplementationindependend.SortResourcesYaml("")
		require.NoError(t, err)

		require.EqualValues(t, "\n", sorted)
	})

	t.Run("single deployment", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "")
		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment)
		require.NoError(t, err)

		require.EqualValues(t, exampleDeployment, sorted)
	})

	t.Run("two ordered deployments", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "")
		exampleDeployment1 := getNginxDeployment("example1", "")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment + "\n---\n" + exampleDeployment1)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(t, []string{exampleDeployment, exampleDeployment1}, splitted)
	})

	t.Run("two ordered deployments 2", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("a", "")
		exampleDeployment1 := getNginxDeployment("b", "")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment + "\n---\n" + exampleDeployment1)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(t, []string{exampleDeployment, exampleDeployment1}, splitted)
	})

	t.Run("two unordered deployments", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "")
		exampleDeployment1 := getNginxDeployment("example1", "")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment1 + "\n---\n" + exampleDeployment)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(t, []string{exampleDeployment, exampleDeployment1}, splitted)
	})

	t.Run("two unordered deployments 2", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("a", "")
		exampleDeployment1 := getNginxDeployment("b", "")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment1 + "\n---\n" + exampleDeployment)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(t, []string{exampleDeployment, exampleDeployment1}, splitted)
	})

	t.Run("with namespaces", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "namespace")
		exampleDeployment1 := getNginxDeployment("example", "")
		exampleDeployment2 := getNginxDeployment("example", "aaaa")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment + "\n---\n" + exampleDeployment1 + "\n---\n" + exampleDeployment2)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(
			t,
			[]string{exampleDeployment1, exampleDeployment2, exampleDeployment},
			splitted,
		)
	})

	t.Run("already sorted by kind", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "namespace")
		exampleReplicaSet := getReplicaSet("example", "namespace")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment + "\n---\n" + exampleReplicaSet)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(
			t,
			[]string{exampleDeployment, exampleReplicaSet},
			splitted,
		)
	})

	t.Run("Unordered by kind", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "namespace")
		exampleReplicaSet := getReplicaSet("example", "namespace")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleReplicaSet + "\n---\n" + exampleDeployment)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(
			t,
			[]string{exampleDeployment, exampleReplicaSet},
			splitted,
		)
	})

	t.Run("Unordered by kind and namespace", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "namespace")
		exampleReplicaSet := getReplicaSet("example", "")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleReplicaSet + "\n---\n" + exampleDeployment)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(
			t,
			[]string{exampleDeployment, exampleReplicaSet},
			splitted,
		)
	})

	t.Run("with namespaces and kind", func(t *testing.T) {
		exampleDeployment := getNginxDeployment("example", "namespace")
		exampleDeployment1 := getNginxDeployment("example", "")
		exampleDeployment2 := getNginxDeployment("example", "aaaa")
		exampleReplicaSet := getReplicaSet("example", "aaaa")

		sorted, err := kubernetesimplementationindependend.SortResourcesYaml(exampleDeployment + "\n---\n" + exampleDeployment1 + "\n---\n" + exampleDeployment2 + "\n---\n" + exampleReplicaSet)
		require.NoError(t, err)

		splitted := yamlutils.SplitMultiYaml(sorted)

		require.EqualValues(
			t,
			[]string{exampleDeployment1, exampleDeployment2, exampleReplicaSet, exampleDeployment},
			splitted,
		)
	})
}
