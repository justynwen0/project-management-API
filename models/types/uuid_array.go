package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(value interface{}) error {

	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
		return errors.New("failed to parse for UUIDArray; unsupported type")
	}
	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")
	parts := strings.Split(str, ",")

	//make([],lenght, capacity)
	*a = make(UUIDArray, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(strings.Trim(s, `"`)) //akan menghapus spasi dan tanda kutip
		if s == "" {
			continue
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid uuid in array: %v", err)
		}
		*a = append(*a, u)
	}
	return nil
}

func (a UUIDArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	postgreformat := make([]string, 0, len(a))
	for _, u := range a {
		postgreformat = append(postgreformat, fmt.Sprintf(`"%s"`, u.String()))
	}
	return "{" + strings.Join(postgreformat, ",") + "}", nil

}

func (UUIDArray) GormDataType() string {
	return "uuid[]"
}
