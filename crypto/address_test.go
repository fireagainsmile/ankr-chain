package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContractAddress(t *testing.T) {
	contAddr1 := CreateContractAddress("0669BA04564603E98DD248185A7440532311CC990DAB69", 0)
	contAddr2 := CreateContractAddress("0669BA04564603E98DD248185A7440532311CC990DAB69", 1)
	t.Logf("contAddr1=%v, contAddr1=%v", contAddr1, contAddr2)
	assert.NotEqual(t, contAddr1, contAddr2)

	contAddr3 := CreateContractAddress("0669BA04564603E98DD248185A7440532311CC990DAB69", 1)
	contAddr4 := CreateContractAddress("0669BA04564603E98DD248185A7440532311CC990DAB69", 1)
	t.Logf("contAddr3=%v, contAddr4=%v", contAddr1, contAddr2)
	assert.Equal(t, contAddr3, contAddr4)
}
