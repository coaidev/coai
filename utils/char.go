package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goccy/go-json"
)

func GetRandomInt(min int, max int) int {
	return Intn(max-min) + min
}

func GenerateCode(length int) string {
	seq := IntnSeq(10, length)

	var code string
	for i := 0; i < length; i++ {
		code += strconv.Itoa(seq[i])
	}
	return code
}

func GenerateChar(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seq := IntnSeq(len(charset), length)

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[seq[i]]
	}
	return string(result)
}

// deprecated: use sql.NullTime directly
func ConvertTime(t []uint8) *time.Time {
	val, err := time.Parse("2006-01-02 15:04:05", string(t))
	if err != nil {
		return nil
	}
	return &val
}

func Unmarshal[T interface{}](data []byte) (form T, err error) {
	err = json.Unmarshal(data, &form)
	return form, err
}

func UnmarshalString[T interface{}](data string) (form T, err error) {
	err = json.Unmarshal([]byte(data), &form)
	return form, err
}

func UnmarshalForm[T interface{}](data string) *T {
	form, err := UnmarshalString[T](data)
	if err != nil {
		return nil
	}
	return &form
}

func Marshal[T interface{}](data T) string {
	res, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(res)
}

func MarshalWithIndent[T interface{}](data T, length ...int) string {
	var indent string
	if len(length) > 0 {
		indent = strings.Repeat(" ", length[0])
	} else {
		indent = "  "
	}

	res, err := json.MarshalIndent(data, "", indent)
	if err != nil {
		return ""
	}
	return string(res)
}

func MapToStruct[T any](data interface{}) *T {
	val := Marshal(data)
	if form, err := Unmarshal[T]([]byte(val)); err == nil {
		return &form
	} else {
		return nil
	}
}

func MapToRawStruct[T any](data interface{}) (*T, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	form, err := Unmarshal[T](val)
	if err != nil {
		return nil, err
	}

	return &form, nil
}

func ParseInt(value string) int {
	if res, err := strconv.Atoi(value); err == nil {
		return res
	} else {
		return 0
	}
}

func ParseInt64(value string) int64 {
	if res, err := strconv.ParseInt(value, 10, 64); err == nil {
		return res
	} else {
		return 0
	}
}

func ParseFloat32(value string) float32 {
	if res, err := strconv.ParseFloat(value, 32); err == nil {
		return float32(res)
	} else {
		return 0
	}
}

func ParseBool(value string) bool {
	if res, err := strconv.ParseBool(value); err == nil {
		return res
	} else {
		return false
	}
}

func ConvertSqlTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetImageMarkdown(url string) string {
	return fmt.Sprintf("![image](%s)", url)
}

func GetBase64ImageMarkdown(b64 string, _desc ...string) string {
	// Extracts the image type from base64 string (e.g., "data:image/png;base64,...") or defaults to png
	var imageType = "png"
	if strings.HasPrefix(b64, "data:image/") {
		parts := strings.Split(b64[11:], ";")
		if len(parts) > 0 {
			imageType = parts[0]
		}
	}

	desc := "image"
	if len(_desc) > 0 && _desc[0] != "" {
		desc = _desc[0]
	}

	return fmt.Sprintf("![%s](data:image/%s;base64,%s)", desc, imageType, b64)
}

// SplitItem is the split function for strings.Split
// e.g.
// SplitItem("a,b,c", ",") => ["a,", "b,", "c"]
func SplitItem(data string, sep string) []string {
	if data == "" {
		return []string{}
	}

	result := strings.Split(data, sep)
	length := len(result)
	for i, item := range result {
		if i == length-1 {
			break
		}
		result[i] = item + sep
	}
	return result
}

func SplitItems(data string, seps []string) []string {
	if len(seps) == 0 {
		return []string{}
	}

	result := []string{data}
	for _, sep := range seps {
		var temp []string
		for _, item := range result {
			temp = append(temp, SplitItem(item, sep)...)
		}
		result = temp
	}
	return result
}

func SplitLangItems(data string) []string {
	return SplitItems(data, []string{",", "，", " ", "\n"})
}

func Extract(data string, length int, flow_ ...string) string {
	flow := ""
	if len(flow_) > 0 {
		flow = flow_[0]
	}

	value := []rune(data)
	if len(value) > length {
		return string(value[:length]) + flow
	} else {
		return data
	}
}

func ExtractUrls(data string) []string {
	re := regexp.MustCompile(`(https?://\S+)`)
	return re.FindAllString(data, -1)
}

func ExtractImages(data string, includeBase64 bool) (content string, images []string) {
	ext := ExtractExternalImages(data)
	if includeBase64 {
		images = append(ext, ExtractBase64Images(data)...)
	} else {
		images = ext
	}

	content = data
	for _, image := range images {
		content = strings.ReplaceAll(content, image, "")
	}

	return content, images
}

