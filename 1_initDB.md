# MongoDB
## Docker start
> docker run -d --name my-mongodb   -e MONGO_INITDB_ROOT_USERNAME=clay   -e MONGO_INITDB_ROOT_PASSWORD=password   -p 27017:27017   mongo:6.0
## Driven
引入Mongodb的驱动
```golang
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
```
## Connect
连接本地的mongodb数据库
> mongodb://[username:password@]host1[:port1][/[database][?options]]

```golang
ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()
//client, err := mongo.NewClient((options.Client().ApplyURI("mongodb://localhost:27017")))
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://clay:password@localhost:27017"))
if err != nil { return err }
```
## Database Selected
```golang
collection := client.Database("baz").Collection("qux")
```

## Insert
> db.COLLECTION_NAME.insert(document)

根据  COLLECTION_NAME 自动创建或插入 表格。

1. use DATABASE_NAME  // 创建或切换数据库
2.  将json数据插入到数据库
```golang
db.clayte.insertOne({"name":"clayer"})
db.clayte.insertOne({"name":"clayer","interest":"football","age":18})
db.clayte.insertOne({"interest":"baseball","age":18})
```

```golang
res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
if err != nil { return err }
id := res.InsertedID
```

## Update
```golang
db.collection.update(
   <query>, // update的查询条件
   <update>, // update的对象。当前插入或更新的内容
   {
     upsert: <boolean>,  //在不存在query条件时，是否插入 update 的数据记录
     multi: <boolean>, // 默认只更新找到的第一条记录，为true时 就会把全部数据进行更新
     writeConcern: <document>
   }
)
```

Example
```golang
db.col.update({'title':'MongoDB 教程'},{$set:{'title':'MongoDB'}},{upsert:true,multi:true});
```

## Query *
查询一条记录的所有数据

> db.col.find({key1:value1, key2:value2}).pretty() // AND
> db.col.find({$or:[{key1:value1, key2:value2}]}).pretty() // OR
> db.col.find({"likes": {$gt:50}, $or: [{"by": "菜鸟教程"},{"title": "MongoDB 教程"}]}).pretty() // AND + OR

## Query xx
查询特定列数数据：1表示显示，0表示不显示，默认0

sort排序，1表示正序， -1表示逆序，从大到小排列。
> db.col.find({},{"title":1,_id:0}).sort({"likes":-1}).limit(1) // 限制一条数据

```golang
cur, err := collection.Find(context.Background(), bson.D{})
if err != nil { log.Fatal(err) }
defer cur.Close(context.Background())
for cur.Next(context.Background()) {
  // To decode into a struct, use cursor.Decode()
  result := struct{
    Foo string
    Bar int32
  }{}
  err := cur.Decode(&result)
  if err != nil { log.Fatal(err) }
  // do something with result...

  // To get the raw bson bytes use cursor.Current
  raw := cur.Current
  // do something with raw...
}
if err := cur.Err(); err != nil {
  return err
}
```
一次性解析所有的数据
```golang
var results []struct{
  Foo string
  Bar int32
}
if err = cur.All(context.Background(), &results); err != nil {
  log.Fatal(err)
}
// do something with results...
```
查询单条数据
```golang
result := struct{
  Foo string
  Bar int32
}{}
filter := bson.D{{"hello", "world"}}
err := collection.FindOne(context.Background(), filter).Decode(&result)
if err != nil { return err }
// do something with result...
```