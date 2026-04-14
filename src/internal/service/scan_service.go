package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
)

type ScanService interface {
	ScanDailyLog(ctx context.Context, pondId int, month string, files []*multipart.FileHeader, username string) (*dto.ScanDailyLogResponse, error)
}

type scanService struct {
	claudeClient   ClaudeClient
	scanLogRepo    repository.ScanLogRepository
	activePondRepo repository.ActivePondRepository
	pondRepo       repository.PondRepository
	farmRepo       repository.FarmRepository
	uploadPath     string
}

func NewScanService(
	claudeClient ClaudeClient,
	scanLogRepo repository.ScanLogRepository,
	activePondRepo repository.ActivePondRepository,
	pondRepo repository.PondRepository,
	farmRepo repository.FarmRepository,
	conf *config.Config,
) ScanService {
	return &scanService{
		claudeClient:   claudeClient,
		scanLogRepo:    scanLogRepo,
		activePondRepo: activePondRepo,
		pondRepo:       pondRepo,
		farmRepo:       farmRepo,
		uploadPath:     conf.App.DailyLogUploadPath,
	}
}

type claudeExtractedData struct {
	Entries []struct {
		Day            int      `json:"day"`
		FreshMorning   *float64 `json:"freshMorning"`
		FreshEvening   *float64 `json:"freshEvening"`
		PelletMorning  *float64 `json:"pelletMorning"`
		PelletEvening  *float64 `json:"pelletEvening"`
		DeathFishCount *int     `json:"deathFishCount"`
	} `json:"entries"`
	Confidence []struct {
		Day            int     `json:"day"`
		FreshMorning   float64 `json:"freshMorning"`
		FreshEvening   float64 `json:"freshEvening"`
		PelletMorning  float64 `json:"pelletMorning"`
		PelletEvening  float64 `json:"pelletEvening"`
		DeathFishCount float64 `json:"deathFishCount"`
	} `json:"confidence"`
	Notes string `json:"notes"`
}

