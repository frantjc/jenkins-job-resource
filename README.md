# jenkins-job-resource

A Concourse resource for jobs on Jenkins.  Written in Go.

## Example

```yaml
resource_types:
  - name: jenkins-job-resource
    type: registry-image
    source:
      repository: logsquaredn/jenkins-job-resource
      tag: latest # TODO: version the resource for reproducible builds

resources:
  - name: my-jenkins-job
    type: jenkins-job-resource
    source:
      ...

jobs:
  - name: some-job
    plan:
      ...
      - put: my-jenkins-job
        params:
          cause: some-job in Concourse caused this build
          build_params:
            foo: bar

  - name: another-job
    plan:
      - get: my-jenkins-job
```

## Source Configuration

| Parameter   | Required            | Example                  | Description                                                            |
| ----------- | ------------------- | ------------------------ | ---------------------------------------------------------------------- |
| `url`       | yes                 | `https://my-jenkins.com` | the base url for a Jenkins deployment                                  |
| `job`       | yes                 |                          | the job name in the Jenkins deployment at the given url                |
| `token`     | for `put`           | see guide below (TODO)   | the `Authentication Token` for the given job                           |
| `username`  | unknown, but likely |                          | a username that can be used to authorize with the given Jenkins url    |
| `login`     | unknown, but likely |                          | either the password or an api token associated with the given username |

> _Note: if_ `login` _is set to a password, then surely_ `username` _must be set to a corresponding username for it to have any effect. If_ `login` _is set to an api token, it is unknown if_ `username` _need be set to a corresponding username. It is unknown if a Jenkins deployment can be configured such that both_ `username` _and_ `login` _are not required_

TODO: Jenkins configuration guide

## Behavior

### `check`

Produces new versions for all builds (after the last version) ordered by the build number.

### `in`

TODO: `skip_download`, make some metadata available via some file(s), more?

Gets the artifacts outputted by the most recent build

### `out`

| Parameter      | Required | Example    | Description                                                           |
| -------------- | -------- | ---------- | --------------------------------------------------------------------- |
| `cause`        | no       |            | the cause of the build being triggered                                |
| `description`  | no       |            | the description of the build being triggered                          |
| `build_params` | no       | any object | name-value pairs that will be passed to the build as build parameters |

Triggers a new build of the target job.
