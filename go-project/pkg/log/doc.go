package log

/*
1.一个算bug的非bug就是如果使用默认的log输出path,那么error级别以上的日志会重复打印两次
	这是因为错误日志会打印到stderr,同时作为日志又会打印到stdout
*/
