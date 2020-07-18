# ping-tool-go

## 設定ファイルのフォーマットについて

下記の形式で追加する。

```
プラント名, IPアドレス
プラント名, IPアドレス
プラント名, IPアドレス
```


## 取りたい情報

- 疎通確認
- Latency
- Packet Loss

## どのくらいの頻度で

オンプレだったら5分に一回くらい情報を収集してもよいか。仮にDBに保存しようとしたときのデータ量はどのくらいか。

5分に一回 = 1日288回
対象が200程度だとすると、

1日57600レコード
1年で2102400レコード程度

## Flow

1. csvからデータ読み込み
2. 各行に対してPING
3. エラーが起きたら通知
   1. Email
   2. Webhook -> chatbot

- 実行速度はそれほど大事ではないので平行実行に対応する必要は基本的にはないが、勉強のためGoroutineを使用して実装する。

- 各拠点に対するPingの結果をDBなどに保存しておくとあとでグラフ化して通信状況を可視化などできて面白そう。Fluentd, Grafara, Prometheusあたりの勉強の題材に使えそう。

- ガッツリDBを使うと運用面で面倒になりそうなのでsqliteあたりで軽く済ませたい。AWSがOKならDynamoDBなどが良さそう。

## 運用方法

実行マシン
 - Synology上に実装
 - AWSのLambda

1. Cronで定期的に実行
2. 常駐型にする

## format

- [ ] PacketRecv : int
- [ ] PacketSent : int
- [ ] PacketLoss : float64
- [ ] IPAddr     : *net.IPAddr
- [ ] Addr       : string
- [ ] Rtts       : []time.duration
- [x] MinRtt     : time.duration
- [x] MaxRtt     : time.duration
- [x] AvgRtt     : time.duration
- [x] StdDevRtt  : time.duration