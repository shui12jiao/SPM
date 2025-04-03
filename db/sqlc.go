package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Nullable interface {
	sql.NullInt32 | sql.NullInt64 | sql.NullString | sql.NullBool | sql.NullTime | uuid.NullUUID
}

func ToNull[T Nullable](value any) T {
	if value == nil || reflect.ValueOf(value).IsNil() {
		// 返回T的零值
		var zeroValue T
		return zeroValue
	}

	switch v := value.(type) {
	case *int32:
		return any(sql.NullInt32{Int32: *v, Valid: true}).(T)
	case *int64:
		return any(sql.NullInt64{Int64: *v, Valid: true}).(T)
	case *string:
		return any(sql.NullString{String: *v, Valid: true}).(T)
	case *bool:
		return any(sql.NullBool{Bool: *v, Valid: true}).(T)
	case *time.Time:
		return any(sql.NullTime{Time: *v, Valid: true}).(T)
	case *uuid.UUID:
		return any(uuid.NullUUID{UUID: *v, Valid: true}).(T)
	default:
		panic(fmt.Sprintf("unsupported type: %T", value))
	}
}

// func ToNullTime(t *time.Time) sql.NullTime {
// 	if t == nil {
// 		return sql.NullTime{Valid: false}
// 	}
// 	return sql.NullTime{Time: *t, Valid: true}
// }

// func ToNullInt32(i *int32) sql.NullInt32 {
// 	if i == nil {
// 		return sql.NullInt32{Valid: false}
// 	}
// 	return sql.NullInt32{Int32: *i, Valid: true}
// }

// func ToNullString(s *string) sql.NullString {
// 	if s == nil {
// 		return sql.NullString{Valid: false}
// 	}
// 	return sql.NullString{String: *s, Valid: true}
// }

// func ToNullBool(b *bool) sql.NullBool {
// 	if b == nil {
// 		return sql.NullBool{Valid: false}
// 	}
// 	return sql.NullBool{Bool: *b, Valid: true}
// }
