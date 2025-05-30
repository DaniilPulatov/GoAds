package utils

import (
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	_ = os.Setenv("JWT_SECRET_KEY", "testsecret")
	token, err := GenerateToken("user-1", 10)
	if err != nil {
		t.Errorf("ошибка при генерации токена: %v", err)
	}
	if token == "" {
		t.Error("токен не должен быть пустым")
	}
}

func TestIsValidPhone(t *testing.T) {
	t.Run("valid phone number", func(t *testing.T) {
		phone := "+998910000000"
		if !IsValidPhone(phone) {
			t.Errorf("expected phone %s to be valid", phone)
		}
	})

	t.Run("invalid phone number", func(t *testing.T) {
		phone := "12345"
		if IsValidPhone(phone) {
			t.Errorf("expected phone %s to be invalid", phone)
		}
	})

	t.Run("empty phone number", func(t *testing.T) {
		phone := ""
		if IsValidPhone(phone) {
			t.Errorf("expected empty phone to be invalid")
		}
	})
}

func TestIsValidPassword(t *testing.T) {
	t.Run("valid password", func(t *testing.T) {
		password := "StrongPass123!"
		if !IsValidPassword(password) {
			t.Errorf("expected password %s to be valid", password)

		}
	})

	t.Run("invalid password", func(t *testing.T) {
		password := "short"
		if IsValidPassword(password) {
			t.Errorf("expected password %s to be invalid", password)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		password := ""
		if IsValidPassword(password) {
			t.Errorf("expected empty password to be invalid")
		}
	})
}

func TestIsSafeLogPath(t *testing.T) {
	t.Run("valid log path", func(t *testing.T) {
		path := "storage/logs/app.log"
		if !IsSafeLogPath(path) {
			t.Errorf("ожидалось, что путь %s будет безопасным", path)
		}
	})

	t.Run("invalid log path (directory traversal)", func(t *testing.T) {
		path := "../logs/app.log"
		if IsSafeLogPath(path) {
			t.Errorf("ожидалось, что путь %s будет небезопасным", path)
		}
	})

	t.Run("invalid log path (outside logs dir)", func(t *testing.T) {
		path := "/tmp/app.log"
		if IsSafeLogPath(path) {
			t.Errorf("ожидалось, что путь %s будет небезопасным", path)
		}
	})

	t.Run("empty path", func(t *testing.T) {
		path := ""
		if IsSafeLogPath(path) {
			t.Errorf("ожидалось, что пустой путь будет небезопасным")
		}
	})
}
