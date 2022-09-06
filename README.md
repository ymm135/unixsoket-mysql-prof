- # unixsoket-mysql

[测试数据](https://github.com/ymm135/test_db) 

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

长时间存储时`fd`耗尽
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

加锁失败:  
```shell
2022/09/02 18:06:28 /data/jenkins-audit/audit/server/service/flowAudit/flow_audit.go:65 Error 1205: Lock wait timeout exceeded; try restarting transaction
[50813.294ms] [rows:0] update icdevicetraffics set ic_device_ip = 'NULL' where ic_device_mac = '64:ae:0c:34:ab:98'

2022/09/02 18:06:28 /data/jenkins-audit/audit/server/service/flowAudit/flow_audit.go:65 Error 1205: Lock wait timeout exceeded; try restarting transaction
[50810.334ms] [rows:0] update icdevicetraffics set ic_device_ip = 'NULL' where ic_device_mac = '00:c0:a8:f2:61:fb'

2022/09/02 18:06:28 /data/jenkins-audit/audit/server/service/flowAudit/flow_audit.go:65 Error 1205: Lock wait timeout exceeded; try restarting transaction
[50809.743ms] [rows:0] update icdevicetraffics set ic_device_ip = 'NULL' where ic_device_mac = '00:0c:29:6b:2a:28'

2022/09/02 18:06:28 /data/jenkins-audit/audit/server/service/flowAudit/flow_audit.go:65 Error 1205: Lock wait timeout exceeded; try restarting transaction
[50808.233ms] [rows:0] update icdevicetraffics set ic_device_ip = '192.168.1.30' where ic_device_mac = '00:e0:ab:01:10:18'
```


### pprof 分析

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


### mysql socket连接/长连接








