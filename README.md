# gen_columns

### 概述
有一个结构作为数据库表结构如下:
```
type User struct {
    ID int `json:"id" bson:"id"`
    Name string `json:"name" bson:"name"`
}
```

当使用这个model里的字段进行sql查询时, 通常使用:

```
map[string]interface{}{
    "id":123456,
}

```

作为查询条件, 如果当字段名更改时, 不得不修改这个map里的key值
如果能够自动生成一个结构体, 用于表示这些column name值, 那么只需修改一处:

```
map[string]interface{}{
    UserColumns.ID: 123456
}
```

### 使用方法
```
gen_columns -tag="bson" -path="./models/user.go"
```

会生成一个独立的文件, 里面的内容为:
```
package models

type _UserColumn struct {
	ID   string
	Name string
}

var UserColumns _UserColumn

func init() {
	UserColumns.ID = "id"
	UserColumns.Name = "name"

}
```