func ExtractImagesFromMarkdown(data string) (images []string) {
	// extract images like `![image](https://xxx.com/xxx?xxx=xxx&xxx=xxx)` and return urls
	re := regexp.MustCompile(`!\[.*\]\((https?://\S+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		images = append(images, match[1])
	}

	return images
}

func ExtractBase64FromMarkdown(data string) (images []string) {
	// extract base64 images like `![image](data:image/png;base64,xxxxxx)`
	re := regexp.MustCompile(`!\[.*?\]\((data:image/\w+;base64,[\w+/=]+)\)`)
	matches := re.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		// We only need the base64 data part
		if len(match) > 1 {
			images = append(images, match[1])
		}
	}

	return images
}

func ExtractBase64Images(data string) []string {
	// get base64 images from data (data:image/png;base64,xxxxxx) (\n \\n [space] \\t \\r \\v \\f break the base64 string)
	re := regexp.MustCompile(`(data:image/\w+;base64,[\w+/=]+)`)
	return re.FindAllString(data, -1)
}

func ExtractExternalImages(data string) []string {
	// https://platform.openai.com/docs/guides/vision/what-type-of-files-can-i-upload

	re := regexp.MustCompile(`(https?://\S+\.(?:png|jpg|jpeg|gif|webp|heif|heic|bmp|svg|ico)(?:\?\S+)?)`)
	return re.FindAllString(data, -1)
}

func ContainUnicode(data string) bool {
	// like `hi\\u2019s` => true
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	return re.MatchString(data)
}

func DecodeUnicode(data string) string {
	// like `hi\\u2019s` => `hi's`
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	return re.ReplaceAllStringFunc(data, func(s string) string {
		unicode, err := strconv.ParseInt(s[2:], 16, 32)
		if err != nil {
			return s
		}

		return string(rune(unicode))
	})
}

func EscapeChar(data string) string {
	// like `\\n` => `\n`, `\\t` => `\t` and so on
	re := regexp.MustCompile(`\\([nrtvfb])`)

	mapper := map[string]string{
		"n": "\n",
		"r": "\r",
		"t": "\t",
		"v": "\v",
		"f": "\f",
		"b": "\b",
	}

	return re.ReplaceAllStringFunc(data, func(s string) string {
		return mapper[s[1:]]
	})
}

func ProcessRobustnessChar(data string) string {
	// like `hi\\u2019s` => `hi's`
	if ContainUnicode(data) {
		data = DecodeUnicode(data)
	}

	// like `\\n` => `\n`, `\\t` => `\t` and so on
	return EscapeChar(data)
}

func SortString(arr []string) []string {
	// sort string array by first char
	// e.g. ["a", "b", "c", "ab", "ac", "bc"] => ["a", "ab", "ac", "b", "bc", "c"]

	if len(arr) <= 1 {
		return arr
	}

	var result []string
	var temp []string
	var first string

	for _, item := range arr {
		if first == "" {
			first = item
			continue
		}

		if strings.HasPrefix(item, first) {
			temp = append(temp, item)
		} else {
			result = append(result, first)
			result = append(result, SortString(temp)...)
			first = item
			temp = []string{}
		}
	}

	if len(temp) > 0 {
		result = append(result, first)
		result = append(result, SortString(temp)...)
	} else {
		result = append(result, first)
	}

	return result
}

func SafeSplit(data string, sep string, seglen int) (res []string) {
	// split string by sep, and each segment has seglen length
	// e.g. SafeSplit("abc,def,ghi", ",", 2) => ["abc", "def,ghi"]

	if data == "" {
		for i := 0; i < seglen; i++ {
			res = append(res, "")
		}
		return res
	}

	arr := strings.Split(data, sep)
	length := len(arr)

	if length == seglen {
		return arr
	}

	if length < seglen {
		for i := 0; i < seglen-length; i++ {
			arr = append(arr, "")
		}
		return arr
	} else {
		offset := length - seglen
		for i := 0; i < offset; i++ {
			arr[seglen-1] += sep + arr[seglen+i]
		}
		return arr[:seglen]
	}
}

func ToSecret(raw string) string {
	// like `axVbeixvN` => `axVb*****`

	data := []rune(raw)
	length := len(data)

	if length < 4 {
		return "****"
	} else {
		suffix := len(data) - 4
		return string(data[:4]) + strings.Repeat("*", suffix)
	}
}

func ToMarkdownCode(lang string, code string) string {
	return fmt.Sprintf("```%s\n%s\n```", lang, code)
}

func ToMarkdownError(err error, body string) error {
	return fmt.Errorf("%s\n%s", err.Error(), ToMarkdownCode("html", body))
}
