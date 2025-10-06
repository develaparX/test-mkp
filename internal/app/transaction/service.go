package transaction

type TransactionService struct {
	Repo *TransactionRepo
}

func NewTransactionService(repo *TransactionRepo) *TransactionService {
	return &TransactionService{Repo: repo}
}

func (s *TransactionService) GetTransactionSummary() ([]TransactionSummary, error) {
	return s.Repo.GetTransactionSummary()
}
