# toukibo_parser
商業登記簿パーサー

## 動作確認
```
 % make run TARGET=sample2
go run main.go -path=sample2
Header
PDF作成日時: 2021-03-29 17:58:00 +0000 UTC
商号: ****株式会社  
本店住所: ****

Body
法人番号 : 0500-01-*****
法人名  : {value: ****, isValid: true, registerAt: },
法人住所 : {value: ****, isValid: true, registerAt: },
公告   : 官報に掲載してする
成立年月日: 昭和**年*月*日
資本金  : {value: 金****万円, isValid: true, registerAt: },
登記記録 : {value: 平成元年法務省令第１５号附則第３項の規定により, isValid: true, registerAt: 平成**年**月**日},
役員   : [{name: ****, position: 取締役, address: , isValid: true, registerAt: 平成**年**月**日, resignedAt: },],
[{name: ****, position: 代表取締役, address: ****, isValid: false, registerAt: 平成**年**月**日, resignedAt: },{name: ****, position: 代表取締役, address: ****, isValid: true, registerAt: 平成**年**月**日, resignedAt: },],

```

## ライセンス
このソフトウェアを利用するには、。ライセンス契約が必要です。詳細はcontact@tychy.jpまでお問い合わせください。

```
Copyright 2023, Ryu Yamada.

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
