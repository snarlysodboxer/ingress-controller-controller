# ingress-controller-controller

#### Why
* TODO

#### How it works
* Watch `Service`s and at each event loop:
  * List all `Service`s, find the list that have the right annotation
  * Calculate a list of `Ingress`s from the annotations
  * List all `Ingress`s from the API, find the list that have the right annotation
  * Find and delete any orphaned `Ingress`s
    * (ingress-controller-controller annotates `Ingress`s which it created)
  * Apply the calculated `Ingress`s to the API

#### Roadmap
* Config via flags
* Docs
* DRY up tests
* E2E tests
