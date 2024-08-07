# Server

The server package is responsible for handling webserver requests for map tiles and various JSON endpoints describing the configured server. Example config:

```toml
[webserver]
port = ":9090"              # set something different than default ":8080"
ssl_cert = "fullchain.pem"  # ssl cert for serving by https
ssl_key = "privkey.pem"     # ssl key for serving by https


[webserver.headers]
Access-Control-Allow-Origin = "*"
```

### Config properties

- `port` (string): [Optional] Port and bind string. For example ":9090" or "127.0.0.1:9090". Defaults to ":8080"
- `hostname` (string): [Optional] The hostname to use in the various JSON endpoints. This is useful if tegola is behind a proxy and can't read the API consumer's request host directly.
- `uri_prefix` (string): [Optional] A prefix to add to all API routes. This is useful when tegola is behind a proxy (i.e. example.com/tegola). The prexfix will be added to all URLs included in the capabilities endpoint responses.
- `ssl_cert` (string): [Optional, unless ssl_key provided] Path to a certificate file for serving through HTTPS
- `ssl_key` (string): [Optional, unless ssl_cert provided] Path to a private key file for serving through HTTPS

## Local development of the embedded viewer

Tegola's built in viewer code is stored in the `ui/` directory. To build the ui `npm` must be installed. Once `npm` is installed the following command can be run from the repository root to generate a .go file for inclusion in the tegola binary:

```
go generate ./server
```

## Disabling the viewer

The viewer can be excluded during building by using the build flag `noViewer`. For example, building tegola from the `cmd/tegola` directory:

```bash
go build -tags "noViewer"
```
