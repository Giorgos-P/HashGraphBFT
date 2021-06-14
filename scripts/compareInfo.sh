# Compare info logs
# Info logs are used for debugging  

cd ~/go/bin/logs/out

echo ""
echo "------Info------"






echo "Compare info_0_ wit info_1_"
cmp info_0_.txt info_1_.txt 
echo "Compare info_0_ wit info_2_"
cmp info_0_.txt info_2_.txt 
echo "Compare info_0_ wit info_3_"
cmp info_0_.txt info_3_.txt 

echo "Compare info_1_ wit info_2_"
cmp info_1_.txt info_2_.txt 
echo "Compare info_1_ wit info_3_"
cmp info_1_.txt info_3_.txt 

echo "Compare info_2_ wit info_3_"
cmp info_2_.txt info_3_.txt 


