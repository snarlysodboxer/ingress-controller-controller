# ingress-controller-controller

#### Roadmap
* Watch `service`s, at each event loop:
  * List all services: `getAllServices`
  * Find the `service`s that have the right annotation: `getAnnotatedServices`
  * Calculate the current `ingress`s: `newIngressList`

  * List all `ingress`s: `getAllIngresses`
  * Find the `ingress`s that have the right annotation: `getAnnotatedIngresses`
  * Delete any orphaned `ingress`s: (use SDK)

  * Apply the `ingress`s: (use SDK)

