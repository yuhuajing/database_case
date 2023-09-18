# MongoDB
## Docker start
> docker run -itd --name mongo-test -p 27017:27017 mongo:6.0

## Driven
引入Mongodb的驱动
```golang
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
```
## Connect
连接本地的mongodb数据库
```golang
ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
defer cancel()
//client, err := mongo.NewClient((options.Client().ApplyURI("mongodb://localhost:27017")))
client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
if err != nil { return err }
```
## Database Selected
```golang
collection := client.Database("baz").Collection("qux")
```
## Insert
```golang
res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
if err != nil { return err }
id := res.InsertedID
```
## Query
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