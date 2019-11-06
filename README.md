# 3scale OPA

[Open policy agent](https://www.openpolicyagent.org) built for 3scale.

Not much to see here yet. Just getting started.

## Instructions

- Compile:

```bash
go build
```

- Start the OPA server with the rules in `example.rego`:

```bash
./3scale-opa run --server --set=plugins.envoy_ext_authz_grpc.addr=:9191 --set=plugins.envoy_ext_authz_grpc.query=data.envoy.authz.allow --set=decision_logs.console=true --ignore=.* example.rego
```

- Start Envoy with the given config file that authorizes using the server above:

```bash
envoy -c envoy_config.yaml
```

