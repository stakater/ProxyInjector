package constants

const (
	EnabledAnnotation = "authproxy.stakater.com/enabled"

	SourceServiceNameAnnotation = "authproxy.stakater.com/source-service-name"

	ImageNameAnnotation       = "authproxy.stakater.com/image-name"
	ImageTagAnnotation        = "authproxy.stakater.com/image-tag"
	ImagePullPolicyAnnotation = "authproxy.stakater.com/image-pull-policy"

	ConfigAnnotation         = "authproxy.stakater.com/config"
	UpstreamUrlAnnotation    = "authproxy.stakater.com/upstream-url"
	RedirectionUrlAnnotation = "authproxy.stakater.com/redirection-url"

	EnableAuthorizationAnnotation = "authproxy.stakater.com/enable-authorization-header"

	ResponseHeaderTimeoutAnnotation = "authproxy.stakater.com/upstream-response-header-timeout"
	TimeoutAnnotation               = "authproxy.stakater.com/upstream-timeout"
	KeepaliveTimeoutAnnotation      = "authproxy.stakater.com/upstream-keepalive-timeout"
	ServerReadTimeoutAnnotation     = "authproxy.stakater.com/server-read-timeout"
	ServerWriteTimeoutAnnotation    = "authproxy.stakater.com/server-write-timeout"
)
