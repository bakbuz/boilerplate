package text

import (
	"github.com/gosimple/slug"
)

func Slugify(input string) string {
	// Özel ayar gerekirse (örneğin Türkçe için 'ı' -> 'i' garantilemek isterseniz)
	slug.CustomSub = map[string]string{
		"ı": "i",
		"İ": "i",
	}
	return slug.Make(input)
}
