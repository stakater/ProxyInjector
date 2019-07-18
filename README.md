# ![](assets/web/proxyinjector-round-100px.png) Proxy Injector
A Kubernetes controller to inject an authentication proxy container to relevant pods

[![Get started with Stakater](https://stakater.github.io/README/stakater-github-banner.png)](http://stakater.com/?utm_source=ProxyInjector&utm_medium=github)

## Problem Statement

We want to automatically inject an authentication proxy container in a pod, for any deployment that requires to connect
 to our SSO provider, instead of manually adding a sidecar container with each deployment 

## Solution

This controller will continuously watch deployments in specific or all namespaces, and automatically add a sidecar container
 for the authentication proxy. Configuration for the proxy is managed through annotations of the respective deployment
 or with ConfigMap of the ProxyInjector.

### Supported proxies

For now the ProxyInjector only supports [Keycloak Gatekeeper](https://github.com/keycloak/keycloak-gatekeeper)
 as the authentication proxy, to work with [Keycloak Server](https://github.com/keycloak/keycloak)


## Usage

The following quickstart let's you set up ProxyInjector:

1. Add configuration to the ProxyInjector
    The following arguments can either be added to the proxy injector `config.yaml` in the ConfigMap for centralized configuration,
     or as annotations on the individual target deployments with a `authproxy.stakater.com/` prefix. In case of both,
     the deployment annotation values will override the central configuration. 

    | Key              | Description                                                               |
    |------------------|---------------------------------------------------------------------------|
    | listen           | the interface address and port the proxy should be listening on           |
    | upstream-url     | url for the upstream endpoint you wish to proxy                           |
    | resources        | list of resources to proxy uri, methods, roles                            |
    | client-id        | client id used to authenticate to the oauth service                       |
    | client-secret    | client secret used to authenticate to the oauth service                   |
    | gatekeeper-image | Keycloak Gatekeeper image e.g. `keycloak/keycloak-gatekeeper:6.0.1` |

The rest of the available options can be found at the [Keycloak Gatekeeper documentation](https://www.keycloak.org/docs/latest/securing_apps/index.html#configuration-options)

2. Deploy the controller by running the following command:

    For Kubernetes Cluster using kubectl
   ```bash
   kubectl apply -f https://raw.githubusercontent.com/stakater/ProxyInjector/master/deployments/kubernetes/proxyinjector.yaml -n default

3. When deploying any application that needs Keycloak authentication, add the following annotations to the `deployment`. The `service` will not need changes as such, all configuration can be provided as annotations in the deployment for the app. And proxy injector automatically modifies the service when injecting the sidecar container.
  
    | Key                                        | Description                                                                                                                                       |
    |--------------------------------------------|--------------------------------------------------------|
    | authproxy.stakater.com/enabled             | (true/false, default=false) Enables Keycloak gatekeeper configuration |
    | authproxy.stakater.com/source-service-name | Name of service that needs to be reconfigured to connect to the proxy. instead of the service directly routing to the app container, it will now route to the proxy sidecar instead. |
    | authproxy.stakater.com/target-port         | (default=80) the port on the pod where the proxy sidecar (keycloak gatekeeper) will be listening. If not specified, the default value of 80 is used. This port should match the `listen` configuration |

    The `authproxy.stakater.com/listen` annotation or the `listen` property in the ProxyInjector ConfigMap should
    specify where the proxy sidecar will listen for incoming requests, e.g. "0.0.0.0:80" i.e. local port 80
 
## Help

### Documentation
You can find more documentation [here](docs/)

### Have a question?
File a GitHub [issue](https://github.com/stakater/ProxyInjector/issues), or send us an [email](mailto:hello@stakater.com).

### Talk to us on Slack
Join and talk to us on the #tools-proxyinjector channel for discussing the ProxyInjector

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://stakater-slack.herokuapp.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater.slack.com/messages/CFCP3MUR4/)

## License

Apache2 © [Stakater](http://stakater.com)

## About

The `ProxyInjector` is maintained by [Stakater][website]. Like it? Please let us know at <hello@stakater.com>

See [our other projects][community]
or contact us in case of professional services and queries on <hello@stakater.com>

  [website]: http://stakater.com/
  [community]: https://www.stakater.com/projects-overview.html

## Contributers

Stakater Team and the Open Source community! :trophy:
