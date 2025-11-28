# Gen

把csv文件转成phone.dat

### 使用

```bash
go run cmd/gen/csv2dat.go <input_csv> <output_dat>
```

或者编译后使用:

```bash
go build -o csv2dat cmd/gen/csv2dat.go
./csv2dat input.csv output.dat
```

### CSV格式

输入的CSV文件必须包含以下表头:

```
prefix,phone,province,city,isp,tel_code,postal_code,area_code
```

例子:

```csv
prefix,phone,province,city,isp,tel_code,postal_code,area_code
130,1300000,山东,济南,中国联通,0531,250000,370100
130,1300001,江苏,常州,中国联通,0519,213000,320400
130,1300002,安徽,合肥,中国联通,0551,230000,340100
```

### 运营商

支持以下运营商的字符串:
- 中国移动
- 中国联通
- 中国电信
- 中国广电
- 中国移动虚拟运营商
- 中国联通虚拟运营商
- 中国电信虚拟运营商
- 中国广电虚拟运营商
