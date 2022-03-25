## logger
logger包是日志管理包
```yaml
base:
  logger:
    # 日志root级别：trace/debug/info/warn/error/fatal/panic，默认：info
    level: info
    time:
      # 时间格式，time包中的内容
      format: time.RFC3339
    # 日志颜色
    color:
      # 启用：true/false，默认：false
      enable: false
    split:
      # 日志是否启用切分：true/false，默认false
      enable: false
      # 日志拆分的单位：MB
      size: 300
    ## 日志文件目录，默认工程目录的logs文件夹
    dir: ./logs/
    
    
```
