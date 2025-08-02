package mox

import (
	"bytes"
	"fmt"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"
)

// generateQueryString uses reflection to generate URL query string from struct
func generateQueryString(obj interface{}) string {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	var params []string

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		formTag := field.Tag.Get("form")
		if formTag == "" {
			continue
		}

		// Extract the form name (first part before comma)
		formName := strings.Split(formTag, ",")[0]

		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				value := formatValue(elem)
				if value != "" {
					params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
				}
			}
		case reflect.Ptr:
			if !fieldValue.IsNil() {
				value := formatValue(fieldValue.Elem())
				if value != "" {
					params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
				}
			}
		default:
			// Handle Option types
			if field.Type.PkgPath() == "github.com/samber/mo" && strings.HasPrefix(field.Type.Name(), "Option[") {
				// Check if Option has value
				isPresent := fieldValue.FieldByName("isPresent")
				if isPresent.IsValid() && isPresent.Bool() {
					optionValue := fieldValue.FieldByName("value")
					if optionValue.IsValid() {
						// Handle different Option value types
						switch optionValue.Kind() {
						case reflect.Slice, reflect.Array:
							// Handle Option[[]T] and Option[[2]T]
							for j := 0; j < optionValue.Len(); j++ {
								elem := optionValue.Index(j)
								value := formatValue(elem)
								if value != "" {
									params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
								}
							}
						case reflect.Ptr:
							// Handle Option[*T]
							if !optionValue.IsNil() {
								value := formatValue(optionValue.Elem())
								if value != "" {
									params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
								}
							}
						default:
							// Handle Option[T] for basic types
							value := formatValue(optionValue)
							if value != "" {
								params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
							}
						}
					}
				}
			} else {
				value := formatValue(fieldValue)
				if value != "" {
					params = append(params, fmt.Sprintf("%s=%s", formName, url.QueryEscape(value)))
				}
			}
		}
	}

	return strings.Join(params, "&")
}

// formatValue converts reflect.Value to string
func formatValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return ""
	}
}

