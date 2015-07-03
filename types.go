package main

type IbexType interface {}

type IbexTupleType struct {
    ElementTypes []IbexType
}

type IbexArrayType struct {
    ElementType IbexType
    Dimensions int
}

type IbexNamedTupleEntry struct {
    Name string
    Type IbexType
}
type IbexNamedTupleType struct {
    types []*IbexNamedTupleEntry
}

type IbexFunctionType struct {
    Argument IbexType
    Return IbexType
}

type IbexSimpleType struct {
    Name string
}
