package k8sutil

import (
	"fmt"
	"k8s.io/api/core/v1"
)

const TlsPort = 443

// -----------------------------------------------------------------------------

func stunnelEnvVars(
	verifyChain string,
	listenPort int32,
	destHostAndPort string,
	checkHost string,
	isNginxIngressStyleCertSecret bool,
	isClientMode bool,
	clientModeSpringBoardDelaySeconds int32,
	index int,
) []v1.EnvVar {

	springboardListenPort := fmt.Sprintf("%d", 6789 + index)
	stunnelEnv := []v1.EnvVar {
		{
			Name: "STUNNEL_VERIFY_CHAIN",
			Value: verifyChain,
		},
		{
			Name: "STUNNEL_ACCEPT_PORT",
			Value: fmt.Sprintf("%d", listenPort),
		},
		{
			Name: "STUNNEL_CONNECT",
			Value: destHostAndPort,
		},
		{
			Name: "SPRINGBOARD_LISTEN_PORT",
			Value: springboardListenPort,
		},
		{
			Name: "STUNNEL_CERT_FILE",
			Value: "/etc/stunnel/certs/tls.crt",
		},
		{
			Name: "STUNNEL_KEY_FILE",
			Value: "/etc/stunnel/certs/tls.key",
		},
	}
	if isClientMode {
		stunnelEnv = append(stunnelEnv, v1.EnvVar{
			Name: "STUNNEL_CLIENT_MODE",
			Value: "yes",
		})
		if clientModeSpringBoardDelaySeconds > 0 {
			stunnelEnv = append(stunnelEnv, v1.EnvVar{
				Name: "SPRINGBOARD_DELAY_SECONDS",
				Value: fmt.Sprintf("%d", clientModeSpringBoardDelaySeconds),
			})
		}
	}
	if checkHost != "" {
		stunnelEnv = append(stunnelEnv, v1.EnvVar{
			Name: "STUNNEL_CHECKHOST_LINE",
			Value: fmt.Sprintf("checkHost=%s", checkHost),
		})
	}
	if isNginxIngressStyleCertSecret {
		// The server cert file names are different because they follow
		// the nginx ingress controller conventions (tls.crt and tls.key).
		// There is no CA for client certificate verification (for now).
		stunnelEnv = append(stunnelEnv,
			v1.EnvVar{
				Name: "STUNNEL_CERT_FILE",
				Value: "/etc/stunnel/certs/tls.crt",
			},
			v1.EnvVar{
				Name: "STUNNEL_KEY_FILE",
				Value: "/etc/stunnel/certs/tls.key",
			},
		)
	}
	return stunnelEnv
}

// -----------------------------------------------------------------------------

func InsertStunnel(
	containerName string,
	listenPort int32,
	verifyChain string,
	destHostAndPort string,
	checkHost string,
	tlsSecretName string,
	isNginxIngressStyleCertSecret bool,
	isClientMode bool,
	volumes []v1.Volume,
	containers []v1.Container,
	clientModeSpringBoardDelaySeconds int32,
	index int,
) ([]v1.Volume, []v1.Container) {

	stunnelEnv := stunnelEnvVars(verifyChain, listenPort, destHostAndPort,
		checkHost, isNginxIngressStyleCertSecret, isClientMode,
		clientModeSpringBoardDelaySeconds, index)

	volumeName := fmt.Sprintf("%s-certs", containerName)
	volumes = append(volumes, v1.Volume{
		Name: volumeName,
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: tlsSecretName,
			},
		},
	})
	containers = append(containers, v1.Container{
		Name: containerName,
		Image: "platform9systems/stunnel-with-springboard:instrumented-3",
		Ports: []v1.ContainerPort{
			{
				ContainerPort: listenPort,
			},
		},
		Env: stunnelEnv,
		VolumeMounts: []v1.VolumeMount{
			{
				Name: volumeName,
				ReadOnly: true,
				MountPath: "/etc/stunnel/certs",
			},
		},
	})
	return volumes, containers
}


