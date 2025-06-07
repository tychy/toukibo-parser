# CLAUDE.md

このファイルは、このリポジトリで作業する際のClaude Code (claude.ai/code) への指針を提供します。

## プロジェクト概要

toukibo-parserは日本の商業登記簿PDFから構造化された法人情報を抽出するツールです。会社情報、役員、資本金、株式情報、重要な日付などを公式登記書類から解析します。

## よく使うコマンド

```bash
# テストデータのダウンロード（テスト実行に必要）
make get/sample

# テスト実行
make test

# パーサーツールのビルド
make build

# PDFファイルの解析実行
make run TARGET=path/to/file.pdf

# サンプルファイルで実行
make run/sample TARGET=sample1

# PDF内のテキスト検索
make find/sample TARGET=sample1 FIND="検索文字列"

# ベンチマーク実行
make bench
```

## アーキテクチャ概要

コードベースは階層化されたアーキテクチャに従っています：

1. **パブリックAPI** (`toukibo_parser.go`): PDF解析のメインエントリーポイント
   - `ParseByPDFPath()`, `ParseByPDFReader()`, `ParseByPDFRawData()`
   - 解析済みデータを含む `*toukibo.Houjin` 構造体を返す

2. **PDFライブラリ** (`internal/pdf/`): 低レベルPDF解析
   - Adobe PDF仕様に基づくカスタム実装
   - PDF構造からテキストコンテンツを抽出

3. **ビジネスロジック** (`internal/toukibo/`): 登記簿固有の解析処理
   - `parse.go`: コンテンツをセクションに分割するメインオーケストレーター
   - `parse_header.go` & `parse_body.go`: セクション別の解析
   - 日本語フォーマット用の特殊パーサー：
     - 通貨変換（円 → 数値）
     - 日付変換（令和 → 標準日付）
     - 株式表記の解析
     - 全角文字の正規化

## 主要なデータ構造

- `Houjin`: ヘッダーとボディを含むルートオブジェクト
- `HoujinHeader`: 文書メタデータ（作成日、基本的な会社情報）
- `HoujinBody`: 詳細な法人情報
- `HoujinExecutiveValue`: 役員・取締役情報
- `HoujinStock`: 優先株式タイプを含む株式構造

## テスト戦略

テストはR2（Cloudflareストレージ）に保存された実際のPDFを使用します。各PDFには期待される出力を含むYAMLファイルが対応しています。テストスイートは1522のサンプル文書で解析精度を検証します。

## 重要なコンテキスト

- 商業登記簿PDFは罫線文字を区切りとした標準化されたレイアウトに従っている
- パーサーはこれらの一貫したパターンを利用して確実な抽出を行う
- 全角日本語文字は処理前に正規化が必要
- 本番環境での使用には商用ライセンスが必要