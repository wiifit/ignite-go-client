package ignite

import (
	"testing"
	"time"
)

func Test_client_ResourceClose(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()
	var status int32

	// insert test values
	tm := time.Date(2018, 4, 3, 14, 25, 32, int(time.Millisecond*123+time.Microsecond*456+789), time.UTC)
	_, err = c.QuerySQLFields("TestDB", false, QuerySQLFieldsData{
		PageSize: 10,
		Query: "INSERT INTO Organization(_key, name, foundDateTime) VALUES" +
			"(?, ?, ?)," +
			"(?, ?, ?)," +
			"(?, ?, ?)",
		QueryArgs: []interface{}{
			int64(1), "Org 1", tm,
			int64(2), "Org 2", tm,
			int64(3), "Org 3", tm},
	}, &status)
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err.Error())
	}
	defer c.CacheRemoveAll("TestDB", false, nil)
	// select test values
	res, err := c.QuerySQLFields("TestDB", false, QuerySQLFieldsData{
		PageSize: 2,
		Query:    "SELECT name, foundDateTime FROM Organization ORDER BY name ASC",
		Timeout:  10000,
	}, &status)
	if err != nil {
		t.Fatalf("failed to select test data: %s", err.Error())
	}

	type args struct {
		id     int64
		status *int32
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
				id:     res.ID,
				status: &status,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ResourceClose(tt.args.id, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("client.ResourceClose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
