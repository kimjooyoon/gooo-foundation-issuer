package issuer

import (
	"fmt"
	"os"
	"strings"
)

const inventoryAuthorityLine = "inventory root_readme=README.md files=excluded physical_lines=excluded other_readmes=retained"

func ReadInventoryAuthority(path string) (InventoryAuthority, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return InventoryAuthority{}, err
	}
	if !strings.Contains(string(bytes), inventoryAuthorityLine) {
		return InventoryAuthority{}, fmt.Errorf("root README inventory authority is not exact")
	}
	return InventoryAuthority{
		Schema:                "gooo/foundation-issuer/inventory-authority/v1",
		RootReadmePath:        "README.md",
		RootReadmeExcluded:    true,
		PhysicalLinesExcluded: true,
		OtherReadmesRetained:  true,
		PolicyDigest:          Digest(bytes),
	}, nil
}
