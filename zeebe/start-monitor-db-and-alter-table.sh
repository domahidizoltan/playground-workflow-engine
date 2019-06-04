#!/bin/sh

(sleep 30 && \
echo "Alter RESOURCES_ column to text in WORKFLOW table" && \
java -cp /opt/h2/bin/h2-1.4.197.jar org.h2.tools.Shell \
-url "jdbc:h2:tcp://localhost:1521/zeebe-monitor" \
-user "sa" \
-password "" \
-sql "ALTER TABLE WORKFLOW ALTER COLUMN RESOURCE_ TEXT;") &

java -cp /opt/h2/bin/h2*.jar org.h2.tools.Server  \
-web \
-webAllowOthers \
-webPort 81 \
-tcp \
-tcpAllowOthers \
-tcpPort 1521 \
-baseDir ${DATA_DIR}
