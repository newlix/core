package server

import "context"

// Server implements the todo RPC methods.
type Server struct{}

func (s *Server) AddItem(ctx context.Context, in AddItemInput) (AddItemOutput, error) {
	return AddItemOutput{}, nil
}

func (s *Server) GetItems(ctx context.Context, in GetItemsInput) (GetItemsOutput, error) {
	return GetItemsOutput{}, nil
}

func (s *Server) RemoveItem(ctx context.Context, in RemoveItemInput) (RemoveItemOutput, error) {
	return RemoveItemOutput{}, nil
}
