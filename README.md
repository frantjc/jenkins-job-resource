# jenkins-job-resource

A Concourse resource for jobs on Jenkins.  Written in Go.

## Example

```yaml
resource_types:
- name: jenkins-job-resource
  type: registry-image
  source:
    repository: logsquaredn/jenkins-job-resource
    tag: latest

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
| `token`     | for `put`           | _see below (1)_          | the `Authentication Token` for the given job                           |
| `username`  | unknown, but likely |                          | a username that can be used to authorize with the given Jenkins url    |
| `login`     | unknown, but likely |                          | either the password or an api token associated with the given username |

> _Note: if_ `login` _is set, then surely_ `username` _must be set to a corresponding username for it to have any effect. It is unknown if a Jenkins deployment can be configured such that both_ `username` _and_ `login` _are not required_

![build-triggers](https://user-images.githubusercontent.com/39865011/100497098-2ccc3c80-3127-11eb-984b-ce09b1681ab1.png)

> _(1) the above can be found at_ `$YOUR_JENKINS` \> `$YOUR_JOB` \> Configure

## Behavior

### `check`

Produces new versions for all builds (after the last version) ordered by the build number

### `in`

| Parameter        | Required | Example              | Description                                                                                                               |
| ---------------- | -------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `accept_results` | no       | `["SUCCESS"]`        | array of acceptable results of the build. The step will fail if none match. Empty arrays are ignored                      |
| `regexp`         | no       | `[path/to/output.*]` | limits downloaded artifacts to only those that match one of the given patterns. Hidden files (ex: `.git`) will be ignored |
| `skip_download`  | no       |                      | whether or not to download any of the artifacts at all. Overrides `regexp`. Default `false`                               |

Optionally gets the specified artifacts outputted by the most recent build

Metadata is made available at `./.metadata/name`. ex:

```
$ cat ./my-jenkins-job/.metadata/result
SUCCESS
```

### `out`

| Parameter           | Required | Example               | Description                                                                                    |
| ------------------- | -------- | --------------------- | ---------------------------------------------------------------------------------------------- |
| `cause`             | no       |                       | the cause of the build being triggered. Default `caused by $ATC_EXTERNAL_URL/builds/$BUILD_ID` |
| `cause_file`        | no       | `path/to/cause`       | path to a file containing the cause of the build being triggered                               |
| `description`       | no       |                       | the description of the build. Default `build triggered by $ATC_EXTERNAL_URL/builds/$BUILD_ID`  |
| `description_file`  | no       | `path/to/description` | path to a file containing the description of the build                                         |
| `build_params`      | no       | any object            | name-value pairs that will be passed to the build as build parameters                          |

Triggers a new build of the target job and gets the result

## Development

### Prerequisites

* golang is *required* - version 1.11.x or above is required for go mod to work
* docker is *required* - version 20.10.x is tested; earlier versions may also work.
* go mod is used for dependency management of the golang packages.

### Tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing environment is consistent across any `docker` enabled platform. When the docker image builds, the test are run inside the docker container, on failure they will stop the build.

The tests can be ran for both the `ubuntu` and `alpine` images a number of ways:

* Against your own Jenkins job in some Jenkins deployment:

```sh
docker build \
  -t jenkins-job-resource \
  --target tests \
  -f dockerfiles/alpine/Dockerfile \
  --build-arg JENKINS_URL=http://example.jenkins.com \
  --build-arg JENKINS_JOB=my-jenkins-job \
  --build-arg JENKINS_JOB_TOKEN=my-jenkins-job-token \
  --build-arg JENKINS_USERNAME=my-username \
  --build-arg JENKINS_API_TOKEN=my-api-token \
  .

docker build \
  -t jenkins-job-resource \
  --target tests \
  -f dockerfiles/ubuntu/Dockerfile \
  --build-arg JENKINS_URL=http://example.jenkins.com \
  --build-arg JENKINS_JOB=my-jenkins-job \
  --build-arg JENKINS_JOB_TOKEN=my_jenkins_job_token \
  --build-arg JENKINS_USERNAME=my-username \
  --build-arg JENKINS_API_TOKEN=my_api_token \
  .
```

* Against a Jenkins automatically spun up in a container using `docker`:

```sh
bin/test
```

* Against a Jenkins automatically spun up in a container using `docker compose` or `docker-compose`:

```sh
docker compose up --build ## docker-compose up --build
```

### Contributing

Please make all pull requests to the `main` branch and ensure tests pass locally.
