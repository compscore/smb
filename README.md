# SMB

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: smb
    tag: latest
  credentials:
    username:
    password:
  target:
  command:
  expectedOutput:
  weight:
  options:
    domain:
    share:
    exists:
    match:
    substring_match:
    regex_match:
    sha256:
    md5:
    sha1:
```

## Parameters

|     parameter     |            path            |   type   | default  | required | description                                           |
| :---------------: | :------------------------: | :------: | :------: | :------: | :---------------------------------------------------- |
|      `name`       |          `.name`           | `string` |   `""`   |  `true`  | `name of check (must be unique)`                      |
|       `org`       |       `.release.org`       | `string` |   `""`   |  `true`  | `organization that check repository belongs to`       |
|      `repo`       |      `.release.repo`       | `string` |   `""`   |  `true`  | `repository of the check`                             |
|       `tag`       |       `.release.tag`       | `string` | `latest` | `false`  | `tagged version of check`                             |
|    `username`     |  `.credentials.username`   | `string` |   `""`   | `false`  | `smb username`                                        |
|    `password`     |  `.credentials.password`   | `string` |   `""`   | `false`  | `default smb password`                                |
|     `target`      |         `.target`          | `string` |   `""`   |  `true`  | `network target for smb server`                       |
|     `command`     |         `.command`         | `string` |   `""`   | `false`  | `file tp check against expectedOutput`                |
| `expectedOutput`  |     `.expectedOutput`      | `string` |   `""`   | `false`  | `expected output for check to measured against`       |
|     `weight`      |         `.weight`          |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`        |
|   `status_code`   |   `.options.status_code`   |  `int`   |   `0`    | `false`  | `check status_code of the response`                   |
|      `match`      |      `.options.match`      |  `bool`  | `false`  | `false`  | `check contents of targeted file are exact match`     |
| `substring_match` | `.options.substring_match` |  `bool`  | `false`  | `false`  | `check contents of targeted file are substring match` |
|   `regex_match`   |   `.options.regex_match`   |  `bool`  | `false`  | `false`  | `check contents of targeted file are regex match`     |

## Examples

```yaml
- name: host_a-smb
  release:
    org: compscore
    repo: smb
    tag: latest
  credentials:
    username: Administrator
    password: changeme
  target: 10.{{ .Team }}.1.1:445
  command: readme.txt
  weight: 2
  options:
    domain: example.local
    share: C$
    exists:
```

```yaml
- name: host_a-smb
  release:
    org: compscore
    repo: smb
    tag: latest
  credentials:
    username: Administrator
    password: changeme
  target: 10.{{ .Team }}.1.1:445
  expectedOutput: ^According to all known laws of aviation
  weight: 1
  command: bee_movie_script.txt
  options:
    domain: example.local
    share: C$
    regex_match:
```
