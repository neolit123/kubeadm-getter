### tokenized-getter

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
tokenized-getter --listen --address=<server-ip> --port=11000 --ttl=240 \
--token=abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 \
--input-path=/some-path

# client example:
tokenized-getter --address=<server-ip> --port=11000 --ttl=240 \
--token=abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 \
--output-path=/some-path --files=ca.crt,ca.key
```

in case you don't have a token call:
```
tokenized-getter --create-token
```

for the list of available command line flags see `--help`.
