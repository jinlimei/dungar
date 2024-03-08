package triggers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
)

func TestMain(m *testing.M) {
	db.FlagWritingIssues(false)
	os.Exit(m.Run())
}

func TestIsJulietNumberOutput(t *testing.T) {
	str := "501050 (9761, 5358, 2291, 3356, 9115, 9918, 5203, 8280, 6908, 2712, 7153, 6889, " +
		"1665, 1333, 9645, 3232, 8905, 542, 1736, 5143, 8301, 2424, 5710, 8597, 8641, 31, 1462, " +
		"6779, 4128, 2791, 5051, 4065, 5351, 8621, 8461, 6006, 9413, 9718, 584, 2701, 4539, 6326, " +
		"1302, 7976, 62, 2851, 5191, 3932, 5765, 9549, 186, 3767, 3068, 425, 8145"

	assert.True(t, isJulietNumberSequence(str))

	str = "Hello, World! 2341, 4291, 2014, 1024, 5812!"

	assert.False(t, isJulietNumberSequence(str))

	str = "0"

	assert.False(t, isJulietNumberSequence(str))
}
