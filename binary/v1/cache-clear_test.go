package ignite

import "testing"

func Test_client_CacheClear(t *testing.T) {
	// get test data
	c, err := getTestClient()
	if err != nil {
		t.Fatalf("failed to open test connection: %s", err.Error())
	}
	defer c.Close()
	if err = c.CacheCreateWithName("TestCache1", nil); err != nil {
		t.Fatalf("failed to create test cache: %s", err.Error())
	}
	defer c.CacheDestroy("TestCache1", nil)
	var status int32

	// put test values
	if err = c.CachePut("TestCache1", false, "key", "value", &status); err != nil {
		t.Fatalf("failed to put test pair: %s", err.Error())
	}

	type args struct {
		cache  string
		binary bool
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
				cache:  "TestCache1",
				binary: false,
				status: &status,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CacheClear(tt.args.cache, tt.args.binary, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("client.CacheClear() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}