package ignite

import (
	"reflect"
	"testing"
	"time"
)

/*
import (
	"reflect"
	"testing"
	"time"
)

func Test_client_QuerySQL(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()
	// insert test values
	tm := time.Date(2018, 4, 3, 14, 25, 32, int(time.Millisecond*123+time.Microsecond*456+789), time.UTC)
	_, err = c.QuerySQLFields("TestDB1", false, QuerySQLFieldsData{
		PageSize: 10,
		Query: "INSERT INTO Organization(_key, name, foundDateTime) VALUES" +
			"(?, ?, ?)," +
			"(?, ?, ?)," +
			"(?, ?, ?)",
		QueryArgs: []interface{}{
			int64(1), "Org 1", tm,
			int64(2), "Org 2", tm,
			int64(3), "Org 3", tm},
	})
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err.Error())
	}
	defer c.CacheRemoveAll("TestDB1", false)

	type args struct {
		cache  string
		binary bool
		data   QuerySQLData
		status *int32
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    QuerySQLResult
		wantErr bool
	}{
		{
			name: "success test 1",
			c:    c,
			args: args{
				cache: "TestDB1",
				data: QuerySQLData{
					Table:    "Organization",
					Query:    `SELECT * FROM Organization ORDER BY name ASC`,
					PageSize: 10,
					Timeout:  10000,
				},
			},
			want: QuerySQLResult{
				Keys:   []interface{}{},
				Values: []interface{}{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.QuerySQL(tt.args.cache, tt.args.binary, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.QuerySQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Keys, tt.want.Keys) ||
				!reflect.DeepEqual(got.Values, tt.want.Values) ||
				!reflect.DeepEqual(got.HasMore, tt.want.HasMore) {
				t.Errorf("client.QuerySQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func Test_client_QuerySQLFields(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()
	defer c.CacheRemoveAll("TestDB1", false)
	tm := time.Date(2018, 4, 3, 14, 25, 32, int(time.Millisecond*123+time.Microsecond*456+789), time.UTC)

	type args struct {
		cache  string
		binary bool
		data   QuerySQLFieldsData
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    QuerySQLFieldsResult
		wantErr bool
	}{
		{
			name: "success test 1",
			c:    c,
			args: args{
				cache: "TestDB1",
				data: QuerySQLFieldsData{
					PageSize: 10,
					Query: "INSERT INTO Organization(_key, name, foundDateTime) VALUES" +
						"(?, ?, ?)," +
						"(?, ?, ?)," +
						"(?, ?, ?)",
					QueryArgs: []interface{}{
						int64(1), "Org 1", tm,
						int64(2), "Org 2", tm,
						int64(3), "Org 3", tm},
				},
			},
			want: QuerySQLFieldsResult{
				FieldCount: 1,
				QuerySQLFieldsPage: QuerySQLFieldsPage{
					Rows: [][]interface{}{[]interface{}{int64(3)}},
				},
			},
		},
		{
			name: "success test 2",
			c:    c,
			args: args{
				cache: "TestDB1",
				data: QuerySQLFieldsData{
					PageSize: 10,
					Query: "INSERT INTO Person(_key, orgId, firstName, lastName, resume, salary) VALUES" +
						"(?, ?, ?, ?, ?, ?)," +
						"(?, ?, ?, ?, ?, ?)," +
						"(?, ?, ?, ?, ?, ?)," +
						"(?, ?, ?, ?, ?, ?)," +
						"(?, ?, ?, ?, ?, ?)",
					QueryArgs: []interface{}{
						int64(4), int64(1), "First name 1", "Last name 1", "Resume 1", float64(100.0),
						int64(5), int64(1), "First name 2", "Last name 2", "Resume 2", float64(200.0),
						int64(6), int64(2), "First name 3", "Last name 3", "Resume 3", float64(300.0),
						int64(7), int64(2), "First name 4", "Last name 4", "Resume 4", float64(400.0),
						int64(8), int64(3), "First name 5", "Last name 5", "Resume 5", float64(500.0)},
				},
			},
			want: QuerySQLFieldsResult{
				FieldCount: 1,
				QuerySQLFieldsPage: QuerySQLFieldsPage{
					Rows: [][]interface{}{[]interface{}{int64(5)}},
				},
			},
		},
		{
			name: "success test 3",
			c:    c,
			args: args{
				cache: "TestDB1",
				data: QuerySQLFieldsData{
					PageSize: 10,
					Query: "SELECT " +
						"o.name AS Name, " +
						"o.foundDateTime AS Found, " +
						"p.firstName AS FirstName, " +
						"p.lastName AS LastName, " +
						"p.salary AS Salary " +
						"FROM Person p INNER JOIN Organization o ON p.orgId = o._key " +
						"WHERE o._key = ? " +
						"ORDER BY p.firstName",
					QueryArgs: []interface{}{
						int64(2)},
					Timeout:           10000,
					IncludeFieldNames: true,
				},
			},
			want: QuerySQLFieldsResult{
				FieldCount: 5,
				Fields:     []string{"NAME", "FOUND", "FIRSTNAME", "LASTNAME", "SALARY"},
				QuerySQLFieldsPage: QuerySQLFieldsPage{
					Rows: [][]interface{}{
						[]interface{}{"Org 2", tm, "First name 3", "Last name 3", float64(300.0)},
						[]interface{}{"Org 2", tm, "First name 4", "Last name 4", float64(400.0)},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.QuerySQLFields(tt.args.cache, tt.args.binary, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.QuerySQLFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (tt.want.Fields != nil && !reflect.DeepEqual(got.Fields, tt.want.Fields)) ||
				(tt.want.Rows != nil && !reflect.DeepEqual(got.Rows, tt.want.Rows)) ||
				!reflect.DeepEqual(got.FieldCount, tt.want.FieldCount) ||
				!reflect.DeepEqual(got.HasMore, tt.want.HasMore) {
				t.Errorf("client.QuerySQLFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_QuerySQLFieldsCursorGetPage(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()

	// insert test values
	tm := time.Date(2018, 4, 3, 14, 25, 32, int(time.Millisecond*123+time.Microsecond*456+789), time.UTC)
	_, err = c.QuerySQLFields("TestDB1", false, QuerySQLFieldsData{
		PageSize: 10,
		Query: "INSERT INTO Organization(_key, name, foundDateTime) VALUES" +
			"(?, ?, ?)," +
			"(?, ?, ?)," +
			"(?, ?, ?)",
		QueryArgs: []interface{}{
			int64(1), "Org 1", tm,
			int64(2), "Org 2", tm,
			int64(3), "Org 3", tm},
	})
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err.Error())
	}
	defer c.CacheRemoveAll("TestDB1", false)
	// select test values
	res, err := c.QuerySQLFields("TestDB1", false, QuerySQLFieldsData{
		PageSize: 2,
		Query:    "SELECT name, foundDateTime FROM Organization ORDER BY name ASC",
		Timeout:  10000,
	})
	if err != nil {
		t.Fatalf("failed to select test data: %s", err.Error())
	}

	type args struct {
		id         int64
		fieldCount int
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    QuerySQLFieldsPage
		wantErr bool
	}{
		{
			name: "success test 1",
			c:    c,
			args: args{
				id:         res.ID,
				fieldCount: res.FieldCount,
			},
			want: QuerySQLFieldsPage{
				Rows: [][]interface{}{
					[]interface{}{"Org 3", tm},
				},
				HasMore: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.QuerySQLFieldsCursorGetPage(tt.args.id, tt.args.fieldCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.QuerySQLFieldsCursorGetPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.QuerySQLFieldsCursorGetPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_ResourceClose(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()

	// insert test values
	tm := time.Date(2018, 4, 3, 14, 25, 32, int(time.Millisecond*123+time.Microsecond*456+789), time.UTC)
	_, err = c.QuerySQLFields("TestDB1", false, QuerySQLFieldsData{
		PageSize: 10,
		Query: "INSERT INTO Organization(_key, name, foundDateTime) VALUES" +
			"(?, ?, ?)," +
			"(?, ?, ?)," +
			"(?, ?, ?)",
		QueryArgs: []interface{}{
			int64(1), "Org 1", tm,
			int64(2), "Org 2", tm,
			int64(3), "Org 3", tm},
	})
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err.Error())
	}
	defer c.CacheRemoveAll("TestDB1", false)
	// select test values
	res, err := c.QuerySQLFields("TestDB1", false, QuerySQLFieldsData{
		PageSize: 2,
		Query:    "SELECT name, foundDateTime FROM Organization ORDER BY name ASC",
		Timeout:  10000,
	})
	if err != nil {
		t.Fatalf("failed to select test data: %s", err.Error())
	}

	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		{
			name: "success test 1",
			c:    c,
			args: args{
				id: res.ID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ResourceClose(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("client.ResourceClose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
