# Compile Hashgraph Client and
# execute clients in the background

#!/bin/bash

./compileHashGraphClient.sh

cd ~/go/bin/logs/client
rm *
cd ~/go/bin/

N=4
CLIENTS=10
REM=0


for (( ID=0; ID<$CLIENTS; ID++ ))
do
	HashGraph_Client $ID $N $CLIENTS $REM &
done
