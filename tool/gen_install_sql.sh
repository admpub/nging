db_host="localhost"
db_username="root"
db_password="root"
db_name="nging"
db_charset="utf8"
db_port=3306
mysqldump -d $db_name -h$db_host -P$db_port -u$db_username -p$db_password --default-character-set=$db_charset --single-transaction --set-gtid-purged=OFF | sed 's/ AUTO_INCREMENT=[0-9]*\s*//g' > ../config/install.sql
#echo 按任意键继续
#read -n 1
