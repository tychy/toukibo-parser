# toukibo-parser
toukibo-parserは商業登記簿のPDFを解析し、法人情報を抽出するツールです。PDFファイルから企業名、住所、役員名、代表者名、設立日などの情報を取得できます。

## 動作イメージ
```
 % make run TARGET=<PDFのパス>
mkdir -p bin
go build -o bin/toukibo-parser main.go
./bin/toukibo-parser -path="testdata/pdf/sample1.pdf"
HoujinKaku: 株式会社
HoujinName: 株式会社近畿商事
HoujinAddress: 東京都Sample区Sample１丁目１番地１
HoujinExecutiveValues: 
  - Name: 大門十三
    Position: 取締役
  - Name: 壹岐正
    Position: 取締役
  - Name: 壹岐正
    Position: 代表取締役
HoujinExecutiveNames: [大門十三,壹岐正]
HoujinRepresentativeNames: [壹岐正]
HoujinCapital:  1000000
HoujinStock: 25642
HoujinPreferredStock: 
  - Type: 普通株式
    Amount: 13106
  - Type: A種種類株式
    Amount: 2332
  - Type: AA種種類株式
    Amount: 2221
  - Type: B種種類株式
    Amount: 3468
  - Type: BB種種類株式
    Amount: 4515
HoujinCreatedAt: 令和3年1月6日
HoujinBankruptedAt: 
HoujinDissolvedAt: 
HoujinContinuedAt: 
```

#### テスト実行
```
% make get/sample
% make test
```

## デモ
[こちら](https://toukibo-parser-demo.tychy.jp/)のページでパーサーの性能をお試しいただくことが可能です。

![デモ](https://github.com/tychy/toukibo-parser/blob/main/misc/parser-demo.jpg)

## 商業登記簿API
このパーサーを利用した商業登記簿APIは[こちら](https://toukikun-docs.tychy.jp/)

## ライセンス
このソフトウェアを利用するには、ライセンス契約が必要です。
このソフトウェアを利用した登記簿に関わるプロダクトを販売する場合は月50万円+税、それ以外の場合には月20万円+税となります。
詳細はcontact@tychy.jpまでお問い合わせください。

```
Copyright 2023 - 2024, Ryu Yamada.

(JA) ソースコードおよびバイナリ形式の改変および再配布、許諾なき商用利用を禁止します。
(EN) We prohibit the modification and redistribution of the source code
and binary format, as well as any commercial use without permission.

---

This repository contains a code from github.com/rsc/pdf.

---
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
## 免責事項
サンプルとして使用される登記簿の個人情報は公開情報です。したがって、以下の点についてご理解ください：

1. 情報の正確性と利用目的について、一切の保証をいたしません。
2. 個人情報の保護に関する法的責任を負いません。
3. 削除要請には速やかに対応いたします。削除を希望する場合は、contact@tychy.jpまでご連絡ください。
