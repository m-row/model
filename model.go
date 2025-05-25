package model

import (
	"context"
	"fmt"
	"net/url"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/m-row/finder"
)

type Dependencies struct {
	DB     finder.Connection
	QB     *squirrel.StatementBuilderType
	FB     *firebase.App
	FBM    *messaging.Client
	PGInfo map[string][]string
}

func SelectSeqID(
	modelID *int,
	modelName string,
	conn finder.Connection,
) {
	query := fmt.Sprintf(
		`
        SELECT setval('%s_id_seq',(SELECT MAX(id) FROM %s))
        `,
		modelName,
		modelName,
	)
	if err := conn.GetContext(
		context.Background(),
		modelID,
		query,
	); err != nil {
		*modelID = 0
		return
	}
}

func InputOrNewUUID(modelUUID *uuid.UUID, v url.Values) {
	id, err := uuid.Parse(v.Get("id"))
	if err != nil {
		*modelUUID = uuid.New()
		return
	}
	*modelUUID = id
}

func BoolParser(str string) bool {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True":
		return true
	case "0", "f", "F", "false", "FALSE", "False":
		return false
	default:
		return false
	}
}

// SQLSelectURLPath returns a conditionally selected image path column as in
//
//		 CASE WHEN nullif(banners.img, '') is not null
//	     THEN FORMAT('http://localhost:8056/%s', banners.img)
//		 ELSE null
//		 END as img
func SQLSelectURLPath(domain, tableName, colName, aliasName string) string {
	return `
    CASE
        WHEN nullif(` + tableName + `.` + colName + `, '') is not null
        THEN FORMAT('` + domain + `/%s', ` + tableName + `.` + colName + `)
        ELSE null
    END as ` + aliasName
}
