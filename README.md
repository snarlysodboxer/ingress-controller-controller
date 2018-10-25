# ingress-controller-controller

#### Why
* In a microservices environment, it's common to distribute Kubernetes manifests across many Git repositories.
* It may be desired to manage one or a handful of `Ingress` objects, rather than one per repo or microservice.
* It might not be desired to manage them in a centralized location, since among other things, this requires tight coordination between the centralized and the distributed.
* ingress-controller-controller allows you to manage one or many `Ingress`s from n number of `Service` objects, using Annotations.
* Used it in conjunction with an Ingress Controller and an [external-dns](https://github.com/kubernetes-incubator/external-dns) controller, Load Balancers and relevant DNS records can be fully automated using simple annotations on `Service`s, GitOps'd in diverse locations.

#### How it works
* Watch `Service`s with a certain label via the Kubernetes API, and at each event loop:
  * Use the API to find the list of `Service`s that have the right annotation
  * Calculate a list of desired `Ingress`s from the annotations
  * Use the API to find the list of `Ingress`s that have the right annotation
  * Delete any orphaned `Ingress`s
    * (ingress-controller-controller annotates `Ingress`s which it created)
  * Apply the desired `Ingress`s to the API

#### Example
See [examples](examples)

#### Roadmap
* Add rate limiter
* Config via flags
* Fix flaky tests (object comparison sometimes fails when the objects are arranged in different orders)
* Docs
* Test for conflicting paths
* Validations for Annotations
* DRY up tests
* E2E tests
* Metrics
