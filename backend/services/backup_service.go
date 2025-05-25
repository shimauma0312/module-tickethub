package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// BackupService はバックアップサービス
type BackupService struct {
	backupRepo   repositories.BackupRepository
	backupDir    string
	dbType       string
	dbConnection string
}

// NewBackupService は新しいBackupServiceを作成します
func NewBackupService(backupRepo repositories.BackupRepository, backupDir, dbType, dbConnection string) *BackupService {
	return &BackupService{
		backupRepo:   backupRepo,
		backupDir:    backupDir,
		dbType:       dbType,
		dbConnection: dbConnection,
	}
}

// CreateBackup はデータベースバックアップを作成します
func (s *BackupService) CreateBackup(ctx context.Context, userID int64, description string) (*models.BackupInfo, error) {
	// バックアップディレクトリが存在しない場合は作成
	if err := os.MkdirAll(s.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	// バックアップファイル名を生成
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("tickethub_backup_%s.sql", timestamp)
	filepath := filepath.Join(s.backupDir, filename)

	// バックアップ情報をデータベースに記録
	backupInfo := models.NewBackupInfo(filename, filepath, userID, description)
	if err := s.backupRepo.Create(ctx, backupInfo); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	// バックアップを非同期で実行
	go func() {
		var err error
		defer func() {
			if err != nil {
				backupInfo.Fail()
				s.backupRepo.Update(context.Background(), backupInfo)
			}
		}()

		// データベースタイプに応じてバックアップコマンドを実行
		switch s.dbType {
		case "sqlite":
			err = s.backupSQLite(filepath)
		case "postgres":
			err = s.backupPostgreSQL(filepath)
		default:
			err = fmt.Errorf("unsupported database type: %s", s.dbType)
		}

		if err == nil {
			// ファイルサイズを取得
			if fileInfo, statErr := os.Stat(filepath); statErr == nil {
				backupInfo.Complete(fileInfo.Size())
			} else {
				backupInfo.Complete(0)
			}
			s.backupRepo.Update(context.Background(), backupInfo)
		}
	}()

	return backupInfo, nil
}

// backupSQLite はSQLiteデータベースのバックアップを作成します
func (s *BackupService) backupSQLite(outputPath string) error {
	// SQLiteの場合、単純にファイルをコピー
	cmd := exec.Command("cp", s.dbConnection, outputPath)
	return cmd.Run()
}

// backupPostgreSQL はPostgreSQLデータベースのバックアップを作成します
func (s *BackupService) backupPostgreSQL(outputPath string) error {
	cmd := exec.Command("pg_dump", s.dbConnection, "-f", outputPath)
	return cmd.Run()
}

// GetBackups はバックアップ一覧を取得します
func (s *BackupService) GetBackups(ctx context.Context, page, limit int) ([]*models.BackupInfo, int, error) {
	return s.backupRepo.List(ctx, page, limit)
}

// GetBackup はバックアップを取得します
func (s *BackupService) GetBackup(ctx context.Context, id int64) (*models.BackupInfo, error) {
	return s.backupRepo.GetByID(ctx, id)
}

// DeleteBackup はバックアップを削除します
func (s *BackupService) DeleteBackup(ctx context.Context, id int64) error {
	backup, err := s.backupRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// ファイルシステムからファイルを削除
	if err := os.Remove(backup.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}

	// データベースからレコードを削除
	return s.backupRepo.Delete(ctx, id)
}

// RestoreBackup はバックアップからデータベースを復元します
func (s *BackupService) RestoreBackup(ctx context.Context, id int64) error {
	backup, err := s.backupRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if backup.Status != "completed" {
		return fmt.Errorf("backup is not completed")
	}

	// データベースタイプに応じて復元コマンドを実行
	switch s.dbType {
	case "sqlite":
		return s.restoreSQLite(backup.FilePath)
	case "postgres":
		return s.restorePostgreSQL(backup.FilePath)
	default:
		return fmt.Errorf("unsupported database type: %s", s.dbType)
	}
}

// restoreSQLite はSQLiteデータベースを復元します
func (s *BackupService) restoreSQLite(backupPath string) error {
	cmd := exec.Command("cp", backupPath, s.dbConnection)
	return cmd.Run()
}

// restorePostgreSQL はPostgreSQLデータベースを復元します
func (s *BackupService) restorePostgreSQL(backupPath string) error {
	cmd := exec.Command("psql", s.dbConnection, "-f", backupPath)
	return cmd.Run()
}

// CleanOldBackups は古いバックアップを削除します
func (s *BackupService) CleanOldBackups(ctx context.Context, retentionDays int) error {
	return s.backupRepo.DeleteOldBackups(ctx, retentionDays)
}
