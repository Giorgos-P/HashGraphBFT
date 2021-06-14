#Count The ACKs the clients received

cd ~/go/bin/logs/client
grep ACK *out* | wc -l
