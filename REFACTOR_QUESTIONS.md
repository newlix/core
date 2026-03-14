## [type.go / method.go] `log.Fatalf` in `InitTypes` / `InitMethods`
- **問題描述**：`InitTypes` 和 `InitMethods` 在遇到重複名稱時呼叫 `log.Fatalf`，直接終止程式。作為 library，呼叫端無法 recover；也無法測試錯誤路徑。
- **可能的做法**：
  A. 改為回傳 `error`（`func InitTypes(...) ([]Type, error)`）— 最正確但是 public API breaking change，所有呼叫端需更新
  B. 改用 `panic` — 至少可以 `recover`，適合「programmer error」語意
  C. 維持現狀 — 這些只在啟動時呼叫，且重複名稱確實是程式設計錯誤
- **目前狀態**：未處理（跳過）

## [field.go] `initFields` 為每個 field 呼叫 `InitTypes` 初始化單一 type
- **問題描述**：`initFields` 中 `f.Type = InitTypes(f.Type)[0]` 對單一 type 執行了不必要的排序與重複檢查。且這造成遞迴呼叫（`InitTypes` → `initFields` → `InitTypes`），雖然目前結構上不會無限遞迴，但語意不清。
- **可能的做法**：
  A. 抽取 `initType(t Type) Type` 單一型別初始化函式，只做 camel name 設定與 field 初始化
  B. 維持現狀 — 效能影響微乎其微，且行為正確
- **目前狀態**：未處理（跳過）

## [generators/sqlc] 測試讀寫同一個 `.gen.sql` 檔案
- **問題描述**：`query_test.go` 和 `schema_test.go` 從 `testdata/todo_queries.gen.sql` / `testdata/todo_schema.gen.sql` 讀取預期輸出，同時又覆寫同檔案。第二次執行測試必定通過，無法偵測 regression。其他 generator（kotlin、swift、golang）都有獨立的 golden file。
- **可能的做法**：
  A. 建立獨立的 golden file（`todo_queries.sql` / `todo_schema.sql`），與 `.gen.sql` 分開
  B. 維持現狀 — sqlc 的 golden file 本身就是 generated，只要 CI 跑在 clean checkout 就不會有問題
- **目前狀態**：未處理（跳過）

## [generators/golang/type.go] `FieldGoType` 與 MarshalJSON 不支援 `[]time.Time`
- **問題描述**：`FieldGoType` 在 `int64Time && GoType == "*time.Time"` 時直接回傳 `"int64"`，忽略 `IsArray`。若某個 field 是 `[]time.Time`，alias struct 會產生 `int64` 而非 `[]int64`，且 MarshalJSON/UnmarshalJSON 的逐欄位邏輯也不支援陣列。目前沒有 spec 使用此組合，但這是潛在的正確性問題。
- **可能的做法**：
  A. 在 `FieldGoType` 中加入 `IsArray` 判斷，並在 marshal/unmarshal 產生迴圈邏輯
  B. 明確文件化不支援 `[]time.Time`，在 `GenerateJSONMarshalerIfNeeded` 加入 guard/panic
- **目前狀態**：未處理（跳過）

## [generators/golang/type.go] UnmarshalJSON 用 0 作為 nil sentinel，無法區分 epoch 與未設定
- **問題描述**：生成的 `UnmarshalJSON` 用 `if tmp.CreatedAt != 0` 判斷是否設定 `*time.Time`。Unix timestamp 0（1970-01-01T00:00:00Z）會被當作「未設定」而保持 nil。MarshalJSON 也將 nil 編碼為 0，所以 round-trip 對稱，但無法表示 epoch zero。
- **可能的做法**：
  A. 改用 `*int64`（pointer）作為 alias field type，JSON `null` 代表未設定，`0` 代表 epoch
  B. 文件化 epoch zero 不可表示（最小侵入）
  C. 維持現狀 — 實務上 epoch zero 極少出現
- **目前狀態**：未處理（跳過）
