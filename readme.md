# LineBot

LineBot + Wit.AI + Wolfram = My Line bot

- 透過 Wit AI 分析語意
- 透過 Wolfram 自動問答

## Quick start

申請好

1. [Line Message API](https://developers.line.biz/console/)
   1. Channel access token
   2. Channel secret
   3. webhooks on, 設定 callback url, 當 line bot 收到訊息時知道要往那轉發
2. [Wit.ai](https://wit.ai/) token
3. [Wolfram](https://developer.wolframalpha.com/) AppID


running

```sh
PORT=8000 ChannelSecret={line secret} ChannelAccessToken={Channel access token} WitToken={wit token} WolframID={wolfram AppID} ./linebot
```


## Friend this bot

QRCODE: ![](image/qrcode.png)

