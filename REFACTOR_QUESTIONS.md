## [type.go / method.go] `log.Fatalf` in `InitTypes` / `InitMethods`
- **問題描述**：`InitTypes` 和 `InitMethods` 在遇到重複名稱時呼叫 `log.Fatalf`，直接終止程式。作為 library，呼叫端無法 recover；也無法測試錯誤路徑。
- **可能的做法**：
  A. 改為回傳 `error`（`func InitTypes(...) ([]Type, error)`）— 最正確但是 public API breaking change，所有呼叫端需更新
  B. 改用 `panic` — 至少可以 `recover`，適合「programmer error」語意
  C. 維持現狀 — 這些只在啟動時呼叫，且重複名稱確實是程式設計錯誤
- **目前狀態**：未處理（跳過）

## [generators/golang/type.go] UnmarshalJSON 用 0 作為 nil sentinel，無法區分 epoch 與未設定
- **問題描述**：生成的 `UnmarshalJSON` 用 `if tmp.CreatedAt != 0` 判斷是否設定 `*time.Time`。Unix timestamp 0（1970-01-01T00:00:00Z）會被當作「未設定」而保持 nil。MarshalJSON 也將 nil 編碼為 0，所以 round-trip 對稱，但無法表示 epoch zero。
- **可能的做法**：
  A. 改用 `*int64`（pointer）作為 alias field type，JSON `null` 代表未設定，`0` 代表 epoch
  B. 文件化 epoch zero 不可表示（最小侵入）
  C. 維持現狀 — 實務上 epoch zero 極少出現
- **目前狀態**：未處理（跳過）
