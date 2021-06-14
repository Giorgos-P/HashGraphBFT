#Compile Hashgraph Server and client

cd ~/go/bin/
rm HashGraph_Client
rm HashGraphBFT

cd ~/go/src/
go install HashGraph_Client/
echo "" 
echo "" 
go install HashGraphBFT/
