package entity

import (
	"fmt"
	"sync"
)

type (
	Symbol string

	Enum struct {
		value   interface{}
		desc    string
		symbol  Symbol
		mapping interface{}
	}

	enumLists []*Enum

	EnumMgr struct {
		enumType Symbol
		lists    *enumLists
		safe     sync.RWMutex
		loader   func(mgr *EnumMgr)
		once     sync.Once
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
	for _, v := range *enums {
		if v.value == e {
			return *v, true
		}
	}
	return *nullEnum, false
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

func (mgr *EnumMgr) load(loader func(*EnumMgr)) {
	if mgr.loader == nil {
		mgr.loader = loader
	}
	mgr.once.Do(func() {
		mgr.loader(mgr)
	})
}
