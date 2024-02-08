# X-Line-Signatureの署名検証
GolangでのLINEBotの署名検証の実装例です。

## 概要
LINEBotのWebhookで受け取ったリクエストの署名検証を行うための実装例です。

## 必要なもの
- LINEBotのチャンネルシークレット

## 流れ
1. リクエストヘッダーからX-Line-Signatureを取得
2. リクエストボディを取得
3. チャネルシークレットを秘密鍵として、HMAC-SHA256アルゴリズムを使用してリクエストボディのダイジェスト値を取得
4. ダイジェスト値をBase64エンコード
5. エンコードした署名とX-Line-Signatureを比較
