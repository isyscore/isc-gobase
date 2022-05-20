package test

import (
	"testing"

	. "github.com/isyscore/isc-gobase/i18n"
)

func TestI18N(t *testing.T) {
	_ = InitI18N("zh-CN")
	_ = LoadI18NLanguage("en-US")

	t.Logf(T("msgid"))
	t.Logf(T("msgstr"))
	t.Logf(T("msgfmt"))
	t.Logf(T("msgcn"))
	t.Logf(Tf("msgf2", "rarnu", 2333))

}