func TestGin(t *testing.T) {
	var err error

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Create pointer values
	ptrStr := "ptrTest"
	ptrBool := true
	ptrInt := rand.Intn(100)
	ptrInt8 := int8(rand.Intn(100))
	ptrInt16 := int16(rand.Intn(100))
	ptrInt32 := int32(rand.Intn(100))
	ptrInt64 := int64(rand.Intn(100))
	ptrUInt := uint(rand.Intn(100))
	ptrUInt8 := uint8(rand.Intn(100))
	ptrUInt16 := uint16(rand.Intn(100))
	ptrUInt32 := uint32(rand.Intn(100))
	ptrUInt64 := uint64(rand.Intn(100))
	ptrFloat32 := rand.Float32() * 100
	ptrFloat64 := rand.Float64() * 100

	v := GinDto{
		// Basic types
		Str:      "test",
		Bool:     true,
		Intv:     rand.Intn(100),
		Int8v:    int8(rand.Intn(100)),
		Int16v:   int16(rand.Intn(100)),
		Int32v:   int32(rand.Intn(100)),
		Int64v:   int64(rand.Intn(100)),
		UIntv:    uint(rand.Intn(100)),
		UInt8v:   uint8(rand.Intn(100)),
		UInt16v:  uint16(rand.Intn(100)),
		UInt32v:  uint32(rand.Intn(100)),
		UInt64v:  uint64(rand.Intn(100)),
		Float32V: rand.Float32() * 100,
		Float64V: rand.Float64() * 100,

		// Pointer types
		PtrStr:     &ptrStr,
		PtrBool:    &ptrBool,
		PtrInt:     &ptrInt,
		PtrInt8:    &ptrInt8,
		PtrInt16:   &ptrInt16,
		PtrInt32:   &ptrInt32,
		PtrInt64:   &ptrInt64,
		PtrUInt:    &ptrUInt,
		PtrUInt8:   &ptrUInt8,
		PtrUInt16:  &ptrUInt16,
		PtrUInt32:  &ptrUInt32,
		PtrUInt64:  &ptrUInt64,
		PtrFloat32: &ptrFloat32,
		PtrFloat64: &ptrFloat64,

		// Slice types
		SliceStr:     []string{"slice1", "slice2"},
		SliceBool:    []bool{true, false},
		SliceInt:     []int{rand.Intn(100), rand.Intn(100)},
		SliceInt8:    []int8{int8(rand.Intn(100)), int8(rand.Intn(100))},
		SliceInt16:   []int16{int16(rand.Intn(100)), int16(rand.Intn(100))},
		SliceInt32:   []int32{int32(rand.Intn(100)), int32(rand.Intn(100))},
		SliceInt64:   []int64{int64(rand.Intn(100)), int64(rand.Intn(100))},
		SliceUInt:    []uint{uint(rand.Intn(100)), uint(rand.Intn(100))},
		SliceUInt8:   []uint8{uint8(rand.Intn(100)), uint8(rand.Intn(100))},
		SliceUInt16:  []uint16{uint16(rand.Intn(100)), uint16(rand.Intn(100))},
		SliceUInt32:  []uint32{uint32(rand.Intn(100)), uint32(rand.Intn(100))},
		SliceUInt64:  []uint64{uint64(rand.Intn(100)), uint64(rand.Intn(100))},
		SliceFloat32: []float32{rand.Float32() * 100, rand.Float32() * 100},
		SliceFloat64: []float64{rand.Float64() * 100, rand.Float64() * 100},

		// Array types
		ArrayStr:     [2]string{"array1", "array2"},
		ArrayBool:    [2]bool{true, false},
		ArrayInt:     [2]int{rand.Intn(100), rand.Intn(100)},
		ArrayInt8:    [2]int8{int8(rand.Intn(100)), int8(rand.Intn(100))},
		ArrayInt16:   [2]int16{int16(rand.Intn(100)), int16(rand.Intn(100))},
		ArrayInt32:   [2]int32{int32(rand.Intn(100)), int32(rand.Intn(100))},
		ArrayInt64:   [2]int64{int64(rand.Intn(100)), int64(rand.Intn(100))},
		ArrayUInt:    [2]uint{uint(rand.Intn(100)), uint(rand.Intn(100))},
		ArrayUInt8:   [2]uint8{uint8(rand.Intn(100)), uint8(rand.Intn(100))},
		ArrayUInt16:  [2]uint16{uint16(rand.Intn(100)), uint16(rand.Intn(100))},
		ArrayUInt32:  [2]uint32{uint32(rand.Intn(100)), uint32(rand.Intn(100))},
		ArrayUInt64:  [2]uint64{uint64(rand.Intn(100)), uint64(rand.Intn(100))},
		ArrayFloat32: [2]float32{rand.Float32() * 100, rand.Float32() * 100},
		ArrayFloat64: [2]float64{rand.Float64() * 100, rand.Float64() * 100},

		// Option types - Basic
		StrOption:     mo.Some("optionStr"),
		BoolOption:    mo.Some(true),
		IntOption:     mo.Some(rand.Intn(100)),
		Int8Option:    mo.Some(int8(rand.Intn(100))),
		Int16Option:   mo.Some(int16(rand.Intn(100))),
		Int32Option:   mo.Some(int32(rand.Intn(100))),
		Int64Option:   mo.Some(int64(rand.Intn(100))),
		UIntOption:    mo.Some(uint(rand.Intn(100))),
		UInt8Option:   mo.Some(uint8(rand.Intn(100))),
		UInt16Option:  mo.Some(uint16(rand.Intn(100))),
		UInt32Option:  mo.Some(uint32(rand.Intn(100))),
		UInt64Option:  mo.Some(uint64(rand.Intn(100))),
		Float32Option: mo.Some(rand.Float32() * 100),
		Float64Option: mo.Some(rand.Float64() * 100),

		// Option types - Slice
		SliceStrOption:     mo.Some([]string{"optSlice1", "optSlice2"}),
		SliceBoolOption:    mo.Some([]bool{true, false}),
		SliceIntOption:     mo.Some([]int{rand.Intn(100), rand.Intn(100)}),
		SliceInt8Option:    mo.Some([]int8{int8(rand.Intn(100)), int8(rand.Intn(100))}),
		SliceInt16Option:   mo.Some([]int16{int16(rand.Intn(100)), int16(rand.Intn(100))}),
		SliceInt32Option:   mo.Some([]int32{int32(rand.Intn(100)), int32(rand.Intn(100))}),
		SliceInt64Option:   mo.Some([]int64{int64(rand.Intn(100)), int64(rand.Intn(100))}),
		SliceUIntOption:    mo.Some([]uint{uint(rand.Intn(100)), uint(rand.Intn(100))}),
		SliceUInt8Option:   mo.Some([]uint8{uint8(rand.Intn(100)), uint8(rand.Intn(100))}),
		SliceUInt16Option:  mo.Some([]uint16{uint16(rand.Intn(100)), uint16(rand.Intn(100))}),
		SliceUInt32Option:  mo.Some([]uint32{uint32(rand.Intn(100)), uint32(rand.Intn(100))}),
		SliceUInt64Option:  mo.Some([]uint64{uint64(rand.Intn(100)), uint64(rand.Intn(100))}),
		SliceFloat32Option: mo.Some([]float32{rand.Float32() * 100, rand.Float32() * 100}),
		SliceFloat64Option: mo.Some([]float64{rand.Float64() * 100, rand.Float64() * 100}),
	}

	// Generate query string using reflection
	queryString := generateQueryString(v)
	fmt.Printf("Generated query string: %s\n", queryString)

	engine := gin.Default()
	engine.GET("/bind/query", func(c *gin.Context) {
		var value GinDto
		require.NoError(t, c.ShouldBindWith(&value, OptionQueryBinding))
		require.Equal(t, v, value)
	})

	engine.POST("/bind/form", func(c *gin.Context) {
		var value GinDto
		require.NoError(t, c.ShouldBindWith(&value, OptionFormBinding))
		require.Equal(t, v, value)
	})

	server := httptest.NewServer(engine)
	client := server.Client()

	_, err = client.Get(server.URL + "/bind/query?" + queryString)
	require.NoError(t, err)

	_, err = client.Post(server.URL+"/bind/form", "application/x-www-form-urlencoded", strings.NewReader(queryString))
	require.NoError(t, err)

	var buf bytes.Buffer
	multipartWriter := multipart.NewWriter(&buf)
	query, err := url.ParseQuery(queryString)
	require.NoError(t, err)
	for name, vs := range query {
		for _, v := range vs {
			require.NoError(t, multipartWriter.WriteField(name, v))
		}
	}
	require.NoError(t, multipartWriter.Close())

	req, err := http.NewRequest("POST", server.URL+"/bind/form", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	_, err = client.Do(req)
	require.NoError(t, err)
}

func TestGinQuery(t *testing.T) {
	var err error
	type User struct {
		Name string `form:"name" binding:"max=1"`
	}
	type User2 struct {
		Name string `form:"name"`
	}
	type User3 struct {
		Name string `form:"name" binding:""`
	}
	engine := gin.Default()
	engine.GET("/bind", func(c *gin.Context) {
		var value User
		require.ErrorContains(t, c.ShouldBindQuery(&value), "failed on the 'max' tag")
		require.Equal(t, "sb", value.Name)
	})
	engine.GET("/bind/valid", func(c *gin.Context) {
		var value User
		require.NoError(t, c.ShouldBindQuery(&value))
		require.Equal(t, "s", value.Name)
	})
	engine.GET("/bind/valid/empty", func(c *gin.Context) {
		var value User
		require.NoError(t, c.ShouldBindQuery(&value))
		require.Equal(t, "", value.Name)
	})
	engine.GET("/bind/valid/empty2", func(c *gin.Context) {
		var value User2
		require.NoError(t, c.ShouldBindQuery(&value))
		require.Equal(t, "", value.Name)
	})
	engine.GET("/bind/valid/empty3", func(c *gin.Context) {
		var value User3
		require.NoError(t, c.ShouldBindQuery(&value))
		require.Equal(t, "", value.Name)
	})
	type User4 struct {
		Name []string `form:"name" binding:""`
	}
	engine.GET("/bind/slice", func(c *gin.Context) {
		var value User3
		require.NoError(t, c.ShouldBindQuery(&value))
		require.Equal(t, "", value.Name)
	})

	server := httptest.NewServer(engine)
	client := server.Client()

	_, err = client.Get(server.URL + "/bind?name=sb")
	require.NoError(t, err)

	_, err = client.Get(server.URL + "/bind/valid?name=s")
	require.NoError(t, err)

	_, err = client.Get(server.URL + "/bind/valid/empty")
	require.NoError(t, err)
	_, err = client.Get(server.URL + "/bind/valid/empty2")
	require.NoError(t, err)
	_, err = client.Get(server.URL + "/bind/valid/empty3")
	require.NoError(t, err)
}
