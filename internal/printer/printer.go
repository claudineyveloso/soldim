package printer

import (
	"fmt"
	"os"

	"github.com/claudineyveloso/soldim.git/internal/types"
)

func GenerateZPL(product types.ProductTag) string {
	return fmt.Sprintf(`^XA
^FO50,50^A0N,25,25^FDNome: %s^FS
^FO50,100^A0N,25,25^FDSKU: %s^FS
^FO50,150^A0N,25,25^FDGTIN: %s^FS
^FO50,200^BY2^BCN,50,Y,N,N^FD%s^FS
^XZ`, product.Nome, product.SKU, product.GTIN, product.GTIN)
}

func SaveZPLToFile(zpl string, filename string) error {
	return os.WriteFile(filename, []byte(zpl), 0644)
}
