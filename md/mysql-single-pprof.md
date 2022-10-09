- # mysql 单表性能  

## 性能数据
测试工具为`sysbench` [详情](../README.md#sysbench性能监测)  
```shell
# 本地测试
HOST=localhost
LUA_FILE=/usr/share/sysbench/tests/include/oltp_legacy/oltp.lua
SIZE=100000000

sysbench $LUA_FILE \
 --mysql-host=$HOST \
 --mysql-port=3306 \
 --mysql-socket=/tmp/mysql.sock \
 --mysql-user=root --mysql-password=root \
 --oltp-test-mode=complex --oltp-tables-count=10 \
 --oltp-table-size=100000 --threads=10 --time=120 \
 --report-interval=10 \
 prepare
 # 运行
 run >> mysysbench.log
 # 清理
 cleanup
```
### 10w  

```shell
Threads started!

[ 10s ] thds: 10 tps: 795.62 qps: 15922.34 (r/w/o: 11147.01/3183.09/1592.24) lat (ms,95%): 55.82 err/s: 0.00 reconn/s: 0.00
[ 20s ] thds: 10 tps: 966.38 qps: 19335.32 (r/w/o: 13534.96/3867.60/1932.75) lat (ms,95%): 44.17 err/s: 0.00 reconn/s: 0.00
[ 30s ] thds: 10 tps: 1016.04 qps: 20321.42 (r/w/o: 14224.58/4064.76/2032.08) lat (ms,95%): 41.85 err/s: 0.00 reconn/s: 0.00
[ 40s ] thds: 10 tps: 938.32 qps: 18764.03 (r/w/o: 13135.53/3751.86/1876.63) lat (ms,95%): 41.85 err/s: 0.00 reconn/s: 0.00
[ 50s ] thds: 10 tps: 909.27 qps: 18182.29 (r/w/o: 12728.44/3635.30/1818.55) lat (ms,95%): 47.47 err/s: 0.00 reconn/s: 0.00
[ 60s ] thds: 10 tps: 945.80 qps: 18916.88 (r/w/o: 13241.75/3783.52/1891.61) lat (ms,95%): 45.79 err/s: 0.00 reconn/s: 0.00
[ 70s ] thds: 10 tps: 974.55 qps: 19495.29 (r/w/o: 13645.49/3900.70/1949.10) lat (ms,95%): 44.17 err/s: 0.00 reconn/s: 0.00
[ 80s ] thds: 10 tps: 979.94 qps: 19593.09 (r/w/o: 13716.33/3916.88/1959.89) lat (ms,95%): 44.17 err/s: 0.00 reconn/s: 0.00
[ 90s ] thds: 10 tps: 930.37 qps: 18613.17 (r/w/o: 13027.33/3725.09/1860.75) lat (ms,95%): 44.17 err/s: 0.00 reconn/s: 0.00
[ 100s ] thds: 10 tps: 904.62 qps: 18087.19 (r/w/o: 12662.75/3615.20/1809.25) lat (ms,95%): 43.39 err/s: 0.00 reconn/s: 0.00
[ 110s ] thds: 10 tps: 977.44 qps: 19554.81 (r/w/o: 13686.90/3913.04/1954.87) lat (ms,95%): 44.98 err/s: 0.00 reconn/s: 0.00
[ 120s ] thds: 10 tps: 1020.66 qps: 20407.27 (r/w/o: 14286.29/4079.65/2041.33) lat (ms,95%): 41.85 err/s: 0.00 reconn/s: 0.00
SQL statistics:
    queries performed:
        read:                            1590470
        write:                           454420
        other:                           227210
        total:                           2272100
    transactions:                        113605 (946.22 per sec.)
    queries:                             2272100 (18924.41 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          120.0570s
    total number of events:              113605

Latency (ms):
         min:                                    1.96
         avg:                                   10.56
         max:                                  335.25
         95th percentile:                       44.17
         sum:                              1200098.51

Threads fairness:
    events (avg/stddev):           11360.5000/21.67
    execution time (avg/stddev):   120.0099/0.02
```

### 100w
```shell
SQL statistics:
    queries performed:
        read:                            321552
        write:                           91872
        other:                           45936
        total:                           459360
    transactions:                        22968  (191.32 per sec.)
    queries:                             459360 (3826.37 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          120.0462s
    total number of events:              22968

Latency (ms):
         min:                                    2.71
         avg:                                   52.26
         max:                                  382.62
         95th percentile:                      134.90
         sum:                              1200279.45

Threads fairness:
    events (avg/stddev):           2296.8000/11.69
    execution time (avg/stddev):   120.0279/0.00
```

### 1000w

```shell
SQL statistics:
    queries performed:
        read:                            134162
        write:                           38332
        other:                           19166
        total:                           191660
    transactions:                        9583   (79.79 per sec.)
    queries:                             191660 (1595.78 per sec.)
    ignored errors:                      0      (0.00 per sec.)
    reconnects:                          0      (0.00 per sec.)

General statistics:
    total time:                          120.0997s
    total number of events:              9583

Latency (ms):
         min:                                    3.95
         avg:                                  125.31
         max:                                 3273.44
         95th percentile:                      303.33
         sum:                              1200799.42

Threads fairness:
    events (avg/stddev):           958.3000/10.29
    execution time (avg/stddev):   120.0799/0.02
```


### 10000w
```shell
mysql> select count(*) from sbtest1;
+-----------+
| count(*)  |
+-----------+
| 100000000 |
+-----------+
1 row in set (14.72 sec)

# count cpu占用
2012 mysql     20   0 3968776 405244   5564 S  56.5  2.5 141:19.38 mysqld
```

性能数据位置，数据准备5个小时，最后只有三张表有1亿，脚本退出了。  

## 原理分析

