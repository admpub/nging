for i in {1..100};do 
sleep 1; 
if [ $i != 5 ];then 
echo "abc # handle $i"
else
echo "abc $i"
fi
done;