package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTag(t *testing.T) {
	tests := []struct {
		name    string
		tagName string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常系：有効なタグ名",
			tagName: "Go",
			wantErr: false,
		},
		{
			name:    "正常系：タグ名が50文字ちょうど",
			tagName: strings.Repeat("あ", 50),
			wantErr: false,
		},
		{
			name:    "正常系：タグ名が1文字",
			tagName: "A",
			wantErr: false,
		},
		{
			name:    "正常系：日本語のタグ名",
			tagName: "プログラミング",
			wantErr: false,
		},
		{
			name:    "正常系：英数字混在のタグ名",
			tagName: "Go1.21",
			wantErr: false,
		},
		{
			name:    "異常系：空文字列のタグ名",
			tagName: "",
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "異常系：タグ名が51文字超過",
			tagName: strings.Repeat("あ", 51),
			wantErr: true,
			errMsg:  "name must be 50 characters or less",
		},
		{
			name:    "異常系：空白のみのタグ名",
			tagName: "   ",
			wantErr: true,
			errMsg:  "name cannot be only whitespace",
		},
		{
			name:    "異常系：タブのみのタグ名",
			tagName: "\t\t",
			wantErr: true,
			errMsg:  "name cannot be only whitespace",
		},
		{
			name:    "異常系：改行のみのタグ名",
			tagName: "\n\n",
			wantErr: true,
			errMsg:  "name cannot be only whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beforeTime := time.Now()
			tag, err := NewTag(tt.tagName)
			afterTime := time.Now()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
				require.NotNil(t, tag)

				assert.Equal(t, tt.tagName, tag.Name)
				assert.Equal(t, int64(0), tag.ID)

				assert.True(t, tag.CreatedAt.After(beforeTime) || tag.CreatedAt.Equal(beforeTime))
				assert.True(t, tag.CreatedAt.Before(afterTime) || tag.CreatedAt.Equal(afterTime))

				assert.Equal(t, tag.CreatedAt, tag.UpdatedAt)
			}
		})
	}
}
