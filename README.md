# yaml-jsonnet-bootstrap

This utility parses a YAML document stream and generates Jsonnet boilerplate to generate the same YAML.

yaml-jsonnet-bootstrap reads a YAML document stream from stdin and writes a valid Jsonnet program to stdout.

I wrote this because I couldn't find an online converter from a YAML stream to a JSON array.

## Usage

Given a YAML document like this:

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: myacct
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: myrole
rules:
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: mybinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: myrole
subjects:
- kind: ServiceAccount
  name: myacct
...
```

then if you run `yaml-jsonnet-bootstrap < /path/to/example.yml > /path/to/example.jsonnet`,
you will generate an unformatted Jsonnet file:

```jsonnet
local var0 = {
  "apiVersion": "v1",
  "kind": "ServiceAccount",
  "metadata": {
    "name": "myacct"
  }
};

local var1 = {
  "apiVersion": "rbac.authorization.k8s.io/v1",
  "kind": "Role",
  "metadata": {
    "name": "myrole"
  },
  "rules": [
    {
      "apiGroups": [
        ""
      ],
      "resources": [
        "pods"
      ],
      "verbs": [
        "get",
        "list",
        "watch"
      ]
    }
  ]
};

local var2 = {
  "apiVersion": "rbac.authorization.k8s.io/v1",
  "kind": "RoleBinding",
  "metadata": {
    "name": "mybinding"
  },
  "roleRef": {
    "apiGroup": "rbac.authorization.k8s.io",
    "kind": "Role",
    "name": "myrole"
  },
  "subjects": [
    {
      "kind": "ServiceAccount",
      "name": "myacct"
    }
  ]
};

{
Objects(conf):: [
var0,
var1,
var2,
],
}
```

Then, you can run jsonnetfmt on the file to get something nicely formatted like:

```
local var0 = {
  apiVersion: 'v1',
  kind: 'ServiceAccount',
  metadata: {
    name: 'myacct',
  },
};

local var1 = {
  apiVersion: 'rbac.authorization.k8s.io/v1',
  kind: 'Role',
  metadata: {
    name: 'myrole',
  },
  rules: [
    {
      apiGroups: [
        '',
      ],
      resources: [
        'pods',
      ],
      verbs: [
        'get',
        'list',
        'watch',
      ],
    },
  ],
};

local var2 = {
  apiVersion: 'rbac.authorization.k8s.io/v1',
  kind: 'RoleBinding',
  metadata: {
    name: 'mybinding',
  },
  roleRef: {
    apiGroup: 'rbac.authorization.k8s.io',
    kind: 'Role',
    name: 'myrole',
  },
  subjects: [
    {
      kind: 'ServiceAccount',
      name: 'myacct',
    },
  ],
};

{
  Objects(conf):: [
    var0,
    var1,
    var2,
  ],
}
```

So the normal pipeline is typically `yaml-jsonnet-bootstrap < file.yml | jsonnetfmt - > out.jsonnet`.
