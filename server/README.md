## server project structure

```shell
├── api
│   └── v1
├── config
├── core
├── docs
├── global
├── initialize
│   └── internal
├── middleware
├── model
│   ├── request
│   └── response
├── packfile
├── resource
│   ├── excel
│   ├── page
│   └── template
├── router
├── service
├── source
└── utils
    ├── timer
    └── upload
```

| Folder       | Name                    | Description                        |
| ------------ | ----------------------- | --------------------------- |
| `api`        | api layer               | api layer |
| `--v1`       | v1 endpoints            | v1 endpoints                |
| `config`     | config package          | config structs matching config.yaml |
| `core`       | core files              | initialization of core components (zap, viper, server) |
| `docs`       | swagger docs            | swagger documentation |
| `global`     | globals                 | global objects |
| `initialize` | initialization          | initialization of router, redis, gorm, validator, timer |
| `--internal` | internal init helpers   | gorm logger customization; functions here are called only by the `initialize` layer |
| `middleware` | middleware layer        | gin middleware code |
| `model`      | model layer             | models mapping to database tables |
| `--request`  | request structs         | data received from the frontend |
| `--response` | response structs        | data returned to the frontend |
| `packfile`   | static asset packaging  | static asset packaging |
| `resource`   | static resources        | static files |
| `--excel` | excel default path      | default path for excel import/export |
| `--page` | form generator          | form generator packaged dist |
| `--template` | templates               | templates used by the code generator |
| `router`     | router layer            | router layer |
| `service`    | service layer           | business logic |
| `source`     | source layer            | functions that seed initial data |
| `utils`      | utilities               | utility helpers |
| `--timer` | timer                   | timer interface wrappers |
| `--upload`   | oss                     | oss interface wrappers |

