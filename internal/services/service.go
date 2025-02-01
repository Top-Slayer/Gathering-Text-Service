package services

import "Text-Gathering-Service/internal/repository"

func CheckText(text string) bool {
	repo := repository.New()
	repo.StoreIntoDB(text)
	repo.Close()

	return false
}

// filter receive lao langage only
