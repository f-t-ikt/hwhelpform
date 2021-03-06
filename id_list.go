package main

import (
    "container/list"
    "sync"
)

type IdList struct {
    list *list.List
    sync.Mutex
}

func NewIdList() *IdList {
    return &IdList {
        list: list.New(),
    }
}

func (il *IdList) Add(v interface{}) *list.Element {
    il.Lock()
    defer il.Unlock()
    return il.list.PushBack(v)
}

func (il *IdList) Remove(v interface{}) interface{} {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value == v {
            return il.list.Remove(e)
        }
    }
    return nil
}

func (il *IdList) Contains(v interface{}) bool {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value == v {
            return true
        }
    }
    return false
}

func (il *IdList) Each(f func(v interface{}) bool) {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if !f(e.Value) {
            return
        }
    }
}

func (il *IdList) ContainsId(v interface{}) bool {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value.(*Post).Id == v.(*Post).Id {
            return true
        }
    }
    return false
}

func (il *IdList) RemoveById(v interface{}) interface{} {
    il.Lock()
    defer il.Unlock()
    for e := il.list.Front(); e != nil; e = e.Next() {
        if e.Value.(*Post).Id == v.(*Post).Id {
            return il.list.Remove(e)
        }
    }
    return nil
}