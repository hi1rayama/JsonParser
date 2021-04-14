# JSONParser&整形ツール

* test.json
```
{
    "color_list": [ "red", "green", "blue" ],
    "num_list": [ 123, 456, 789 ],
    "mix_list": [ "red", 456, true ],
    "array_list": [ [ 12, 23 ], [ 34, 45 ], [ 56, 67 ] ],
    "object_list": [{ "name": "Tanaka", "age": 26 },{ "name": "Suzuki", "age": 32 }]
  }
```

```
$ go run main.go test.json
{
  "color_list":[
    "red",
    "green",
    "blue"
  ],
  "num_list":[
    123,
    456,
    789
  ],
  "mix_list":[
    "red",
    456,
    true
  ],
  "array_list":[
    [
      12,
      23
    ],
    [
      34,
      45
    ],
    [
      56,
      67
    ]
  ],
  "object_list":[
    {
      "name":"Tanaka",
      "age":26
    },
    {
      "name":"Suzuki",
      "age":32
    }
  ]
}
```


## 未対応
* エスケープ文字
* 実数や浮動小数
* 日本語などアルファベット以外の文字列