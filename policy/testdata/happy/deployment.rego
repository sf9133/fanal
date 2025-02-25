package testdata.xyz_100

import data.services

__rego_metadata__ := {
    "id": "XYZ-100",
    "title": "Bad Deployment",
    "version": "v1.0.0",
    "severity": "HIGH",
    "type": "Kubernetes Security Check",
}

__rego_input__ := {
    "selector": {
        "types": ["kubernetes"]
    },
    "combine": false,
}

deny[msg] {
  input.kind == "Deployment"
  services.ports[_] == 22
  msg := sprintf("deny %s", [input.metadata.name])
}
