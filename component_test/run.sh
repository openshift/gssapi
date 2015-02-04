#!/bin/bash -eu

# boot2docker doesn't seem to like /tmp so use the home direcotry for the build
BASE_DIR=$(cd .. && pwd)
export TEST_DIR="$HOME/tmp/$(uuidgen)"
mkdir -p $TEST_DIR
cp -R $BASE_DIR $TEST_DIR
DOCKER_DIR=$TEST_DIR/gssapi/component_test/docker

if [[ "$OSTYPE" == "darwin"* ]]; then
        DOCKER=docker
else
        DOCKER='sudo docker'
fi

function log() {
	>&2 /bin/echo "go-gssapi-test: $*"
}

function cleanup_containers() {
	log "Clean up running containers"
	running=`$DOCKER ps --all | grep 'go-gssapi-test' | awk '{print $1}'`
	if [[ "$running" != "" ]]; then
		echo $running | xargs $DOCKER stop >/dev/null
		echo $running | xargs $DOCKER rm >/dev/null
	fi
}

function cleanup() {
        set +e

        log "
        
kdc logs:"
        $DOCKER logs kdc 2>&1 

        log "
        
service logs:"
        if [[ "$SERVICE_LOG_FILTER" != "" ]]; then
                $DOCKER logs service 2>&1 | egrep -v "gssapi-sample:\t[0-9 /:]+ ACCESS "
        else
                $DOCKER logs service 2>&1 
        fi

	cleanup_containers

	log "Clean up build directory"
	rm -rf $TEST_DIR
}

function build_image() {
        comp=$1
        name=$2
        func=$3
        img="go-gssapi-test-${name}"
        image=$($DOCKER images --quiet ${img})

        if [[ "$REUSE_DOCKER_IMAGES" != "" && "$image" != "" ]]; then
                log "Reuse cached docker image ${img} ${image}"
        else
                log "Build docker image ${img}"
                if [[ "$func" != "" ]]; then
                        (${func})
                fi
                $DOCKER build \
                        --quiet \
                        --rm \
                        --tag=${img} \
                        $DOCKER_DIR/${comp}
        fi
}

function run_image() {
        comp=$1
        name=$2
        options=$3
        img="go-gssapi-test-${name}"
        log "Run docker image ${img}"
	options="${options} \
                --hostname=${comp} \
                --name=${comp} \
                --env SERVICE_NAME=${SERVICE_NAME} \
                --env USER_NAME=${USER_NAME} \
                --env USER_PASSWORD=${USER_PASSWORD} \
                --env REALM_NAME=${REALM_NAME} \
                --env DOMAIN_NAME=${DOMAIN_NAME}"
	$DOCKER run -P ${options} ${img}
}

function map_ports() {
	comp=$1
        port=$2
        COMP=`echo $comp | tr '[:lower:]' '[:upper:]'`
        if [[ "${OSTYPE}" == "darwin"* ]]; then
                b2d_ip=$(boot2docker ip)
                export ${COMP}_PORT_${port}_TCP_ADDR=${b2d_ip}
        else
                export ${COMP}_PORT_${port}_TCP_ADDR=127.0.0.1
        fi
        export ${COMP}_PORT_${port}_TCP_PORT=$(docker port ${comp} ${port} | cut -f2 -d ':')
}


# Cleanup
trap 'cleanup' INT TERM EXIT

cleanup_containers


# KDC
if [[ "${EXT_KDC_IP}" == "" ]]; then
        cat $DOCKER_DIR/kdc/krb5.conf.template \
                | sed -e "s/KDC_ADDRESS/0.0.0.0:88/g" \
                | sed -e "s/DOMAIN_NAME/${DOMAIN_NAME}/g" \
                | sed -e "s/REALM_NAME/${REALM_NAME}/g" \
                > $DOCKER_DIR/kdc/krb5.conf

        suffix=$(/bin/echo ${REALM_NAME} | shasum | cut -f1 -d ' ')
        build_image "kdc" "kdc-${suffix}" "" >/dev/null
        run_image "kdc" "kdc-${suffix}" "--detach" >/dev/null
        map_ports "kdc" 88
else
        export KDC_PORT_88_TCP_ADDR=${EXT_KDC_IP}
        export KDC_PORT_88_TCP_PORT=${EXT_KDC_PORT}
fi

while ! echo exit | nc $KDC_PORT_88_TCP_ADDR $KDC_PORT_88_TCP_PORT >/dev/null; do
        echo "Waiting for kdc to start"
        sleep 1
done

function keytab_from_kdc() {
        $DOCKER cp kdc:/etc/docker-kdc/krb5.keytab $DOCKER_DIR/service
}

function keytab_from_options() {
        cp ${KEYTAB_FILE} $DOCKER_DIR/service/krb5.keytab
}

if [[ "${EXT_KDC_IP}" == "" ]]; then
        DOCKER_KDC_OPTS='--link=kdc:kdc'
        KEYTAB_FUNCTION='keytab_from_kdc'
else
        DOCKER_KDC_OPTS="--env KDC_PORT_88_TCP_ADDR=${EXT_KDC_IP} \
                --env KDC_PORT_88_TCP_PORT=${EXT_KDC_PORT}"
        KEYTAB_FUNCTION='keytab_from_options'
fi

# GSSAPI service
log "Build and unit-test gssapi on host"
go test github.com/apcera/gssapi

build_image "service" "service" "$KEYTAB_FUNCTION" >/dev/null
run_image "service" \
        "service" \
        "--detach \
        $DOCKER_KDC_OPTS \
        --volume $TEST_DIR/gssapi:/opt/go/src/github.com/apcera/gssapi" >/dev/null
map_ports "service" 80

while ! echo exit | nc $SERVICE_PORT_80_TCP_ADDR $SERVICE_PORT_80_TCP_PORT >/dev/null; do
        echo "Waiting for service to start"
        sleep 1
done

# GSSAPI client
if [[ "$OSTYPE" != "darwin"* || "$CLIENT_IN_CONTAINER" != "" ]]; then
        build_image "client" "client" "" >/dev/null
        run_image "client" \
                "client" \
                "--link=service:service \
                $DOCKER_KDC_OPTS \
                --volume $TEST_DIR/gssapi:/opt/go/src/github.com/apcera/gssapi" \
                >/dev/null
else
        log "Run gssapi sample client on host"
        KRB5_CONFIG_TEMPLATE=${DOCKER_DIR}/client/krb5.conf.template \
                DOMAIN_NAME=${DOMAIN_NAME} \
                GSSAPI_PATH=/opt/local/lib/libgssapi_krb5.dylib \
                KRB5_CONFIG=${TEST_DIR}/krb5.conf \
                REALM_NAME=${REALM_NAME} \
                SERVICE_NAME=${SERVICE_NAME} \
                USER_NAME=${USER_NAME} \
                USER_PASSWORD=${USER_PASSWORD} \
                ${DOCKER_DIR}/client/entrypoint.sh
fi
echo "OK TEST PASSED"
