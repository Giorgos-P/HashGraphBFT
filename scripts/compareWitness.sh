# Compare Witnesses

cd ~/go/bin/logs/out

echo ""
echo "------Witness------"

wc -l witness_0_.txt
wc -l witness_1_.txt
wc -l witness_2_.txt
wc -l witness_3_.txt

echo "Compare witness_0_ wit witness_1_"
cmp witness_0_.txt witness_1_.txt 
echo "Compare witness_0_ wit witness_2_"
cmp witness_0_.txt witness_2_.txt 
echo "Compare witness_0_ wit witness_3_"
cmp witness_0_.txt witness_3_.txt 

echo "Compare witness_1_ wit witness_2_"
cmp witness_1_.txt witness_2_.txt 
echo "Compare witness_1_ wit witness_3_"
cmp witness_1_.txt witness_3_.txt 

echo "Compare witness_2_ wit witness_3_"
cmp witness_2_.txt witness_3_.txt 

