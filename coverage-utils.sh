#!/bin/sh

# Barbara S1 Req5 - Employees, S2 Req6 - Purchase Order
if [[ $1 == "Barbara" ]]; then
    echo "Starting the tests for Barbara ..."
    go test -coverprofile=coverage.out ./internal/employee/... ./internal/purchase_order/... ./cmd/server/handler/... && go tool cover -html=coverage.out
fi

# Debora S1 Req1 - Sellers, S2 Req3 - Product Batch
if [[ $1 == "Debora" ]]; then
    echo "Starting the tests for Debora ..."
    go test -coverprofile=coverage.out ./internal/seller/... ./internal/product_batch/... ./cmd/server/handler/... && go tool cover -html=coverage.out
fi
# Pedro S1 Req6 - Buyers, S2 Req5 - Inbound Order
if [[ $1 == "Pedro" ]]; then
    echo "Starting the tests for Pedro ..."
    go test -coverprofile=coverage.out ./internal/buyer/... ./internal/inbound_order/... ./cmd/server/handler/... && go tool cover -html=coverage.out
fi

# Matheus S1 Req4 - Products, S2 Req2 - Carry
if [[ $1 == "Pedro" ]]; then
    echo "Starting the tests for Pedro ..."
    go test -coverprofile=coverage.out ./internal/product/... ./internal/carry/... ./cmd/server/handler/... && go tool cover -html=coverage.out
fi

# João S1 Req2 - Warehouse, S2 Req4 - Product Record
if [[ $1 == "João" ]]; then
    echo "Starting the tests for João ..."
    go test -coverprofile=coverage.out ./internal/warehouse/... ./internal/product_record/... ./cmd/server/handler/... && go tool cover -html=coverage.out
fi

# Yuri S1 Req3 - Sections, S2 Req1 - Locality
if [[ $1 == "Yuri" ]]; then
    echo "Starting the tests for Yuri ..."
    go test -coverprofile=coverage.out ./internal/section/... ./internal/locality/... ./cmd/server/handler/... -run 'Unit' && go tool cover -html=coverage.out
fi
