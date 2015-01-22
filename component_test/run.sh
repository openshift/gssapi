#!/bin/bash -e

function print() {
	echo "apcera/go-gssapi-test: $*"
}

function cleanup_containers() {
	print "Clean up running containers"
	running=`sudo docker ps --all | grep 'apcera/go-gssapi-test' | awk '{print $1}'`
	if [[ "$running" != "" ]]; then
		echo $running | xargs sudo docker stop >/dev/null
		echo $running | xargs sudo docker rm >/dev/null
	fi
}

function cleanup() {
	cleanup_containers

	print "Clean up build directory"
	rm -rf $TMP_DIR
}

function build_image() {
	comp=$1
	check=$2
	func=$3
	img="apcera/go-gssapi-test-${comp}"
	image=$(sudo docker images --quiet ${img})

	if [[ "$check" == "" || "$image" == "" ]]; then
		print "Build ${img}"
		cp -R ./docker/${1} $TMP_DIR
		if [[ "$func" != "" ]]; then
			(${func})
		fi
		sudo docker build \
			--quiet \
			--rm \
			--tag=${img} \
			$TMP_DIR/${1} >/dev/null
	else
		print "Reuse ${img} cached image ${image}"
	fi
}

function run_image() {
	comp=$1
	options=$2
	img="apcera/go-gssapi-test-${comp}"
	print "Run ${img}"
	options="${options} --hostname=${comp} --name=${comp}"
	sudo docker run ${options} ${img}
}

TMP_DIR="/tmp/$(uuidgen)"
mkdir -p $TMP_DIR
trap 'cleanup' INT TERM EXIT

print "Build gssapi, test client, and server"
go test github.com/apcera/gssapi
cp -R ./docker/service $TMP_DIR/
(cd $TMP_DIR/service && go build github.com/apcera/gssapi/component_test/service)
cp -R ./docker/client $TMP_DIR/
(cd $TMP_DIR/client && go test -c github.com/apcera/gssapi/component_test/client)

cleanup_containers

build_image "kdc" "unless exists"
run_image "kdc" "--detach" >/dev/null

function copy_keytab() {
	sudo docker cp kdc:/etc/docker-kdc/krb5.keytab $TMP_DIR/base-service
}
build_image "base-service" "unless exists" "copy_keytab"
build_image "service"
run_image "service" "--detach --link=kdc:kdc" >/dev/null


build_image "base-client" "unless exists"
build_image "client"
set +e
run_image "client" "--link=service:service --link=kdc:kdc"
if [[ "$?" != "0" ]]; then
	print "Client failed, see service logs below:\n\n"
	sudo docker logs service 2>&1
	exit 1
fi
set -e
echo "TEST PASSED"
