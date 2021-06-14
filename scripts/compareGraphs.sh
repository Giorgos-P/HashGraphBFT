# Show the Hashgraph sizes
# Compare the graphs  

cd ~/go/bin/logs/out

echo ""
echo "------Graph------"

grep Elements graph_0_.txt
grep Elements graph_1_.txt
grep Elements graph_1_.txt
grep Elements graph_1_.txt

echo "Compare graph_0_ with graph_1_"
cmp graph_0_.txt graph_1_.txt 
echo "Compare graph_0_ with graph_2_"
cmp graph_0_.txt graph_2_.txt 
echo "Compare graph_0_ with graph_3_"
cmp graph_0_.txt graph_3_.txt 

echo "Compare graph_1_ with graph_2_"
cmp graph_1_.txt graph_2_.txt 
echo "Compare graph_1_ with graph_3_"
cmp graph_1_.txt graph_3_.txt 

echo "Compare graph_2_ with graph_3_"
cmp graph_2_.txt graph_3_.txt 


