# Pretty default backend
![GitHub last commit](https://img.shields.io/github/last-commit/bsord/pretty-default-backend.svg)
![Build and Publish](https://github.com/bsord/pretty-default-backend/workflows/Build%20and%20Publish/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsord/pretty-default-backend)](https://goreportcard.com/report/github.com/bsord/pretty-default-backend)
![License](https://img.shields.io/github/license/bsord/pretty-default-backend.svg?style=flat)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-green.svg)

An easily customized pretty default backend replacement for Kubernetes Nginx ingress controller with a neutral default configuration.

![pretty-default-backend](https://raw.githubusercontent.com/bsord/pretty-default-backend/master/cover.png)

## Getting Started
A default installation will deploy a single instance of pretty-default-backend to the same namespace of the ingress that will utilize it.

### Requirements
* [Kubernetes](https://microk8s.io/docs)
* [Helm v3](https://helm.sh/docs/intro/install/)

### Installation
Replace values **[namespace-of-ingress]**, and **[ingress-name]** in the commands below according to your environment
1. Add Helm Repository
    ```sh
    helm add repo bsord https://h.cfhr.io/bsord/charts
    ```
2. Install the helm chart (to same namespace as ingress)
    ```sh
    helm install bsord/pretty-default-backend --set bgColor="#443322" --set brandingText="YourBrandingText" bsord/pretty-default-backend -n [namespace-of-ingress]
    ```
3. Configure an ingress to use pretty-default-backend with one of the following two options:

    * Patch the Ingress Directly
        ```sh
        kubectl annotate ingress [ingress-name] -n [namespace-of-ingress] \
        nginx.ingress.kubernetes.io/default-backend=pretty-default-backend --overwrite
    
        kubectl annotate ingress [ingress-name] -n [namespace-of-ingress] \
        nginx.ingress.kubernetes.io/custom-http-errors="404,503" --overwrite
        ```
    * Pass as parameters to standard templated Helm v3 chart
        ```sh
        Helm install [release-name] \
        --set "ingress.annotations.nginx\\.ingress\\.kubernetes\\.io/default-backend=pretty-default-backend" \
        --set  "ingress.annotations.nginx\\.ingress\\.kubernetes\\.io/custom-http-errors=404\\,503\\,501" \
        [repo/chart-name]
        ```
        _Please note you must escape any special characters including commas and periods with a `\` backslash when passing complex strings to `--set`._
        
### Verifying
Browse to a location that does not exist and would trigger a 404 on the ingress you annotated above.

### Parameters
The parameters below can be passed using `--set KEY=VALUE` in the helm install/upgrade command above.

| Key | Value | Default |
| ------------- | ------------- | ------------- |
| `bgColor` | Background color of the page in hex value | #334455 |
| `brandingText` | Branding text at bottom of error box | BrandingText(2020) |

# Todo:
- [x] Write a functional ReadMe
- [ ] Fix workflow so it only triggers on succesful merge
- [ ] Use seperate writer stream before sending response (prevent broken responses)
- [ ] Add support for rich html variable input
