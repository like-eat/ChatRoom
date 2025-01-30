echo Starting Zookeeper service...
set ZOOKEEPER_HOME=G:\\GoogleDownloads\\zookeeper-3.4.10
set ZOOCFG=%ZOOKEEPER_HOME%\\conf\\zoo.cfg
call "%ZOOKEEPER_HOME%\\bin\\zkServer.cmd" start