#!/bin/bash -e

while getopts ":s:i:e:m:u:p:a:c:l:h" arg
do
  case "$arg" in
    s)
      search=$OPTARG         # 'officer'
    ;;
    i)
      index=$OPTARG          # 'officer-mapping'
    ;;
    e)
      es_url=$OPTARG         # 'chs-pp-es1.ch.gov.uk:9200'
    ;;
    m)
      mongo_url=$OPTARG      # 'chs-pp-mes-sl2.ch.gov.uk:27019'
    ;;
    u)
      username=$OPTARG       # 'admin'
    ;;
    p)
      password=$OPTARG       # 'admin'
    ;;
    a)
      alphakey_url=$OPTARG   # 'http://chs-alphakey-pp.internal.ch'
    ;;
    c)
      create_mapping=$OPTARG # 'true' or 'false'
    ;;
    l)
      company_limit=$OPTARG  # limit number of companies to transfer
    ;;
    h)
      echo ""
      echo "Elastic Loadinator Script"
      echo ""
      echo " * Empties elastic search index of data"
      echo " * Creates new mapping to elastic search index"
      echo " * Calls the relevant bindex to load in data to elastic search index"
      echo ""
      echo "Options are:"
      echo ""
      echo "      OPTION    ENV-VAR        DESCRIPTION                              EXAMPLE ('' does NOT indicate default value)"
      echo "        -s      search         The search type you intend to write to.  officer or disqualified-officer or company"
      echo "        -e      es_url         The elastic search url.                  chs-pp-es1.ch.gov.uk:9200"
      echo "        -i      index          The name of the elastic search index.    test-company"
      echo "        -m      mongo_url      The mongo db url.                        chs-pp-mes-sl2.ch.gov.uk:27019"
      echo "        -u      username       Username for authentication.             admin"
      echo "        -p      password       Password for authentication.             admin"
      echo "        -a      alphakey_url   The alphakey service url.                http://chs-alphakey-pp.internal.ch"
      echo "        -c      create_mapping Boolean flag to create mapping or not.   true or false (defaults to false)"
      echo "        -l      company_limit  Limit the number of companies to load.   100 (0 = no limit, defaults to 0)"
      echo ""
      exit 0
    ;;
    \?)
      echo "ERROR: Unknown option $OPTARG"
    ;;
  esac
done

search=${search:?ERROR: var not set [-s search]}
index=${index:?ERROR: var not set [-i index]}
es_url=${es_url:?ERROR: var not set [-e es_url]}
mongo_url=${mongo_url:?ERROR: var not set [-m mongo_url]}
alphakey_url=${alphakey_url:?ERROR: var not set [-a alphakey_url]}
create_mapping=${create_mapping:-false}
company_limit=${company_limit:-0}

full_es_url="$es_url/"$index
echo "            search: $search"
echo "elastic search url: $full_es_url"
echo "      mongo db url: $mongo_url"
echo " mongo db username: $username"
echo " mongo db password: $password"
echo "      alphakey url: $alphakey_url"
echo "    company limit: $company_limit"

# Check load type
if [ $search = "company" ]
then
    bindex="./companybindex/companybindex"
    scheme="search_scheme.json"
else
    echo "Incorrect search - use company"
    echo "Use -h for further options"
    exit 1
fi

echo "-----------------------------------"
echo "STEP 1: Delete existing index if -c flag set to true"
if [ $create_mapping = "true" ]
then
    echo "DELETING INDEX $full_es_url"
    delete_index="curl -XDELETE $full_es_url"
    echo $delete_index
    delete_index_response=`$delete_index`
    echo $delete_index_response
else
    echo "NOT DELETING INDEX"
fi

echo "-----------------------------------"
echo "STEP 2: Create index with new mapping if -c flag set to true"
if [ $create_mapping = "true" ]
then
    echo "CREATING INDEX WITH NEW MAPPING $full_es_url"
    curl -XPUT -H "Content-Type: application/json" $full_es_url -d@./config/$scheme
else
    echo "NOT CREATING INDEX WITH NEW MAPPING"
fi

# Check for authentication fields and build full mongo url
if [ -z "$username" ] && [ -z "$password" ]
then
    full_mongo_url="$mongo_url"
    echo "mongo url is: $full_mongo_url"
elif [ -z "$username" ]
then
    echo "Missing username!"
    exit 1
elif [ -z "$password" ]
then
    echo "Missing password!"
    exit 1
else
    echo "username: $username"
    echo "password: $password"

    full_mongo_url="$username:$password@$mongo_url"
    echo "mongo url is $full_mongo_url"
fi

echo "bindex: $bindex"

echo "-----------------------------------"
echo "STEP 3: Start $search load"
upload="$bindex -mongo-url=$full_mongo_url -es-dest-url=$es_url -es-dest-type=alpha_search -alphakey-url=$alphakey_url -es-dest-index=$index -company-limit=$company_limit"
echo $upload
exec $upload
