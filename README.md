# certdump

This simple little utility dumps out some useful information from X509 certificates.

## Example on a UCP node

If you want to dump out the cert information from a UCP deployment,
you can use the following:

```bash
docker run --rm \
    -v ucp-auth-api-certs:/ucp-auth-api-certs \
    -v ucp-auth-store-certs:/ucp-auth-store-certs \
    -v ucp-auth-worker-certs:/ucp-auth-worker-certs \
    -v ucp-auth-worker-data:/ucp-auth-worker-data \
    -v ucp-client-root-ca:/ucp-client-root-ca \
    -v ucp-cluster-root-ca:/ucp-cluster-root-ca \
    -v ucp-controller-client-certs:/ucp-controller-client-certs \
    -v ucp-controller-server-certs:/ucp-controller-server-certs \
    -v ucp-kv-certs:/ucp-kv-certs \
    -v ucp-node-certs:/ucp-node-certs \
	-v /var/lib/docker/swarm/certificates:/swarm \
    dhiltgen/certdump
```

## Common failure scenarios to look out for

* Due to clock skew or just generally having clocks set wrong, the cert timestamps may be out-of-whack.
* Regenerated root CAs.  If the root CAs have been regenerated, and old certs are still lying around they
  may casually look OK, but not be trusted.  You can cross-reference the Issuer ID's to make sure they match.
