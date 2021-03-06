package unifiedpushserver

import (
	"fmt"
	"os"
	"strings"

	aerogearv1alpha1 "github.com/aerogear/unifiedpush-operator/pkg/apis/aerogear/v1alpha1"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type image struct {
	defaultImage string
	envVarName   string
}

func (i image) image() string {
	customImage, found := os.LookupEnv(i.envVarName)
	if found {
		return customImage
	}
	return i.defaultImage
}

var (
	proxy       = image{"docker.io/openshift/oauth-proxy:v1.1.0", "OAUTH_PROXY_IMAGE"}
	postgresql  = image{"docker.io/centos/postgresql-96-centos7:9.6", "POSTGRESQL_IMAGE"}
	unifiedpush = image{"docker.io/aerogear/unifiedpush-wildfly-plain:2.2.1.Final", "UNIFIEDPUSH_IMAGE"}
)

func labels(cr *aerogearv1alpha1.UnifiedPushServer, suffix string) map[string]string {
	return map[string]string{
		"app":     cr.Name,
		"service": fmt.Sprintf("%s-%s", cr.Name, suffix),
	}
}

// objectMeta returns the default ObjectMeta for all the other objects here
func objectMeta(cr *aerogearv1alpha1.UnifiedPushServer, suffix string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%s", cr.Name, suffix),
		Namespace: cr.Namespace,
		Labels:    labels(cr, suffix),
	}
}

func generatePassword() (string, error) {
	generatedPassword, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "error generating password")
	}
	return strings.Replace(generatedPassword.String(), "-", "", -1), nil
}
