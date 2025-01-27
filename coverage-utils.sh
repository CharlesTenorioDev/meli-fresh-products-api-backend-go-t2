#!/bin/sh

# Barbara S1 Req5 - Employees, S2 Req6 - Purchase Order
if [[ $1 == "Barbara" ]]; then
    echo "Starting the tests for Barbara ..."
    go test -coverprofile=coverage.out ./internal/employee/... ./internal/purchase_order/... ./cmd/server/handler/... 

    echo "mode: set" > ./coverage_filter.out
    sed '/locality_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/employees_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/purchase_order\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/employee\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi

# Debora S1 Req1 - Sellers, S2 Req3 - Product Batch
if [[ $1 == "Debora" ]]; then
    echo "Starting the tests for Debora ..."
    go test -coverprofile=coverage.out ./internal/seller/... ./internal/product_batch/... ./cmd/server/handler/... 

    echo "mode: set" > ./coverage_filter.out
    sed '/sellers_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/product_batch_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/seller\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/product_batch\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi

# João S1 Req2 - Warehouse, S2 Req4 - Product Record
if [[ $1 == "João" ]]; then
    echo "Starting the tests for João ..."
    go test -coverprofile=coverage.out ./internal/warehouse/... ./internal/product_record/... ./cmd/server/handler/...

    echo "mode: set" > ./coverage_filter.out
    sed '/warehouse_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/product_records_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/warehouse\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/product_record\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi

# Yuri S1 Req3 - Sections, S2 Req1 - Locality
if [[ $1 == "Yuri" ]]; then
    echo "Starting the tests for Yuri ..."
    go test -coverprofile=coverage.out ./internal/section/... ./internal/locality/... ./cmd/server/handler/... -run 'Unit'

    echo "mode: set" > ./coverage_filter.out
    sed '/sections_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/locality_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/section\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/locality\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi



# Matheus S1 Req4 - Products, S2 Req2 - Carry
if [[ $1 == "Matheus" ]]; then
    echo "Starting the tests for Matheus ..."
    go test -coverprofile=coverage.out ./internal/product/... ./internal/carry/... ./cmd/server/handler/...

    echo "mode: set" > ./coverage_filter.out
    sed '/product_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/carry_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/product\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/carry\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi

# Pedro S1 Req6 - Buyers, S2 Req5 - Inbound Order
if [[ $1 == "Pedro" ]]; then
    echo "Starting the tests for Pedro ..."
    go test -coverprofile=coverage.out ./internal/buyer/... ./internal/inbound_order/... ./cmd/server/handler/...

    echo "mode: set" > ./coverage_filter.out
    sed '/buyer_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/inbound_orders_handler/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/buyer\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    sed '/internal\/inbound_order\/service_default.go/!d' ./coverage.out >> ./coverage_filter.out
    go tool cover -html=coverage_filter.out
fi