package model

import (
	"testing"
)

func TestCommand_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cmd     Command
		wantErr bool
	}{
		{
			name: "valid command",
			cmd: Command{
				Name:            "ls",
				Category:        "操作系统/通用Linux命令",
				InstallRequired: false,
				Description:     "列出目录内容",
				Usage:           []string{"ls [选项] [目录]"},
				Platforms:       []string{"linux", "darwin"},
				Examples: []Example{
					{Command: "ls -l", Description: "长格式列表"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			cmd: Command{
				Category:    "操作系统/通用Linux命令",
				Description: "列出目录内容",
			},
			wantErr: true,
		},
		{
			name: "missing category",
			cmd: Command{
				Name:        "ls",
				Description: "列出目录内容",
			},
			wantErr: true,
		},
		{
			name: "missing description",
			cmd: Command{
				Name:     "ls",
				Category: "操作系统/通用Linux命令",
			},
			wantErr: true,
		},
		{
			name: "invalid risk level",
			cmd: Command{
				Name:        "test",
				Category:    "test",
				Description: "test",
				Platforms:   []string{"linux"},
				Risks: []Risk{
					{Level: "invalid", Description: "test"},
				},
			},
			wantErr: true,
		},
		{
			name: "valid risk levels",
			cmd: Command{
				Name:        "test",
				Category:    "test",
				Description: "test",
				Platforms:   []string{"linux"},
				Usage:       []string{"test"},
				Examples:    []Example{{Command: "test", Description: "test"}},
				Risks: []Risk{
					{Level: RiskLevelLow, Description: "low risk"},
					{Level: RiskLevelMedium, Description: "medium risk"},
					{Level: RiskLevelHigh, Description: "high risk"},
					{Level: RiskLevelCritical, Description: "critical risk"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cmd.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Command.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRiskLevel_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		level RiskLevel
		want  bool
	}{
		{"low", RiskLevelLow, true},
		{"medium", RiskLevelMedium, true},
		{"high", RiskLevelHigh, true},
		{"critical", RiskLevelCritical, true},
		{"invalid", RiskLevel("invalid"), false},
		{"empty", RiskLevel(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.level.IsValid(); got != tt.want {
				t.Errorf("RiskLevel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_GetRiskLevel(t *testing.T) {
	tests := []struct {
		name string
		cmd  Command
		want RiskLevel
	}{
		{
			name: "no risks",
			cmd:  Command{},
			want: RiskLevelLow,
		},
		{
			name: "low risk",
			cmd: Command{
				Risks: []Risk{{Level: RiskLevelLow}},
			},
			want: RiskLevelLow,
		},
		{
			name: "critical risk",
			cmd: Command{
				Risks: []Risk{
					{Level: RiskLevelLow},
					{Level: RiskLevelCritical},
				},
			},
			want: RiskLevelCritical,
		},
		{
			name: "high risk",
			cmd: Command{
				Risks: []Risk{
					{Level: RiskLevelMedium},
					{Level: RiskLevelHigh},
				},
			},
			want: RiskLevelHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cmd.GetRiskLevel(); got != tt.want {
				t.Errorf("Command.GetRiskLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_SupportsPlatform(t *testing.T) {
	cmd := Command{
		Platforms: []string{"linux", "darwin"},
	}

	tests := []struct {
		platform string
		want     bool
	}{
		{"linux", true},
		{"darwin", true},
		{"windows", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.platform, func(t *testing.T) {
			if got := cmd.SupportsPlatform(tt.platform); got != tt.want {
				t.Errorf("Command.SupportsPlatform(%v) = %v, want %v", tt.platform, got, tt.want)
			}
		})
	}
}
