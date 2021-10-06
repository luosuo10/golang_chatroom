package utils

import "sync"

type IDGenerator struct {
	globalID int
	idLocker sync.Mutex
}

var (
	IDGeneratorOnce sync.Once
	globalIDGen IDGenerator
)

func CreateIDGenerator() {
	IDGeneratorOnce.Do(func() {
		globalIDGen = IDGenerator {
			globalID: 0,
		}
	})
}

func GenUserID() int {
	globalIDGen.idLocker.Lock()
	defer globalIDGen.idLocker.Unlock()

	globalIDGen.globalID++
	return globalIDGen.globalID
}
