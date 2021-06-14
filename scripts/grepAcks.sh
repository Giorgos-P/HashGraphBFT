cd ~/go/bin/logs/client
egrep Total *out* > totalTime.txt
egrep Time *out* > TimePerMsg.txt

egrep ACK *out* > Acks.txt
cat Acks.txt | awk '{ print $10 }' > Acks2.txt



CLIENTS=20


for (( ID=0; ID<$CLIENTS; ID++ ))
do
  egrep "^"$ID"_" TimePerMsg.txt | tail -n 1 >> LastMsgTime.txt
done

cat LastMsgTime.txt | awk '{ print $8 }' > LastMsgTime2.txt
cat LastMsgTime.txt | awk '{ print $6 }' > AckPerNode.txt

cd ~/go/bin/logs/out
NODES=4

for (( ID=0; ID<$NODES; ID++ ))
do
  egrep "Size:" "info_"$ID* | tail -n -1 >> Gossip.txt
  egrep Gossip "info_"$ID* | tail -n -1 >> Gossip.txt
done

