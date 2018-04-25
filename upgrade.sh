#!/bin/bash

# COMMAND line parameters
get_existing_servers() {
    EXISTING_SERVERS=$(docker ps | grep webapp | tr -s " "| cut -d" " -f1 | xargs -I{} docker inspect {} --format  '{{.Name}}'| tr -d '/')
}

# helper to get container names
get_existing_servers
ORIG_EXISTING_SERVERS=$EXISTING_SERVERS
echo "Found existing servers: ${ORIG_EXISTING_SERVERS}"

# build new instance of webapp
docker-compose build webapp

NUM_SERVERS=${#ORIG_EXISTING_SERVERS[@]}

for server in $ORIG_EXISTING_SERVERS; do
    echo "Upgrading ${server}"

    # scale up to have one new node
    echo "Scaling up"
    docker-compose up -d --scale webapp=$((NUM_SERVERS + 1))

    # remove old container
    docker stop $server
    docker rm $server

    # write new config file with new host
    servers_config=""
    get_existing_servers
    for new_server in $EXISTING_SERVERS; do
        servers_config="server $new_server:3000;${servers_config}";
    done
    echo "Generated config: ${servers_config}"

    sed -e  "s/{{SERVERS}}/${servers_config}/" server/load_balancer/nginx.conf.template> server/load_balancer/nginx.conf
    cat server/load_balancer/nginx.conf

    echo "Restarting lb"
    # reload lb
    docker exec letstalk_lb_1 nginx -s reload
done
