# elasticsearch-data-loader
A tool that Loads data from MongoDB into ElasticSearch.

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
./run-elastic-search -s company -e enva.es.ch.gov.uk:9400 -i alpha_search -m chs-pp-mes-sl2.ch.gov.uk:27019 -u admin -p admin -a http://chs-alphakey-pp.internal.ch -c false
```

* Anything to be executed should be executed from the project root — ie this directory.

* The script `run-elastic-search` sets up an index with the correct settings and mappings required for
search using the config/all_scheme.json file. It then calls the relevant go scripts (companybindex) which will copy and transforms data from mongo db to ElasticSearch.

* `run-elastic-search` will ask for several parameters, to view these use the help parameter `-h`


