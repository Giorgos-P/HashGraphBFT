# Compile Hashgraph Server and
# execute N servers in the background

./compileHashGraphBFT.sh

cd ~/go/bin/logs/error/
rm *
cd ~/go/bin/logs/out/
rm *
cd ~/go/bin/

N=4
CLIENTS=10
SCE=0
REM=0

for (( ID=0; ID<$N; ID++ ))
do
	HashGraphBFT $ID $N $CLIENTS $SCE $REM &
done