func (s *scanService) ScanDailyLog(ctx context.Context, pondId int, month string, files []*multipart.FileHeader, username string) (*dto.ScanDailyLogResponse, error) {
	// Validate month format
	_, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("invalid month format, expected YYYY-MM: %w", err))
	}

	// Validate pond access
	data, err := s.pondRepo.GetByIDWithFarmAndActivePond(ctx, pondId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if data == nil || data.Pond == nil {
		return nil, errors.ErrPondNotFound
	}
	if data.ClientId == 0 {
		return nil, errors.ErrFarmNotFound
	}
	ok, err := utils.CanAccessClient(ctx, data.ClientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if !ok {
		return nil, errors.ErrAuthPermissionDenied
	}

	// Get active pond — try active first, then latest (for closed ponds)
	activePond := data.ActivePond
	if activePond == nil {
		activePond, err = s.activePondRepo.GetLatestByPondID(ctx, pondId)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		if activePond == nil {
			return nil, errors.ErrPondNotActive
		}
	}

	// Save images to disk and encode to base64
	var imagesBase64 []string
	var mimeTypes []string
	var imagePaths []string

	scanDir := filepath.Join(s.uploadPath, "scans", fmt.Sprintf("pond_%d", pondId), month)
	if err := os.MkdirAll(scanDir, 0755); err != nil {
		return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to create scan directory: %w", err))
	}

	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to open file %s: %w", fh.Filename, err))
		}

		fileBytes, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to read file %s: %w", fh.Filename, err))
		}

		// Save to disk
		ts := time.Now().UnixMilli()
		saveName := fmt.Sprintf("%d_%s", ts, fh.Filename)
		savePath := filepath.Join(scanDir, saveName)
		if err := os.WriteFile(savePath, fileBytes, 0644); err != nil {
			return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to save file: %w", err))
		}
		imagePaths = append(imagePaths, savePath)

		imagesBase64 = append(imagesBase64, EncodeImageToBase64(fileBytes))
		mimeTypes = append(mimeTypes, DetectMimeType(fh.Filename))
	}

	// Build prompt
	prompt := buildScanPrompt(month)

	// Call Claude Vision API
	claudeResp, err := s.claudeClient.SendVisionRequest(imagesBase64, mimeTypes, prompt)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(fmt.Errorf("AI scan failed: %w", err))
	}

	// Extract text from response
	var responseText string
	for _, c := range claudeResp.Content {
		if c.Type == "text" {
			responseText = c.Text
			break
		}
	}

	// Parse JSON from response (strip markdown code fences if present)
	jsonStr := strings.TrimSpace(responseText)
	if strings.HasPrefix(jsonStr, "```") {
		lines := strings.Split(jsonStr, "\n")
		// Remove first and last lines (```json and ```)
		if len(lines) > 2 {
			lines = lines[1 : len(lines)-1]
		}
		jsonStr = strings.Join(lines, "\n")
	}

	var extracted claudeExtractedData
	if err := json.Unmarshal([]byte(jsonStr), &extracted); err != nil {
		return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to parse AI response as JSON: %w", err))
	}

	// Save scan log
	imagePathsJSON, _ := json.Marshal(imagePaths)
	extractedJSON, _ := json.Marshal(extracted.Entries)
	confidenceJSON, _ := json.Marshal(extracted.Confidence)

	scanLog := &model.ScanLog{
		ActivePondId:     activePond.Id,
		Month:            month,
		ImagePaths:       imagePathsJSON,
		RawResponse:      responseText,
		ExtractedData:    extractedJSON,
		ConfidenceScores: confidenceJSON,
		Status:           "pending_review",
	}
	scanLog.CreatedBy = username
	scanLog.UpdatedBy = username

	if err := s.scanLogRepo.Create(ctx, scanLog); err != nil {
		return nil, errors.ErrGeneric.Wrap(fmt.Errorf("failed to save scan log: %w", err))
	}

	// Build response
	resp := &dto.ScanDailyLogResponse{
		ScanLogId: scanLog.Id,
		Month:     month,
		ImageUrls: imagePaths,
		Notes:     extracted.Notes,
	}

	for _, e := range extracted.Entries {
		resp.Entries = append(resp.Entries, dto.ScanEntry{
			Day:            e.Day,
			FreshMorning:   e.FreshMorning,
			FreshEvening:   e.FreshEvening,
			PelletMorning:  e.PelletMorning,
			PelletEvening:  e.PelletEvening,
			DeathFishCount: e.DeathFishCount,
		})
	}

	for _, c := range extracted.Confidence {
		resp.Confidence = append(resp.Confidence, dto.ScanConfidence{
			Day:            c.Day,
			FreshMorning:   c.FreshMorning,
			FreshEvening:   c.FreshEvening,
			PelletMorning:  c.PelletMorning,
			PelletEvening:  c.PelletEvening,
			DeathFishCount: c.DeathFishCount,
		})
	}

	return resp, nil
}

func buildScanPrompt(month string) string {
	return fmt.Sprintf(`Extract daily fish feeding data from this handwritten paper for month: %s.

The paper records daily data with these columns:
- Day number (วันที่)
- Fresh feed morning (สดเช้า) - kg
- Fresh feed evening (สดเย็น) - kg
- Pellet feed morning (เม็ดเช้า) - kg
- Pellet feed evening (เม็ดเย็น) - kg
- Death fish count (ปลาตาย) - integer (optional)

Return JSON in this exact format:
{
  "entries": [
    {
      "day": 1,
      "freshMorning": 50.0,
      "freshEvening": 30.0,
      "pelletMorning": null,
      "pelletEvening": null,
      "deathFishCount": 2
    }
  ],
  "confidence": [
    {
      "day": 1,
      "freshMorning": 0.95,
      "freshEvening": 0.90,
      "pelletMorning": 0.0,
      "pelletEvening": 0.0,
      "deathFishCount": 0.85
    }
  ],
  "notes": "observations about image quality or unrecognized content"
}

Rules:
- Convert Thai numerals ๐-๙ to Arabic 0-9
- Feed amounts are typically 0-200 kg per feeding
- Death count is typically 0-100 per day
- If value is unclear, use best guess with confidence below 0.7
- If a field is not present on the paper, set it to null with confidence 0.0
- Only include days that have data on the paper
- Return ONLY valid JSON, no markdown formatting`, month)
}
