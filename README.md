# sleuth-client-kube

A service intended to be run in kubernetes that listens for new replica sets in all namespaces, extracts deployment information from the annotations on the replica sets, and then published the deployment information to Sleuth.

## References
- https://github.com/sleuth-io/sleuth-client
- https://github.com/MatthewDolan/sleuth-client-go
- https://help.sleuth.io/sleuth-api#rest-api-details

