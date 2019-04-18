#!/bin/sh

IMAGE=gcr.io/gapic-images/gapic-showcase:0.1.1
SERVER=host.docker.internal:7469
FIXTURES=`pwd`/test/fixtures
VOLUME="$FIXTURES":/root/fixtures

OUTPUT=`mktemp`

### Start server
node build/src/index.js &
server_pid=$!
echo "Server PID is $server_pid"

retval=0

### 1. Test echo

docker run --rm --network host $IMAGE --address $SERVER \
  echo echo --response content --response.content test -j >> $OUTPUT

### 2. Test expand

docker run --rm --network host $IMAGE --address $SERVER \
  echo expand --content 'ab cd ef gh' -j >> $OUTPUT

### 3. Test collect

docker run -v $VOLUME --rm --network host $IMAGE --address $SERVER \
  echo collect --from_file /root/fixtures/collect.in -j >> $OUTPUT

### 4. Test chat

# note: commented out because of client bug
# docker run -v $VOLUME --rm --network host $IMAGE --address $SERVER \
#   echo chat --from_file /root/fixtures/chat.in --out_file /dev/stdout

### 5. Test pagedExpand

docker run --rm --network host $IMAGE --address $SERVER \
  echo paged-expand --content 'ab cd ef gh ij kl mn' --page_size 2 -j >> $OUTPUT

docker run --rm --network host $IMAGE --address $SERVER \
  echo paged-expand --content 'ab cd ef gh ij kl mn' --page_size 2 --page_token 4 -j >> $OUTPUT

docker run --rm --network host $IMAGE --address $SERVER \
  echo paged-expand --content 'ab cd ef gh ij kl mn' --page_size 2 --page_token 5 -j >> $OUTPUT

docker run --rm --network host $IMAGE --address $SERVER \
  echo paged-expand --content 'ab cd ef gh ij kl mn' --page_size 2 --page_token 6 -j >> $OUTPUT

### 6. Test wait
start_time=`date +%s`
docker run --rm --network host $IMAGE --address $SERVER \
  echo wait --end ttl --end.ttl.seconds 5 --follow --response success --response.success.content okay -j >> $OUTPUT
end_time=`date +%s`
if [ $(( end_time - start_time )) -lt 5 ]; then
  echo "Wait returned too fast"
  retval=1
fi

### Kill server
kill %1

### Check output
if ! diff "$FIXTURES/test.baseline" "$OUTPUT"; then
  echo "Differences in output"
  retval=2
fi

if [ $retval == 0 ]; then
  echo "Tests passed!"
fi

exit $retval
