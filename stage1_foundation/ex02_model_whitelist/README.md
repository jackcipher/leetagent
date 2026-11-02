# ex02 模型白名单 + 默认值（Easy）

## 题目
前端可以传任意 `model` 字符串，直接透传给上游是安全隐患（用户可能调到很贵或不存在的模型）。
请实现：从环境变量加载白名单，并对请求里的 model 做校验。

## 接口
```go
type Env interface{ Get(key string) string }
func LoadModelConfig(env Env) (ModelConfig, error)
func (c ModelConfig) Resolve(requested string) (string, error)
```

## 要求
- `AGENT_MODELS`：逗号分隔，去首尾空格、忽略空项、去重并保持首次出现顺序。列表为空 → `ErrNoModels`。
- `AGENT_DEFAULT_MODEL`：为空则用列表第一个；不在列表里 → `ErrBadDefault`。
- `Resolve`：空串（含纯空格）→ 默认模型；在白名单 → 原样；否则 `ErrModelNotAllowed`（errors.Is 可识别）。

## 验收
`go test ./stage1_foundation/ex02_model_whitelist/`

## 思考
为什么用 `Env` 接口而不是直接调 `os.Getenv`？（提示：测试、以及"配置来源可以换成文件/配置中心"。）
生产里可以写 `type OSEnv struct{}; func (OSEnv) Get(k string) string { return os.Getenv(k) }`。
