package dto

type LangNew struct {
  Locale string `json:"lg" validate:"required"`
  Group string `json:"group" validate:"required"`
  Key string `json:"key" validate:"required"`
  Val string `json:"val" validate:"required"`
}

type LangDel struct {
  Version int `json:"version" validate:"required"`
  ForceDel bool `json:"force_del" default:"false"`
}

type LangEdit struct {
  Locale string `json:"lg" validate:"required"`
  Group string `json:"group" validate:"required"`
  Key string `json:"key" validate:"required"`
  Val string `json:"val" validate:"required"`
  Version int `json:"version" validate:"required"`
}
