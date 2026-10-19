# elasticsearch-data-loader
A tool that Loads data from MongoDB into ElasticSearch

## Requirements
---------------
You must download and make the [alpha key service](https://github.com/companieshouse/alpha-key-service). Running the start script locally should start it on port 4025.

In order to run this service you need:

* alpha-key-service running
* Mongo DB instance
* Elastic Search 7+

## Getting Started
-------------------
When loading indices for development environments, the number of replica sets must be set as appropriate.
This can be tweaked in `config/search_scheme.json` but **do not** check in this change.

When the alpha key service is running then run the following command to drop the index and repopulate.

## Examples 
```bash
./run-elastic-search.sh -s company -e http://enva.es.ch.gov.uk:9400 -i alpha_search -m chs-pp-mes-sl2.ch.gov.uk:27019 -u admin -p admin -a http://chs-alphakey-pp.internal.ch -c false
```

## Docker / ECS

`ecs-image-build/docker_start.sh` maps environment variables to loader flags:

- `SEARCH` -> `-s` (default: `company`)
- `INDEX` -> `-i` (required)
- `ES_URL` -> `-e` (required)
- `MONGO_URL` -> `-m` (required)
- `ALPHAKEY_URL` -> `-a` (required)
- `USERNAME` -> `-u` (optional, must be set with `PASSWORD`)
- `PASSWORD` -> `-p` (optional, must be set with `USERNAME`)
- `CREATE_MAPPING` -> `-c` (default: `false`, valid: `true|false`)
- `COMPANY_LIMIT` -> `-l` (default: `0`)

Example run:

```bash
docker run --rm \
  -e SEARCH=company \
  -e INDEX=alpha_search \
  -e ES_URL=http://enva.es.ch.gov.uk:9400 \
  -e MONGO_URL=chs-pp-mes-sl2.ch.gov.uk:27019 \
  -e ALPHAKEY_URL=http://chs-alphakey-pp.internal.ch \
  -e CREATE_MAPPING=false \
  <your-image>:<tag>
```

## Running in AWS ECS

Use the ECS service to run this as a one-off batch load.

1. Set the ECS service task count to `1` to start a run.
2. Monitor task logs in CloudWatch until the load completes successfully.
3. After a successful run, the service will automatically scale back (roll back) to task count `0`.

### Staging permissions

For the staging environment, a member of the FESS team must start the run because they have the required permissions.

