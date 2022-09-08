- # unixsoket-mysql

[测试数据](https://github.com/ymm135/test_db) 

- [基本信息](#基本信息)
  - [pprof 分析](#pprof-分析)
  - [mysql 不同连接方式性能表现](#mysql-不同连接方式性能表现)
  - [性能表现](#性能表现)
  - [`sysbench`性能监测](#sysbench性能监测)
    - [准备数据](#准备数据)
    - [执行测试](#执行测试)
    - [清理数据](#清理数据)
    - [测试结果](#测试结果)
  - [sql耗时分析](#sql耗时分析)
## 基本信息

高并发存储时，mysql的资源占用  
```shell
top - 14:16:57 up 21:20,  1 user,  load average: 2.72, 1.58, 0.84
Tasks: 158 total,   1 running, 157 sleeping,   0 stopped,   0 zombie
%Cpu(s):  7.3 us,  1.3 sy,  0.0 ni, 78.1 id, 13.2 wa,  0.0 hi,  0.1 si,  0.0 st
KiB Mem : 16178732 total, 14373064 free,  1296160 used,   509508 buff/cache
KiB Swap: 19914748 total, 19369212 free,   545536 used. 14595852 avail Mem 

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND                                                                                                                                                         
 1961 mysql     20   0 4755208 377580   3700 S  55.8  2.3 598:21.99 mysqld                                                                                                                                                          
23315 root      20   0  565200 145756   3248 S  14.3  0.9  80:40.11 Southwest_Engin                                                                                                                                                 
21668 root      20   0 1429800 181912  32944 S   1.7  1.1   0:15.85 audit-server                                                                                                                                                    
21691 root      20   0 1525192  51080  17760 S   0.3  0.3   0:02.85 watcher                                                                                                                                                         
    1 root      20   0  191152   2664   1492 S   0.0  0.0   0:03.15 systemd                                                                                                                                                         
    2 root      20   0       0      0      0 S   0.0  0.0   0:00.01 kthreadd     
```

长时间存储时`fd`浏览  
```shell
$ ls -l /proc/21668/fd
总用量 0
lr-x------ 1 root root 64 9月   2 18:03 0 -> pipe:[3182739]
l-wx------ 1 root root 64 9月   2 18:03 1 -> pipe:[3182740]
lrwx------ 1 root root 64 9月   2 18:03 10 -> socket:[4182657]
lrwx------ 1 root root 64 9月   2 18:03 100 -> socket:[3936054]
lrwx------ 1 root root 64 9月   2 18:03 101 -> socket:[3937797]
lrwx------ 1 root root 64 9月   2 18:03 102 -> socket:[3933068]
lrwx------ 1 root root 64 9月   2 18:03 103 -> socket:[3936762]
lrwx------ 1 root root 64 9月   2 18:03 104 -> socket:[3871143]
lrwx------ 1 root root 64 9月   2 18:03 105 -> socket:[3937798]
lrwx------ 1 root root 64 9月   2 18:03 106 -> socket:[3936763]
lrwx------ 1 root root 64 9月   2 18:03 107 -> socket:[3933069]
lrwx------ 1 root root 64 9月   2 18:03 108 -> socket:[3936764]
lrwx------ 1 root root 64 9月   2 18:03 109 -> socket:[3933070]
lrwx------ 1 root root 64 9月   2 18:03 11 -> socket:[3179927]
lrwx------ 1 root root 64 9月   2 18:03 110 -> socket:[3936055]
lrwx------ 1 root root 64 9月   2 18:03 111 -> socket:[3936765]
lrwx------ 1 root root 64 9月   2 18:03 112 -> socket:[3933071]
lrwx------ 1 root root 64 9月   2 18:03 113 -> socket:[4189297]
lrwx------ 1 root root 64 9月   2 18:03 114 -> socket:[4189296]
lrwx------ 1 root root 64 9月   2 18:03 115 -> socket:[4190406]
lrwx------ 1 root root 64 9月   2 18:03 116 -> socket:[4189655]
lrwx------ 1 root root 64 9月   2 18:03 117 -> socket:[4192422]
lrwx------ 1 root root 64 9月   2 18:03 118 -> socket:[4186937]
lrwx------ 1 root root 64 9月   2 18:03 119 -> socket:[4192320]
lrwx------ 1 root root 64 9月   2 18:03 12 -> socket:[3176177]
lrwx------ 1 root root 64 9月   2 18:03 120 -> socket:[4184954]
lrwx------ 1 root root 64 9月   2 18:03 121 -> socket:[4187288]
lrwx------ 1 root root 64 9月   2 18:03 122 -> socket:[4193313]
```



### pprof 分析
`http://127.0.0.1:6060/debug/pprof/`  

协程分析
```shell
goroutine profile: total 355731
320510 @ 0x8be025 0x8cfd65 0x8cfd4e 0x8f1687 0x911d85 0x16b06b7 0x16b06b8 0x9f2f23 0x9f2e5f 0x8f5601
#	0x8f1686	sync.runtime_SemacquireMutex+0x46					/usr/lib/golang/src/runtime/sema.go:71
#	0x911d84	sync.(*Mutex).lockSlow+0x104						/usr/lib/golang/src/sync/mutex.go:138
#	0x16b06b6	sync.(*Mutex).Lock+0xf6							/usr/lib/golang/src/sync/mutex.go:81
#	0x16b06b7	audit/server/socket.ObtainProtoSockHandler.func1+0xf7			/data/jenkins-audit/audit/server/socket/socket_handler.go:83
#	0x9f2f22	audit/server/utils/socket.(*UnixSocket).HandleServerContext+0x42	/data/jenkins-audit/audit/server/utils/socket/unix_socket.go:62
#	0x9f2e5e	audit/server/utils/socket.(*UnixSocket).HandleServerConn+0x3e		/data/jenkins-audit/audit/server/utils/socket/unix_socket.go:52

28793 @ 0x8be025 0x8cfd65 0x8cfd4e 0x8f1687 0x911d85 0x16b0c4e 0x16b0c4f 0x9f2f23 0x9f2e5f 0x8f5601
#	0x8f1686	sync.runtime_SemacquireMutex+0x46					/usr/lib/golang/src/runtime/sema.go:71
#	0x911d84	sync.(*Mutex).lockSlow+0x104						/usr/lib/golang/src/sync/mutex.go:138
#	0x16b0c4d	sync.(*Mutex).Lock+0x4ad						/usr/lib/golang/src/sync/mutex.go:81
#	0x16b0c4e	audit/server/socket.EventSockHandler.func1+0x4ae			/data/jenkins-audit/audit/server/socket/socket_handler.go:121
#	0x9f2f22	audit/server/utils/socket.(*UnixSocket).HandleServerContext+0x42	/data/jenkins-audit/audit/server/utils/socket/unix_socket.go:62
#	0x9f2e5e	audit/server/utils/socket.(*UnixSocket).HandleServerConn+0x3e		/data/jenkins-audit/audit/server/utils/socket/unix_socket.go:52

6387 @ 0x8be025 0x8cf297 0xe7b635 0x8f5601
#	0xe7b634	database/sql.(*DB).connectionOpener+0xb4	/usr/lib/golang/src/database/sql/sql.go:1133
```

数据处理不过来，阻塞在`auditMutex.Lock()`,`behaviorMutex.Lock()`,`connectionOpener select`  

`/usr/lib/golang/src/database/sql/sql.go`
```
func OpenDB(c driver.Connector) *DB {
	ctx, cancel := context.WithCancel(context.Background())
	db := &DB{
		connector:    c,
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		lastPut:      make(map[*driverConn]string),
		connRequests: make(map[uint64]chan connRequest),
		stop:         cancel,
	}

	go db.connectionOpener(ctx)  // 调用connectionOpener

	return db
}

// Runs in a separate goroutine, opens new connections when requested.
func (db *DB) connectionOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-db.openerCh:
			db.openNewConnection(ctx)
		}
	}
}
```

```shell
goroutine 21 [select, 154 minutes]:
database/sql.(*DB).connectionOpener(0xc000270680, 0x1bf5f98, 0xc000235000)
	/usr/lib/golang/src/database/sql/sql.go:1133 +0xb5
created by database/sql.OpenDB
	/usr/lib/golang/src/database/sql/sql.go:740 +0x12a
```

pprof heap分析
```shell
      File: audit-server
Build ID: 3f40d65b9579003996e94713e1ad33737a96d818
Type: alloc_space
Time: Sep 5, 2022 at 4:29pm (CST)
Showing nodes accounting for 8229.34MB, 100% of 8229.34MB total
----------------------------------------------------------+-------------
      flat  flat%   sum%        cum   cum%   calls calls% + context 	 	 
----------------------------------------------------------+-------------
                                         1539.50MB   100% |   github.com/go-sql-driver/mysql.(*connector).Connect /root/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/connector.go:74 (inline)
 1539.50MB 18.71% 18.71%  1539.50MB 18.71%                | github.com/go-sql-driver/mysql.newBuffer /root/go/pkg/mod/github.com/go-sql-driver/mysql@v1.5.0/buffer.go:38
----------------------------------------------------------+-------------
                                          712.22MB   100% |   github.com/gosnmp/gosnmp.(*GoSNMP).Connect /root/go/pkg/mod/github.com/gosnmp/gosnmp@v1.35.0/gosnmp.go:272
  712.22MB  8.65% 27.36%   712.22MB  8.65%                | github.com/gosnmp/gosnmp.(*GoSNMP).connect /root/go/pkg/mod/github.com/gosnmp/gosnmp@v1.35.0/gosnmp.go:315
----------------------------------------------------------+-------------
                                          615.82MB   100% |   audit/server/utils/socket.(*UnixSocket).StartServer /data/jenkins-audit/audit/server/utils/socket/unix_socket.go:69
  615.82MB  7.48% 34.85%   615.82MB  7.48%                | audit/server/utils/socket.(*UnixSocket).createServer /data/jenkins-audit/audit/server/utils/socket/unix_socket.go:45
----------------------------------------------------------+-------------
                                          426.16MB   100% |   runtime.newproc1 /usr/lib/golang/src/runtime/proc.go:4065
  426.16MB  5.18% 40.02%   426.16MB  5.18%                | runtime.malg /usr/lib/golang/src/runtime/proc.go:3988
----------------------------------------------------------+-------------
                                          336.11MB   100% |   os.(*File).Readdirnames /usr/lib/golang/src/os/dir.go:70
  336.11MB  4.08% 44.11%   336.11MB  4.08%                | os.(*File).readdir /usr/lib/golang/src/os/dir_unix.go:35
----------------------------------------------------------+-------------
                                          101.37MB 53.30% |   database/sql.(*DB).execDC.func2 /usr/lib/golang/src/database/sql/sql.go:1573
                                           86.81MB 45.65% |   database/sql.resultFromStatement /usr/lib/golang/src/database/sql/sql.go:2521
                                               1MB  0.53% |   database/sql.(*DB).queryDC.func1 /usr/lib/golang/src/database/sql/sql.go:1646
                                               1MB  0.53% |   database/sql.rowsiFromStatement /usr/lib/golang/src/database/sql/sql.go:2691
  190.18MB  2.31% 46.42%   190.18MB  2.31%                | database/sql.driverArgsConnLocked /usr/lib/golang/src/database/sql/convert.go:108
----------------------------------------------------------+-------------
                                          164.64MB   100% |   bufio.NewReader /usr/lib/golang/src/bufio/bufio.go:63 (inline)
  164.64MB  2.00% 48.42%   164.64MB  2.00%                | bufio.NewReaderSize /usr/lib/golang/src/bufio/bufio.go:57
```


### mysql 不同连接方式性能表现  

[mysql8 不同连接方式性能表现](https://blog.herecura.eu/blog/2021-03-03-mysql-local-vs-remote/)  

| type | transactions/sec | queries/sec | 95% latency(ms)| percentage |  
| ---- | ---- | ---- | ---- |  ---- | 
| local socket | 1159.13 | 33614.79 | 2.18 | 100% |  
| local tcp | 900.29 | 26108.54 | 2.81 | 77.7% |  
| remote tcp | 326.32 | 9463.20 | 7.84 | 28.2% |  
| remote socket | 259.16 | 7515.59 | 9.73 | 22.4% |  

### 性能表现
如果远程连接，mysql每秒处理10条数据，本地连接，处理20条

```
# remote 
paeseDataAndStore handler data 6 pps
paeseDataAndStore handler data 8 pps
paeseDataAndStore handler data 10 pps

# local
paeseDataAndStore handler data 20 pps
paeseDataAndStore handler data 26 pps
paeseDataAndStore handler data 24 pps
paeseDataAndStore handler data 21 pps
```

### `sysbench`性能监测

sysbench是跨平台的基准测试工具，支持多线程，支持多种数据库；主要包括以下几种测试：

- cpu性能
- 磁盘io性能
- 调度程序性能
- 内存分配及传输速度
- POSIX线程性能
- 数据库性能(OLTP基准测试)  


本文主要介绍对数据库性能的测试。  

安装
```shell
# centos
curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.rpm.sh | sudo bash
sudo yum -y install sysbench

# ubuntu 
curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.deb.sh | sudo bash
sudo apt -y install sysbench

# macos
brew install sysbench
```

使用说明
```shell
Usage:
  sysbench [options]... [testname] [command]

Commands implemented by most tests: prepare run cleanup help

General options:
  --threads=N                     number of threads to use [1]
  --events=N                      limit for total number of events [0]
  --time=N                        limit for total execution time in seconds [10]
  --forced-shutdown=STRING        number of seconds to wait after the --time limit before forcing shutdown, or 'off' to disable [off]
  --thread-stack-size=SIZE        size of stack per thread [64K]
  --rate=N                        average transactions rate. 0 for unlimited rate [0]
  --report-interval=N             periodically report intermediate statistics with a specified interval in seconds. 0 disables intermediate reports [0]
  --report-checkpoints=[LIST,...] dump full statistics and reset all counters at specified points in time. The argument is a list of comma-separated values representing the amount of time in seconds elapsed from start of test when report checkpoint(s) must be performed. Report checkpoints are off by default. []
  --debug[=on|off]                print more debugging info [off]
  --validate[=on|off]             perform validation checks where possible [off]
  --help[=on|off]                 print help and exit [off]
  --version[=on|off]              print version and exit [off]
  --config-file=FILENAME          File containing command line options
  --tx-rate=N                     deprecated alias for --rate [0]
  --max-requests=N                deprecated alias for --events [0]
  --max-time=N                    deprecated alias for --time [0]
  --num-threads=N                 deprecated alias for --threads [1]
```

mysql参数
```shell
mysql options:
  --mysql-host=[LIST,...]          MySQL server host [localhost]
  --mysql-port=[LIST,...]          MySQL server port [3306]
  --mysql-socket=[LIST,...]        MySQL socket
  --mysql-user=STRING              MySQL user [sbtest]
  --mysql-password=STRING          MySQL password []
  --mysql-db=STRING                MySQL database name [sbtest]
  --mysql-ssl[=on|off]             use SSL connections, if available in the client library [off]
  --mysql-ssl-cipher=STRING        use specific cipher for SSL connections []
  --mysql-compression[=on|off]     use compression, if available in the client library [off]
  --mysql-debug[=on|off]           trace all client library calls [off]
  --mysql-ignore-errors=[LIST,...] list of errors to ignore, or "all" [1213,1020,1205]
  --mysql-dry-run[=on|off]         Dry run, pretend that all MySQL client API calls are successful without executing them [off]
```

lua脚本位置
```shell
▶ ls -l /usr/local/Cellar/sysbench/1.0.20_2/share/sysbench/tests/include/oltp_legacy 
total 104
-rw-r--r--  1 ymm  admin  1195  4 24  2020 bulk_insert.lua
-rw-r--r--  1 ymm  admin  4696  4 24  2020 common.lua
-rw-r--r--  1 ymm  admin   366  4 24  2020 delete.lua
-rw-r--r--  1 ymm  admin  1171  4 24  2020 insert.lua
-rw-r--r--  1 ymm  admin  3004  4 24  2020 oltp.lua
-rw-r--r--  1 ymm  admin   368  4 24  2020 oltp_simple.lua
-rw-r--r--  1 ymm  admin   527  4 24  2020 parallel_prepare.lua
-rw-r--r--  1 ymm  admin   369  4 24  2020 select.lua
-rw-r--r--  1 ymm  admin  1448  4 24  2020 select_random_points.lua
-rw-r--r--  1 ymm  admin  1556  4 24  2020 select_random_ranges.lua
-rw-r--r--  1 ymm  admin   369  4 24  2020 update_index.lua
-rw-r--r--  1 ymm  admin   578  4 24  2020 update_non_index.lua
```

注意事项
在执行sysbench时，应该注意：  
1. 尽量不要在MySQL服务器运行的机器上进行测试，一方面可能无法体现网络（哪怕是局域网）的影响，另一方面，sysbench的运行（尤其是设置的并发数较高时）会影响MySQL服务器的表现。
2. 可以逐步增加客户端的并发连接数（--thread参数），观察在连接数不同情况下，MySQL服务器的表现；如分别设置为10,20,50,100等。
3. 一般执行模式选择complex即可，如果需要特别测试服务器只读性能，或不使用事务时的性能，可以选择simple模式或nontrx模式。
4. 如果连续进行多次测试，注意确保之前测试的数据已经被清理干净。  

#### 准备数据

创建数据库
```shell
CREATE DATABASE IF NOT EXISTS sbtest DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
```

```shell
# 远程测试
HOST=10.25.10.125
LUA_FILE=/usr/local/Cellar/sysbench/1.0.20_2/share/sysbench/tests/include/oltp_legacy/oltp.lua

# 本地测试
HOST=localhost
LUA_FILE=/usr/share/sysbench/tests/include/oltp_legacy/oltp.lua

sysbench $LUA_FILE \
 --mysql-host=$HOST \
 --mysql-port=3306 \
 --mysql-user=root --mysql-password=root \
 --oltp-test-mode=complex --oltp-tables-count=10 \
 --oltp-table-size=100000 --threads=10 --time=120 \
 --report-interval=10 \
 prepare
```

> 执行模式为complex，使用了10个表，每个表有10万条数据，客户端的并发线程数为10，执行时间为120秒，每10秒生成一次报告  

输出日志， 创建10张还有10万数据的表  
```shell
Creating table 'sbtest1'...
Inserting 100000 records into 'sbtest1'
Creating secondary indexes on 'sbtest1'...
Creating table 'sbtest2'...
Inserting 100000 records into 'sbtest2'
Creating secondary indexes on 'sbtest2'...
Creating table 'sbtest3'...
Inserting 100000 records into 'sbtest3'
Creating secondary indexes on 'sbtest3'...
Creating table 'sbtest4'...
Inserting 100000 records into 'sbtest4'
Creating secondary indexes on 'sbtest4'...
Creating table 'sbtest5'...
Inserting 100000 records into 'sbtest5'
Creating secondary indexes on 'sbtest5'...
Creating table 'sbtest6'...
Inserting 100000 records into 'sbtest6'
Creating secondary indexes on 'sbtest6'...
Creating table 'sbtest7'...
Inserting 100000 records into 'sbtest7'
Creating secondary indexes on 'sbtest7'...
Creating table 'sbtest8'...
Inserting 100000 records into 'sbtest8'
Creating secondary indexes on 'sbtest8'...
Creating table 'sbtest9'...
Inserting 100000 records into 'sbtest9'
Creating secondary indexes on 'sbtest9'...
Creating table 'sbtest10'...
Inserting 100000 records into 'sbtest10'
Creating secondary indexes on 'sbtest10'..
```

表结构
```shell
mysql> select * from sbtest1 limit 1;
+----+-------+-------------------------------------------------------------------------------------------------------------------------+-------------------------------------------------------------+
| id | k     | c                                                                                                                       | pad                                                         |
+----+-------+-------------------------------------------------------------------------------------------------------------------------+-------------------------------------------------------------+
|  1 | 49929 | 83868641912-28773972837-60736120486-75162659906-27563526494-20381887404-41576422241-93426793964-56405065102-33518432330 | 67847967377-48000963322-62604785301-91415491898-96926520291 |
+----+-------+-------------------------------------------------------------------------------------------------------------------------+-------------------------------------------------------------+
1 row in set (0.00 sec)
```

查看lua脚本，数据库操作包含增删改查   
```shell
for i=1, oltp_point_selects do
      rs = db_query("SELECT c FROM ".. table_name .." WHERE id=" ..
                       sb_rand(1, oltp_table_size))
   end

   if oltp_range_selects then

   for i=1, oltp_simple_ranges do
      rs = db_query("SELECT c FROM ".. table_name .. get_range_str())
   end

   for i=1, oltp_sum_ranges do
      rs = db_query("SELECT SUM(K) FROM ".. table_name .. get_range_str())
   end

   for i=1, oltp_order_ranges do
      rs = db_query("SELECT c FROM ".. table_name .. get_range_str() ..
                    " ORDER BY c")
   end

   for i=1, oltp_distinct_ranges do
      rs = db_query("SELECT DISTINCT c FROM ".. table_name .. get_range_str() ..
                    " ORDER BY c")
   end

   for i=1, oltp_index_updates do
      rs = db_query("UPDATE " .. table_name .. " SET k=k+1 WHERE id=" .. sb_rand(1, oltp_table_size))
   end

   for i=1, oltp_non_index_updates do
      c_val = sb_rand_str("###########-###########-###########-###########-###########-###########-###########-###########-###########-###########")
      query = "UPDATE " .. table_name .. " SET c='" .. c_val .. "' WHERE id=" .. sb_rand(1, oltp_table_size)
      rs = db_query(query)
      if rs then
        print(query)
      end
   end

   for i=1, oltp_delete_inserts do

   i = sb_rand(1, oltp_table_size)

   rs = db_query("DELETE FROM " .. table_name .. " WHERE id=" .. i)
   
   c_val = sb_rand_str([[
###########-###########-###########-###########-###########-###########-###########-###########-###########-###########]])
   pad_val = sb_rand_str([[
###########-###########-###########-###########-###########]])

   rs = db_query("INSERT INTO " .. table_name ..  " (id, k, c, pad) VALUES " .. string.format("(%d, %d, '%s', '%s')",i, sb_rand(1, oltp_table_size) , c_val, pad_val))

   end
```


> 提示文件找不到 FATAL: error 2002: Can't connect to local MySQL server through socket '/var/lib/mysql/mysql.sock' (2)  
> 本地文件在`/tmp/mysql.sock`  

```shell
mkdir -p /var/lib/mysql/

# 注意用户权限
chown -R mysql:mysql /var/lib/mysql/

# /etc/my.cnf
[mysqld]
socket=/var/lib/mysql/mysql.sock
```



#### 执行测试
将测试结果导出到文件中，便于后续分析。  

```shell
sysbench $LUA_FILE \
 --mysql-host=$HOST \
 --mysql-port=3306 \
 --mysql-user=root --mysql-password=root \
 --oltp-test-mode=complex --oltp-tables-count=10 \
 --oltp-table-size=100000 --threads=10 --time=120 \
 --report-interval=10 \
 run >> mysysbench.log
```


#### 清理数据
```shell
sysbench $LUA_FILE \
 --mysql-host=$HOST \
 --mysql-port=3306 \
 --mysql-user=root --mysql-password=root \
  --oltp-test-mode=complex --oltp-tables-count=10 \
 --oltp-table-size=100000 --threads=10 --time=120 \
 --report-interval=10 \
 cleanup
```


#### 测试结果
远程连接测试
```shell
SQL statistics:
    queries performed:
        read:                            36372
        write:                           10392
        other:                           5196
        total:                           51960
    transactions:                        2598   (21.56 per sec.)
    queries:                             51960  (431.13 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          120.5195s
    total number of events:              2598

Latency (ms):
         min:                                  181.88
         avg:                                  463.26
         max:                                 1084.22
         95th percentile:                      719.92
         sum:                              1203547.73

Threads fairness:
    events (avg/stddev):           259.8000/2.60
    execution time (avg/stddev):   120.3548/0.09
```

mysql本地连接测试
```shell
SQL statistics:
    queries performed:
        read:                            125552
        write:                           35872
        other:                           17936
        total:                           179360
    transactions:                        8968   (74.54 per sec.)
    queries:                             179360 (1490.70 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          120.3180s
    total number of events:              8968

Latency (ms):
         min:                                   36.81
         avg:                                  133.91
         max:                                  788.42
         95th percentile:                      376.49
         sum:                              1200926.35

Threads fairness:
    events (avg/stddev):           896.8000/5.44
    execution time (avg/stddev):   120.0926/0.08
```

### sql耗时分析
```shell
mysql> show variables like 'profiling';
+---------------+-------+
| Variable_name | Value |
+---------------+-------+
| profiling     | OFF   |
+---------------+-------+
1 row in set (0.00 sec)

mysql> set GLOBAL profiling = on;
mysql> set profiling = on;
Query OK, 0 rows affected, 1 warning (0.00 sec)

mysql> show variables like 'profiling';
+---------------+-------+
| Variable_name | Value |
+---------------+-------+
| profiling     | ON    |
+---------------+-------+
1 row in set (0.00 sec)

mysql> show profiles;
+----------+------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Query_ID | Duration   | Query                                                                                                                                                        |
+----------+------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
|        1 | 0.00281875 | show variables like 'profiling'                                                                                                                              |
|        2 | 0.03362250 | INSERT INTO `employees` (`birth_date`,`first_name`,`last_name`,`gender`,`hire_date`) VALUES ('1953-09-02 00:00:00','G12','Tester','M','1986-06-26 00:00:00') |
+----------+------------+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
2 rows in set, 1 warning (0.00 sec)

mysql> show profile cpu, block io for query 2;
+----------------------+----------+----------+------------+--------------+---------------+
| Status               | Duration | CPU_user | CPU_system | Block_ops_in | Block_ops_out |
+----------------------+----------+----------+------------+--------------+---------------+
| starting             | 0.000116 | 0.000051 |   0.000061 |            0 |             0 |
| checking permissions | 0.000014 | 0.000005 |   0.000007 |            0 |             0 |
| Opening tables       | 0.000029 | 0.000014 |   0.000016 |            0 |             0 |
| init                 | 0.000031 | 0.000013 |   0.000017 |            0 |             0 |
| System lock          | 0.000013 | 0.000006 |   0.000007 |            0 |             0 |
| update               | 0.000089 | 0.000041 |   0.000049 |            0 |             0 |
| end                  | 0.000010 | 0.000003 |   0.000005 |            0 |             0 |
| query end            | 0.033245 | 0.000106 |   0.000127 |            0 |             8 |
| closing tables       | 0.000025 | 0.000010 |   0.000012 |            0 |             0 |
| freeing items        | 0.000030 | 0.000013 |   0.000017 |            0 |             0 |
| cleaning up          | 0.000022 | 0.000010 |   0.000012 |            0 |             0 |
+----------------------+----------+----------+------------+--------------+---------------+
11 rows in set, 1 warning (0.00 sec)
```

正常模式:  
2000条数据耗时`144`S  

预处理模式:
2000条数据耗时`97`S

批量处理模式:
10000条数据耗时`1`S  

如果批量插入超过1w提示`prepared statement contains too many placeholders gorm`   

这种情况是只有一个任务执行，如果有多个任务对同一个表格执行呢？  
测试情况为20个协程同时往一张表中插入1w数据?  

```shell
6 save data end 耗时: 1 s
2 save data end 耗时: 2 s
7 save data end 耗时: 2 s
14 save data end 耗时: 2 s
11 save data end 耗时: 3 s
13 save data end 耗时: 3 s
16 save data end 耗时: 3 s
19 save data end 耗时: 3 s
15 save data end 耗时: 3 s
5 save data end 耗时: 3 s
3 save data end 耗时: 3 s
10 save data end 耗时: 3 s
17 save data end 耗时: 3 s
4 save data end 耗时: 3 s
8 save data end 耗时: 3 s
1 save data end 耗时: 3 s
12 save data end 耗时: 3 s
9 save data end 耗时: 3 s
0 save data end 耗时: 3 s
18 save data end 耗时: 3 s
batch insert data,cost  3 s,avg 0.15 s
```

测试情况为20个协程同时往不同张表中插入1w数据?  





















