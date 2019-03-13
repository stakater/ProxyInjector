# ![](assets/web/proxyinjector-round-100px.png) Proxy Injector
A Kubernetes controller to inject an authentication proxy container to relevant pods

[![Get started with Stakater](https://stakater.github.io/README/stakater-github-banner.png)](http://stakater.com/?utm_source=ProxyInjector&utm_medium=github)

## Problem Statement

We want to automatically inject a [keycloak gatekeeper](https://github.com/keycloak/keycloak-gatekeeper) container in a pod, for any deployment that requires to connect
 to keycloak, instead of manually adding a sidecar container with each deployment 

## Solution

This controller will continuously watch deployments in specific or all namespaces, and automatically add a sidecar container for [keycloak gatekeeper](https://github.com/keycloak/keycloak-gatekeeper). Configuration for the keycloak gatekeeper is done through annotations of the respective deployment or with ConfigMap of the ProxyInjector.


## Usage

The following quickstart let's you set up ProxyInjector:

1. Deploy the controller by running the following command:

    For Kubernetes Cluster
   ```bash
   kubectl apply -f https://raw.githubusercontent.com/stakater/ProxyInjector/master/deployments/kubernetes/proxyinjector.yaml -n default

2. When deploying any application that needs Keycloak authentication, add the following annotations to the deployment.
  
    | Key                                        | Description                                                                                                                                       |
    |--------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
    | authproxy.stakater.com/enabled             | Enables Keycloak gatekeeper configuration                                                                                                         |
    | authproxy.stakater.com/source-service-name | Name of service that needs to be reconfigured to connect to the proxy                                                                             |
    | authproxy.stakater.com/target-port         | (default=80) the port number that should be changed for the target port of the service, i.e. the port where keycloak gatekeeper will be listening |
    
    The following arguments can either be added to the proxy injector `config.yaml` in the ConfigMap, or as annotations
    on the target deployments with a `authproxy.stakater.com/` prefix.

    | Key              | Description                                                               |
    |------------------|---------------------------------------------------------------------------|
    | listen           | the interface the proxy should be listening on                            |
    | upstream-url     | url for the upstream endpoint you wish to proxy                           |
    | resources        | list of resources to proxy uri, methods, roles                            |
    | client-id        | client id used to authenticate to the oauth service                       |
    | client-secret    | client secret used to authenticate to the oauth service                   |
    | gatekeeper-image | Keycloak Gatekeeper image e.g. `keycloak/keycloak-gatekeeper:4.6.0.Final` |

The rest of the available options can be found at the [Keycloak Gatekeeper documentation](https://github.com/keycloak/keycloak-gatekeeper#keycloak-proxy)
 
## Help

### Have a question?
File a GitHub [issue](https://github.com/stakater/ProxyInjector/issues), or send us an [email](mailto:hello@stakater.com).

### Talk to us on Slack
Join and talk to us on the #tools-proxyinjector channel for discussing the Ingress Monitor Controller

[![Join Slack](https://stakater.github.io/README/stakater-join-slack-btn.png)](https://stakater-slack.herokuapp.com/)
[![Chat](https://stakater.github.io/README/stakater-chat-btn.png)](https://stakater.slack.com/messages/)

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
