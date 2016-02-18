package parser

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBlockify(t *testing.T) {
	str := `a
    b
    c
        d
    e
f`

	body, err := Blockify(str)
	assert.Nil(t, err)
	assert.NotNil(t, body)

	lineType := GeneralLine{}
	bodyType := &GeneralBody{}

	assert.Len(t, body.children, 3)
	assert.IsType(t, lineType, body.children[0])
	assert.Equal(t, "a", body.children[0].(GeneralLine).line)

	assert.IsType(t, bodyType, body.children[1])
	sub1 := body.children[1].(*GeneralBody)
	assert.Len(t, sub1.children, 4)
	assert.IsType(t, lineType, sub1.children[0])
	assert.Equal(t, "b", sub1.children[0].(GeneralLine).line)
	assert.IsType(t, lineType, sub1.children[1])
	assert.Equal(t, "c", sub1.children[1].(GeneralLine).line)

	assert.IsType(t, bodyType, sub1.children[2])
	sub2 := sub1.children[2].(*GeneralBody)
	assert.Len(t, sub2.children, 1)
	assert.IsType(t, lineType, sub2.children[0])
	assert.Equal(t, "d", sub2.children[0].(GeneralLine).line)

	assert.IsType(t, lineType, sub1.children[3])
	assert.Equal(t, "e", sub1.children[3].(GeneralLine).line)

	assert.IsType(t, lineType, body.children[2])
	assert.Equal(t, "f", body.children[2].(GeneralLine).line)
}
