### kubeadm-getter

this project uses TCP sockets to transfer files from a server machine
to a client machine.

the data transfer is encrypted using AES-256 with Galois/Counter Mode (GCM).
the encryption key is a Kubernetes bootstrap token that both the client
and server should know.

A handshake is performed to ensure that the client knows the same token
as the server.

the client can only request files from the server input folder and its children.
the list of files can be separated using `,`. wildcards are not allowed.

the connection automatically terminates after a period of time (TTL).

```
# server example:
sudo kubeadm-getter --listen --address=<server-ip> --port=11000 --ttl=240 \
--token=abcdef.1234567890abcdef --input-path=/etc/kubernetes

# client example:
sudo kubeadm-getter --address=<server-ip> --port=11000 --ttl=240 \
--token=abcdef.1234567890abcdef --output-path=/etc/kubernetes \
--files=pki/ca.crt,pki/ca.key,admin.conf
```

in case you don't have a token call:
```
kubeadm-getter --create-token
```

for the list of available command line flags see `--help`.
