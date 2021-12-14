go get github.com/webx-top/db
go install github.com/webx-top/db/cmd/dbgenerator
dbgenerator -h 127.0.0.1 -d nging -p root -o ../application/dbschema -match "^(nging_alert_recipient|nging_alert_topic|nging_cloud_backup|nging_cloud_storage|nging_code_invitation|nging_code_verification|nging_config|nging_file|nging_file_embedded|nging_file_moved|nging_file_thumb|nging_kv|nging_login_log|nging_sending_log|nging_task|nging_task_group|nging_task_log|nging_user|nging_user_role|nging_user_u2f)$" -backup "../config/install.sql" -charset utf8mb4

# &dbschema\.([a-zA-Z0-9]+)\{\}    =>    dbschema.New$1(ctx)