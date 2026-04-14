package system

import (
	"context"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"reflect"
	"testing"
)

func Test_autoCodeTemplate_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		info request.AutoCode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &autoCodeTemplate{}
			if err := s.Create(tt.args.ctx, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_autoCodeTemplate_Preview(t *testing.T) {
	if global.GVA_DB == nil {
		t.Skip("requires initialized DB (global.GVA_DB)")
	}
	type args struct {
		ctx  context.Context
		info request.AutoCode
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test package",
			args: args{
				ctx:  context.Background(),
				info: request.AutoCode{},
			},
			wantErr: false,
		},
		{
			name: "test plugin",
			args: args{
				ctx:  context.Background(),
				info: request.AutoCode{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJson := `{"structName":"SysUser","tableName":"sys_users","packageName":"sysUsers","package":"gva","abbreviation":"sysUsers","description":"sysUsersTable","businessDB":"","autoCreateApiToSql":true,"autoCreateMenuToSql":true,"autoMigrate":true,"gvaModel":true,"autoCreateResource":false,"fields":[{"fieldName":"Uuid","fieldDesc":"user UUID","fieldType":"string","dataType":"varchar","fieldJson":"uuid","primaryKey":false,"dataTypeLong":"191","columnName":"uuid","comment":"user UUID","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"Username","fieldDesc":"username","fieldType":"string","dataType":"varchar","fieldJson":"username","primaryKey":false,"dataTypeLong":"191","columnName":"username","comment":"username","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"Password","fieldDesc":"user login password","fieldType":"string","dataType":"varchar","fieldJson":"password","primaryKey":false,"dataTypeLong":"191","columnName":"password","comment":"user login password","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"NickName","fieldDesc":"user nickname","fieldType":"string","dataType":"varchar","fieldJson":"nickName","primaryKey":false,"dataTypeLong":"191","columnName":"nick_name","comment":"user nickname","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"SideMode","fieldDesc":"user side theme","fieldType":"string","dataType":"varchar","fieldJson":"sideMode","primaryKey":false,"dataTypeLong":"191","columnName":"side_mode","comment":"user side theme","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"HeaderImg","fieldDesc":"user avatar","fieldType":"string","dataType":"varchar","fieldJson":"headerImg","primaryKey":false,"dataTypeLong":"191","columnName":"header_img","comment":"user avatar","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"BaseColor","fieldDesc":"BasicColor","fieldType":"string","dataType":"varchar","fieldJson":"baseColor","primaryKey":false,"dataTypeLong":"191","columnName":"base_color","comment":"BasicColor","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"AuthorityId","fieldDesc":"user role ID","fieldType":"int","dataType":"bigint","fieldJson":"authorityId","primaryKey":false,"dataTypeLong":"20","columnName":"authority_id","comment":"user role ID","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"Phone","fieldDesc":"user phone number","fieldType":"string","dataType":"varchar","fieldJson":"phone","primaryKey":false,"dataTypeLong":"191","columnName":"phone","comment":"user phone number","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"Email","fieldDesc":"user email","fieldType":"string","dataType":"varchar","fieldJson":"email","primaryKey":false,"dataTypeLong":"191","columnName":"email","comment":"user email","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}},{"fieldName":"Enable","fieldDesc":"UserYesNofrozen 1Normal 2Freeze","fieldType":"int","dataType":"bigint","fieldJson":"enable","primaryKey":false,"dataTypeLong":"19","columnName":"enable","comment":"UserYesNofrozen 1Normal 2Freeze","require":false,"errorText":"","clearable":true,"fieldSearchType":"","fieldIndexType":"","dictType":"","front":true,"dataSource":{"association":1,"table":"","label":"","value":""}}],"humpPackageName":"sys_users"}`
			err := json.Unmarshal([]byte(testJson), &tt.args.info)
			if err != nil {
				t.Error(err)
				return
			}
			err = tt.args.info.Pretreatment()
			if err != nil {
				t.Error(err)
				return
			}
			got, err := AutoCodeTemplate.Preview(tt.args.ctx, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("Preview() error = %+v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Preview() got = %v, want %v", got, tt.want)
			}
		})
	}
}
