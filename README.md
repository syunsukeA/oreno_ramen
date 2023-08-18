# oreno_ramen

## 実行方法
```
git clone https://github.com/syunsukeA/oreno_ramen.git
cd oreno_ramen
docker compose up --build
```
上記コマンドを実行しコンテナが立ち上がると下記提供機能を使用可能。
## 提供機能
### App
http://127.0.0.1:8080 にアクセスすることで接続可能。
俺のラーメンAPIを提供。エンドポイントは下記のSwagger UIを参照。

### Swagger UI
http://127.0.0.1:8081 にアクセスすることで接続可能。
APIのエンドポイント詳細について記載。
該当エンドポイントをクリックし、"Try it out"を押すことで実際にHTTPリクエストを送ることが可能。
認証が必要なエンドポイントには右上に鍵マークが付いており、そこから認証情報を記載する必要あり。（詳細は決定次第記載）
#### Basic認証
Valueに下記の情報を入力
Basic <”username:pass”をbase64符号化した文字列>
(usernameとpassの区切りは”:”でしきり、その仕切りを入れたままbase64符号化する)

### Air
ホットリロードツールのサービス名。
コードを変更した際は自動で再実行がなされるようになっており、再度buildする必要がない。

## Tips
### DBの変更反映
Docker volumeでローカルでDBデータを永続化させている。手元にDBデータがある場合はinit.sqlが実行されないので変更を反映したい場合には /.data 配下の mysqlを削除する必要がある。