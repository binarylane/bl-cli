/*
Copyright 2018 The Doctl Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bl

import (
	"context"
	"github.com/binarylane/go-binarylane"
)

// Balance is a wrapper for binarylane.Balance.
type Balance struct {
	*binarylane.Balance
}

// BalanceService is an interface for interacting with BinaryLane's balance api.
type BalanceService interface {
	Get() (*Balance, error)
}

type balanceService struct {
	client *binarylane.Client
}

var _ BalanceService = &balanceService{}

// NewBalanceService builds an BalanceService instance.
func NewBalanceService(client *binarylane.Client) BalanceService {
	return &balanceService{
		client: client,
	}
}

func (as *balanceService) Get() (*Balance, error) {
	binarylaneBalance, _, err := as.client.Balance.Get(context.TODO())
	if err != nil {
		return nil, err
	}

	balance := &Balance{Balance: binarylaneBalance}
	return balance, nil
}
