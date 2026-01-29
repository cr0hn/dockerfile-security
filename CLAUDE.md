# dockerfile-sec

Herramienta CLI que analiza Dockerfiles en busca de vulnerabilidades de seguridad y credenciales expuestas, usando reglas regex definidas en YAML.

**Autor:** Daniel Garcia (cr0hn@cr0hn.com)
**Licencia:** BSD

## Estado actual

Reescritura de Python a Go. El PRD completo está en `PRD.md`.

## Estructura objetivo (Go)

```
cmd/dockerfile-sec/main.go       # Entry point
internal/
  analyzer/                      # Motor de análisis (regex matching)
  rules/                         # Carga de reglas (embed, file, HTTP)
    embedded/                    # core.yaml, credentials.yaml, security.yaml, packages.yaml, configuration.yaml
  output/                        # Formateadores (tabla ASCII, JSON)
  ignore/                        # Sistema de reglas ignoradas
  config/                        # Configuración CLI
testdata/                        # Fixtures de tests
```

## Comandos

```bash
go build -o bin/dockerfile-sec ./cmd/dockerfile-sec/   # Build
go test ./...                                           # Tests
go test -cover ./...                                    # Cobertura
go vet ./...                                            # Análisis estático
```

## Interfaz CLI (mantener compatibilidad con versión Python)

```
dockerfile-sec [DOCKERFILE] [-F ignore-file] [-i rule-ids] [-r rules-file]
               [-R core|credentials|security|packages|configuration|all|none] [-o output.json] [-q] [-E]
```

- Sin archivo: lee de stdin
- `-E`: exit code 1 si hay issues (CI/CD)
- `-q`: modo silencioso
- `-R`: soporta categorías separadas por comas (e.g., `-R core,security`)
- Salida: tabla ASCII en terminal, JSON en pipe

## Reglas

**35 reglas** en 5 categorías:
- `core` (10 reglas): buenas prácticas
- `credentials` (11 reglas): secretos expuestos, incluyendo GitHub PAT, RSA keys, OpenAI, Stripe, Docker registry auth
- `security` (7 reglas): Docker socket mounts, permisos peligrosos, SUID/SGID, sudo, BuildKit secrets
- `packages` (4 reglas): apt-get, pip, npm, curl|bash
- `configuration` (3 reglas): --privileged, puertos peligrosos, STOPSIGNAL

Formato YAML con campos: `id`, `description`, `regex`, `reference`, `severity` (Low, Medium, High, Critical).

## Convenciones

- Go idiomático: interfaces, error handling explícito, `embed` para reglas
- Tests: table-driven, `testdata/`, golden files
- Sin CGO. Binario estático.
- Commits en nombre de Dani (cr0hn@cr0hn.com), no de Claude
