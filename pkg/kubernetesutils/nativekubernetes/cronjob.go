package nativekubernetes

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CronJobExists(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, cronJobName string) (bool, error) {
	if clientset == nil {
		return false, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return false, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if cronJobName == "" {
		return false, tracederrors.TracedErrorEmptyString("cronJobName")
	}

	logging.LogInfoByCtxf(ctx, "Check if CronJob '%s' in namespace '%s' exists.", cronJobName, namespaceName)

	_, err := clientset.BatchV1().CronJobs(namespaceName).Get(ctx, cronJobName, metav1.GetOptions{})
	if err != nil {
		logging.LogInfoByCtxf(ctx, "CronJob '%s' in namespace '%s' does not exist.", cronJobName, namespaceName)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "CronJob '%s' in namespace '%s' exists.", cronJobName, namespaceName)
	return true, nil
}

func DeleteCronJob(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, cronJobName string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if cronJobName == "" {
		return tracederrors.TracedErrorEmptyString("cronJobName")
	}

	logging.LogInfoByCtxf(ctx, "Delete CronJob '%s' in namespace '%s' started.", cronJobName, namespaceName)

	exists, err := CronJobExists(ctx, clientset, namespaceName, cronJobName)
	if err != nil {
		return err
	}

	if exists {
		err = clientset.BatchV1().CronJobs(namespaceName).Delete(ctx, cronJobName, metav1.DeleteOptions{})
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete CronJob '%s' in namespace '%s'.", cronJobName, namespaceName)
		}

		logging.LogChangedByCtxf(ctx, "CronJob '%s' in namespace '%s' deleted.", cronJobName, namespaceName)
	} else {
		logging.LogInfoByCtxf(ctx, "CronJob '%s' in namespace '%s' does not exist. Skip delete.", cronJobName, namespaceName)
	}

	logging.LogInfoByCtxf(ctx, "Delete CronJob '%s' in namespace '%s' finished.", cronJobName, namespaceName)

	return nil
}

func CreateCronJob(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string, cronJobName string, schedule string, image string, command []string, labels map[string]string) error {
	if clientset == nil {
		return tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return tracederrors.TracedErrorEmptyString("namespaceName")
	}

	if cronJobName == "" {
		return tracederrors.TracedErrorEmptyString("cronJobName")
	}

	if schedule == "" {
		return tracederrors.TracedErrorEmptyString("schedule")
	}

	if image == "" {
		return tracederrors.TracedErrorEmptyString("image")
	}

	logging.LogInfoByCtxf(ctx, "Create CronJob '%s' in namespace '%s' started.", cronJobName, namespaceName)

	cronJob := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:   cronJobName,
			Labels: labels,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    cronJobName,
									Image:   image,
									Command: command,
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	_, err := clientset.BatchV1().CronJobs(namespaceName).Create(ctx, cronJob, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			logging.LogInfoByCtxf(ctx, "CronJob '%s' in namespace '%s' already exists.", cronJobName, namespaceName)
			return nil
		}
		return tracederrors.TracedErrorf("Failed to create CronJob '%s' in namespace '%s': %w", cronJobName, namespaceName, err)
	}

	logging.LogChangedByCtxf(ctx, "CronJob '%s' in namespace '%s' created.", cronJobName, namespaceName)
	logging.LogInfoByCtxf(ctx, "Create CronJob '%s' in namespace '%s' finished.", cronJobName, namespaceName)

	return nil
}

func ListCronJobs(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	if clientset == nil {
		return nil, tracederrors.TracedErrorNil("clientset")
	}

	if namespaceName == "" {
		return nil, tracederrors.TracedErrorEmptyString("namespaceName")
	}

	cronJobList, err := clientset.BatchV1().CronJobs(namespaceName).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list CronJobs in namespace '%s'.", namespaceName)
	}

	cronJobNames := []string{}
	for _, cj := range cronJobList.Items {
		cronJobNames = append(cronJobNames, cj.Name)
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' CronJobs in namespace '%s'.", len(cronJobNames), namespaceName)

	return cronJobNames, nil
}

func ListCronJobNames(ctx context.Context, clientset *kubernetes.Clientset, namespaceName string) ([]string, error) {
	return ListCronJobs(ctx, clientset, namespaceName)
}
