package mapblockparser

import (
	"errors"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func Parse(data []byte) (*MapBlock, error) {
	mapblock := NewMapblock()
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	log.WithFields(log.Fields{
		"data-length": len(data),
	}).Debug("Parsing mapblock")

	offset := 0

	// version
	mapblock.Version = data[0]

	//flags
	flags := data[1]
	mapblock.Underground = (flags & 0x01) == 0x01

	content_width := data[4]
	params_width := data[4]

	if content_width != 2 {
		return nil, errors.New("content_width = " + strconv.Itoa(int(content_width)))
	}

	if params_width != 2 {
		return nil, errors.New("params_width = " + strconv.Itoa(int(params_width)))
	}

	//mapdata (blocks)
	offset = 6

	//metadata
	count, err := parseMapdata(&mapblock, data[offset:])
	if err != nil {
		return nil, err
	}

	offset += count

	count, err = parseMetadata(&mapblock, data[offset:])
	if err != nil {
		return nil, err
	}

	offset += count

	//static objects

	offset++ //static objects version
	staticObjectsCount := readU16(data, offset)
	offset += 2
	for i := 0; i < staticObjectsCount; i++ {
		offset += 13
		dataSize := readU16(data, offset)
		offset += dataSize + 2
	}

	//timestamp
	offset += 4

	//mapping version
	offset++

	numMappings := readU16(data, offset)
	offset += 2
	for i := 0; i < numMappings; i++ {
		nodeId := readU16(data, offset)
		offset += 2

		nameLen := readU16(data, offset)
		offset += 2

		blockName := string(data[offset : offset+nameLen])
		offset += nameLen

		mapblock.BlockMapping[nodeId] = blockName
	}

	return &mapblock, nil
}