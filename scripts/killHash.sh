#Kill Hashgraph processes

kill $(ps aux | grep 'Hash*' | awk '{print $2}')
