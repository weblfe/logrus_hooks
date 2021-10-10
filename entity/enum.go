package entity

import (
	"encoding/json"
	"fmt"
	"sync"
)

type (
	// Symbol 符号类型
	Symbol string

	// Enum 枚举类型
	Enum struct {
		value   interface{} // 枚举值
		desc    string      // 枚举描述
		symbol  Symbol      // 枚举类型
		mapping interface{} // 枚举自定义映射值
	}

	// enumLists 枚举列表
	enumLists []*Enum

	// EnumMgr 分类枚举列表
	EnumMgr struct {
		enumType Symbol
		lists    *enumLists
		safe     sync.RWMutex
		loader   func(mgr *EnumMgr)
		once     sync.Once
	}

	exportEnum struct {
		Symbol  Symbol
		Value   interface{}
		Mapping interface{}
		Desc    string
	}
)

var (
	nullEnum = NewEnum(nil, "nil", "nil enums")
)

func NewEnum(v interface{}, symbol Symbol, desc ...string) *Enum {
	var enum = new(Enum)
	enum.symbol = symbol
	enum.value = v
	desc = append(desc, fmt.Sprintf("%v", v))
	enum.desc = desc[0]
	return enum
}

func (enum *Enum) Symbol() Symbol {
	return enum.symbol
}

func (enum *Enum) String() string {
	return fmt.Sprintf("%s:%v", enum.symbol, enum.value)
}

func (enum *Enum) Equal(e *Enum) bool {
	if enum == nil || e == nil {
		return false
	}
	if enum.value == e.value && enum.symbol == e.symbol {
		return true
	}
	return false
}

func (enum *Enum) SetCustom(m interface{}) *Enum {
	if enum == nil {
		return nil
	}
	enum.mapping = m
	return enum
}

func (enum *Enum) GetCustom() interface{} {
	if enum == nil {
		return nil
	}
	return enum.mapping
}

func (enum *Enum) IsNull() bool {
	if enum == nullEnum {
		return true
	}
	if enum.symbol == nullEnum.symbol && enum.value == nullEnum.value {
		return true
	}
	return false
}

func (enum *Enum) Export() *exportEnum {
	return &exportEnum{
		Symbol:  enum.symbol,
		Value:   enum.value,
		Desc:    enum.desc,
		Mapping: enum.mapping,
	}
}

func (enums *enumLists) In(e *Enum) bool {
	if enums == nil || e == nil {
		return false
	}
	for _, v := range *enums {
		if v.Equal(e) {
			return true
		}
	}
	return false
}

func (enums *enumLists) Len() int {
	if enums == nil {
		return 0
	}
	return len(*enums)
}

func (enums *enumLists) Register(e *Enum) bool {
	if enums == nil {
		return false
	}
	*enums = append(*enums, e)
	return true
}

func (enums *enumLists) Get(e interface{}) (Enum, bool) {
	var enum *Enum
	switch e.(type) {
	case Enum:
		var v = e.(Enum)
		enum = &v
	case *Enum:
		enum = e.(*Enum)
	}
	for _, v := range *enums {
		if enum != nil {
			if enum.Equal(v) {
				return *v, true
			}
			continue
		}
		if v.value == e || v.mapping == e {
			return *v, true
		}
	}
	return *nullEnum, false
}

func (enums *enumLists) Index(e interface{}) (int, bool) {
	var enum *Enum
	switch e.(type) {
	case Enum:
		var v = e.(Enum)
		enum = &v
	case *Enum:
		enum = e.(*Enum)
	}
	for i, v := range *enums {
		if enum != nil {
			if enum.Equal(v) {
				return i, false
			}
			continue
		}
		if v.value == e {
			return i, true
		}
	}
	return -1, false
}

func NewEnumMgr(symbol Symbol) *EnumMgr {
	var enumMgr = new(EnumMgr)
	enumMgr.enumType = symbol
	enumMgr.safe = sync.RWMutex{}
	enumMgr.lists = new(enumLists)
	enumMgr.once = sync.Once{}
	return enumMgr
}

func (mgr *EnumMgr) Add(enum *Enum) *EnumMgr {
	mgr.safe.Lock()
	defer mgr.safe.Unlock()
	if !mgr.enable(enum) {
		return mgr
	}
	if mgr.lists.In(enum) {
		return mgr
	}
	mgr.lists.Register(enum)
	return mgr
}

func (mgr *EnumMgr) Index(v interface{}) (int, bool) {
	mgr.safe.Lock()
	defer mgr.safe.Unlock()
	return mgr.lists.Index(v)
}

func (mgr *EnumMgr) In(enum *Enum) bool {
	mgr.safe.Lock()
	defer mgr.safe.Unlock()
	if mgr.enable(enum) {
		return false
	}
	if mgr.lists.In(enum) {
		return true
	}
	return false
}

func (mgr *EnumMgr) enable(enum *Enum) bool {
	if enum == nil {
		return false
	}
	if mgr.enumType != enum.symbol {
		return false
	}
	return true
}

func (mgr *EnumMgr) Get(v interface{}) (Enum, bool) {
	return mgr.lists.Get(v)
}

func (mgr *EnumMgr) Symbol() Symbol {
	mgr.safe.Lock()
	defer mgr.safe.Unlock()
	return mgr.enumType
}

func (mgr *EnumMgr) String() string {
	mgr.safe.Lock()
	defer mgr.safe.Unlock()
	var mapping []*exportEnum
	for _, v := range *mgr.lists {
		var m = v.Export()
		mapping = append(mapping, m)
	}
	if bytes, err := json.Marshal(mapping); err == nil {
		return string(bytes)
	}
	return ""
}

func (mgr *EnumMgr) load(loader func(*EnumMgr)) {
	if mgr.loader == nil {
		mgr.loader = loader
	}
	mgr.once.Do(func() {
		mgr.loader(mgr)
	})
}
