#set -x
set -e

function print() {
	echo "apcera/go-gssapi-test: $*"
}

print "Run the gssapi unit tests and build the sample app"
TMP_DIR="/tmp/$(uuidgen)"
mkdir -p $TMP_DIR

go test github.com/apcera/gssapi
go build -o $TMP_DIR/sample github.com/apcera/gssapi/sample

print "Clean up apcera/go-gssapi-test containers"
running=`sudo docker ps -a | grep 'apcera/go-gssapi-test' | awk '{print $1}'`
if [[ "$running" != "" ]]; then
	echo $running | xargs sudo docker stop >/dev/null
	echo $running | xargs sudo docker rm >/dev/null
fi

if [[ `sudo docker images -q apcera/go-gssapi-test-kdc` == "" ]]; then
	print "Build apcera/go-gssapi-test-kdc"
	cp -R ./kdc $TMP_DIR
	sudo docker build \
		-q \
		-t apcera/go-gssapi-test-kdc \
		$TMP_DIR/kdc \
		>/dev/null
fi

print "Run apcera/go-gssapi-test-kdc"
sudo docker run \
	-d \
	--hostname=kdc \
	--name=kdc \
	apcera/go-gssapi-test-kdc \
	>/dev/null

print "Build apcera/go-gssapi-test-service"
cp -R ./service $TMP_DIR
cp $TMP_DIR/sample $TMP_DIR/service/service
sudo docker cp kdc:/etc/docker-kdc/krb5.keytab $TMP_DIR/service
sudo docker build \
	-q \
	-t apcera/go-gssapi-test-service \
	$TMP_DIR/service >/dev/null

print "Run apcera/go-gssapi-test-service"
sudo docker run \
	-d \
	--hostname=service \
	--name=service \
	--link kdc:kdc \
	apcera/go-gssapi-test-service \
	>/dev/null

print "Build apcera/go-gssapi-test-client"
cp -R ./client $TMP_DIR
cp $TMP_DIR/sample $TMP_DIR/client/client
sudo docker build \
	-q \
	-t apcera/go-gssapi-test-client \
	$TMP_DIR/client \
	>/dev/null

rm -Rf $TMP_DIR

set +e
print "Run apcera/go-gssapi-test-client"
sudo docker run \
	--hostname=client \
	--name=client \
	--link kdc:kdc \
	--link service:service \
	apcera/go-gssapi-test-client

if [[ "$?" != "0" ]]; then
	print "Client failed, see service logs"
	sudo docker logs service 2>&1
	exit 1
fi
set -e
