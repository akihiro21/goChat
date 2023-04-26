# goChat
研究で使っていた実験用アプリをpublicにできるように対策したリポジトリ
# 概要
- 人間同士のチャットや、対話システムを模した実験をするためのアプリ
- 特定の実験のため、csvでシナリオを追加している
# 実行
ルート直下に.envファイルを追加。以下の環境変数について記述   
  
SESSION_AUTHENTICATION_KEY  
SESSION_ENCRYPTION_KEY  
TZ   
MYSQL_ROOT_PASSWORD  
MYSQL_DATABASE  
MYSQL_USER  
MYSQL_PASSWORD  

# 問題
- 実験に間に合わせるためやっつけな部分が多数(csvの読み込みとか)
- なんちゃってアーキテクチャ（ちゃんと設計学んでないので変）
- 要リファクタリング（goやクリーンコードをちゃんと学んでいなくて書いた部分がたくさんあるので色々ひどい）
- interfaceを無駄な使い方してる。
