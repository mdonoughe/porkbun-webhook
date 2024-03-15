package porkbun

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook"
	acme "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/nrdcg/porkbun"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	klog "k8s.io/klog/v2"
)

type PorkbunSolver struct {
	kube *kubernetes.Clientset
}

func (e *PorkbunSolver) Name() string {
	return "porkbun"
}

type Config struct {
	ApiKeySecretRef    corev1.SecretKeySelector `json:"apiKeySecretRef"`
	SecretKeySecretRef corev1.SecretKeySelector `json:"secretKeySecretRef"`
}

func (e *PorkbunSolver) readConfig(request *acme.ChallengeRequest) (*porkbun.Client, error) {
	config := Config{}

	if request.Config != nil {
		if err := json.Unmarshal(request.Config.Raw, &config); err != nil {
			return nil, errors.Wrap(err, "config error")
		}
	}

	apiKey, err := e.resolveSecretRef(config.ApiKeySecretRef, request)
	if err != nil {
		return nil, err
	}

	secretKey, err := e.resolveSecretRef(config.SecretKeySecretRef, request)
	if err != nil {
		return nil, err
	}

	return porkbun.New(secretKey, apiKey), nil
}

func (e *PorkbunSolver) resolveSecretRef(selector corev1.SecretKeySelector, ch *acme.ChallengeRequest) (string, error) {
	secret, err := e.kube.CoreV1().Secrets(ch.ResourceNamespace).Get(context.Background(), selector.Name, metav1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "get error for secret %q %q", ch.ResourceNamespace, selector.Name)
	}

	bytes, ok := secret.Data[selector.Key]
	if !ok {
		return "", errors.Errorf("secret %q %q does not contain key %q", ch.ResourceNamespace, selector.Name, selector.Key)
	}

	return string(bytes), nil
}

func (e *PorkbunSolver) Present(ch *acme.ChallengeRequest) error {
	klog.Infof("Handling present request for %q %q", ch.ResolvedFQDN, ch.Key)

	client, err := e.readConfig(ch)
	if err != nil {
		return errors.Wrap(err, "initialization error")
	}

	domain := strings.TrimSuffix(ch.ResolvedZone, ".")
	entity := strings.TrimSuffix(ch.ResolvedFQDN, "."+ch.ResolvedZone)
	name := strings.TrimSuffix(ch.ResolvedFQDN, ".")
	records, err := client.RetrieveRecords(context.Background(), domain)
	if err != nil {
		return errors.Wrap(err, "retrieve records error")
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Name == name && record.Content == ch.Key {
			klog.Infof("Record %s is already present", record.ID)
			return nil
		}
	}

	id, err := client.CreateRecord(context.Background(), domain, porkbun.Record{
		Name:    entity,
		Type:    "TXT",
		Content: ch.Key,
		TTL:     "60",
	})
	if err != nil {
		return errors.Wrap(err, "create record error")
	}

	klog.Infof("Created record %v", id)
	return nil
}

func (e *PorkbunSolver) CleanUp(ch *acme.ChallengeRequest) error {
	klog.Infof("Handling cleanup request for %q %q", ch.ResolvedFQDN, ch.Key)

	client, err := e.readConfig(ch)
	if err != nil {
		return errors.Wrap(err, "initialization error")
	}

	domain := strings.TrimSuffix(ch.ResolvedZone, ".")
	name := strings.TrimSuffix(ch.ResolvedFQDN, ".")
	records, err := client.RetrieveRecords(context.Background(), domain)
	if err != nil {
		return errors.Wrap(err, "retrieve records error")
	}

	for _, record := range records {
		if record.Type == "TXT" && record.Name == name && record.Content == ch.Key {
			id, err := strconv.ParseInt(record.ID, 10, 32)
			if err != nil {
				return errors.Wrap(err, "found TXT record, but it's ID is malformed")
			}

			record.Content = ch.Key
			err = client.DeleteRecord(context.Background(), domain, int(id))
			if err != nil {
				return errors.Wrap(err, "delete record error")
			}

			klog.Infof("Deleted record %v", id)
			return nil
		}
	}

	klog.Info("No matching record to delete")

	return nil
}

func (e *PorkbunSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	klog.Info("Initializing")

	kube, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return errors.Wrap(err, "kube client creation error")
	}

	e.kube = kube
	return nil
}

func New() webhook.Solver {
	return &PorkbunSolver{}
}
