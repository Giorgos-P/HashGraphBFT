#Compare servers Ordered transactions

cd ~/go/bin/logs/out

echo ""
echo "------Order------"

wc -l order_0_.txt
wc -l order_1_.txt
wc -l order_2_.txt
wc -l order_3_.txt

echo "Compare order_0_ wit order_1_"
cmp order_0_.txt order_1_.txt 
echo "Compare order_0_ wit order_2_"
cmp order_0_.txt order_2_.txt 
echo "Compare order_0_ wit order_3_"
cmp order_0_.txt order_3_.txt 

echo "Compare order_1_ wit order_2_"
cmp order_1_.txt order_2_.txt 
echo "Compare order_1_ wit order_3_"
cmp order_1_.txt order_3_.txt 

echo "Compare order_2_ wit order_3_"
cmp order_2_.txt order_3_.txt 


