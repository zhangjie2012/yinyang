# yinyang

## 背景

想做一个节日倒计时的小程序，需要农历计算，发现竟然没有一个库可以直接来使用。所以就自己写了一个。

## 实现

中国的农历并不是通过算法算出来的，而是天文台预测的。

香港天文台公共和农历对照表：https://www.hko.gov.hk/tc/gts/time/conversion1_text.htm

思路就是从网站上抓数据到 rawdata 目录下，然后解析 rawdata 的数据，建立一个映射关系，提供对外服务。

*阴历和农历并不完全对等，这里只是为了好听，取了 Yin 和 Yang。*

## 功能

以 API 服务的方式暴露出去。

- [X] 支持公历和农历相互转换
- [X] 农历日历

## API 列表

- `/api/v1/years` 年份列表
- `/api/v1/years/${year}/months` 某年的月份列表
- `/api/v1/years/${year}/months/${month}/${leap}/days` 某年某月的农历日列表
- `/api/v1/conv/yang-yin/${year}/${month}/${day}` 公历转农历
- `/api/v1/conv/yin-yang/${year}/${month}/${leap}/${day}` 农历转公历

## 开发 & 测试

```
make build-local && ./bin/yinyang -rawdata ./rawdata -port 8080
```

```
➜  ~ curl -s localhost:8080/api/v1/conv/yang-yin/2021/4/28 | jq .
{
  "code": 0,
  "data": {
    "year_num": 2021,
    "year_tian": "辛",
    "year_di": "丑",
    "year_zodiac": "牛",
    "month_num": 3,
    "month_leap": false,
    "month_name": "三月",
    "day_num": 17,
    "day_name": "十七",
    "weekday": "星期三",
    "solarterm": ""
  },
  "error": ""
}
```
