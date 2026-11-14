/-
Copyright 2026 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-/

import Aip4222.Syntax

namespace Aip4222

/--
RFC 6570 §1.5 unreserved character set:
`unreserved = ALPHA / DIGIT / "-" / "." / "_" / "~"`.
-/
def isUnreservedByte (b : UInt8) : Bool :=
  (b >= 65 && b <= 90) ||
  (b >= 97 && b <= 122) ||
  (b >= 48 && b <= 57) ||
  b == 45 || b == 46 || b == 95 || b == 126

def hexDigitUpper (n : UInt8) : Char :=
  let v := n.toNat % 16
  if v < 10 then Char.ofNat (48 + v) else Char.ofNat (55 + v)

def percentEncodeByte (b : UInt8) : List Char :=
  if isUnreservedByte b then
    [Char.ofNat b.toNat]
  else
    ['%', hexDigitUpper (b >>> 4), hexDigitUpper (b &&& 0x0F)]

/-- Kernel-reducible UTF-8 encoding of a Lean `String`. -/
def utf8Bytes (s : String) : List UInt8 :=
  s.toList.flatMap String.utf8EncodeChar

/-- RFC 6570 §3.2.2 percent-encoding over UTF-8 bytes. -/
def percentEncode (s : String) : String :=
  String.ofList ((utf8Bytes s).flatMap percentEncodeByte)

def formatKeyValue (k v : String) : String :=
  percentEncode k ++ "=" ++ percentEncode v

/-- Join percent-encoded `key=value` routing header entries with `&`. -/
def formatHeaderPairs (pairs : List (String × String)) : String :=
  joinStringsWith "&" (pairs.map fun (k, v) => formatKeyValue k v)

end Aip4222
