# Dockerfile-sec

Dockerfile-sec is a simple but powerful rules-based checker for Dockerfiles.

## Install

```bash
> pip install dockerfile-sec 
```

## Quick start

Analyze a Dockerfile

```bash
> dockerfile-sec examples/Dockerfile-example
+----------+-------------------------------------------+----------+
| Rule Id  | Description                               | Severity |
+----------+-------------------------------------------+----------+
| core-002 | Missing USER sentence in dockerfile       | Medium   |
| core-003 | Posible text plain password in dockerfile | High     |
| core-005 | Recursive copy found                      | Medium   |
| core-006 | Use of COPY instead of ADD                | Low      |
| core-007 | Use image tag instead of SHA256 hash      | Medium   |
| cred-001 | Generic credential                        | Medium   |
+----------+-------------------------------------------+----------+  
```

## Using docker

```bash
> cat Dockerfile | docker run --rm -t cr0hn/dockerfile-sec  
```

    IMPORTANT: By using docker you can pass a rules file or a docker file as paramenter. You need to use a pipe or mount a volume

## Usage

### With remote rules

```bash
> dockerfile-sec -r http://127.0.0.1:9999/rules/credentials.yaml Dockerfile 
```

### With built-in rules

**All rules**

All rules are enabled by default:

```bash
> dockerfile-sec Dockerfile
```

**Core rules only**

https://github.com/cr0hn/dockerfile-security/blob/master/dockerfile_sec/rules/core.yaml

```bash
> dockerfile-sec -R core Dockerfile
```

**Credentials rules only**

https://github.com/cr0hn/dockerfile-security/blob/master/dockerfile_sec/rules/credentials.yaml

```bash
> dockerfile-sec -R credentials Dockerfile
```

**Disabling built-in rules**

```bash
> dockerfile-sec -R none Dockerfile
```

### With user defined rules

```bash
> dockerfile-sec -r my-rules.yaml Dockerfile
```

### Export results as json

```bash
> dockerfile-sec -o results.json Dockerfile 
```

### Quiet mode

Not writing anything in the console:

```bash
> dockerfile-sec -q -o results.json Dockerfile 
```


### Filtering false positives

**By ignore file**

Dockerfile-sec allows to ignore rules by using a file that contains the rules you want to ignore.

```bash
> dockerfile-sec -F ignore-rules.text Dockerfile 
```

Ignore file format contains the *IDs* of rules you want to ignore. **one ID per line**. Example:

```bash
> ls ignore-rules.text
core-001
core-007
```

**By using the cli**

You also can use cli to ignore specific *IDs*:

```bash
> dockerfile-sec -i core-001,core007 Dockerfile 
```

## Using as a pipeline

You also can use dockerfile-sec as UNIX pipeline.

Loading Dockerfile from stdin:

```bash
> cat Dockerfile | dockerfile-sec -i core-001,core007 
```

Exposing results via pipe:

```bash
> cat Dockerfile | dockerfile-sec -i core-001,core007 | jq 
```

## Output formats

### JSON Output format

```json

[
  {
    "description": "Missing USER sentence in dockerfile",
    "id": "core-002",
    "reference": "https://snyk.io/blog/10-docker-image-security-best-practices/",
    "severity": "Medium"
  }
]

```

## References

- https://snyk.io/blog/10-docker-image-security-best-practices/
- https://medium.com/microscaling-systems/dockerfile-security-tuneup-166f1cdafea1
- https://medium.com/@tariq.m.islam/container-deployments-a-lesson-in-deterministic-ops-a4a467b14a03
- https://spacelift.io/blog/docker-security
