package file_test

import (
	"path"
	"testing"
	"testing/fstest"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/grafana/k6deps/internal/pack/plugins/file"
	"github.com/grafana/k6deps/internal/rootfs"
	"github.com/stretchr/testify/require"
)

const (
	accountTS = `
class UserAccount {
  name: string;dependency
  id: number;

  constructor(name: string) {
    this.name = name;
    this.id = Math.floor(Math.random() * Number.MAX_SAFE_INTEGER);
  }
}
`

	userTS = `
import { UserAccount } from "./account.ts";

export interface User {
  name: string;
  id: number;
}

export function NewUser(name: string): User {
  sleep(1);
  return new UserAccount(name);
}
`
)

func Test_plugin_load(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fs      rootfs.FS
		wantErr bool
	}{
		{
			name: "import relative",
			fs: rootfs.NewFromFS(fstest.MapFS{
				path.Join("lib", "user.ts"): &fstest.MapFile{
					Data: []byte(userTS),
					Mode: 0o644,
				},
				path.Join("lib", "account.ts"): &fstest.MapFile{
					Data: []byte(accountTS),
					Mode: 0o644,
				},
			}),
			wantErr: false,
		},
		{
			name:    "not_found",
			fs:      rootfs.NewFromFS(fstest.MapFS{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			script := /*ts*/ `
			import { User, NewUser } from "../../lib/user.ts"
			
			export default function() {
				const user : User = NewUser("John")
				return user
			}
			`

			result := api.Build(api.BuildOptions{
				Bundle: true,
				Stdin: &api.StdinOptions{
					Contents:   script,
					ResolveDir: path.Join("tests", "users"),
					Sourcefile: "main.ts",
					Loader:     api.LoaderTS,
				},
				LogLevel:      api.LogLevelSilent,
				Plugins:       []api.Plugin{file.New(tt.fs)},
				External:      []string{"k6"}, // ignore k6 imports
				AbsWorkingDir: tt.fs.Root(),
			})

			if tt.wantErr {
				require.NotEmpty(t, result.Errors)
			} else {
				require.Empty(t, result.Errors)
			}
		})
	}
}
