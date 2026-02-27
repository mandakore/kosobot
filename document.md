1. 概要
Raspberry Pi Zero 2 W で 24 時間運用する、身内サーバー（約 10 人）向けの Discord ボット。ユーザーの GitHub での活動や AtCoder への参加、Discord 内での交流を「経験値（XP）」として数値化し、自動的にレベル（ロール）を付与する。

2. 技術スタック（指定）
言語: Go (Golang)

Discord ライブラリ: github.com/bwmarrin/discordgo

データベース: SQLite3 (CGO 不要な modernc.org/sqlite を推奨)

定期実行: github.com/robfig/cron/v3

実行環境: Raspberry Pi Zero 2 W (Linux / ARM64)

3. 主要機能
3.1 経験値 (XP) システム
ユーザーごとに累計 XP を保持する。

XP に応じて「レベル」を計算し、特定のレベルに達したら Discord の「ロール」を自動的に付与・更新する。

3.2 GitHub 連動 (自動)
仕様: 12時間ごとに全ユーザーの GitHub コントリビューション数を確認。

判定基準: 「前回確認時より総コントリビューション数が増えているか（＝今日、草が生えたか）」を基準とする。

報酬: 草が生えていれば一律で一定の XP を付与する（回数ではなく、継続を重視）。

API: GitHub GraphQL API (v4) を使用し、totalContributions を取得する。

3.3 AtCoder 連動 (手動)
仕様: スラッシュコマンド /contest を実装。

入力項目: ユーザーが自身の「パフォーマンス」または「成績」を数値で入力。

報酬: 入力された数値に基づき、計算式（例: パフォーマンス/10）に従って XP を付与する。

3.4 リアクション加算 (自動)
仕様: 他のユーザーからのリアクション（絵文字）を検知。

報酬: メッセージの投稿者に XP を付与する。

制限: 自分のメッセージに対する自己リアクションは除外する。

4. データベース設計 (SQLite)
以下のカラムを持つ users テーブルを作成すること。

discord_id (TEXT, PK): ユーザーの Discord ID

github_id (TEXT): ユーザーの GitHub ユーザー名

current_xp (INTEGER): 現在の累計 XP

current_level (INTEGER): 現在のレベル

last_github_contributions (INTEGER): 前回確認時の草の総数

updated_at (DATETIME): 最終更新日時

5. 実装フェーズ（AI へのステップ指示）
AI には以下の順番でコードを書かせてください。

フェーズ 1: Go のプロジェクト初期化と、Discord への接続、SQLite のテーブル作成コードの作成。

フェーズ 2: /set-github コマンドの実装（ユーザーの Discord ID と GitHub ID を紐付けて DB に保存）。

フェーズ 3: GitHub GraphQL API を叩き、草の有無を判定して XP を加算する定期実行ロジックの作成。

フェーズ 4: /contest コマンドの実装と、リアクション検知（MessageReactionAdd）による XP 加算ロジックの作成。

フェーズ 5: XP が増えるたびにレベル判定を行い、Discord ロールを付与・更新する関数の作成。

6. 非機能要件
省メモリ: Raspberry Pi Zero 2 W のメモリ (512MB) を圧迫しないよう、キャッシュは最小限にする。

耐障害性: ボットが再起動してもデータが消えないよう、XP 加算のたびに DB をコミットする。

環境変数: トークン類（DISCORD_TOKEN, GITHUB_TOKEN）は .env ファイルまたは環境変数から読み込む形式にすること。