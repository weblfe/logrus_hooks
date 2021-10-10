package utils

import (
		"os"
		"testing"
		"time"
)

type (

	user struct {
		User string `json:"user" env:"user_name"`
		Sex  uint   `json:"sex" env:"user_sex"`
	}

	info struct {
		ID     int    `json:"id" env:"id,1"`
		Avatar string `json:"avatar" env:"avatar"`
		Images images `json:"images"`
	}

	images struct {
		Tag      string        `json:"tag" env:"image_tag,default"`
		CreateAt time.Time     `json:"create_at" env:"create_at"`
		Duration time.Duration `json:"duration" env:"duration,10s"`
	}

	Data struct {
		Code uint    `json:"code" env:"data_code"`
		Msg  string `json:"msg" env:"data_msg"`
		Info info   `json:"info"`
	}

	testEnv struct {
		Name     string                 `json:"name" env:"name"`
		Password string                 `json:"password" env:"password,123"`
		Number   int                    `json:"number" env:"number"`
		Boolean  bool                   `json:"boolean" env:"boolean"`
		Arr      []int                  `json:"arr" env:"arr"`
		Maps     map[string]interface{} `json:"maps" env:"maps"`
		User     user                   `json:"user"`
		Data     *Data                  `json:"data"`
	}

)

func  initTestData()  {
		_ = os.Setenv("NAME", "test")
		_ = os.Setenv("PASSWORD", "test1111")
		_ = os.Setenv("NUMBER", "11")
		_ = os.Setenv("BOOLEAN", "true")
		_ = os.Setenv("ARR", "[1,1,1]")
		_ = os.Setenv("USER_NAME", "env")
		_ = os.Setenv("USER_SEX", "1")
		_ = os.Setenv("DATA_CODE", "200")
		_ = os.Setenv("DATA_MSG", "OK")
		_ = os.Setenv("ID", "20")
		_ = os.Setenv("CREATE_AT", "2006-01-02 15:04:05")
		_ = os.Setenv("AVATAR", "http://127.0.0.1/image.png")
		_ = os.Setenv("MAPS", `{"name":"123","num":1,"bool":true,"nil":null}`)
}

func TestEnvTagLoader_Marshal(t *testing.T) {
	var (
		loader      = NewEnvDecoder()
		env         = new(testEnv)
		duration, _ = time.ParseDuration("10s")
		dateTime, _ = time.Parse(DateTimeLayout, "2006-01-02 15:04:05")
	)
	initTestData()
	env.Data = new(Data)
	if err := loader.Marshal(env); err != nil {
		t.Error(err)
	}
	if env.Password != "test1111" {
		t.Error("解析环境变量失败")
	}
	if env.Name != "test" {
		t.Error("解析环境变量失败")
	}
	if env.Number != 11 {
		t.Error("解析环境变量失败")
	}
	if !env.Boolean {
		t.Error("解析环境变量失败")
	}
	if env.Maps == nil {
		t.Error("解析环境变量失败")
	}
	if env.Arr == nil {
		t.Error("解析环境变量失败")
	}
	if env.Data.Code != 200 {
		t.Error("解析环境变量失败")
	}
	if env.Data.Msg != "OK" {
		t.Error("解析环境变量失败")
	}
	if env.Data.Info.Avatar != "http://127.0.0.1/image.png" {
		t.Error("解析环境变量 嵌套结构失败")
	}
	if env.Data.Info.ID != 20 {
		t.Error("解析环境变量 嵌套结构失败")
	}
	if !env.Data.Info.Images.CreateAt.Equal(dateTime) {
		t.Error("解析环境变量 时间类型失败")
	}
	if env.Data.Info.Images.Duration != duration {
		t.Error("解析环境变量 时间类型失败")
	}
}
