# Pretty default backend
![GitHub last commit](https://img.shields.io/github/last-commit/bsord/pretty-default-backend.svg)
![Build and Publish](https://github.com/bsord/pretty-default-backend/workflows/Build%20and%20Publish/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/bsord/pretty-default-backend)](https://goreportcard.com/report/github.com/bsord/pretty-default-backend)
![License](https://img.shields.io/github/license/bsord/pretty-default-backend.svg?style=flat)

An easily customized pretty default backend replacement for kubernetes nginx ingress controller with a neutral default configuration.

![Docker+Node](https://raw.githubusercontent.com/bsord/pretty-default-backend/master/pretty-default-backend.png)

## Requirements
* Kubernetes with nginx ingress installed
* Helm v3 installed and configured

## Getting Started
You can use either installation method below, or both together. However, if custom-http-errors are defined (see configuration) for both the global installation and an ingress specific one, the global backend will be used for any conflicting error code definitions.

### Installation
Replace values [namespace-of-ingress], and [ingress-name] in the commands below according to your environment
1. Add Helm Repository
```sh
helm add repo bsord https://h.cfhr.io/bsord/charts
```
2. Install the helm chart (to same namespace as ingress)
```sh
helm install bsord/pretty-default-backend --set bgColor="#443322" --set brandingText="YourBrandingText" ./chart -n [namespace-of-ingress]
```
3. Patch Annotations on existing ingress
```sh
kubectl annotate ingress [ingress-name] -n [namespace-of-ingress] ingress.annotations.nginx.ingress.kubernetes.io/default-backend pretty-default-backend
kubectl annotate ingress [ingress-name] -n [namespace-of-ingress] ingress.annotations.nginx.ingress.kubernetes.io/custom-http-errors "404,503"
```

### Parameters
The parameters below can be passed using `--set KEY=VALUE` in the helm install/upgrade command above.

| Key | Value | Default |
| ------------- | ------------- | ------------- |
| `bgColor` | Background color of the page in hex value | #334455 |
| `brandingText` | Branding text at bottom of error box | BrandingText(2020) |
# Todo:
:all-the-things:
Write a how to
fix workflow so it only triggers on succesful merge
Update this readme.
Use seperate writer stream before sending response (prevent broken responses)
