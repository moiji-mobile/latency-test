#!/usr/bin/env bash

# vars
SIZES="16 32 1024"
SERVERS="C go python ruby pharo"
SERVERS="C go python ruby"
INTERVAL=0.25ms
NUM_MSG=100000


clean_all() {
	echo "T: Killing parent $$ and all children. Bye"
	pkill -9 -P $$
}
trap clean_all EXIT

# build everything we need
build_all() {
	echo "T: Building software"
	CFLAGS="-O2 -Wall" make -C contrib/ server
	set -e
	go build
	go build contrib/server.go
	set +e
}


test_one() {
	server=$1
	size=$2

	echo "T: Testing $server with $size"

	case $server in
	"go")
		DEST="localhost:7999"
		;;
	"python")
		DEST="localhost:8000"
		;;
	"C")
		DEST="localhost:8001"
		;;
	"ruby")
		DEST="localhost:8002"
		;;
	"pharo")
		DEST="localhost:8003"
		;;
	esac

	# warm up...
	echo "T: Result in result_${server}_${size}.csv"
	./latency-test  -messages $NUM_MSG -interval $INTERVAL -destination $DEST -packet-size $size > /dev/null 2>&1
	./latency-test  -messages $NUM_MSG -interval $INTERVAL -destination $DEST -packet-size $size -write-result result_${server}_${size}.csv > /dev/null
	RES=$?
	if [ $RES -gt 0 ]; then
		echo "T: Failed to run properly. TCP buffering?"
		exit 23
	fi
}

start_one() {
	server=$1
	case $server in
	"go")
		EXEC="./server"
		;;
	"python")
		EXEC="python3 ./contrib/server.py"
		;;
	"C")
		EXEC="./contrib/server 8001"
		;;
	"ruby")
		EXEC="./contrib/server.rb 8002"
		;;
	"pharo")
		EXEC="./contrib/pharo-vm..."
		;;
	esac

	echo "T: Launching $EXEC"
	$EXEC &
	sleep 1s
}

start_all() {
	for server in $SERVERS
	do
		start_one $server
	done
}

test_all() {
	for size in $SIZES
	do
		for server in $SERVERS
		do
			test_one $server $size
		done
	done
}

plot_all() {
	OUT="set datafile separator ','; plot "
	for size in $SIZES
	do
		for server in $SERVERS
		do
			OUT="$OUT \"result_${server}_${size}.csv\" using (\$2):(\$3/1000000) title \"$server ($size)\","
		done
	done

	# Print without the last ,
	echo ${OUT%?}
}


build_all
start_all
test_all
plot_all
