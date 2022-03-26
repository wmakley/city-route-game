package admin

import (
	"github.com/gorilla/schema"
)

type Controller struct {
	FormDecoder *schema.Decoder
	TemplateRoot string
	AssetHost    string

}

type ControllerConfig struct {
	FormDecoder *schema.Decoder
	TemplateRoot string
	AssetHost    string
}
