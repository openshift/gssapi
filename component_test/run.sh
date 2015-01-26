#!/bin/bash -eu

#TODO: Add a -c flag to clear out old images, or rebuild by default, but have a flag to reuse

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

function print() {
	>&2 echo "apcera/go-gssapi-test: $*"
}

function cleanup_containers() {
	print "Clean up running containers"
	running=`$DOCKER ps --all | grep 'apcera/go-gssapi-test' | awk '{print $1}'`
	if [[ "$running" != "" ]]; then
		echo $running | xargs $DOCKER stop >/dev/null
		echo $running | xargs $DOCKER rm >/dev/null
	fi
}

function cleanup() {
        set +e
        print "Service logs:"
        $DOCKER logs service 2>&1

	cleanup_containers

	print "Clean up build directory"
	rm -rf $TEST_DIR
}

function build_image() {
	comp=$1
	check=$2
	func=$3
	img="apcera/go-gssapi-test-${comp}"
	image=$($DOCKER images --quiet ${img})

	if [[ "$check" == "" || "$image" == "" ]]; then
		print "Build docker image ${img}"
		if [[ "$func" != "" ]]; then
			(${func})
		fi
		$DOCKER build \
			--quiet \
			--rm \
			--tag=${img} \
			$DOCKER_DIR/${comp}
	else
		print "Reuse cached docker image ${img} ${image}"
	fi
}

function run_image() {
	comp=$1
	options=$2
	img="apcera/go-gssapi-test-${comp}"
	print "Run docker image ${img}"
	options="${options} --hostname=${comp} --name=${comp}"
	$DOCKER run -P ${options} ${img}
}


# Cleanup
trap 'cleanup' INT TERM EXIT
cleanup_containers


# KDC
build_image "kdc" "unless exists" "" >/dev/null
run_image "kdc" "--detach" >/dev/null


# GSSAPI service
print "Build and unit-test gssapi on host"
go test github.com/apcera/gssapi
function copy_keytab() {
	$DOCKER cp kdc:/etc/docker-kdc/krb5.keytab $DOCKER_DIR/service
}
build_image "service" "unless exists" "copy_keytab" >/dev/null
run_image "service" \
        "--detach \
        --link=kdc:kdc\
        --volume $TEST_DIR/gssapi:/opt/go/src/github.com/apcera/gssapi" >/dev/null


# GSSAPI client
if [[ "${OSTYPE}" == "darwin"* ]]; then
        print "Build gssapi sample client on host"
        B2D_IP=$(boot2docker ip)
        export SERVICE_PORT_80_TCP_ADDR=${B2D_IP}
        export SERVICE_PORT_80_TCP_PORT=$(docker port service 80/tcp | cut -f2 -d ':')
        export KDC_PORT_88_TCP_ADDR=${B2D_IP}
        export KDC_PORT_88_TCP_PORT=$(docker port kdc 88/tcp | cut -f2 -d ':')
        export KRB5_CONFIG_TEMPLATE=${DOCKER_DIR}/client/krb5.conf.template
        export KRB5_CONFIG=${TEST_DIR}/krb5.conf
        export GSSAPI_PATH=/opt/local/lib/libgssapi_krb5.dylib

        ${DOCKER_DIR}/client/entrypoint.sh
else
        build_image "client" "unless exists" "" >/dev/null
        #set +e
        run_image "client" \
                "--link=service:service \
                --link=kdc:kdc \
                --volume $TEST_DIR/gssapi:/opt/go/src/github.com/apcera/gssapi" \
                >/dev/null
        #if [[ "$?" != "0" ]]; then
        #        print "Client failed, see service logs below:\n\n"
        #        $DOCKER logs service 2>&1
        #        exit 1
        #fi
        #set -e
fi
echo "TEST PASSED"
